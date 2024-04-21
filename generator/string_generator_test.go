package generator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShouldGenerateString(t *testing.T) {
	InitDumb()
	res := Data{}

	sut := GenerateString("teststring")
	sut(res, &Context{})

	assert.Equal(t, 1, len(res))
	// 'those' is what gofakeit generates as a word with Dumb random source and seed == 1
	assert.Equal(t, "those", res["teststring"])
}
