package responder

import (
	"macaw/data"
	"macaw/template"
	"math/rand"
	"os"
)

// Responder generates responses
type Responder interface {
	Generate(request data.Request) []data.Response
}

// PolicyResponder generates policy responses
type PolicyResponder struct {
	TemplatePath string
}

// Generate creates a random amount of
func (pr *PolicyResponder) Generate(request template.Request) []string {
	var responses []string

	amount := rand.Intn(10) + 1

	base, err := os.ReadFile(pr.TemplatePath)
	if err != nil {
		panic(err)
	}
	response := template.NewResponse(string(base), amount, &request)

	for i := 0; i < amount; i++ {
		responses = append(responses, response.Create())
	}

	return responses
}
