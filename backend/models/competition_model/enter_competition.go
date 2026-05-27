package competition_model

type EnterCompetitionReq struct {
	CompetitionID int64 `json:"competition_id" uri:"competition_id"`
}
type EnterCompetitionResp struct{}
