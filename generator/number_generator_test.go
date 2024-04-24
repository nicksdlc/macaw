package generator

import (
	"testing"
)

func TestShouldGenerateFloat(t *testing.T) {
	// Seed == 1 produces 0.0, so need some other value to produce higher number
	InitDumbWithSeed(9999999999999999999)

	data := Data{}
	sut := GenerateFloat("float-test")
	sut(data, &Context{})

	AssertFloat32Field(t, data, "float-test", 0.77787805)

}
