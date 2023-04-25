package number

import (
	"fmt"

	"github.com/nicksdlc/macaw/template/types"
)

type IncrementalNumber struct {
	value int
}

func incrementalNumber(params []string) types.Type {
	startingValue := 0
	if len(params) == 1 {
		fmt.Sscanf(params[0], "%d", &startingValue)
	}

	return &IncrementalNumber{
		value: startingValue,
	}
}

func (in *IncrementalNumber) Value() string {
	in.value++
	return fmt.Sprintf("%d", in.value)
}
