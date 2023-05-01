package prototype

import (
	"os"

	"github.com/nicksdlc/macaw/config"
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

// PrototypeBuilder builds message templates
type PrototypeBuilder interface {
	Build() []MessagePrototype
}

func readMessageTemplate(filePath string) []byte {
	base, err := os.ReadFile(filePath)
	if err != nil {
		panic(err)
	}
	return base
}

func buildMediators(bodyTemplate string, options config.Options) mediator.MediatorChain {
	var chain mediator.MediatorChain

	chain.Append(mediator.NewGeneratingMediator(options.Quantity, bodyTemplate))
	chain.Append(mediator.NewDelayingMediator(options.Delay))

	return chain
}

func buildMatcher(cfg config.Matchers) []matchers.Matcher {
	messageMatchers := []matchers.Matcher{}

	if cfg.Field.Name != "" {
		messageMatchers = append(messageMatchers, &matchers.FieldMatcher{
			Field: cfg.Field.Name,
			Value: cfg.Field.Value,
		})
	}

	if cfg.Contains.Value != "" {
		messageMatchers = append(messageMatchers, &matchers.BodyContainsMatcher{
			Contains: cfg.Contains.Value,
		})
	}

	return messageMatchers
}

func buildBodyTemplate(body config.Body) string {
	if body.File != "" {
		return string(readMessageTemplate(body.File))
	}
	return body.String
}
