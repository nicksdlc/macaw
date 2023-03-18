package sender

import (
	"time"

	"github.com/nicksdlc/macaw/connectors"
	"github.com/nicksdlc/macaw/generator"
)

// Sender sends message to the external interface
type Sender interface {
	Send() error
}

// MessageSender generates and sends defined amount of messages
type MessageSender struct {
	connector connectors.Connector
	requester *generator.JSONRequester
}

// Send generates and send requests to connector
func (rs *MessageSender) Send() error {
	for _, req := range rs.requester.Generate() {
		err := rs.connector.Post(req)
		if err != nil {
			return err
		}
		time.Sleep(time.Duration(rs.requester.Request.Delay) * time.Millisecond)
	}

	return nil
}
