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
	ConnectionRetries config.Retry

	responsePrototypes []model.MessagePrototype

	sendChannel    *amqp.Channel
	receiveChannel *amqp.Channel
	connection     *amqp.Connection
	queues         []amqp.Queue
	exchanges      []exchange

	retrierPolicy *backoff.ExponentialBackOff
}

type exchange struct {
	name       string
	routingKey string
}

// NewRMQExchangeCommunicator creates new communicator with default connection
func NewRMQExchangeCommunicator(connectionString string, retries config.Retry, exchanges []config.Exchange, queues []config.Queue) *RMQExchangeCommunicator {
	rc := &RMQExchangeCommunicator{
		ConnectionString:  connectionString,
		ConnectionRetries: retries,
	}

	var err error
	connectToRabbit := func() error {
		rc.connection, err = amqp.Dial(connectionString)
		if err != nil {
			log.Printf("Not able to connect to rabbit error: %s", err.Error())
		} else {
			log.Printf("Connection succeded")
		}
		return err
	}

	rc.retrierPolicy = backoff.NewExponentialBackOff()
	rc.retrierPolicy.MaxInterval = time.Duration(rc.ConnectionRetries.Interval) * time.Second
	rc.retrierPolicy.MaxElapsedTime = time.Duration(rc.ConnectionRetries.ElapsedTime) * time.Minute

	err = backoff.Retry(connectToRabbit, rc.retrierPolicy)
	failOnError(err, "Failed to connect to RabbitMQ")

	err = backoff.Retry(func() error {
		rc.sendChannel, err = rc.connection.Channel()
		return err
	}, rc.retrierPolicy)
	failOnError(err, "Failed to open a channel")

	err = backoff.Retry(func() error {
		rc.receiveChannel, err = rc.connection.Channel()
		return err
	}, rc.retrierPolicy)
	failOnError(err, "Failed to open a channel")

	for _, queue := range queues {
		err = backoff.Retry(func() error {
			q, err := rc.sendChannel.QueueDeclare(
				queue.Name, // name
				true,       // durable
				false,      // delete when unused
				false,      // exclusive
				false,      // no-wait
				queue.Args, // arguments
			)
			rc.queues = append(rc.queues, q)
			return err
		}, rc.retrierPolicy)
		failOnError(err, "Failed to declare a queue")
	}

	for _, ex := range exchanges {
		rc.exchanges = append(rc.exchanges, exchange{name: ex.Name, routingKey: ex.RoutingKey})
	}

	return rc
}

// RespondWith defines responses to be sent to RMQ
func (rc *RMQExchangeCommunicator) RespondWith(response []model.MessagePrototype) {
	rc.responsePrototypes = response
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
func (rc *RMQExchangeCommunicator) Post(message model.RequestMessage) error {
	return rc.post(rc.exchanges[0], message)
}

// ConsumeMediateReplyWithResponse opens a channel to wait for the information from rabbit mq
func (rc *RMQExchangeCommunicator) ConsumeMediateReplyWithResponse() {
	for _, prototype := range rc.responsePrototypes {
		go rc.consume(prototype)
	}
}

func (rc *RMQExchangeCommunicator) post(exchange exchange, message model.RequestMessage) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := backoff.Retry(func() error {
		err := rc.receiveChannel.PublishWithContext(ctx,
			exchange.name,       // exchange
			exchange.routingKey, // routing key
			false,               // mandatory
			false,               // immediate
			amqp.Publishing{
				ContentType: "application/json",
				Body:        message.Body,
			})
		return err
	}, rc.retrierPolicy)
	if err != nil {
		log.Panicf("%s: %s", "Failed to publish a message", err)
		return err
	}
	log.Printf(" [x] Sent %s to exchange %s\n", message, rc.exchanges[0].name)
	return nil
}

func (rc *RMQExchangeCommunicator) consume(prototype model.MessagePrototype) {
	amqpMsgs, err := rc.sendChannel.Consume(
		prototype.From, // queue
		"",             // consumer
		true,           // auto-ack
		false,          // exclusive
		false,          // no-local
		false,          // no-wait
		nil,            // args
	)
	failOnError(err, "Failed to register a consumer")

	go func() {
		for delivery := range amqpMsgs {
			message := model.RequestMessage{
				Body: delivery.Body,
			}

			resp := model.ResponseMessage{}
			for _, mediator := range prototype.Mediators {
				mediator(message, &resp)
			}

			for _, msg := range resp.Responses {
				if matchAny(prototype, message) {
					rc.post(rc.getExchange(prototype.To), model.RequestMessage{Body: []byte(msg)})
				}
			}
		}
	}()
}

func (rc *RMQExchangeCommunicator) getExchange(name string) exchange {
	for _, ex := range rc.exchanges {
		if ex.name == name {
			return ex
		}
	}
	return exchange{}
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}
