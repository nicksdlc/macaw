package builder

import (
	"os"

	"github.com/nicksdlc/macaw/config"
	"github.com/nicksdlc/macaw/generator"
	"github.com/nicksdlc/macaw/model"
	"github.com/nicksdlc/macaw/template"
)

// PrototypeBuilder builds message templates
type PrototypeBuilder interface {
	Build() []model.MessagePrototype
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
func (rtb *ResponsePrototypeBuilder) Build() []model.MessagePrototype {
	var templates []model.MessagePrototype

	for _, response := range rtb.responseConfig {
		templates = append(templates, model.MessagePrototype{
			BodyTemplate: rtb.buildBodyTemplate(response),
			Mediators:    rtb.buildMediators(response),
			From:         response.ResponseRequest.To,
			To:           response.To,
			Matcher:      rtb.buildMatcher(response),
		})
	}

	return templates
}

func (rtb *ResponsePrototypeBuilder) buildMediators(responseConfig config.Response) []model.Mediator {
	var mediators []model.Mediator

	mediators = append(mediators, func(message model.RequestMessage, response *model.ResponseMessage) error {
		// TODO: should be moved once figure out how to do that
		responder := &generator.GenericResponder{
			Response:     responseConfig,
			RespTemplate: rtb.buildBodyTemplate(responseConfig),
		}
		response.Responses = responder.Generate(template.Serialize(message.Headers, message.Body))
		return nil
	})

	return mediators
}

func (rtb *ResponsePrototypeBuilder) buildBodyTemplate(response config.Response) string {
	if response.File != "" {
		return string(readResponseTemplate(response.File))
	}
	return response.String
}

func (rtb *ResponsePrototypeBuilder) buildMatcher(response config.Response) []model.Matcher {
	if response.ResponseRequest.Field.Name != "" {
		return []model.Matcher{&model.FieldMatcher{
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
