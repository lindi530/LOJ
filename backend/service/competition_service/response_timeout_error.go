package competition_service

import "GO1/pkg/constants"

type responseTimeoutError struct{}

func (responseTimeoutError) Error() string {
	return constants.CompetitionSubmitMessageJudgeTimeout
}
