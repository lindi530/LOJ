package constants

import (
	"errors"
	"time"
)

var (
	LoginUserIDKey   = "login_user_id"
	LoginUserNameKey = "login_user_name"
	NoUserID         = int64(-1)
	NoUserName       = ""
	NotLogin         = errors.New("not login")
	AccessUserIDKey  = "access_user_id"
)

const (
	CompetitionSubmitQueue             = "competition_submit_queue"
	CompetitionSubmitConsumerTagPrefix = "competition-submit-"
	CompetitionSubmitContentTypeJSON   = "application/json"
	CompetitionSubmitTimeout           = 2 * time.Minute
)

const (
	JudgeStatusPending             = "Pending"
	JudgeStatusCompiling           = "Compiling"
	JudgeStatusRunning             = "Running"
	JudgeStatusAccepted            = "Accepted"
	JudgeStatusWrongAnswer         = "Wrong Answer"
	JudgeStatusTimeLimitExceeded   = "Time Limit Exceeded"
	JudgeStatusMemoryLimitExceeded = "Memory Limit Exceeded"
	JudgeStatusRuntimeError        = "Runtime Error"
	JudgeStatusCompileError        = "Compile Error"
	JudgeStatusOutputLimitExceeded = "Output Limit Exceeded"
	JudgeStatusSystemError         = "System Error"
)

const (
	CompetitionOver                               = "该竞赛已经结束！"
	CompetitionSubmitMessageUnsafeCode            = "unsafe code"
	CompetitionSubmitMessageJudgeQueueUnavailable = "judge queue unavailable"
	CompetitionSubmitMessageJudgeTimeout          = "judge timeout"
	CompetitionSubmitMessageDataQueryError        = "data query error"
	CompetitionSubmitMessageNotStarted            = "competition has not started"
	CompetitionSubmitMessageEnded                 = "competition has ended"
	CompetitionSubmitMessageTestCasesNotFound     = "test cases not found"
	CompetitionSubmitMessageSaveSubmissionFailed  = "save submission failed"
)
