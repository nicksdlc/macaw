package number

import (
	"github.com/nicksdlc/macaw/template/types"
)

var numberTypes = make(map[string]func(params []string) types.Type)

func init() {
	numberTypes["incremental"] = incrementalNumber
	numberTypes["any"] = any
	numberTypes["between"] = betweenNumber
}

// Number returns a number type
// In template it is used like this:
// {{.Number "incremental" "10"}}
// {{.Number "any"}}
// or
// {{.Number "any" "100"}} for seed to be 100
// {{.Number "between" "10" "20"}}
func Number(params []string) types.Type {
	if len(params) == 0 {
		return any([]string{})
	}
	return numberTypes[params[0]](params[1:])
}
