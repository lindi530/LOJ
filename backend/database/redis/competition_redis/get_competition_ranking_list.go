package competition_redis

import (
	"context"
	"sort"
	"strconv"
	"strings"
	"time"

	"GO1/global"
	"GO1/models/competition_model"
)

func GetCompetitionRankingList(competitionID int64) (data []competition_model.CompetitionRankingListItem, err error) {
	ctx := context.Background()
	rankingData, err := global.RedisClient.
		ZRevRangeWithScores(ctx, competitionRankingListKey(competitionID), 0, -1).
		Result()
	if err != nil {
		return nil, err
	}
	if len(rankingData) == 0 {
		return []competition_model.CompetitionRankingListItem{}, nil
	}

	userIDs := make([]string, 0, len(rankingData))
	for _, item := range rankingData {
		userIDs = append(userIDs, item.Member.(string))
	}

	userNames, err := global.RedisClient.
		HMGet(ctx, competitionRankingUserNameKey(competitionID), userIDs...).
		Result()
	if err != nil {
		return nil, err
	}

	data = make([]competition_model.CompetitionRankingListItem, 0, len(rankingData))
	for i, item := range rankingData {
		userIDStr := item.Member.(string)
		userID, err := strconv.ParseInt(userIDStr, 10, 64)
		if err != nil {
			return nil, err
		}

		score := int64(item.Score)
		solvedCount := int((score + competitionRankingScoreBase - 1) / competitionRankingScoreBase)
		penalty := int(int64(solvedCount)*competitionRankingScoreBase - score)

		userName := ""
		if i < len(userNames) && userNames[i] != nil {
			userName = userNames[i].(string)
		}

		submissions, err := getCompetitionRankingSubmissionInfos(ctx, competitionID, userID)
		if err != nil {
			return nil, err
		}

		data = append(data, competition_model.CompetitionRankingListItem{
			UserID:      userID,
			UserName:    userName,
			SolvedCount: solvedCount,
			Penalty:     penalty,
			Submissions: submissions,
		})
	}

	return data, nil
}

func getCompetitionRankingSubmissionInfos(ctx context.Context, competitionID int64, userID int64) ([]competition_model.CompetitionRankingSubmissionInfo, error) {
	userProblemFields, err := global.RedisClient.
		HGetAll(ctx, competitionRankingUserProblemKey(competitionID, userID)).
		Result()
	if err != nil {
		return nil, err
	}
	if len(userProblemFields) == 0 {
		return []competition_model.CompetitionRankingSubmissionInfo{}, nil
	}

	problemFields := make(map[string]map[string]string)
	for field, value := range userProblemFields {
		problemID, infoField, ok := strings.Cut(field, ":")
		if !ok {
			continue
		}
		fields := problemFields[problemID]
		if fields == nil {
			fields = make(map[string]string)
			problemFields[problemID] = fields
		}
		fields[infoField] = value
	}

	submissions := make([]competition_model.CompetitionRankingSubmissionInfo, 0, len(problemFields))
	for _, fields := range problemFields {
		if len(fields) == 0 {
			continue
		}

		info := competition_model.CompetitionRankingSubmissionInfo{
			ProblemNumber: fields["problem_number"],
			IsAC:          parseRedisBool(fields["is_ac"]),
			WrongCount:    parseRedisInt(fields["wrong_count"]),
		}
		if info.IsAC && fields["first_ac_time"] != "" {
			firstACTime, err := time.Parse(time.RFC3339Nano, fields["first_ac_time"])
			if err != nil {
				return nil, err
			}
			info.FirstACTime = &firstACTime
		}
		submissions = append(submissions, info)
	}

	sort.SliceStable(submissions, func(i, j int) bool {
		return submissions[i].ProblemNumber < submissions[j].ProblemNumber
	})

	return submissions, nil
}

func parseRedisBool(value string) bool {
	return value == "1" || value == "true"
}

func parseRedisInt(value string) int {
	number, err := strconv.Atoi(value)
	if err != nil || number < 0 {
		return 0
	}

	return number
}
