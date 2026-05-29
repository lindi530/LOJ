package competition_service

import (
	"GO1/global"
	"GO1/pkg/constants"
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
			ContentType:   constants.CompetitionSubmitContentTypeJSON,
			CorrelationId: correlationID,
			Body:          body,
		},
	)
}
