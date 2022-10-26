package types

import (
	"math/rand"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

var stringTypes = make(map[string]func([]string) Type)

func init() {
	stringTypes["variant"] = variantString
	stringTypes["any"] = anyString
}

func String(parameters []string) Type {
	if len(parameters) == 0 {
		return &AnyString{}
	}
	return stringTypes[parameters[0]](parameters[1:])
}

type AnyString struct {
}

func anyString([]string) Type {
	return &AnyString{}
}

func (as *AnyString) Value() string {
	b := make([]rune, 12)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

type VariantString struct {
	options []string
}

func variantString(variants []string) Type {
	return &VariantString{
		options: variants,
	}
}

func (os *VariantString) Value() string {
	return os.options[rand.Intn(len(os.options))]
}
