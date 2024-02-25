package prototype

import (
	"testing"

	"github.com/nicksdlc/macaw/config"
	"github.com/stretchr/testify/assert"
)

func TestOnEmptyRequestConfigNoProtypesIsBuilt(t *testing.T) {
	// Given
	requestConfig := []config.Request{}

	// When
	requestPrototypes := NewRequestPrototypeBuilder(requestConfig).Build()

	// Then
	assert.Equal(t, 0, len(requestPrototypes))
}

func TestOnValidRequestConfigRequestPrototypeIsBuilt(t *testing.T) {
	// Given
	requestConfig := []config.Request{
		{
			Alias: "testAlias",
			To:    "testTo",
			Body: config.Body{
				String: []string{"testBody"},
			},
		},
	}

	// When
	requestPrototypes := NewRequestPrototypeBuilder(requestConfig).Build()

	// Then
	assert.Equal(t, 1, len(requestPrototypes))
	assert.Equal(t, "testAlias", requestPrototypes[0].Alias)
	assert.Equal(t, "testTo", requestPrototypes[0].To)
	assert.Equal(t, "testBody", requestPrototypes[0].BodyTemplate[0])
}
