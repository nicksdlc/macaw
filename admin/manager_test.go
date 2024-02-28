package admin

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestManagerCreatedWithHealthEndpoint(t *testing.T) {
	// Given
	healthEndpoint := Endpoint{
		Path: "/health",
	}

	// When
	manager := NewManager(1234, healthEndpoint)

	// Then
	assert.NotNil(t, manager)
	assert.Equal(t, 1, len(manager.endpoints))
	assert.Equal(t, "/health", manager.endpoints[0].Path)
}

func TestHealthEndpointReturnsOk(t *testing.T) {
	// Given
	healthEndpoint := HealthEndpoint()

	// When
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rr := httptest.NewRecorder()

	http.HandlerFunc(healthEndpoint.Function).ServeHTTP(rr, req)

	// Then
	assert.Equal(t, http.StatusOK, rr.Code, "Expected HTTP status code 200 OK")

}
