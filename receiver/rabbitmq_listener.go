package receiver

import (
	"log"

	"github.com/nicksdlc/macaw/config"
	"github.com/nicksdlc/macaw/connectors"
	"github.com/nicksdlc/macaw/generator"
	"github.com/nicksdlc/macaw/template"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RMQReceiver struct {
	rmqConnector *connectors.RMQExchangeConnector
	responder    *generator.GenericResponder
	deliveries   <-chan amqp.Delivery
}

func NewRMQReceiver(connector *connectors.RMQExchangeConnector, resp config.Response) RMQReceiver {
	return RMQReceiver{
		rmqConnector: connector,
		responder: &generator.GenericResponder{
			Response: resp,
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
