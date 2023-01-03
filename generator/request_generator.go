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

	amount := pr.Request.Amount

	base, err := os.ReadFile(pr.Request.File)
	if err != nil {
		panic(err)
	}
	request := template.NewRequest(string(base))

	for i := 0; i < amount; i++ {
		requests = append(requests, request.Create())
	}

	return requests
}
