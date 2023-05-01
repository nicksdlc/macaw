package prototype

import (
	"github.com/nicksdlc/macaw/config"
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
	}

	return templates
}
