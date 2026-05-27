package competition_model

type GetCompetitionProblemsReq struct {
	CompetitionID int64 `json:"competition_id" uri:"competition_id" binding:"required"`
}

type GetCompetitionProblemsResp struct {
}
