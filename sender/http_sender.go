package sender

import (
	"github.com/nicksdlc/macaw/config"
	"github.com/nicksdlc/macaw/connectors"
	"github.com/nicksdlc/macaw/generator"
)

// NewHTTPSender creates new MessageSender
func NewHTTPSender(connector *connectors.HTTPConnector, request config.Request) *MessageSender {
	return &MessageSender{
		connector: connector,
		requester: &generator.JSONRequester{
			Request: request,
		},
	}
}
