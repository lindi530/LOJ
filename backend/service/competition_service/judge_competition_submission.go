package competition_service

import (
	"GO1/database/mysql/competition_mysql"
	"GO1/database/mysql/problem_mysql"
	"GO1/database/redis/calendar_redis"
	"GO1/database/redis/competition_redis"
	"GO1/global"
	"GO1/models/competition_model"
	"GO1/models/problem_model"
	"GO1/models/problem_submission_model"
	"GO1/models/ws_model"
	"GO1/pkg/constants"
	"GO1/service/problem_service"
	"GO1/service/ws_service"
	"time"
)

func judgeCompetitionSubmission(job competitionSubmitJob) competitionSubmitResult {
	message := &ws_model.EditStatus{
		Type: ws_model.MessageTypeEditStatus,
		To:   job.UserID,
	}
	ws_service.WsHub.SendEditData(message, constants.JudgeStatusPending)

	competition, competitionProblem, examples, err := competition_mysql.GetCompetitionProblemForSubmit(job.CompetitionID, job.ProblemNumber)
	if err != nil {
		ws_service.WsHub.SendEditData(message, constants.JudgeStatusSystemError)
		return failedCompetitionSubmitResult(constants.CompetitionSubmitMessageDataQueryError)
	}

	now := time.Now()
	if competition.StartTime.IsZero() || now.Before(competition.StartTime) {
		ws_service.WsHub.SendEditData(message, constants.JudgeStatusSystemError)
		return failedCompetitionSubmitResult(constants.CompetitionSubmitMessageNotStarted)
	}
	if !competition.EndTime.IsZero() && now.After(competition.EndTime) {
		ws_service.WsHub.SendEditData(message, constants.JudgeStatusSystemError)
		return failedCompetitionSubmitResult(constants.CompetitionSubmitMessageEnded)
	}
	if len(examples) == 0 {
		ws_service.WsHub.SendEditData(message, constants.JudgeStatusSystemError)
		return failedCompetitionSubmitResult(constants.CompetitionSubmitMessageTestCasesNotFound)
	}

	var constraints problem_model.ProblemConstraint
	if err := problem_mysql.GetProblemConstraints(int64(competitionProblem.ProblemID), job.Language, &constraints); err != nil {
		ws_service.WsHub.SendEditData(message, constants.JudgeStatusSystemError)
		return failedCompetitionSubmitResult(constants.CompetitionSubmitMessageDataQueryError)
	}

	runResult := problem_service.RunCode(job.UserID, job.Code, job.Language, &examples, constraints.MemoryLimit, constraints.TimeLimit, message)
	execTime, memoryUsage := problem_service.SummarizeRunResultMetrics(runResult)

	state, score, accepted := summarizeCompetitionRunResult(runResult)
	state = NormalizeJudgeStatus(state)

	global.Logger.Error("state, score, accepted: ", state, score, accepted)

	submitTime := time.Now()
	problemSubmission := &problem_submission_model.ProblemSubmission{
		UserId:      job.UserID,
		ProblemId:   uint(competitionProblem.ProblemID),
		Code:        job.Code,
		Language:    job.Language,
		State:       state,
		ExecTime:    execTime,
		MemoryUsage: memoryUsage,
		Score:       score,
		CreatedAt:   submitTime,
	}
	isFirstAC, err := competition_mysql.SaveCompetitionJudgeResult(
		job.CompetitionID,
		competitionProblem.ProblemID,
		job.UserID,
		problemSubmission,
		accepted,
	)
	if err != nil {
		ws_service.WsHub.SendEditData(message, constants.JudgeStatusSystemError)
		return failedCompetitionSubmitResult(constants.CompetitionSubmitMessageSaveSubmissionFailed)
	}

	if err := competition_redis.UpdateCompetitionRankingList(
		job.CompetitionID,
		competition.StartTime,
		competition.EndTime,
		submitTime,
		competitionProblem.ProblemID,
		competitionProblem.ProblemNumber,
		job.UserID,
		job.UserName,
		accepted,
		isFirstAC,
	); err != nil {
		global.Logger.Error("update competition ranking list failed:", err)
		ws_service.WsHub.SendEditData(message, constants.JudgeStatusSystemError)
		return failedCompetitionSubmitResult(constants.CompetitionSubmitMessageSaveSubmissionFailed)
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
