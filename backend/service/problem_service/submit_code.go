package problem_service

import (
	"GO1/database/mysql/problem_mysql"
	"GO1/database/redis/calendar_redis"
	"GO1/middlewares/response"
	"GO1/models/problem_model"
	"GO1/models/problem_submission_model"
	"GO1/models/ws_model"
	"GO1/pkg/constants"
	"GO1/service/ws_service"
	"math"
	"time"
)

func SubmitCode(userid int64, codeSubmission problem_model.CodeSubmission, message *ws_model.EditStatus) (resp response.Response) {

	if message == nil {
		message = &ws_model.EditStatus{
			Type: ws_model.MessageTypeEditStatus,
			To:   userid,
		}
		ws_service.WsHub.SendEditData(message, constants.JudgeStatusPending)
	}

	var examples []problem_model.Example

	err := problem_mysql.GetProblemExamples(codeSubmission.ProblemID, &examples)
	if err != nil {
		resp.Code = 1
		resp.Message = "判定失败，请重试"
		return
	}

	var constraints problem_model.ProblemConstraint
	err = problem_mysql.GetProblemConstraints(int64(codeSubmission.ProblemID), codeSubmission.Language, &constraints)
	if err != nil {
		resp.Code = 1
		resp.Message = err.Error()
		return
	}

	resp.Code = 0
	runResult := RunCode(userid, codeSubmission.Code, codeSubmission.Language, &examples, constraints.MemoryLimit, constraints.TimeLimit, message)
	msgContent := constants.JudgeStatusAccepted
	totalCount := len(runResult)
	acCount := totalCount
	for _, result := range runResult {
		if !result.Passed {
			resp.Code = 1
			msgContent = result.Error
			acCount--
		}
	}

	if totalCount == acCount {
		problem_mysql.SaveAcState(userid, codeSubmission.ProblemID)
		problem_mysql.UpdateAcCount(codeSubmission.ProblemID)
	}
	problem_mysql.UpdateSubmitCount(codeSubmission.ProblemID)
	problem_mysql.SaveSubmission(&problem_submission_model.ProblemSubmission{
		UserId:    userid,
		ProblemId: codeSubmission.ProblemID,
		Code:      codeSubmission.Code,
		Language:  codeSubmission.Language,
		State:     msgContent,
		Score:     math.Round(float64(acCount)/float64(totalCount)*10000) / 100,
		CreatedAt: time.Now(),
	})

	if msgContent == constants.JudgeStatusAccepted {
		_ = calendar_redis.SaveACProblem(userid, codeSubmission.ProblemID)
	}

	ws_service.WsHub.SendEditData(message, msgContent)
	resp.Data = runResult

	return
}
