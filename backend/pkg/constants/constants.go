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
	SaberInitialRating     = 1500
	SaberInitialLevelIndex = 3
	SaberPKRatingDelta     = 25
)

const (
	SaberLevelBlackIron = "黑铁"
	SaberLevelBronze    = "黄铜"
	SaberLevelSilver    = "白银"
	SaberLevelGold      = "黄金"
	SaberLevelPlatinum  = "铂金"
	SaberLevelEmerald   = "翡翠"
	SaberLevelDiamond   = "钻石"
	SaberLevelMaster    = "大师"
	SaberLevelGrand     = "宗师"
	SaberLevelKing      = "王者"
)

var SaberLevels = []string{
	SaberLevelBlackIron,
	SaberLevelBronze,
	SaberLevelSilver,
	SaberLevelGold,
	SaberLevelPlatinum,
	SaberLevelEmerald,
	SaberLevelDiamond,
	SaberLevelMaster,
	SaberLevelGrand,
	SaberLevelKing,
}

func SaberLevelByRating(rating int) string {
	levelIndex := SaberInitialLevelIndex + (rating-SaberInitialRating)/100
	if levelIndex < 0 {
		return SaberLevels[0]
	}
	if levelIndex >= len(SaberLevels) {
		return SaberLevels[len(SaberLevels)-1]
	}
	return SaberLevels[levelIndex]
}

const (
	CompetitionSubmitQueue             = "competition_submit_queue"
	CompetitionSubmitConsumerTagPrefix = "competition-submit-"
	CompetitionSubmitContentTypeJSON   = "application/json"
	CompetitionSubmitTimeout           = 2 * time.Minute
)

const (
	VideoTranscodeQueue           = "video_transcode_queue"
	VideoTranscodeContentTypeJSON = "application/json"
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
