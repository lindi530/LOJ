package competition_service

import (
	"GO1/global"
	"GO1/pkg/constants"
	"encoding/json"
)

func StartCompetitionSubmitConsumer() {
	if global.MQChannel == nil {
		global.Logger.Error("competition submit queue unavailable")
		return
	}

	msgs, err := global.MQChannel.Consume(
		constants.CompetitionSubmitQueue,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		global.Logger.Error("failed to register competition submit consumer:", err)
		return
	}

	go func() {
		for msg := range msgs {
			var job competitionSubmitJob
			if err := json.Unmarshal(msg.Body, &job); err != nil {
				global.Logger.Error("invalid competition submit message:", err)
				_ = msg.Nack(false, false)
				continue
			}

			result := judgeCompetitionSubmission(job)
			if msg.ReplyTo != "" {
				if err := publishCompetitionSubmitResult(msg.ReplyTo, msg.CorrelationId, result); err != nil {
					global.Logger.Error("publish competition submit result failed:", err)
					_ = msg.Nack(false, true)
					continue
				}
			}

			if err := msg.Ack(false); err != nil {
				global.Logger.Error("ack competition submit message failed:", err)
			}
		}
	}()
}
