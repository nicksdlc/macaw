package connectors

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// HTTPConnector is a connector to provide http infrastructure for the
type HTTPConnector struct {
	endpoint    string
	httpClient  *http.Client
	requestType string
}

// NewHTTPConnector creates new connector to send requests
func NewHTTPConnector(endpoint, requestType string, client *http.Client) HTTPConnector {
	return HTTPConnector{
		endpoint:    endpoint,
		requestType: requestType,
		httpClient:  client,
	}
}

// Close closes connection
func (m HTTPConnector) Close() error {

	panic("not implemented") // TODO: Implement
}

// Post posts requests to endpoint
func (m HTTPConnector) Post(body string) error {
	_, err := m.sendRequest(body)

	return err
}

// PostWithResponse posts requests to endpoint and gets response
func (m HTTPConnector) PostWithResponse(body string) (string, error) {
	response, err := m.sendRequest(body)
	if err != nil {
		return "", err
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(response.Body)

	return buf.String(), nil
}

func (m HTTPConnector) sendRequest(body string) (*http.Response, error) {
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
