package receiver

import (
	"log"
	"macaw/connectors"
	"macaw/responder"
	"macaw/template"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RMQReceiver struct {
	rmqConnector *connectors.RMQExchangeConnector
	responder    *responder.PolicyResponder
	deliveries   <-chan amqp.Delivery
}

func NewRMQReceiver(connector *connectors.RMQExchangeConnector, responsePath string) RMQReceiver {
	return RMQReceiver{
		rmqConnector: connector,
		responder: &responder.PolicyResponder{
			TemplatePath: responsePath,
		},
	}
}

func (r *RMQReceiver) Listen() {
	r.deliveries = r.rmqConnector.Consume()

	go func() {
		log.Printf("Started listening")
		for d := range r.deliveries {
			r.Notify(d.Body)
		}
	}()
}

// Notify - sends a message to the external source
func (r *RMQReceiver) Notify(message []byte) {
	log.Printf("Received a message: %s", message)
	req := template.Serialize(message)

	for _, resp := range r.responder.Generate(req) {
		r.rmqConnector.Post(resp)
	}

}
