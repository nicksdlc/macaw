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

func TestManagerStartedWithHealthEndpointReturnsOk(t *testing.T) {
	// Given
	manager := NewManager(1234, HealthEndpoint())

	// When
	manager.Start()

	req, err := http.NewRequest("GET", "/health", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	handler.ServeHTTP(rr, req)

	// Then
	assert.Equal(t, http.StatusOK, rr.Code, "handler returned wrong status code")
}
