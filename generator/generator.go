package generator

import (
	"github.com/brianvoe/gofakeit/v7"
	"github.com/brianvoe/gofakeit/v7/source"
)

// Context holds the context information, like current depth or settings for the process
type Context struct {
}

// Data is map to which generator puts data
type Data map[string]interface{}

// Generator is function that fills provided Data map with generated values
type Generator func(Data, *Context)

// ObjGenerator is a generator for object
type ObjGenerator func() Data

// NewObjGenerator returns an ObjGenerator function. This function creates instance of Data and Context
// Then calls generator provided as argument and returns Data filled with generated values.
func NewObjGenerator(g Generator) ObjGenerator {
	return func() Data {
		data := Data{}
		ctx := &Context{}
		g(data, ctx)
		return data
	}
}

// GenerateUint generates uint
func GenerateUint(field string) Generator {
	return func(data Data, _ *Context) {
		data[field] = gofakeit.Uint()
	}
}

// GenerateString generates a random word
func GenerateString(field string) Generator {
	return func(data Data, _ *Context) {
		data[field] = gofakeit.Word()
	}
}

// Compose is a helper function that composes 2 generators. If one is nil, returns other without changes.
// if both are not nil - creates a wrapper function that calls first generator and then calls second one
func Compose(f Generator, g Generator) Generator {
	if f == nil {
		return g
	}

	if g == nil {
		return f
	}

	return func(data Data, ctx *Context) {
		f(data, ctx)
		g(data, ctx)
	}
}

// InitDumb initializes generator in predictable mode for tests
func InitDumb() {
	gofakeit.GlobalFaker = gofakeit.NewFaker(source.NewDumb(1), false)
}
