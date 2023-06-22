package prototype

import (
	"testing"

	"github.com/nicksdlc/macaw/config"
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
	fieldMatcher := config.Matcher{
		Type:  "field",
		Name:  "id",
		Value: "test",
	}

	responseConfig := []config.Response{
		{
			ResponseRequest: config.ResponseRequest{
				To:       "test",
				Matchers: []config.Matcher{fieldMatcher},
			},
			Body: config.Body{String: "test"},
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
	fieldMatcher1 := config.Matcher{
		Type:  "field",
		Name:  "id",
		Value: "test",
	}
	fieldMatcher2 := config.Matcher{
		Type:  "field",
		Name:  "id2",
		Value: "test2",
	}
	responseConfig := []config.Response{
		{
			ResponseRequest: config.ResponseRequest{
				To:       "test",
				Matchers: []config.Matcher{fieldMatcher1},
			},
			Body: config.Body{String: "test"},
		},
		{
			ResponseRequest: config.ResponseRequest{
				To:       "test2",
				Matchers: []config.Matcher{fieldMatcher2},
			},
			Body: config.Body{String: "test2"},
		},
	}

	// When
	responsePrototypes := NewResponsePrototypeBuilder(responseConfig).Build()

	// Then
	assert.Equal(t, 2, len(responsePrototypes))
	assertPrototype(t, "test", responsePrototypes[0])
	assertPrototype(t, "test2", responsePrototypes[1])
}

func assertPrototype(t *testing.T, name string, responsePrototype MessagePrototype) {
	assert.Equal(t, name, responsePrototype.From)
	assert.Equal(t, name, responsePrototype.BodyTemplate)
	assert.Equal(t, 1, len(responsePrototype.Matcher))
}

func TestOneMediatorCreatedForResponse(t *testing.T) {
	// Given
	fieldMatcher := config.Matcher{
		Type:  "field",
		Name:  "id",
		Value: "test",
	}
	responseConfig := []config.Response{
		{
			ResponseRequest: config.ResponseRequest{
				To:       "test",
				Matchers: []config.Matcher{fieldMatcher},
			},
			Body: config.Body{String: "test"},
		},
	}

	// When
	responsePrototypes := NewResponsePrototypeBuilder(responseConfig).Build()

	// Then
	assert.NotNil(t, responsePrototypes[0].Mediators)
}