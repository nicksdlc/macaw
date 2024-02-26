package communicators

import (
	"github.com/nicksdlc/macaw/model"
	"github.com/nicksdlc/macaw/prototype"
)

// Communicator is an interface the represents communicator to external source
type Communicator interface {
	RespondWith(response []prototype.MessagePrototype)

	RequestWith(request []prototype.MessagePrototype)

	GetResponses() []prototype.MessagePrototype

	UpdateResponse(response prototype.MessagePrototype)

	Close() error

	PostAndListen() (chan model.ResponseMessage, error)

	ConsumeMediateReplyWithResponse()
}
