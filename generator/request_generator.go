package generator

import (
	"os"

	"github.com/nicksdlc/macaw/config"
	"github.com/nicksdlc/macaw/template"
)

// Requester interface to generates requests
type Requester interface {
	Generate() []string
}

// JSONRequester generates json requests
type JSONRequester struct {
	Request config.Request
}

// Generate creates a list of requests
func (pr *JSONRequester) Generate() []string {
	var requests []string

	quantity := pr.Request.Options.Quantity

	base := pr.Request.Body.String
	if base == "" {
		baseBytes, err := os.ReadFile(pr.Request.Body.File)
		if err != nil {
			panic(err)
		}
		base = string(baseBytes)
	}
	request := template.NewRequest(base)

	for i := 0; i < quantity; i++ {
		requests = append(requests, request.Create())
	}

	return requests
}
