package responder

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"testing"

	"github.com/nicksdlc/macaw/communicators"
	"github.com/nicksdlc/macaw/config"
	"github.com/stretchr/testify/assert"
)

func TestShouldRespondWithPreSetResponseToMessage(t *testing.T) {
	expected := "\"name\": 10"

	port, err := getFreePort()
	if err != nil {
		t.Fatalf("No port is available: %s", err.Error())
	}
	configuredResponse := config.Response{
		ResponseRequest: config.ResponseRequest{
			To: "/test",
		},
		Body:    &config.Body{String: []string{"{\"name\": {{.FromRequestHeaders \"requestID\"}}}"}},
		Options: &config.Options{Quantity: 1, Delay: "0"},
	}

	sut := NewMessageResponder(communicators.NewHTTPCommunicator("localhost", uint16(port), nil), []config.Response{configuredResponse})
	sut.Listen()
	response, err := http.Get(fmt.Sprintf("http://localhost:%d/test?requestID=10", port))

	if err != nil {
		t.Fatalf("expected a no error, instead got: %s", err.Error())
	}

	assertResponseContentCorrect(t, response, expected)
}

func TestShouldNotRespondIfDoesNotMatch(t *testing.T) {
	port, err := getFreePort()
	if err != nil {
		t.Fatalf("No port is available: %s", err.Error())
	}
	fieldMatcher := config.Matcher{
		Type:  "field",
		Name:  "id",
		Value: "test",
	}
	configuredResponse := config.Response{
		ResponseRequest: config.ResponseRequest{
			To:       "/test",
			Matchers: []config.Matcher{fieldMatcher},
		},
		Body:    &config.Body{String: []string{"{\"name\": {{.FromRequestHeaders \"requestID\"}}}"}},
		Options: &config.Options{Quantity: 1, Delay: "0"},
	}

	sut := NewMessageResponder(communicators.NewHTTPCommunicator("localhost", uint16(port), nil), []config.Response{configuredResponse})
	sut.Listen()
	response, err := http.Get(fmt.Sprintf("http://localhost:%d/test?requestID=NOT_10", port))
	actual, _ := io.ReadAll(response.Body)

	if err != nil {
		t.Fatalf("expected a no error, instead got: %s", err.Error())
	}

	assert.Equal(t, string(actual), "")
}

func TestShouldRespondToDifferntRequest(t *testing.T) {
	expectedFromTest := "\"name\": 10"
	expectedFromTest2 := "\"name\": 42"

	port, err := getFreePort()
	if err != nil {
		t.Fatalf("No port is available: %s", err.Error())
	}
	configuredResponses := []config.Response{
		{
			ResponseRequest: config.ResponseRequest{
				To: "/test",
			},
			Body:    &config.Body{String: []string{"{\"name\": {{.FromRequestHeaders \"requestID\"}}}"}},
			Options: &config.Options{Quantity: 1, Delay: "0"},
		},
		{
			ResponseRequest: config.ResponseRequest{
				To: "/test2",
			},
			Body:    &config.Body{String: []string{"{\"name\": {{.FromRequestHeaders \"otherID\"}}}"}},
			Options: &config.Options{Quantity: 1, Delay: "0"},
		},
	}

	sut := NewMessageResponder(communicators.NewHTTPCommunicator("localhost", uint16(port), nil), configuredResponses)
	sut.Listen()
	response, _ := http.Get(fmt.Sprintf("http://localhost:%d/test?requestID=10", port))
	otherResponse, err := http.Get(fmt.Sprintf("http://localhost:%d/test2?otherID=42", port))

	if err != nil {
		t.Fatalf("expected a no error, instead got: %s", err.Error())
	}

	assertResponseContentCorrect(t, response, expectedFromTest)
	assertResponseContentCorrect(t, otherResponse, expectedFromTest2)
}

func TestSameEndpointShouldRespondWithDifferentResponses(t *testing.T) {
	expectedFromTest := "\"name\": 10"
	expectedFromTest2 := "\"name\": 42"

	port, err := getFreePort()
	if err != nil {
		t.Fatalf("No port is available: %s", err.Error())
	}
	configuredResponses := []config.Response{
		{
			ResponseRequest: config.ResponseRequest{
				To: "/test",
				Matchers: []config.Matcher{
					{
						Type:  "field",
						Name:  "requestID",
						Value: "10",
					},
				},
			},
			Body:    &config.Body{String: []string{"{\"name\": {{.FromRequestHeaders \"requestID\"}}}"}},
			Options: &config.Options{Quantity: 1, Delay: "0"},
		},
		{
			ResponseRequest: config.ResponseRequest{
				To: "/test",
			},
			Body:    &config.Body{String: []string{"{\"name\": {{.FromRequestHeaders \"otherID\"}}}"}},
			Options: &config.Options{Quantity: 1, Delay: "0"},
		},
	}

	sut := NewMessageResponder(communicators.NewHTTPCommunicator("localhost", uint16(port), nil), configuredResponses)
	sut.Listen()
	response, _ := http.Get(fmt.Sprintf("http://localhost:%d/test?requestID=10", port))
	otherResponse, err := http.Get(fmt.Sprintf("http://localhost:%d/test?otherID=42", port))

	if err != nil {
		t.Fatalf("expected a no error, instead got: %s", err.Error())
	}

	assertResponseContentCorrect(t, otherResponse, expectedFromTest2)
	assertResponseContentCorrect(t, response, expectedFromTest)
}

func assertResponseContentCorrect(t *testing.T, response *http.Response, expectedContent string) {
	results, _ := io.ReadAll(response.Body)
	if !strings.Contains(string(results), expectedContent) {
		t.Fatalf("expected a %s, instead got: %s", expectedContent, string(results))
	}
}

func getFreePort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return 0, err
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port, nil
}
