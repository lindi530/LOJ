package competition_model

import "time"

type CompetitionResult struct {
	Id              int64     `json:"id"`
	CompetitionId   int64     `json:"competition_id"`
	CompetitionName string    `json:"competition_name"`
	UserId          int64     `json:"user_id"`
	CompetitionRank int       `json:"competition_rank"`
	SolvedCount     int       `json:"solved_count"`
	Penalty         int       `json:"penalty"`
	RatingBefore    int       `json:"rating_before"`
	RatingAfter     int       `json:"rating_after"`
	CreatedAt       time.Time `json:"created_at"`
}

/*
{
  "id": 0,
  "competition_id": 0,
  "user_id": 0,
  "competition_rank": 0,
  "solved_count": 0,
  "penalty": 0,
  "rating_before": 0,
  "rating_after": 0,
  "created_at": ""
}
*/
