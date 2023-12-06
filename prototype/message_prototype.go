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

	// Type is a message type
	Type string

	// BodyTemplate is a template for message body
	// It uses text/template syntax
	BodyTemplate []string

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

func buildMediators(bodyTemplate []string, options config.Options) mediator.MediatorChain {
	var chain mediator.MediatorChain

	chain.Append(mediator.NewGeneratingMediator(options.Quantity, bodyTemplate))
	chain.Append(mediator.NewDelayingMediator(options.Delay))
	chain.Append(mediator.NewRandomDelayingMediator(options.RandomDelay.Min, options.RandomDelay.Max))

	return chain
}

func buildMatcher(cfg []config.Matcher) []matchers.Matcher {
	messageMatchers := []matchers.Matcher{}

	for _, matcher := range cfg {
		messageMatchers = append(messageMatchers, matcherTypes[matcher.Type](matcher))
	}

	return messageMatchers
}

func buildBodyTemplate(body config.Body) []string {
	var bodyTemplate []string

	for _, bodyPart := range body.File {
		bodyTemplate = append(bodyTemplate, string(readMessageTemplate(bodyPart)))
	}

	for _, bodyPart := range body.String {
		bodyTemplate = append(bodyTemplate, bodyPart)
	}

	return bodyTemplate
}
