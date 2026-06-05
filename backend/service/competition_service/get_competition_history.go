package competition_service

import (
	"GO1/database/mysql/competition_mysql"
	"GO1/middlewares/response"
	"GO1/models/competition_model"
	"GO1/pkg/constants"
)

func GetCompetitionHistory(req *competition_model.GetCompetitionHistoryReq) (resp response.Response) {
	results, err := competition_mysql.GetCompetitionHistory(req.UserID)
	if err != nil {
		resp.Code = 1
		resp.Message = constants.CompetitionSubmitMessageDataQueryError
		return
	}

	respData := make([]competition_model.GetCompetitionHistoryResp, 0, len(results))
	for _, result := range results {
		respData = append(respData, competition_model.GetCompetitionHistoryResp{
			RankScore:       result.RatingAfter,
			RankScoreChange: result.RatingAfter - result.RatingBefore,
			Rank:            result.CompetitionRank,
			CompetitionId:   result.CompetitionId,
			CompetitionName: result.CompetitionName,
		})
	}

	resp.Data = respData
	return
}
