package sender

import (
	"time"

	"github.com/nicksdlc/macaw/config"
	"github.com/nicksdlc/macaw/connectors"
	"github.com/nicksdlc/macaw/generator"
)

// RMQSender sends messages to rabbitmq
type RMQSender struct {
	rmqConnector *connectors.RMQExchangeConnector
	requester    *generator.JSONRequester
}

// NewRMQSender constructor
func NewRMQSender(connector *connectors.RMQExchangeConnector, request config.Request) RMQSender {
	return RMQSender{
		rmqConnector: connector,
		requester: &generator.JSONRequester{
			Request: request,
		},
	}
}

// Send message to rabbitmq
func (rs *RMQSender) Send() error {
	for _, req := range rs.requester.Generate() {
		rs.rmqConnector.Post(req)
		time.Sleep(time.Duration(rs.requester.Request.Delay) * time.Millisecond)
	}

	return nil
}
