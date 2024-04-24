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

	"github.com/stretchr/testify/require"
)

func TestShouldRespondFromOAPISpecBasicTypes(t *testing.T) {
	schemaPath := "./testdata/single_endpoint_basic_types.yaml"
	// Note, that number is generated as float and requires way higher number as seed to be > 0
	// So, in some cases number fields will be 0, although still present in final data
	gen.InitDumb()
	expected := readExpectedJson("./testdata/expected_basic_types.json")
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
	require.JSONEq(t, expected, resultStr)
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

func readExpectedJson(name string) string {
	file, _ := os.ReadFile(name)
	return string(file)
}

func formatErrors(errs []*errors.ValidationError) string {
	var b bytes.Buffer
	for _, e := range errs {
		b.WriteString(fmt.Sprintf("%s : %s", e.Reason, e.Message))
	}

	return b.String()
}
