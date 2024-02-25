package prototype

import "github.com/nicksdlc/macaw/config"

// RequestPrototypeBuilder is a builder for request prototypes
type RequestPrototypeBuilder struct {
	requestConfig []config.Request
}

// NewRequestPrototypeBuilder creates a new instance of RequestPrototypeBuilder
func NewRequestPrototypeBuilder(requestConfig []config.Request) *RequestPrototypeBuilder {
	return &RequestPrototypeBuilder{
		requestConfig: requestConfig,
	}
}

// Build creates a list of request prototypes
func (b *RequestPrototypeBuilder) Build() []MessagePrototype {
	var prototypes []MessagePrototype

	for _, request := range b.requestConfig {
		bodyTemplate := buildBodyTemplate(request.Body)
		prototypes = append(prototypes, MessagePrototype{
			Alias:        request.Alias,
			Type:         request.Type,
			BodyTemplate: bodyTemplate,
			Mediators:    buildMediators(bodyTemplate, request.Options),
			To:           request.To,
			From:         request.From,
		})
	}

	return prototypes
}
