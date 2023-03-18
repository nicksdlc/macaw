package sender

import (
	"github.com/nicksdlc/macaw/config"
	"github.com/nicksdlc/macaw/connectors"
	"github.com/nicksdlc/macaw/generator"
)

// NewRMQSender constructor
func NewRMQSender(connector *connectors.RMQExchangeConnector, request config.Request) MessageSender {
	return MessageSender{
		connector: connector,
		requester: &generator.JSONRequester{
			Request: request,
		},
	}
}
