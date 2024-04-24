package oapi

import (
	"testing"

	gen "github.com/nicksdlc/macaw/generator"
)

func TestShoudGenerateFloatField(t *testing.T) {
	gen.InitDumbWithSeed(8888888888888888)
	model, _ := LoadModel("./testdata/single_number_field.yaml")

	sut := NewGeneratorFromSchema(model.Model.Components.Schemas.First().Value().Schema())
	res := sut()

	gen.AssertFloat32Field(t, res, "float32-test", 0.12335807)
}
