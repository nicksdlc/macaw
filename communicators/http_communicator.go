package communicators

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/nicksdlc/macaw/model"
)

// HTTPCommunicator is a communicator to provide http infrastructure for the
type HTTPCommunicator struct {
	endpoint   string
	port       uint16
	httpClient *http.Client
}

// NewHTTPCommunicator creates new communicator to send requests
func NewHTTPCommunicator(endpoint string, port uint16, client *http.Client) *HTTPCommunicator {
	return &HTTPCommunicator{
		endpoint:   endpoint,
		port:       port,
		httpClient: client,
	}
}

// Close closes connection
func (m *HTTPCommunicator) Close() error {
	panic("not implemented") // TODO: Implement
}

// Post posts requests to endpoint
func (m *HTTPCommunicator) Post(body model.RequestMessage) error {
	_, err := m.sendRequest(string(body.Body))

	return err
}

// PostWithResponse posts requests to endpoint and gets response
func (m *HTTPCommunicator) PostWithResponse(body model.RequestMessage) (string, error) {
	response, err := m.sendRequest(string(body.Body))
	if err != nil {
		return "", err
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(response.Body)

	return buf.String(), nil
}

// Consume listens as a server and consumes sends messages to the channel
func (m *HTTPCommunicator) Consume() <-chan model.RequestMessage {
	msgs := make(chan model.RequestMessage)

	return msgs
}

// ConsumeMediateReply listens as a server and consumes sends messages to the channel
func (m *HTTPCommunicator) ConsumeMediateReply(mediators []model.Mediator) {
	go func() {
		mux := http.NewServeMux()
		mux.HandleFunc(m.endpoint, func(w http.ResponseWriter, r *http.Request) {
			message := model.RequestMessage{
				Body:    getRequestBody(r),
				Headers: getQueryParams(r.URL),
			}

			resp := model.ResponseMessage{}
			for _, mediator := range mediators {
				resp, _ = mediator(message, resp)
			}
			io.WriteString(w, resp.Responses[0])
		})

		http.ListenAndServe(fmt.Sprintf(":%d", m.port), mux)
	}()
}

func (m HTTPCommunicator) sendRequest(body string) (*http.Response, error) {
	resourceBytes, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	response, err := m.httpClient.Post(m.endpoint, "application/json", bytes.NewBuffer(resourceBytes))
	if err != nil {
		return nil, err
	}

	if response.StatusCode >= 500 {
		return nil, fmt.Errorf("Internal Server Error")
	}

	return response, nil
}

func getQueryParams(url *url.URL) map[string]string {
	var params = make(map[string]string)
	for k, v := range url.Query() {
		params[k] = v[0]
	}
	return params
}

func getRequestBody(r *http.Request) []byte {
	if b, err := io.ReadAll(r.Body); err == nil {
		return b
	}

	return nil
}
