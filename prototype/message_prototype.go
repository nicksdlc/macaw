package prototype

import (
	"github.com/nicksdlc/macaw/mediator"
	"github.com/nicksdlc/macaw/prototype/matchers"
)

// MessagePrototype is a template for message with mediators, bodyTemplate and headers
type MessagePrototype struct {
	// Headers is a map of headers
	// Optional and used mostly in HTTP requests
	Headers map[string]string

	// Parameters is a map of parameters
	// Optional and used mostly in HTTP requests
	Parameters map[string]string

	// BodyTemplate is a template for message body
	// It uses text/template syntax
	BodyTemplate string

	Mediators mediator.MediatorChain

	From string

	To string

	Matcher []matchers.Matcher
}
