package competition_service

type responseTimeoutError struct{}

func (responseTimeoutError) Error() string {
	return "judge timeout"
}
