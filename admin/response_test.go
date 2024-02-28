package admin

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/nicksdlc/macaw/config"
	"github.com/nicksdlc/macaw/prototype"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func TestGetResponsesReturnsListOfResponses(t *testing.T) {
	// Given
	communicator := &MockCommunicator{}
	responsesCfg := []config.Response{
		{
			Alias: "test1",
		},
		{
			Alias: "test2",
		},
	}

	// When
	sut := ResponsesEndpoint(communicator, &responsesCfg)

	req := httptest.NewRequest(http.MethodGet, "/responses", nil)
	rr := httptest.NewRecorder()

	http.HandlerFunc(sut.Function).ServeHTTP(rr, req)

	// Read body of the response and parse it to the list of responses
	body, err := io.ReadAll(rr.Body)
	if err != nil {
		t.Fatalf("Could not read response body: %v", err)
	}
	var responses []config.Response
	err = yaml.Unmarshal(body, &responses)
	if err != nil {
		t.Fatalf("Could not unmarshal response body: %v", err)
	}

	// Then
	assert.Equal(t, http.StatusOK, rr.Code, "Expected HTTP status code 200 OK")
	assert.Equal(t, responsesCfg[0].Alias, responses[0].Alias, "Expected first response to have alias 'test1'")
	assert.Equal(t, responsesCfg[1].Alias, responses[1].Alias, "Expected second response to have alias 'test2'")
}

func TestUpdateAllResponsesInCommunicator(t *testing.T) {
	// Given
	communicator := &MockCommunicator{}
	responsesCfg := []config.Response{
		{
			Alias: "test1",
			Body: &config.Body{
				String: []string{"test1"},
			},
		},
		{
			Alias: "test2",
		},
	}

	sut := ResponsesEndpoint(communicator, &responsesCfg)

	// Read body of the response and parse it to the list of responses

	responseRequestCfg := []config.Response{
		{
			Alias: "test1",
			Body: &config.Body{
				String: []string{"updated-test1", "test2"},
			},
		},
		{
			Alias: "test2",
			Body: &config.Body{
				String: []string{"test3", "test4"},
			},
		},
	}
	rr := httptest.NewRecorder()
	responseRequest, err := yaml.Marshal(responseRequestCfg)
	if err != nil {
		t.Fatalf("Could not marshal response request: %v", err)
	}
	req := httptest.NewRequest(http.MethodPost, "/responses", bytes.NewReader(responseRequest))

	http.HandlerFunc(sut.Function).ServeHTTP(rr, req)

	// Then
	assert.Nil(t, err, "Expected no error")
	assert.Equal(t, http.StatusOK, rr.Code, "Expected HTTP status code 200 OK")
	assert.Equal(t, communicator.GetResponses()[0].BodyTemplate[0], "updated-test1", "Expected first response to have body 'test1'")

}

func TestUpdateResponseInCommunicator(t *testing.T) {
	// Given
	communicator := &MockCommunicator{}
	responsesCfg := []config.Response{
		{
			Alias: "test1",
			Body: &config.Body{
				String: []string{"test1"},
				File:   nil,
			},
		},
		{
			Alias: "test2",
		},
	}
	communicator.RespondWith(prototype.NewResponsePrototypeBuilder(responsesCfg).Build())

	sut := ResponsesEndpoint(communicator, &responsesCfg)

	responseRequestCfg := []config.Response{
		{
			Alias: "test1",
			Body: &config.Body{
				String: []string{"updated-test1", "test2"},
			},
		},
		{
			Alias: "test2",
		},
	}
	rr := httptest.NewRecorder()
	responseRequest, err := yaml.Marshal(responseRequestCfg[0])
	if err != nil {
		t.Fatalf("Could not marshal response request: %v", err)
	}
	req := httptest.NewRequest(http.MethodPatch, "/responses?response-alias=test1", bytes.NewReader(responseRequest))

	http.HandlerFunc(sut.Function).ServeHTTP(rr, req)

	// Then
	assert.Nil(t, err, "Expected no error")
	assert.Equal(t, http.StatusOK, rr.Code, "Expected HTTP status code 200 OK")
	assert.Equal(t, communicator.GetResponses()[0].BodyTemplate[0], "updated-test1", "Expected first response to have body 'test1'")

}
