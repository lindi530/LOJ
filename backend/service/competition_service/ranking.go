package competition_service

import (
	"GO1/database/mysql/competition_mysql"
	"GO1/middlewares/response"
	"GO1/models/competition_model"
	"GO1/models/user_model"
	"GO1/pkg/constants"
)

func GetRanking(req *competition_model.GetRankingReq) (resp response.Response) {
	ranking, err := competition_mysql.GetRanking()
	if err != nil {
		resp.Code = 1
		resp.Message = constants.CompetitionSubmitMessageDataQueryError
		return
	}

	respData := make([]competition_model.GetRankingResp, 0, len(ranking))
	for _, item := range ranking {
		if item.UserAvatar != "" {
			item.UserAvatar = user_model.GetAvatarPath(item.UserAvatar)
		}
		respData = append(respData, item)
	}

	resp.Data = respData
	return
}
