package prototype

import (
	"github.com/nicksdlc/macaw/config"
	"github.com/nicksdlc/macaw/mediator"

	matchers "github.com/nicksdlc/macaw/prototype/matchers"
)

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
		bodyTemplate := buildBodyTemplate(response.Body)
		templates = append(templates, MessagePrototype{
			BodyTemplate: bodyTemplate,
			Mediators:    buildMediators(bodyTemplate, response.Options),
			From:         response.ResponseRequest.To,
			To:           response.To,
			Matcher:      buildMatcher(response.ResponseRequest.Matchers),
		})

		// prepend matchingMediator to the mediators list
		// so that it is the first mediator to be executed
		// and all others can be skipped if it doesn't match
		templates[len(templates)-1].Mediators.Prepend(
			mediator.NewMatchingMediator(
				matchers.ParsePattern(response.ResponseRequest.Match),
				templates[len(templates)-1].Matcher))
	}

	return templates
}
