package generator

import (
	"os"

	"github.com/nicksdlc/macaw/config"
	"github.com/nicksdlc/macaw/data"
	"github.com/nicksdlc/macaw/template"
)

// Responder generates responses
type Responder interface {
	Generate(request data.Request) []string
}

// GenericResponder generates policy responses
type GenericResponder struct {
	Response config.Response
}

// Generate creates a slive of responses
func (pr *GenericResponder) Generate(request template.IncomingRequest) []string {
	var responses []string

	amount := pr.Response.Amount

	base, err := os.ReadFile(pr.Response.File)
	if err != nil {
		panic(err)
	}
	response := template.NewResponse(string(base), amount, &request)

	for i := 0; i < amount; i++ {
		responses = append(responses, response.Create())
	}

	return responses
}
