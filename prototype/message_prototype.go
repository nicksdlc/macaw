package prototype

import (
	"github.com/nicksdlc/macaw/mediator"
	"github.com/nicksdlc/macaw/prototype/matchers"
)

// MessagePrototype is a template for message with mediators, bodyTemplate and headers
type MessagePrototype struct {
	Headers map[string]string

	BodyTemplate string

	Mediators mediator.MediatorChain

	From string

	To string

	Matcher []matchers.Matcher
}
