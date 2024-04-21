package oapi

import (
	"testing"

	gen "github.com/nicksdlc/macaw/generator"
	"github.com/stretchr/testify/assert"
)

func TestShouldGenerateStrings(t *testing.T) {
	gen.InitDumb()
	model, _ := LoadModel("./testdata/multiple_string_fields.yaml")
	sut := NewGeneratorFromSchema(model.Model.Components.Schemas.First().Value().Schema())

	res := sut()

	assert.Equal(t, 3, len(res), "Should generate 1 field")
	gen.AssertStringField(t, res, "strField1", "those")
	gen.AssertStringField(t, res, "strField2", "here")
	gen.AssertStringField(t, res, "strField3", "over")
}
