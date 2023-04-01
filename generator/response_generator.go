package generator

import (
	"os"

	"github.com/nicksdlc/macaw/config"
	"github.com/nicksdlc/macaw/template"
)

// Responder generates responses
type Responder interface {
	Generate(request template.IncomingRequest) []string
}

// GenericResponder generates policy responses
type GenericResponder struct {
	Response     config.Response
	RespTemplate string
}

// Generate creates a slice of responses
func (pr *GenericResponder) Generate(request template.IncomingRequest) []string {
	var responses []string
	amount := pr.Response.Amount

	if pr.RespTemplate == "" {
		pr.RespTemplate = string(pr.readResponseTemplate())
	}
	base := pr.RespTemplate

	response := template.NewResponse(string(base), amount, &request)

	for i := 0; i < amount; i++ {
		responses = append(responses, response.Create())
	}

	return responses
}

func (pr *GenericResponder) readResponseTemplate() []byte {
	base, err := os.ReadFile(pr.Response.File)
	if err != nil {
		panic(err)
	}
	return base
}
