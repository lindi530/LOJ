package competition_service

func failedCompetitionSubmitResult(message string) competitionSubmitResult {
	return competitionSubmitResult{
		Code:    1,
		Message: message,
	}
}
