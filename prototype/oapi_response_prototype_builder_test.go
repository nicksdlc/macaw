package prototype

import (
	"testing"

	"github.com/nicksdlc/macaw/config"
	"github.com/stretchr/testify/assert"
)

type endpoint struct {
	address string
	method  string
}

func TestSingleGetEndpoint(t *testing.T) {
	responseConfig := []config.Response{
		{
			FromOpenAPI: "./testdata/oapi_single_endpoint_number_property.yaml",
			ResponseRequest: config.ResponseRequest{
				To: "/base-addr",
			},
			Options: &config.Options{},
		},
	}

	responsePrototypes := NewResponsePrototypeBuilder(responseConfig).Build()

	assert.Equal(t, 1, len(responsePrototypes))
	assert.Equal(t, "/base-addr/with-int", responsePrototypes[0].From)
	assert.Equal(t, "GET", responsePrototypes[0].Type)
}

func TestMultipleGetEndpoint(t *testing.T) {

	expected := []endpoint{
		{
			address: "/base-addr/endpoint1",
			method:  "GET",
		},
		{
			address: "/base-addr/endpoint2",
			method:  "GET",
		},
		{
			address: "/base-addr/endpoint3",
			method:  "GET",
		},
		{
			address: "/base-addr/endpoint4",
			method:  "GET",
		},
	}

	responseConfig := []config.Response{
		{
			FromOpenAPI: "./testdata/oapi_multiple_endpoint.yaml",
			ResponseRequest: config.ResponseRequest{
				To: "/base-addr",
			},
			Options: &config.Options{},
		},
	}

	responsePrototypes := NewResponsePrototypeBuilder(responseConfig).Build()

	assert.Equal(t, 4, len(responsePrototypes))

	for i, response := range expected {
		assert.Equal(t, response.address, responsePrototypes[i].From)
		assert.Equal(t, response.method, responsePrototypes[i].Type)
	}
}
