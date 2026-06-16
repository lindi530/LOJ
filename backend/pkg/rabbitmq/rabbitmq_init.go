package rabbitmq

import (
	"GO1/global"
	"GO1/pkg/constants"
	"github.com/streadway/amqp"
	"log"
	"os"
)

var MQConn *amqp.Connection
var MQChannel *amqp.Channel

func InitRabbitMQ() {

	rabbitURL := os.Getenv("RABBITMQ_URL")
	if rabbitURL == "" {
		rabbitURL = global.Config.RabbitMQ.URL
	}

	var err error
	MQConn, err = amqp.Dial(rabbitURL)
	if err != nil {
		log.Fatalf("RabbitMQ connect error: %v", err)
	}

	MQChannel, err = MQConn.Channel()
	if err != nil {
		log.Fatalf("RabbitMQ channel error: %v", err)
	}

	for _, queueName := range []string{"match_queue", "ac_calendar", constants.CompetitionSubmitQueue, constants.VideoTranscodeQueue} {
		_, err = MQChannel.QueueDeclare(
			queueName,
			true,
			false,
			false,
			false,
			nil,
		)
		if err != nil {
			log.Fatalf("RabbitMQ queue %q declare error: %v", queueName, err)
		}
	}

	global.MQConn = MQConn
	global.MQChannel = MQChannel

	global.Logger.Info("RabbitMQ initialized.")
}
