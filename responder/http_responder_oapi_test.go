package responder

import (
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/nicksdlc/macaw/communicators"
	"github.com/nicksdlc/macaw/config"
	gen "github.com/nicksdlc/macaw/generator"
	"github.com/stretchr/testify/assert"
)

func TestShouldRespondFromOAPISpec(t *testing.T) {
	gen.InitDumb()
	expected := "{\"testInt\":2}"
	port, err := getFreePort()
	if err != nil {
		t.Fatalf("No port is available: %s", err.Error())
	}
	configuredResponse := config.Response{
		FromOpenAPI: "./testdata/single_endpoint_single_field.yaml",
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
}
