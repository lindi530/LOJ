package competition_model

import "time"

type GetCompetitionRankingListReq struct {
	CompetitionID int64 `json:"competition_id" uri:"competition_id" binding:"required"`
}

type GetCompetitionRankingListResp struct {
	CompetitionID   int64                        `json:"competition_id"`
	CompetitionName string                       `json:"competition_name"`
	RankingList     []CompetitionRankingListItem `json:"ranking_list"`
	ProblemNumbers  []string                     `json:"problem_numbers"`
}

type CompetitionRankingListItem struct {
	CompetitionRank int                                `json:"competition_rank"`
	UserID          int64                              `json:"user_id"`
	UserName        string                             `json:"user_name"`
	SolvedCount     int                                `json:"solved_count"`
	Penalty         int                                `json:"penalty"`
	Submissions     []CompetitionRankingSubmissionInfo `json:"submissions"`
}

type CompetitionRankingSubmissionInfo struct {
	ProblemNumber string     `json:"problem_number"`
	IsAC          bool       `json:"is_ac"`
	WrongCount    int        `json:"wrong_count"`
	FirstACTime   *time.Time `json:"first_ac_time,omitempty"`
}
