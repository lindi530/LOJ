package competition_service

import (
	"GO1/global"
	"GO1/pkg/constants"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/streadway/amqp"
)

func publishCompetitionSubmitJob(job competitionSubmitJob) (competitionSubmitResult, error) {
	replyQueue, err := global.MQChannel.QueueDeclare(
		"",
		false,
		true,
		true,
		false,
		nil,
	)
	if err != nil {
		return competitionSubmitResult{}, err
	}

	correlationID := uuid.NewString()
	consumerTag := constants.CompetitionSubmitConsumerTagPrefix + correlationID
	msgs, err := global.MQChannel.Consume(
		replyQueue.Name,
		consumerTag,
		true,
		true,
		false,
		false,
		nil,
	)
	if err != nil {
		return competitionSubmitResult{}, err
	}
	defer global.MQChannel.Cancel(consumerTag, false)

	body, err := json.Marshal(job)
	if err != nil {
		return competitionSubmitResult{}, err
	}

	if err := global.MQChannel.Publish(
		"",
		constants.CompetitionSubmitQueue,
		false,
		false,
		amqp.Publishing{
			ContentType:   constants.CompetitionSubmitContentTypeJSON,
			CorrelationId: correlationID,
			ReplyTo:       replyQueue.Name,
			Body:          body,
		},
	); err != nil {
		return competitionSubmitResult{}, err
	}

	timer := time.NewTimer(constants.CompetitionSubmitTimeout)
	defer timer.Stop()

	for {
		select {
		case msg, ok := <-msgs:
			if !ok {
				return competitionSubmitResult{}, responseTimeoutError{}
			}
			if msg.CorrelationId != correlationID {
				continue
			}
			var result competitionSubmitResult
			if err := json.Unmarshal(msg.Body, &result); err != nil {
				return competitionSubmitResult{}, err
			}
			return result, nil
		case <-timer.C:
			return competitionSubmitResult{}, responseTimeoutError{}
		}
	}
}
