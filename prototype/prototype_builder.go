package prototype

import (
	"os"

	"github.com/nicksdlc/macaw/config"
	"github.com/nicksdlc/macaw/mediator"
	"github.com/nicksdlc/macaw/prototype/matchers"
)

// PrototypeBuilder builds message templates
type PrototypeBuilder interface {
	Build() []MessagePrototype
}

// ResponsePrototypeBuilder builds response templates
type ResponsePrototypeBuilder struct {
	responseConfig []config.Response
}

// NewResponsePrototypeBuilder creates a new ResponseTemplateBuilder
func NewResponsePrototypeBuilder(responseConfig []config.Response) *ResponsePrototypeBuilder {
	return &ResponsePrototypeBuilder{
		responseConfig: responseConfig,
	}
}

// Build builds response templates from response configuration
func (rtb *ResponsePrototypeBuilder) Build() []MessagePrototype {
	var templates []MessagePrototype

	for _, response := range rtb.responseConfig {
		templates = append(templates, MessagePrototype{
			BodyTemplate: rtb.buildBodyTemplate(response),
			Mediators:    rtb.buildMediators(response),
			From:         response.ResponseRequest.To,
			To:           response.To,
			Matcher:      rtb.buildMatcher(response),
		})
	}

	return templates
}

func (rtb *ResponsePrototypeBuilder) buildMediators(response config.Response) mediator.MediatorChain {
	var chain mediator.MediatorChain

	chain.Append(mediator.NewGeneratingMediator(response.Amount, rtb.buildBodyTemplate(response)))
	chain.Append(mediator.NewDelayingMediator(response.Delay))

	return chain
}

func (rtb *ResponsePrototypeBuilder) buildBodyTemplate(response config.Response) string {
	if response.File != "" {
		return string(readResponseTemplate(response.File))
	}
	return response.String
}

func (rtb *ResponsePrototypeBuilder) buildMatcher(response config.Response) []matchers.Matcher {
	if response.ResponseRequest.Field.Name != "" {
		return []matchers.Matcher{&matchers.FieldMatcher{
			Field: response.ResponseRequest.Field.Name,
			Value: response.ResponseRequest.Field.Value,
		}}
	}
	return nil
}

func readResponseTemplate(filePath string) []byte {
	base, err := os.ReadFile(filePath)
	if err != nil {
		panic(err)
	}
	return base
}
