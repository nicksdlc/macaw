package mediator

import (
	"github.com/nicksdlc/macaw/model"
)

// Mediator is an interface that should be implemented by all mediators
type Mediator interface {
	// Mediate is a function that should be done with message
	Mediate(message model.RequestMessage, _ <-chan model.ResponseMessage) <-chan model.ResponseMessage
}
