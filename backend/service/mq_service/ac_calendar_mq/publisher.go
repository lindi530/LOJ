package ac_calendar_mq

import (
	"GO1/database/redis/calendar_redis"
	"GO1/global"
	"GO1/models/mq_models"
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/goccy/go-json"
	"github.com/streadway/amqp"
)

func SendPendingACProblems(year int, month int, day int) {
	ctx := context.Background()

	date := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Local)
	dateStr := date.Format("2006-01-02")
	key := fmt.Sprintf("pending:user:problem:%s", dateStr)

	var cursor uint64
	const count int64 = 100

	// 记录每个用户 AC 了多少题
	userACCount := make(map[int64]int)

	// 记录每个用户对应的 Redis member，发送成功后用于删除
	userMembers := make(map[int64][]string)

	for {
		members, nextCursor, err := calendar_redis.ReadAllACProblem(ctx, key, cursor, count)
		if err != nil {
			fmt.Println("SScan redis failed:", err)
			return
		}

		for _, member := range members {
			userID, err := parseUserIDFromMember(member)
			if err != nil {
				fmt.Println("invalid redis member:", member, err)

				// 格式错误的数据可以删除，避免之后一直扫到
				_ = calendar_redis.RemoveACProblem(ctx, key, member)
				continue
			}

			userACCount[userID]++
			userMembers[userID] = append(userMembers[userID], member)
		}

		cursor = nextCursor
		if cursor == 0 {
			break
		}
	}

	for userID, acCount := range userACCount {
		err := SendACCountToMQ(userID, acCount, date)
		if err != nil {
			// 发送失败，不删除 Redis，下次继续重试
			fmt.Println("send mq failed:", err)
			continue
		}

		// 这个用户的统计发送成功后，删除他的所有记录
		for _, member := range userMembers[userID] {
			if err := calendar_redis.RemoveACProblem(ctx, key, member); err != nil {
				fmt.Println("SRem redis failed:", err)
			}
		}
	}
}

func SendACCountToMQ(userID int64, acCount int, date time.Time) error {
	msg := mq_models.ACProblem{
		UserID:  userID,
		AcCount: acCount,
		Date:    date,
	}

	body, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	// TODO: 替换成你的消息队列发送逻辑
	return global.MQChannel.Publish(
		"",            // exchange
		"ac_calendar", // routing key (queue name)
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}

func parseUserIDFromMember(member string) (int64, error) {
	parts := strings.Split(member, ":")
	if len(parts) != 2 {
		return 0, fmt.Errorf("invalid member format: %s", member)
	}

	userID, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return 0, err
	}

	return userID, nil
}
