package admin

import "net/http"

// Endpoint configuration
type Endpoint struct {
	Path     string
	Function func(w http.ResponseWriter, r *http.Request)
}

// HealthEndpoint creates a new health endpoint
func HealthEndpoint() Endpoint {
	return Endpoint{
		Path: "/health",
		Function: func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		},
	}
}
