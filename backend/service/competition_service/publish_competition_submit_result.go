package competition_service

import (
	"GO1/global"
	"encoding/json"

	"github.com/streadway/amqp"
)

func publishCompetitionSubmitResult(replyTo, correlationID string, result competitionSubmitResult) error {
	body, err := json.Marshal(result)
	if err != nil {
		return err
	}

	return global.MQChannel.Publish(
		"",
		replyTo,
		false,
		false,
		amqp.Publishing{
			ContentType:   "application/json",
			CorrelationId: correlationID,
			Body:          body,
		},
	)
}
