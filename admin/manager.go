package admin

import (
	"fmt"
	"net/http"
)

// Manager responsible for endpoints of configuration and health
type Manager struct {
	port      uint16
	endpoints []Endpoint
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

	// Start the server
	go http.ListenAndServe(fmt.Sprintf(":%d", m.port), mux)
}
