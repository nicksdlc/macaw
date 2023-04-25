package number

import (
	"fmt"

	"math/rand"

	"github.com/nicksdlc/macaw/template/types"
)

type AnyNumber struct {
	seed int
}

func any(params []string) types.Type {
	defaultSeed := 1000000000 // 1 billion

	if len(params) == 1 { // if seed is provided
		fmt.Sscanf(params[0], "%d", &defaultSeed)
	}

	return &AnyNumber{
		seed: int(defaultSeed),
	}
}

func (an *AnyNumber) Value() string {
	return fmt.Sprintf("%d", rand.Intn(an.seed))
}
