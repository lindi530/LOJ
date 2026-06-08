package competition_model

type EnterCompetitionReq struct {
	CompetitionID       int64  `json:"competition_id" uri:"competition_id"`
	CompetitionPassword string `json:"password" uri:"password"`
}
type EnterCompetitionResp struct{}
