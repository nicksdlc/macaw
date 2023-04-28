package communicators

import (
	"github.com/nicksdlc/macaw/model"
	"github.com/nicksdlc/macaw/prototype"
)

// Communicator is an interface the represents communicator to external source
type Communicator interface {
	RespondWith(response []prototype.MessagePrototype)

	Close() error

	Post(body model.RequestMessage) error

	ConsumeMediateReplyWithResponse()
}
