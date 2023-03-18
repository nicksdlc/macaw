package connectors

import (
	"context"
	"log"
	"time"

	"github.com/cenkalti/backoff"
	"github.com/nicksdlc/macaw/config"
	amqp "github.com/rabbitmq/amqp091-go"
)

// RMQExchangeConnector connector to RabbitMQ
type RMQExchangeConnector struct {
	ConnectionString  string
	Exchange          string
	Queues            []string
	ConnectionRetries config.Retry

	sendChannel    *amqp.Channel
	receiveChannel *amqp.Channel
	connection     *amqp.Connection
	inQ            amqp.Queue
	outQ           amqp.Queue
}

// NewRMQExchangeConnector creates new connector with default connection
func NewRMQExchangeConnector(connectionString string, retries config.Retry, exchange string, queue ...string) *RMQExchangeConnector {
	rc := &RMQExchangeConnector{
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

// Close closes connection to RMQ
func (rc *RMQExchangeConnector) Close() error {
	err := rc.sendChannel.Close()
	if err != nil {
		return err
	}

	return rc.connection.Close()
}

// Post sends request to exchange
func (rc *RMQExchangeConnector) Post(body string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := rc.receiveChannel.PublishWithContext(ctx,
		rc.Exchange, // exchange
		"",          // routing key
		false,       // mandatory
		false,       // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        []byte(body),
		})

	if err != nil {
		log.Panicf("%s: %s", "Failed to publish a message", err)
		return err
	}
	log.Printf(" [x] Sent %s to exchange %s\n", body, rc.Exchange)
	return nil
}

func (rc *RMQExchangeConnector) PostIn(body string) {
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

func (rc *RMQExchangeConnector) Consume() <-chan amqp.Delivery {
	msgs, err := rc.sendChannel.Consume(
		rc.inQ.Name, // queue
		"",          // consumer
		true,        // auto-ack
		false,       // exclusive
		false,       // no-local
		false,       // no-wait
		nil,         // args
	)
	failOnError(err, "Failed to register a consumer")

	return msgs
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}
