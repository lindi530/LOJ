package competition_service

import (
	"GO1/database/mysql/competition_mysql"
	"GO1/database/redis/calendar_redis"
	"GO1/models/competition_model"
	"GO1/models/problem_submission_model"
	"GO1/models/ws_model"
	"GO1/service/problem_service"
	"GO1/service/ws_service"
	"time"
)

func judgeCompetitionSubmission(job competitionSubmitJob) competitionSubmitResult {
	message := &ws_model.EditStatus{
		Type: ws_model.MessageTypeEditStatus,
		To:   job.UserID,
	}
	ws_service.WsHub.SendEditData(message, "Pending")

	competition, competitionProblem, examples, err := competition_mysql.GetCompetitionProblemForSubmit(job.CompetitionID, job.ProblemNumber)
	if err != nil {
		ws_service.WsHub.SendEditData(message, "Wrong Answer")
		return failedCompetitionSubmitResult("data query error")
	}

	now := time.Now()
	if competition.StartTime.IsZero() || now.Before(competition.StartTime) {
		ws_service.WsHub.SendEditData(message, "Wrong Answer")
		return failedCompetitionSubmitResult("competition has not started")
	}
	if !competition.EndTime.IsZero() && now.After(competition.EndTime) {
		ws_service.WsHub.SendEditData(message, "Wrong Answer")
		return failedCompetitionSubmitResult("competition has ended")
	}
	if len(examples) == 0 {
		ws_service.WsHub.SendEditData(message, "Wrong Answer")
		return failedCompetitionSubmitResult("test cases not found")
	}

	runResult := problem_service.RunCode(job.UserID, job.Code, job.Language, &examples, message)
	state, score, accepted := summarizeCompetitionRunResult(runResult)

	problemSubmission := &problem_submission_model.ProblemSubmission{
		UserId:    job.UserID,
		ProblemId: uint(competitionProblem.ProblemID),
		Code:      job.Code,
		Language:  job.Language,
		State:     state,
		Score:     score,
		CreatedAt: time.Now(),
	}
	isFirstAC, err := competition_mysql.SaveCompetitionJudgeResult(
		job.CompetitionID,
		competitionProblem.ProblemID,
		job.UserID,
		problemSubmission,
		accepted,
	)
	if err != nil {
		ws_service.WsHub.SendEditData(message, "Wrong Answer")
		return failedCompetitionSubmitResult("save submission failed")
	}

	if accepted {
		_ = calendar_redis.SaveACProblem(job.UserID, uint(competitionProblem.ProblemID))
	}

	ws_service.WsHub.SendEditData(message, state)
	return competitionSubmitResult{
		Code:    0,
		Message: state,
		Data: competition_model.SubmitCompetitionProblemResp{
			SubmissionID: problemSubmission.Id,
			State:        state,
			Score:        score,
			IsFirstAC:    isFirstAC,
			Results:      runResult,
		},
	}
}
