package sender

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/nicksdlc/macaw/communicators"
	"github.com/nicksdlc/macaw/config"
)

func TestSendSimplePostMessage(t *testing.T) {
	req := "{\"name\": \"test\"}"

	// Send response to be tested
	server := httptest.NewServer(createHandler(http.StatusOK, req))
	defer server.Close()

	sender := prepareSender(req, server)

	err := sender.Send()
	if err != nil {
		t.Fatalf("expected a no error, instead got: %s", err.Error())
	}
}

func TestRecievedServerError(t *testing.T) {
	req := "{\"name\": \"test\"}"

	// Send response to be tested
	server := httptest.NewServer(createHandler(http.StatusInternalServerError, req))
	defer server.Close()

	sender := prepareSender(req, server)

	err := sender.Send()
	if err == nil {
		t.Fatalf("expected an error, instead got: nil")
	}
}

func prepareSender(body string, server *httptest.Server) *MessageSender {
	sendRequest := config.Request{
		Body: config.Body{String: body},
		To:   server.URL,
		Options: config.Options{
			Quantity: 1,
			Delay:    1,
		},
	}

	communicator := communicators.NewHTTPCommunicator(server.URL, 0, server.Client())
	sender := NewMessageSender(communicator, sendRequest)
	return sender
}

func createHandler(status int, body string) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(status)
		b, _ := json.Marshal(body)
		rw.Write(b)
	}
}
