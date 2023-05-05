package sender

import (
	"github.com/nicksdlc/macaw/communicators"
	"github.com/nicksdlc/macaw/config"
	"github.com/nicksdlc/macaw/model"
	"github.com/nicksdlc/macaw/prototype"
)

// Sender sends message to the external interface
type Sender interface {
	Send() error
}

// MessageSender generates and sends defined quantity of messages
type MessageSender struct {
	communicator   communicators.Communicator
	requestBuilder prototype.PrototypeBuilder
}

// NewMessageSender creates a new sender for the provided communicator
func NewMessageSender(communicator communicators.Communicator, request []config.Request) *MessageSender {
	return &MessageSender{
		communicator:   communicator,
		requestBuilder: prototype.NewRequestPrototypeBuilder(request),
	}
}

// SendWithResponse sendr requests and publish responses to the channel
func (rs *MessageSender) SendWithResponse() (chan model.ResponseMessage, error) {
	rs.communicator.RequestWith(rs.requestBuilder.Build())

	return rs.communicator.PostAndListen()
}
