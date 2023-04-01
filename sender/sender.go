package sender

import (
	"time"

	"github.com/nicksdlc/macaw/communicators"
	"github.com/nicksdlc/macaw/config"
	"github.com/nicksdlc/macaw/generator"
	"github.com/nicksdlc/macaw/model"
)

// Sender sends message to the external interface
type Sender interface {
	Send() error
}

// MessageSender generates and sends defined amount of messages
type MessageSender struct {
	communicator communicators.Communicator
	requester    *generator.JSONRequester
}

// NewMessageSender creates a new sender for the provided communicator
func NewMessageSender(communicator communicators.Communicator, request config.Request) *MessageSender {
	return &MessageSender{
		communicator: communicator,
		requester: &generator.JSONRequester{
			Request: request,
		},
	}
}

// Send generates and send requests to communicator
func (rs *MessageSender) Send() error {
	for _, req := range rs.requester.Generate() {
		err := rs.communicator.Post(model.RequestMessage{Body: []byte(req)})
		if err != nil {
			return err
		}
		time.Sleep(time.Duration(rs.requester.Request.Delay) * time.Millisecond)
	}

	return nil
}
