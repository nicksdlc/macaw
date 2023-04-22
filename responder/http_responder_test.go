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
		File:   "",
		Amount: 1,
	}

	sut := NewMessageResponder(communicators.NewHTTPCommunicator("localhost", uint16(port), nil), configuredResponse)
	sut.responder.RespTemplate = "{\"name\": {{.FromRequestHeaders \"requestID\"}}}"
	sut.Listen()
	response, err := http.Get(fmt.Sprintf("http://localhost:%d/test?requestID=10", port))

	if err != nil {
		t.Fatalf("expected a no error, instead got: %s", err.Error())
	}

	results, _ := io.ReadAll(response.Body)
	if !strings.Contains(string(results), expected) {
		t.Fatalf("expected a %s, instead got: %s", expected, string(results))
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
