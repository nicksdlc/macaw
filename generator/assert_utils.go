package generator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func AssertIntField(t *testing.T, data Data, field string, expected int) {
	if value, ok := data[field].(int); ok {
		assert.Equal(t, expected, value, "field '%s' should be %d", field, expected)
	} else {
		assert.Failf(t, "Failed to cast", "field '%s' is not an int", field)
	}
}

func AssertStringField(t *testing.T, data Data, field string, expected string) {
	if value, ok := data[field].(string); ok {
		assert.Equal(t, expected, value, "field '%s' should be %d", field, expected)
	} else {
		assert.Failf(t, "Failed to cast", "field '%s' is not a string", field)
	}
}
