package communicators

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/nicksdlc/macaw/model"
)

// HTTPCommunicator is a communicator used to send and receive HTTP requests
type HTTPCommunicator struct {
	serveEndpoint   string
	port            uint16
	httpClient      *http.Client
	responseHandler map[string]func(w http.ResponseWriter, r *http.Request)
}

// NewHTTPCommunicator creates new communicator to send requests
func NewHTTPCommunicator(endpoint string, port uint16, client *http.Client) *HTTPCommunicator {
	return &HTTPCommunicator{
		serveEndpoint: endpoint,
		port:          port,
		httpClient:    client,
	}
}

// RespondWith sets the response handler
func (m *HTTPCommunicator) RespondWith(responses []model.MessagePrototype) {
	m.responseHandler = make(map[string]func(w http.ResponseWriter, r *http.Request))

	for _, response := range responses {
		// re-assignment is required since, if not done here - will always point to last response
		res := response
		m.responseHandler[response.Destination] = func(w http.ResponseWriter, r *http.Request) {
			message := model.RequestMessage{
				Body:    getRequestBody(r),
				Headers: getQueryParams(r.URL),
			}

			resp := model.ResponseMessage{}
			for _, mediator := range res.Mediators {
				mediator(message, &resp)
			}

			if m.matchAny(res, message) {
				io.WriteString(w, resp.Responses[0])
			}
		}
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
		mux.HandleFunc(m.serveEndpoint, func(w http.ResponseWriter, r *http.Request) {
			message := model.RequestMessage{
				Body:    getRequestBody(r),
				Headers: getQueryParams(r.URL),
			}

			resp := model.ResponseMessage{}
			for _, mediator := range mediators {
				mediator(message, &resp)
			}
			io.WriteString(w, resp.Responses[0])
		})

		http.ListenAndServe(fmt.Sprintf(":%d", m.port), mux)
	}()
}

// ConsumeMediateReplyWithResponse listens as a server and consumes sends messages to the channel
func (m *HTTPCommunicator) ConsumeMediateReplyWithResponse() {
	go func() {
		mux := http.NewServeMux()
		for endpoint, handler := range m.responseHandler {
			log.Printf("Listening on endpoint: %s", endpoint)
			mux.HandleFunc(endpoint, handler)
		}
		http.ListenAndServe(fmt.Sprintf(":%d", m.port), mux)
	}()

}

func (m HTTPCommunicator) sendRequest(body string) (*http.Response, error) {
	resourceBytes, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	response, err := m.httpClient.Post(m.serveEndpoint, "application/json", bytes.NewBuffer(resourceBytes))
	if err != nil {
		return nil, err
	}

	if response.StatusCode >= 500 {
		return nil, fmt.Errorf("internal Server Error")
	}

	return response, nil
}

func (*HTTPCommunicator) matchAny(res model.MessagePrototype, message model.RequestMessage) bool {
	if len(res.Matcher) == 0 {
		return true
	}

	for _, matcher := range res.Matcher {
		if matcher.Match(message) {
			return true
		}
	}
	return false
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
