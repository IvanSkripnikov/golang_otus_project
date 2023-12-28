package queue

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/IvanSkripnikov/golang_otus_project/logger"
	"github.com/IvanSkripnikov/golang_otus_project/models"
	amqp "github.com/rabbitmq/amqp091-go"
)

func SendEventToQueue(eventName string, bannerID, slotID, groupID int) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")

	failOnError(err, "Failed to connect to RabbitMQ")

	defer conn.Close()

	ch, err := conn.Channel()

	failOnError(err, "Failed to open a channel")

	defer ch.Close()

	q, err := ch.QueueDeclare(

		"events-queue",

		false,

		false,

		false,

		false,

		nil,
	)

	failOnError(err, "Failed to declare a queue")

	message := models.Message{Type: eventName, BannerID: bannerID, SlotID: slotID, GroupID: groupID}

	body, err := json.Marshal(message)
	if err != nil {
		logger.SendToErrorLog(fmt.Sprintf("Error: %s", err))

		return
	}

	err = ch.PublishWithContext(
		context.Background(),
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        body,
		})

	failOnError(err, "Failed to publish a message")

	logger.SendToInfoLog(fmt.Sprintf("[x] Congrats, sending message: %s", body))
}

func failOnError(err error, msg string) {
	if err != nil {
		logger.SendToFatalLog(fmt.Sprintf("%s: %s", msg, err))
	}
}
