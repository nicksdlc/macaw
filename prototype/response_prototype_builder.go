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
		templates = append(templates, rtb.BuildResponse(response))
	}

	return templates
}

func (rtb *ResponsePrototypeBuilder) BuildResponse(responseConfig config.Response) MessagePrototype {
	bodyTemplate := buildBodyTemplate(responseConfig.Body)
	prototype := MessagePrototype{
		Alias:        responseConfig.Alias,
		BodyTemplate: bodyTemplate,
		Mediators:    buildMediators(bodyTemplate, responseConfig.Options),
		From:         responseConfig.ResponseRequest.To,
		To:           responseConfig.To,
		Matcher:      buildMatcher(responseConfig.ResponseRequest.Matchers),
	}

	// prepend matchingMediator to the mediators list
	// so that it is the first mediator to be executed
	// and all others can be skipped if it doesn't match
	prototype.Mediators.Prepend(
		mediator.NewMatchingMediator(
			matchers.ParsePattern(responseConfig.ResponseRequest.Match),
			prototype.Matcher))

	return prototype
}
