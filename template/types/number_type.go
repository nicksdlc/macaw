package types

import (
	"fmt"
	"math/rand"
)

var numberTypes = make(map[string]func() Type)

func init() {
	numberTypes["incremental"] = incrementalNumber
	numberTypes["any"] = newNumber
}

func Number(params []string) Type {
	if len(params) == 0 {
		return &AnyNumber{}
	}
	return numberTypes[params[0]]()
}

type AnyNumber struct {
}

func newNumber() Type {
	return &AnyNumber{}
}

func (an *AnyNumber) Value() string {
	return fmt.Sprintf("%d", rand.Intn(10000000))
}

type IncrementalNumber struct {
	value int
}

func incrementalNumber() Type {
	return &IncrementalNumber{
		value: 0,
	}
}

func (in *IncrementalNumber) Value() string {
	in.value++
	return fmt.Sprintf("%d", in.value)
}
