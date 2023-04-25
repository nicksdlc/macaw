package number

import (
	"fmt"
	"math/rand"
	"strconv"

	"github.com/nicksdlc/macaw/template/types"
)

// BetweemNumber type for number in between min and max
type BetweemNumber struct {
	min int
	max int
}

func betweenNumber(params []string) types.Type {
	if len(params) != 2 {
		return &AnyNumber{}
	}

	min, err := strconv.Atoi(params[0])
	if err != nil {
		return &AnyNumber{}
	}

	max, err := strconv.Atoi(params[1])
	if err != nil {
		return &AnyNumber{}
	}

	return &BetweemNumber{
		min: min,
		max: max,
	}
}

// Value for between number returns random number between min and max
func (bn *BetweemNumber) Value() string {
	return fmt.Sprintf("%d", rand.Intn(bn.max-bn.min+1)+bn.min)
}
