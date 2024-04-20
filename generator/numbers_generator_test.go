package generator

import (
	"testing"
)

func TestUintShouldSet2(t *testing.T) {
	InitDumb()
	res := Data{}

	sut := GenerateUint("testuint")
	sut(res, &Context{})

	AssertUintField(t, res, "testuint", 2)
}
