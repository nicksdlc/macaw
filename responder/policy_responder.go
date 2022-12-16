package responder

import (
	"macaw/config"
	"macaw/data"
	"macaw/template"
	"os"
)

// Responder generates responses
type Responder interface {
	Generate(request data.Request) []data.Response
}

// GenericResponder generates policy responses
type GenericResponder struct {
	Response config.Response
}

// Generate creates a random amount of
func (pr *GenericResponder) Generate(request template.Request) []string {
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
