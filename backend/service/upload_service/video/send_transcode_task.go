package video

import (
	"encoding/json"
	"errors"

	"GO1/global"
	"GO1/pkg/constants"

	"github.com/streadway/amqp"
)

func sendVideoTranscodeTask(task videoTranscodeTask) error {
	if global.MQChannel == nil {
		return errors.New("rabbitmq channel unavailable")
	}

	body, err := json.Marshal(task)
	if err != nil {
		return err
	}

	return global.MQChannel.Publish(
		"",
		constants.VideoTranscodeQueue,
		false,
		false,
		amqp.Publishing{
			ContentType:  constants.VideoTranscodeContentTypeJSON,
			DeliveryMode: amqp.Persistent,
			Body:         body,
		},
	)
}
