package competition_service

import (
	"GO1/models/problem_model"
	"math"
)

func summarizeCompetitionRunResult(runResult []problem_model.RunResult) (string, float64, bool) {
	totalCount := len(runResult)
	if totalCount == 0 {
		return "Judge Failed", 0, false
	}

	state := "Accepted"
	acCount := totalCount
	for _, result := range runResult {
		if !result.Passed {
			acCount--
			if result.Error != "" {
				state = result.Error
			} else {
				state = "Wrong Answer"
			}
		}
	}

	score := math.Round(float64(acCount)/float64(totalCount)*10000) / 100
	return state, score, acCount == totalCount
}
