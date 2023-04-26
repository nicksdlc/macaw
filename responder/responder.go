package responder

import (
	"github.com/nicksdlc/macaw/builder"
	"github.com/nicksdlc/macaw/communicators"
	"github.com/nicksdlc/macaw/config"
)

// Responder listens to the incoming message and updates the responder
type Responder interface {
	Listen()

	Notify(message []byte)
}

// MessageResponder listens to http requests and responds with messages
type MessageResponder struct {
	communicator    communicators.Communicator
	responseBuilder builder.PrototypeBuilder
}

// NewMessageResponder creates listener for requests
func NewMessageResponder(communicator communicators.Communicator, resp []config.Response) *MessageResponder {
	return &MessageResponder{
		communicator:    communicator,
		responseBuilder: builder.NewResponsePrototypeBuilder(resp),
	}
}

// Listen for HTTPServer requests
func (h *MessageResponder) Listen() {
	h.communicator.RespondWith(h.responseBuilder.Build())
	h.communicator.ConsumeMediateReplyWithResponse()
}
