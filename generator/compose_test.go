package generator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShouldCallFirstIfSecondNil(t *testing.T) {
	var fCalled bool = false
	f := func(_ Data, _ *Context) {
		fCalled = true
	}

	sut := Compose(f, nil)
	sut(nil, nil)

	assert.True(t, fCalled)
}

func TestShouldCallSecondIfFirstNil(t *testing.T) {
	var gCalled bool = false
	g := func(_ Data, _ *Context) {
		gCalled = true
	}

	sut := Compose(nil, g)
	sut(nil, nil)

	assert.True(t, gCalled)
}

func TestShouldCallBothIfNotNil(t *testing.T) {
	var fCalled bool = false
	var gCalled bool = false
	f := func(_ Data, _ *Context) {
		fCalled = true
	}
	g := func(_ Data, _ *Context) {
		gCalled = true
	}

	sut := Compose(f, g)
	sut(nil, nil)

	assert.True(t, fCalled, "First should be called")
	assert.True(t, gCalled, "Second should be called")

}
