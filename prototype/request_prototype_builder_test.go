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
			To: "test",
			Body: config.Body{
				String: "test",
			},
		},
	}

	// When
	requestPrototypes := NewRequestPrototypeBuilder(requestConfig).Build()

	// Then
	assert.Equal(t, 1, len(requestPrototypes))
	assert.Equal(t, "test", requestPrototypes[0].To)
	assert.Equal(t, "test", requestPrototypes[0].BodyTemplate)
}
