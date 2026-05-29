package competition_service

import (
	"GO1/models/problem_model"
	"GO1/pkg/constants"
	"math"
)

func summarizeCompetitionRunResult(runResult []problem_model.RunResult) (string, float64, bool) {
	totalCount := len(runResult)
	if totalCount == 0 {
		return constants.JudgeStatusSystemError, 0, false
	}

	state := constants.JudgeStatusAccepted
	acCount := totalCount
	for _, result := range runResult {
		if !result.Passed {
			acCount--
			if result.Error != "" {
				state = NormalizeJudgeStatus(result.Error)
			} else {
				state = constants.JudgeStatusWrongAnswer
			}
		}
	}

	score := math.Round(float64(acCount)/float64(totalCount)*10000) / 100
	return state, score, acCount == totalCount
}
