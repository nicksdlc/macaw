package admin

import (
	"net/http"
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

func TestManagerStartedWithHealthEndpointReturnsOk(t *testing.T) {
	// Given
	manager := NewManager(1234, HealthEndpoint())

	// When
	manager.Start()

	resp, err := http.Get("http://localhost:1234/health")
	if err != nil {
		t.Fatalf("Could not make GET request: %v", err)
	}

	// Then
	assert.Equal(t, http.StatusOK, resp.StatusCode, "Expected HTTP status code 200 OK")

}
