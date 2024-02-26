package communicators

import (
	"context"
	"log"
	"time"

	"github.com/cenkalti/backoff"
	"github.com/nicksdlc/macaw/config"
	"github.com/nicksdlc/macaw/model"
	"github.com/nicksdlc/macaw/prototype"
	"github.com/nicksdlc/macaw/prototype/matchers"
	amqp "github.com/rabbitmq/amqp091-go"
)

// RMQExchangeCommunicator communicator to RabbitMQ
type RMQExchangeCommunicator struct {
	ConnectionString  string
	ConnectionRetries config.Retry

	responsePrototypes []prototype.MessagePrototype
	requestPrototypes  []prototype.MessagePrototype

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

	rc.createRetryPolicy()
	rc.connectToRabbit()
	rc.createChannels()
	rc.declareQueues(queues)
	rc.createExchanges(exchanges)

	return rc
}

// GetResponses returns responses
func (rc *RMQExchangeCommunicator) GetResponses() []prototype.MessagePrototype {
	return rc.responsePrototypes
}

// UpdateResponse updates response
func (rc *RMQExchangeCommunicator) UpdateResponse(response prototype.MessagePrototype) {
	for i, resp := range rc.responsePrototypes {
		if resp.Alias == response.Alias {
			rc.responsePrototypes[i] = response
		}
	}
}

// RespondWith defines responses to be sent to RMQ
func (rc *RMQExchangeCommunicator) RespondWith(response []prototype.MessagePrototype) {
	rc.responsePrototypes = response
}

// RequestWith defines requests to be sent to RMQ
func (rc *RMQExchangeCommunicator) RequestWith(request []prototype.MessagePrototype) {
	rc.requestPrototypes = request
}

// Close closes connection to RMQ
func (rc *RMQExchangeCommunicator) Close() error {
	err := rc.sendChannel.Close()
	if err != nil {
		return err
	}

	return rc.connection.Close()
}

// PostAndListen sends request to exchange and waits for response
func (rc *RMQExchangeCommunicator) PostAndListen() (chan model.ResponseMessage, error) {
	responseChannel := make(chan model.ResponseMessage)

	for _, prototype := range rc.requestPrototypes {
		p := prototype
		go func() {
			resp := model.ResponseMessage{}
			for r := range p.Mediators.Run(model.RequestMessage{}, resp) {
				log.Printf("Sending message to %s", p.To)
				rc.post(rc.getExchange(p.To), model.RequestMessage{Body: []byte(r.Body)})
			}
			rc.consumeAndStore(responseChannel, p)
		}()
	}

	return responseChannel, nil
}

// ConsumeMediateReplyWithResponse opens a channel to wait for the information from rabbit mq
func (rc *RMQExchangeCommunicator) ConsumeMediateReplyWithResponse() {
	var listeners []*RMQQueueListener
	for _, queue := range rc.queues {
		listener := NewRMQQueueListener(queue.Name, rc.receiveChannel)
		listeners = append(listeners, listener)
		listener.Listen()
	}

	for _, prototype := range rc.responsePrototypes {
		messageChannel := make(chan model.RequestMessage)
		for _, listener := range listeners {
			if listener.QueueName == prototype.From {
				listener.AddListener(messageChannel)
			}
		}
		go rc.consume(prototype, messageChannel)
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
	failOnError(err, "Failed to publish a message")

	log.Printf(" [rmq-connector] Sent to exchange %s\n", exchange.name)
	return nil
}

func (rc *RMQExchangeCommunicator) consume(messagePrototype prototype.MessagePrototype, messages chan model.RequestMessage) {
	go func() {
		for message := range messages {
			resp := model.ResponseMessage{}
			for r := range messagePrototype.Mediators.Run(message, resp) {
				if matchers.MatchAny(messagePrototype.Matcher, message) {
					rc.post(rc.getExchange(messagePrototype.To), model.RequestMessage{Body: []byte(r.Body)})
				}
			}
		}
	}()
}

func (rc *RMQExchangeCommunicator) consumeAndStore(storeChan chan model.ResponseMessage, messagePrototype prototype.MessagePrototype) {
	amqpMsgs, err := rc.receiveChannel.Consume(
		messagePrototype.From, // queue
		"",                    // consumer
		true,                  // auto-ack
		false,                 // exclusive
		false,                 // no-local
		false,                 // no-wait
		nil,                   // args
	)
	failOnError(err, "Failed to register a consumer")

	go func() {
		for delivery := range amqpMsgs {
			message := model.ResponseMessage{
				Body: string(delivery.Body),
			}
			storeChan <- message
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

func (rc *RMQExchangeCommunicator) connectToRabbit() {
	err := backoff.Retry(func() error {
		var err error
		rc.connection, err = amqp.Dial(rc.ConnectionString)
		if err != nil {
			log.Printf("Not able to connect to rabbit error: %s", err.Error())
		} else {
			log.Printf("Connection succeded")
		}
		return err
	}, rc.retrierPolicy)
	failOnError(err, "Failed to connect to RabbitMQ")
}

func (rc *RMQExchangeCommunicator) createRetryPolicy() {
	rc.retrierPolicy = backoff.NewExponentialBackOff()
	rc.retrierPolicy.MaxInterval = time.Duration(rc.ConnectionRetries.Interval) * time.Second
	rc.retrierPolicy.MaxElapsedTime = time.Duration(rc.ConnectionRetries.ElapsedTime) * time.Minute
}

func (rc *RMQExchangeCommunicator) createChannels() {
	var err error
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
}

func (rc *RMQExchangeCommunicator) declareQueues(queues []config.Queue) {
	var err error
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
}

func (rc *RMQExchangeCommunicator) createExchanges(exchanges []config.Exchange) {
	for _, ex := range exchanges {
		rc.exchanges = append(rc.exchanges, exchange{name: ex.Name, routingKey: ex.RoutingKey})
	}
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}
