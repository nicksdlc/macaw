package responder

import (
	"github.com/nicksdlc/macaw/communicators"
	"github.com/nicksdlc/macaw/config"
	"github.com/nicksdlc/macaw/generator"
	"github.com/nicksdlc/macaw/model"
	"github.com/nicksdlc/macaw/template"
)

// Responder listens to the incoming message and updates the responder
type Responder interface {
	Listen()

	Notify(message []byte)
}

// MessageResponder listens to http requests and responds with messages
type MessageResponder struct {
	communicator communicators.Communicator
	responder    *generator.GenericResponder
}

// NewMessageResponder creates listener for requests
func NewMessageResponder(communicator communicators.Communicator, resp config.Response) *MessageResponder {
	return &MessageResponder{
		communicator: communicator,
		responder: &generator.GenericResponder{
			Response: resp,
		},
	}
}

// Listen for HTTPServer requests
func (h *MessageResponder) Listen() {
	var mediators []model.Mediator
	mediators = append(mediators, func(message model.RequestMessage, response model.ResponseMessage) (model.ResponseMessage, error) {
		response.Responses = h.responder.Generate(template.Serialize(message.Headers, message.Body))
		return response, nil
	})
	h.communicator.RespondWith(h.responder.Response, mediators)
	h.communicator.ConsumeMediateReplyWithResponse()
}
