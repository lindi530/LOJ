package competition_service

import (
	"GO1/database/mysql/competition_mysql"
	"GO1/database/redis/competition_redis"
	"GO1/middlewares/response"
	"GO1/models/competition_model"
	"GO1/pkg/constants"
	"sort"
	"time"
)

func GetCompetitionRankingList(req *competition_model.GetCompetitionRankingListReq) (resp response.Response) {
	competition, err := competition_mysql.GetCompetitionTime(req.CompetitionID)
	if err != nil {
		resp.Code = 1
		resp.Message = constants.CompetitionSubmitMessageDataQueryError
		return
	}

	now := time.Now()
	if competition.StartTime.IsZero() || now.Before(competition.StartTime) {
		resp.Code = 1
		resp.Message = constants.CompetitionSubmitMessageNotStarted
		return
	}

	var rankingList []competition_model.CompetitionRankingListItem
	if !competition.EndTime.IsZero() && now.After(competition.EndTime.Add(time.Hour)) {
		rankingList, err = competition_mysql.GetCompetitionRankingList(req.CompetitionID, competition.StartTime)
	} else {
		rankingList, err = competition_redis.GetCompetitionRankingList(req.CompetitionID)
	}
	if err != nil {
		resp.Code = 1
		resp.Message = constants.CompetitionSubmitMessageDataQueryError
		return
	}
	if rankingList == nil {
		rankingList = []competition_model.CompetitionRankingListItem{}
	}

	problemNumbers, err := getCompetitionProblemNumbers(req.CompetitionID, now)
	if err != nil {
		resp.Code = 1
		resp.Message = constants.CompetitionSubmitMessageDataQueryError
		return
	}

	sortAndRankCompetitionRankingList(rankingList)

	resp.Data = competition_model.GetCompetitionRankingListResp{
		CompetitionID:   competition.ID,
		CompetitionName: competition.Name,
		RankingList:     rankingList,
		ProblemNumbers:  problemNumbers,
	}
	return
}

func sortAndRankCompetitionRankingList(rankingList []competition_model.CompetitionRankingListItem) {
	sort.SliceStable(rankingList, func(i, j int) bool {
		if rankingList[i].SolvedCount != rankingList[j].SolvedCount {
			return rankingList[i].SolvedCount > rankingList[j].SolvedCount
		}
		if rankingList[i].Penalty != rankingList[j].Penalty {
			return rankingList[i].Penalty < rankingList[j].Penalty
		}
		return rankingList[i].UserID < rankingList[j].UserID
	})

	rank := 0
	for i := range rankingList {
		if i == 0 ||
			rankingList[i].SolvedCount != rankingList[i-1].SolvedCount ||
			rankingList[i].Penalty != rankingList[i-1].Penalty {
			rank = i + 1
		}
		rankingList[i].CompetitionRank = rank
	}
}
