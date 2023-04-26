package communicators

import (
	"github.com/nicksdlc/macaw/model"
)

// Communicator is an interface the represents communicator to external source
type Communicator interface {
	RespondWith(response []model.MessagePrototype)

	Close() error

	Post(body model.RequestMessage) error

	Consume() <-chan model.RequestMessage

	ConsumeMediateReply(mediators []model.Mediator)

	ConsumeMediateReplyWithResponse()
}
