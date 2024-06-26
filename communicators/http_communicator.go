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
	"github.com/nicksdlc/macaw/prototype"
	"github.com/nicksdlc/macaw/prototype/matchers"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	responseGauge = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "macaw_response_current",
		Help: "The number of currently processed responses",
	})
)

// HTTPCommunicator is a communicator used to send and receive HTTP requests
type HTTPCommunicator struct {
	serveEndpoint   string
	port            uint16
	httpClient      *http.Client
	responseHandler map[string]*dynamicHandler
	requests        []prototype.MessagePrototype
	responses       []prototype.MessagePrototype
	splitResponses  map[string][]*prototype.MessagePrototype
	mux             *http.ServeMux
}

// NewHTTPCommunicator creates new communicator to send requests
func NewHTTPCommunicator(endpoint string, port uint16, client *http.Client) *HTTPCommunicator {
	return &HTTPCommunicator{
		serveEndpoint: endpoint,
		port:          port,
		httpClient:    client,
	}
}

func (m *HTTPCommunicator) GetResponses() []prototype.MessagePrototype {
	return m.responses
}

func (m *HTTPCommunicator) UpdateResponse(response prototype.MessagePrototype) {
	for i, resp := range m.splitResponses[response.From] {
		if resp.Alias == response.Alias {
			m.splitResponses[response.From][i] = &response
		}
	}
}

// RespondWith sets the response handler
func (m *HTTPCommunicator) RespondWith(responses []prototype.MessagePrototype) {
	m.responses = responses

	if m.responseHandler == nil {
		m.responseHandler = make(map[string]*dynamicHandler)
	}

	m.splitResponses = splitResponsesByEndpoint(m.responses)

	for _, response := range m.splitResponses {
		// re-assignment is required since, if not done here - will always point to last response
		res := response

		handlerFunc := func(w http.ResponseWriter, req *http.Request) {
			responseGauge.Inc()
			defer responseGauge.Dec()

			message := model.RequestMessage{
				Body:    getRequestBody(req),
				Headers: getQueryParams(req.URL),
			}

			for _, respPrototype := range res {
				resp := model.ResponseMessage{}
				for r := range respPrototype.Mediators.Run(message, resp) {
					if matchers.MatchAny(respPrototype.Matcher, message) {
						if contentType, present := r.Metadata["Content-Type"]; present && contentType != "" {
							w.Header().Add("Content-Type", contentType)
						}
						io.WriteString(w, r.Body)
						return
					}
				}
			}
		}

		if m.responseHandler[response[0].From] == nil {
			m.responseHandler[response[0].From] = &dynamicHandler{
				handler: handlerFunc,
			}
		} else {
			m.responseHandler[response[0].From].handler = handlerFunc
		}
	}
}

// RequestWith sets the request handler
func (m *HTTPCommunicator) RequestWith(requests []prototype.MessagePrototype) {
	m.requests = requests
}

// Close closes connection
func (m *HTTPCommunicator) Close() error {
	panic("implement me")
}

// PostAndListen posts requests to endpoint and listens for response
func (m *HTTPCommunicator) PostAndListen() (chan model.ResponseMessage, error) {
	responseChan := make(chan model.ResponseMessage)
	signalChannel := make(chan struct{})

	for _, request := range m.requests {
		go func(request prototype.MessagePrototype, sigChannel chan struct{}) {
			resp := model.ResponseMessage{}
			for r := range request.Mediators.Run(model.RequestMessage{}, resp) {
				result, _ := m.postWithResponse(request, r)
				responseChan <- result
			}
			// upon compeltion of requests send for specific prortotype, send signal to channel
			sigChannel <- struct{}{}
		}(request, signalChannel)
	}

	go m.waitForResponses(signalChannel, responseChan)

	return responseChan, nil
}

// ConsumeMediateReplyWithResponse listens as a server and consumes sends messages to the channel
func (m *HTTPCommunicator) ConsumeMediateReplyWithResponse() {
	go func() {
		m.mux = http.NewServeMux()
		for endpoint, handler := range m.responseHandler {
			log.Printf("Listening on endpoint: %s", endpoint)
			m.mux.Handle(endpoint, handler)
		}
		http.ListenAndServe(fmt.Sprintf(":%d", m.port), m.mux)
	}()
}

func (m *HTTPCommunicator) waitForResponses(signalChannel chan struct{}, responseChan chan model.ResponseMessage) {
	count := 0
	for sig := range signalChannel {
		if sig == struct{}{} {
			count++
		}
		if len(m.requests) == count {
			close(responseChan)
			close(signalChannel)
		}
	}
}

func (m *HTTPCommunicator) postWithResponse(request prototype.MessagePrototype, r model.ResponseMessage) (model.ResponseMessage, error) {
	resp, err := m.sendRequest(request.To, request.Type, r.Body)
	if err != nil {
		log.Printf("Error sending request: %s", err)
		return model.ResponseMessage{
			Body:     "",
			Metadata: model.NewResponseMetadata(resp.Status),
		}, err
	}

	result := model.ResponseMessage{
		Body:     readResponseBody(resp),
		Metadata: model.NewResponseMetadata(resp.Status),
	}

	return result, err
}

func (m HTTPCommunicator) sendRequest(path, requestType, body string) (*http.Response, error) {
	resourceBytes, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	serveUrl, err := url.Parse(m.serveEndpoint + "/" + path)
	if err != nil {
		return nil, err
	}

	response, err := m.httpClient.Do(&http.Request{
		Method: requestType,
		URL:    serveUrl,
		Body:   io.NopCloser(bytes.NewBuffer(resourceBytes)),
	})

	if err != nil {
		return nil, err
	}

	return response, nil
}

func splitResponsesByEndpoint(responses []prototype.MessagePrototype) map[string][]*prototype.MessagePrototype {
	endpointSlices := make(map[string][]*prototype.MessagePrototype)
	for _, response := range responses {
		res := response
		endpointSlices[response.From] = append(endpointSlices[response.From], &res)
	}
	return endpointSlices
}

func readResponseBody(response *http.Response) string {
	buf := new(bytes.Buffer)
	buf.ReadFrom(response.Body)

	return buf.String()
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

type dynamicHandler struct {
	handler func(w http.ResponseWriter, r *http.Request)
}

func (d *dynamicHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	d.handler(w, r)
}
