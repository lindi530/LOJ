package competition_model

type HasEnterCompetitionReq struct {
	CompetitionID int64 `json:"competition_id" uri:"competition_id"`
}

type HasEnterCompetitionResp struct {
	HasEnter bool `json:"has_enter"`
}
