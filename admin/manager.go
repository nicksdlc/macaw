package admin

import (
	"context"
	"fmt"
	"log"
	"net/http"
)

// Manager responsible for endpoints of configuration and health
type Manager struct {
	port      uint16
	endpoints []Endpoint
	server    *http.Server
}

// NewManager creates a new instance of Manager
func NewManager(port uint16, endpoints ...Endpoint) *Manager {
	return &Manager{
		port:      port,
		endpoints: endpoints,
	}
}

// AddEndpoint adds an endpoint to the manager
func (m *Manager) AddEndpoint(endpoint Endpoint) {
	m.endpoints = append(m.endpoints, endpoint)
}

// Start starts the manager
func (m *Manager) Start() {

	// Create a new ServeMux.
	mux := http.NewServeMux()

	// start the manager
	for _, endpoint := range m.endpoints {
		// start the endpoint
		mux.HandleFunc(endpoint.Path, endpoint.Function)
	}

	m.server = &http.Server{
		Addr:    fmt.Sprintf(":%d", m.port),
		Handler: mux,
	}

	// start the server
	if err := m.server.ListenAndServe(); err != nil {
		log.Fatalf("Error starting server: %s", err)
	}
}

// Stop stops the manager
func (m *Manager) Stop() {
	m.server.Shutdown(context.Background())
}
