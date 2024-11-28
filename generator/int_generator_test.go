package generator

import (
	"testing"
)

func TestUintShouldSet2(t *testing.T) {
	InitDumb()
	res := Data{}

	sut := GenerateInt("testuint")
	sut(res, &Context{})

	AssertIntField(t, res, "testuint", 2)
}
