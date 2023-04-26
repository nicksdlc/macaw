package communicators

import (
	"github.com/nicksdlc/macaw/model"
)

// Communicator is an interface the represents communicator to external source
type Communicator interface {
	RespondWith(response []model.MessagePrototype)

	Close() error

	Post(body model.RequestMessage) error

	ConsumeMediateReplyWithResponse()
}

// Should be moved to a mediator maybe
func matchAny(res model.MessagePrototype, message model.RequestMessage) bool {
	if len(res.Matcher) == 0 {
		return true
	}

	for _, matcher := range res.Matcher {
		if matcher.Match(message) {
			return true
		}
	}
	return false
}
