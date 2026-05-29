package competition_service

import (
	"GO1/pkg/constants"
	"strings"
)

func NormalizeJudgeStatus(state string) string {
	state = strings.TrimSpace(state)
	switch state {
	case constants.JudgeStatusPending,
		constants.JudgeStatusCompiling,
		constants.JudgeStatusRunning,
		constants.JudgeStatusAccepted,
		constants.JudgeStatusWrongAnswer,
		constants.JudgeStatusTimeLimitExceeded,
		constants.JudgeStatusMemoryLimitExceeded,
		constants.JudgeStatusRuntimeError,
		constants.JudgeStatusCompileError,
		constants.JudgeStatusOutputLimitExceeded,
		constants.JudgeStatusSystemError:
		return state
	}

	switch {
	case strings.HasPrefix(state, constants.JudgeStatusRuntimeError):
		return constants.JudgeStatusRuntimeError
	case strings.HasPrefix(state, constants.JudgeStatusCompileError):
		return constants.JudgeStatusCompileError
	}

	return constants.JudgeStatusSystemError
}
