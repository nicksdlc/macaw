package builder

import (
	"testing"

	"github.com/nicksdlc/macaw/config"
	"github.com/nicksdlc/macaw/model"
	"github.com/stretchr/testify/assert"
)

func TestNoResponseCreatedOnEmptyConfigurtion(t *testing.T) {
	// Given
	responseConfig := []config.Response{}

	// When
	responsePrototypes := NewResponsePrototypeBuilder(responseConfig).Build()

	// Then
	assert.Equal(t, 0, len(responsePrototypes))
}

func TestResponseCreatedOnValidConfigurtion(t *testing.T) {
	// Given
	responseConfig := []config.Response{
		{
			ResponseRequest: config.ResponseRequest{
				To: "test",
				Field: config.Field{
					Name:  "id",
					Value: "test",
				},
			},
			String: "test",
		},
	}

	// When
	responsePrototypes := NewResponsePrototypeBuilder(responseConfig).Build()

	// Then
	assert.Equal(t, 1, len(responsePrototypes))
	assertPrototype(t, "test", responsePrototypes[0])

}

func TestMultipleResponsesCreatedOnValidConfiguration(t *testing.T) {
	// Given
	responseConfig := []config.Response{
		{
			ResponseRequest: config.ResponseRequest{
				To: "test",
				Field: config.Field{
					Name:  "id",
					Value: "test",
				},
			},
			String: "test",
		},
		{
			ResponseRequest: config.ResponseRequest{
				To: "test2",
				Field: config.Field{
					Name:  "id",
					Value: "test2",
				},
			},
			String: "test2",
		},
	}

	// When
	responsePrototypes := NewResponsePrototypeBuilder(responseConfig).Build()

	// Then
	assert.Equal(t, 2, len(responsePrototypes))
	assertPrototype(t, "test", responsePrototypes[0])
	assertPrototype(t, "test2", responsePrototypes[1])
}

func assertPrototype(t *testing.T, name string, responsePrototype model.MessagePrototype) {
	assert.Equal(t, name, responsePrototype.From)
	assert.Equal(t, name, responsePrototype.BodyTemplate)
	assert.Equal(t, 1, len(responsePrototype.Mediators))
	assert.Equal(t, 1, len(responsePrototype.Matcher))
}

func TestOneMediatorCreatedForResponse(t *testing.T) {
	// Given
	responseConfig := []config.Response{
		{
			ResponseRequest: config.ResponseRequest{
				To: "test",
				Field: config.Field{
					Name:  "id",
					Value: "test",
				},
			},
			String: "test",
		},
	}

	// When
	responsePrototypes := NewResponsePrototypeBuilder(responseConfig).Build()

	// Then
	assert.Equal(t, 1, len(responsePrototypes[0].Mediators))
}
