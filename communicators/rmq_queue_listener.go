package communicators

import (
	"log"

	"github.com/nicksdlc/macaw/model"
	amqp "github.com/rabbitmq/amqp091-go"
)

type RMQQueueListener struct {
	QueueName string

	receiveChannel    *amqp.Channel
	listeningChannels []chan model.RequestMessage
}

func NewRMQQueueListener(queueName string, receiveChannel *amqp.Channel) *RMQQueueListener {
	return &RMQQueueListener{
		QueueName:         queueName,
		receiveChannel:    receiveChannel,
		listeningChannels: make([]chan model.RequestMessage, 0),
	}
}

func (ql *RMQQueueListener) Listen() {
	ampqMsgs, err := ql.receiveChannel.Consume(
		ql.QueueName,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Panicf("Failed to register a consumer: %s", err)
	}

	go func() {
		for delivery := range ampqMsgs {
			message := model.RequestMessage{
				Body: delivery.Body,
			}

			for _, ch := range ql.listeningChannels {
				ch <- message
			}
		}
	}()
}

func (ql *RMQQueueListener) AddListener(ch chan model.RequestMessage) {
	ql.listeningChannels = append(ql.listeningChannels, ch)
}
