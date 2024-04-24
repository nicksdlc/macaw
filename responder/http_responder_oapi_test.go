package responder

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/nicksdlc/macaw/communicators"
	"github.com/nicksdlc/macaw/config"
	gen "github.com/nicksdlc/macaw/generator"
	"github.com/pb33f/libopenapi"
	validator "github.com/pb33f/libopenapi-validator"
	"github.com/pb33f/libopenapi-validator/errors"
	"github.com/stretchr/testify/assert"
)

func TestShouldRespondFromOAPISpecBasicTypes(t *testing.T) {
	schemaPath := "./testdata/single_endpoint_basic_types.yaml"
	gen.InitDumb()
	expected := "{\"testInt\":2,\"testString\":\"it\"}"
	port, err := getFreePort()
	if err != nil {
		t.Fatalf("No port is available: %s", err.Error())
	}
	configuredResponse := config.Response{
		FromOpenAPI: schemaPath,
		ResponseRequest: config.ResponseRequest{
			To: "/test-oapi",
		},
		Options: &config.Options{Quantity: 1, Delay: "0"},
	}

	sut := NewMessageResponder(communicators.NewHTTPCommunicator("localhost", uint16(port), nil), []config.Response{configuredResponse})
	sut.Listen()
	response, err := http.Get(fmt.Sprintf("http://localhost:%d/test-oapi/example", port))

	if err != nil {
		t.Fatalf("expected no error, instead got: %s", err.Error())
	}

	result, _ := io.ReadAll(response.Body)
	resultStr := string(result)
	assert.Equal(t, expected, resultStr)
	validateAgainstSchema(t, schemaPath, response)
}

func validateAgainstSchema(t *testing.T, schemaPath string, response *http.Response) {
	schFile, _ := os.ReadFile(schemaPath)

	document, err := libopenapi.NewDocument(schFile)

	if !assert.NoError(t, err, "Failed to read schema file %s", schemaPath) {
		return
	}

	rbValidator, _ := validator.NewValidator(document)

	_, errors := rbValidator.ValidateHttpResponse(response.Request, response)

	assert.Equal(t, 0, len(errors), "Response has %d errors\n %s", len(errors), formatErrors(errors))

}

func formatErrors(errs []*errors.ValidationError) string {
	var b bytes.Buffer
	for _, e := range errs {
		b.WriteString(fmt.Sprintf("%s : %s", e.Reason, e.Message))
	}

	return b.String()
}
