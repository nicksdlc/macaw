package communicators

import (
	"context"
	"log"
	"time"

	"github.com/cenkalti/backoff"
	"github.com/nicksdlc/macaw/config"
	"github.com/nicksdlc/macaw/model"
	amqp "github.com/rabbitmq/amqp091-go"
)

// RMQExchangeCommunicator communicator to RabbitMQ
type RMQExchangeCommunicator struct {
	ConnectionString  string
	Exchange          string
	Queues            []string
	ConnectionRetries config.Retry

	mediators []model.Mediator

	sendChannel    *amqp.Channel
	receiveChannel *amqp.Channel
	connection     *amqp.Connection
	inQ            amqp.Queue
	outQ           amqp.Queue
}

// NewRMQExchangeCommunicator creates new communicator with default connection
func NewRMQExchangeCommunicator(connectionString string, retries config.Retry, exchange string, queue ...string) *RMQExchangeCommunicator {
	rc := &RMQExchangeCommunicator{
		ConnectionString:  connectionString,
		Exchange:          exchange,
		Queues:            queue,
		ConnectionRetries: retries,
	}

	var err error
	operation := func() error {
		rc.connection, err = amqp.Dial(connectionString)
		if err != nil {
			log.Printf("Not able to connect to rabbit error: %s", err.Error())
		} else {
			log.Printf("Connection succeded")
		}
		return err
	}

	backoffPolicy := backoff.NewExponentialBackOff()
	backoffPolicy.MaxInterval = time.Duration(rc.ConnectionRetries.Interval) * time.Second
	backoffPolicy.MaxElapsedTime = time.Duration(rc.ConnectionRetries.ElapsedTime) * time.Minute

	err = backoff.Retry(operation, backoffPolicy)
	failOnError(err, "Failed to connect to RabbitMQ")

	rc.sendChannel, err = rc.connection.Channel()
	failOnError(err, "Failed to open a channel")

	rc.receiveChannel, err = rc.connection.Channel()
	failOnError(err, "Failed to open a channel")

	rc.inQ, err = rc.sendChannel.QueueDeclare(
		queue[0], // name
		true,     // durable
		false,    // delete when unused
		false,    // exclusive
		false,    // no-wait
		nil,      // arguments
	)
	failOnError(err, "Failed to declare a queue")

	rc.outQ, err = rc.receiveChannel.QueueDeclare(
		queue[1], // name
		true,     // durable
		false,    // delete when unused
		false,    // exclusive
		false,    // no-wait
		nil,      // arguments
	)
	failOnError(err, "Failed to declare a queue")

	return rc
}

// RespondWith defines responses to be sent to RMQ
func (rc *RMQExchangeCommunicator) RespondWith(response config.Response, mediators []model.Mediator) {
	rc.mediators = mediators
}

// Close closes connection to RMQ
func (rc *RMQExchangeCommunicator) Close() error {
	err := rc.sendChannel.Close()
	if err != nil {
		return err
	}

	return rc.connection.Close()
}

// Post sends request to exchange
func (rc *RMQExchangeCommunicator) Post(body model.RequestMessage) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := rc.receiveChannel.PublishWithContext(ctx,
		rc.Exchange, // exchange
		"",          // routing key
		false,       // mandatory
		false,       // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body.Body,
		})

	if err != nil {
		log.Panicf("%s: %s", "Failed to publish a message", err)
		return err
	}
	log.Printf(" [x] Sent %s to exchange %s\n", body, rc.Exchange)
	return nil
}

// PostIn might be soon deleted, un-used
func (rc *RMQExchangeCommunicator) PostIn(body string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := rc.sendChannel.PublishWithContext(ctx,
		rc.Exchange, // exchange
		rc.inQ.Name, // routing key
		false,       // mandatory
		false,       // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        []byte(body),
		})
	failOnError(err, "Failed to publish a message")
	log.Printf(" [x] Sent %s\n", body)
}

// Consume opens a channel to wait for the information from rabbit mq
func (rc *RMQExchangeCommunicator) Consume() <-chan model.RequestMessage {

	amqpMsgs, err := rc.sendChannel.Consume(
		rc.inQ.Name, // queue
		"",          // consumer
		true,        // auto-ack
		false,       // exclusive
		false,       // no-local
		false,       // no-wait
		nil,         // args
	)
	failOnError(err, "Failed to register a consumer")

	msgs := make(chan model.RequestMessage)
	go func() {
		for delivery := range amqpMsgs {
			msgs <- model.RequestMessage{
				Body: delivery.Body,
			}
		}
	}()

	return msgs
}

// ConsumeMediateReply opens a channel to wait for the information from rabbit mq
func (rc *RMQExchangeCommunicator) ConsumeMediateReply(mediators []model.Mediator) {

	amqpMsgs, err := rc.sendChannel.Consume(
		rc.inQ.Name, // queue
		"",          // consumer
		true,        // auto-ack
		false,       // exclusive
		false,       // no-local
		false,       // no-wait
		nil,         // args
	)
	failOnError(err, "Failed to register a consumer")

	go func() {
		for delivery := range amqpMsgs {
			message := model.RequestMessage{
				Body: delivery.Body,
			}

			resp := model.ResponseMessage{}
			for _, mediator := range mediators {
				resp, _ = mediator(message, resp)
			}

			for _, msg := range resp.Responses {
				rc.Post(model.RequestMessage{Body: []byte(msg)})
			}
		}
	}()
}

// ConsumeMediateReplyWithResponse opens a channel to wait for the information from rabbit mq
func (rc *RMQExchangeCommunicator) ConsumeMediateReplyWithResponse() {
	rc.ConsumeMediateReply(rc.mediators)
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}
