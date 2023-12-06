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
	// Given
	req := "{\"name\": \"test\"}"

	server := httptest.NewServer(createHandler(http.StatusOK, req))
	request := []config.Request{createRequest("test", req, 1, server)}
	defer server.Close()

	// When
	sender := prepareSender(request, server)
	result, err := sender.SendWithResponse()

	// Then
	if err != nil {
		t.Fatalf("expected a no error, instead got: %s", err.Error())
	}
	for r := range result {
		if r.Metadata.Status() != "200 OK" {
			t.Fatalf("expected a 200 status, instead got: %s", r.Metadata["status"])
		}
	}
}

func TestSimpleGetMessageWhenResponseIs500(t *testing.T) {
	// Given
	req := "{\"name\": \"test\"}"

	server := httptest.NewServer(createHandler(http.StatusInternalServerError, req))
	request := []config.Request{createRequest("test", req, 1, server)}
	defer server.Close()

	// When
	sender := prepareSender(request, server)
	result, err := sender.SendWithResponse()
	count := 0

	// Then
	if err != nil {
		t.Fatalf("expected a no error, instead got: %s", err.Error())
	}
	for r := range result {
		count++
		if r.Metadata.Status() != "500 Internal Server Error" {
			t.Fatalf("expected a 500 status, instead got: %s", r.Metadata["status"])
		}
	}
	if count != 1 {
		t.Fatalf("expected a 1 result, instead got: %d", count)
	}
}

func TestMultipleGetMessagesWithResponse200Ok(t *testing.T) {
	// Given
	req := "{\"name\": \"test\"}"

	server := httptest.NewServer(createHandler(http.StatusOK, req))
	request := []config.Request{createRequest("test", req, 2, server)}
	defer server.Close()

	// When
	sender := prepareSender(request, server)
	result, err := sender.SendWithResponse()
	count := 0

	// Then
	if err != nil {
		t.Fatalf("expected a no error, instead got: %s", err.Error())
	}
	for r := range result {
		count++
		if r.Metadata.Status() != "200 OK" {
			t.Fatalf("expected a 200 status, instead got: %s", r.Metadata["status"])
		}
	}
	if count != 2 {
		t.Fatalf("expected a 2 result, instead got: %d", count)
	}
}

func TestDifferentRequestsWithBoth200and400(t *testing.T) {
	// Given
	req1 := "{\"name\": \"test\"}"
	req2 := "{\"name\": \"test2\"}"

	server := httptest.NewServer(createHandler(http.StatusOK, req1))
	server2 := httptest.NewServer(createHandler(http.StatusBadRequest, req2))
	request := []config.Request{createRequest("test", req1, 1, server), createRequest("test2", req2, 1, server2)}
	defer server.Close()
	defer server2.Close()

	// When
	sender := prepareSender(request, server)
	result, err := sender.SendWithResponse()
	count := 0

	// Then
	if err != nil {
		t.Fatalf("expected a no error, instead got: %s", err.Error())
	}
	for r := range result {
		count++
		if r.Metadata.Status() != "200 OK" && r.Metadata.Status() != "400 Bad Request" {
			t.Fatalf("expected a 200 or 400 status, instead got: %s", r.Metadata["status"])
		}
	}
	if count != 2 {
		t.Fatalf("expected a 2 result, instead got: %d", count)
	}
}

func prepareSender(requests []config.Request, server *httptest.Server) *MessageSender {
	communicator := communicators.NewHTTPCommunicator(server.URL, 0, server.Client())
	sender := NewMessageSender(communicator, requests)
	return sender
}

func createRequest(path, body string, quantity int, server *httptest.Server) config.Request {
	sendRequest := config.Request{
		Body: config.Body{String: []string{body}},
		Type: "GET",
		To:   path,
		Options: config.Options{
			Quantity: quantity,
			Delay:    "1",
		},
	}
	return sendRequest
}

func createHandler(status int, body string) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(status)
		b, _ := json.Marshal(body)
		rw.Write(b)
	}
}
