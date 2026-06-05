package ac_calendar_mq

import (
	"GO1/database/mysql/calendar_mysql"
	"GO1/global"
	"GO1/models/mq_models"

	"github.com/goccy/go-json"
)

func ConsumeACCalendar() error {

	msgs, err := global.MQChannel.Consume(
		"ac_calendar", // queue
		"",            // consumer
		false,         // autoAck，建议 false，手动 ack
		false,         // exclusive
		false,         // noLocal
		false,         // noWait
		nil,           // args
	)
	if err != nil {
		return err
	}

	for msg := range msgs {
		var data mq_models.ACProblem

		if err := json.Unmarshal(msg.Body, &data); err != nil {
			global.Logger.Error("json unmarshal failed:", err)

			// 解析失败，丢弃这条消息
			msg.Nack(false, false)
			continue
		}

		// TODO: 处理你的业务逻辑

		global.Logger.Info("收到消息:", data)
		if err := calendar_mysql.SaveACCount(&data); err != nil {
			global.Logger.Error("save ac count failed:", err)

			// 		// 消费失败，拒绝消息
			//		// requeue=true 表示重新入队
			//		// requeue=false 表示丢弃，或者进入死信队列
			_ = msg.Nack(false, true)
			continue
		}

		// 处理成功后确认
		if err := msg.Ack(false); err != nil {
			global.Logger.Error("ack failed:", err)
		}
	}

	return nil
}
