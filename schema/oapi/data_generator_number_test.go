package oapi

import (
	"testing"

	gen "github.com/nicksdlc/macaw/generator"
	"github.com/stretchr/testify/assert"
)

func TestShouldGenerateSingleField(t *testing.T) {
	gen.InitDumb()
	model, _ := LoadModel("./testdata/single_int_field.yaml")
	sut := NewGeneratorFromSchema(model.Model.Components.Schemas.First().Value().Schema())

	res := sut()

	assert.Equal(t, 1, len(res), "Should generate 1 field")
	gen.AssertUintField(t, res, "number", 2)
}

func TestShouldGenerateMultipleFields(t *testing.T) {
	gen.InitDumb()
	model, _ := LoadModel("./testdata/multiple_int_fields.yaml")
	sut := NewGeneratorFromSchema(model.Model.Components.Schemas.First().Value().Schema())

	res := sut()

	assert.Equal(t, 3, len(res), "Should generate 3 fields")

	gen.AssertUintField(t, res, "field1", 2)
	gen.AssertUintField(t, res, "field2", 3)
	gen.AssertUintField(t, res, "field3", 4)
}
