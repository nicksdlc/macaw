package admin

import (
	"io"
	"log"
	"net/http"

	"github.com/nicksdlc/macaw/communicators"
	"github.com/nicksdlc/macaw/config"
	"github.com/nicksdlc/macaw/prototype"
	"gopkg.in/yaml.v3"
)

// UpdateEndpoint creates a new update endpoint
func UpdateEndpoint(communicator communicators.Communicator, responsesCfg []config.Response) Endpoint {
	return Endpoint{
		Path: "/update",
		Function: func(w http.ResponseWriter, r *http.Request) {
			handleUpdateRequest(w, r, communicator, responsesCfg)
		},
	}
}

func handleUpdateRequest(w http.ResponseWriter, r *http.Request, communicator communicators.Communicator, responsesCfg []config.Response) {
	if !isPatchRequest(r) {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	alias := getAlias(r)
	if alias == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	responseCfg := findResponseConfig(alias, responsesCfg)
	if responseCfg == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	body, err := readRequestBody(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("Error reading request body: %s", err)
		return
	}

	newResponse, err := unmarshalResponse(body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("Error unmarshalling response: %s", err.Error())
		return
	}

	updateResponseInCommunicator(alias, &newResponse, communicator, responseCfg)

	w.WriteHeader(http.StatusOK)
}

// Check if the request is a Patch
func isPatchRequest(r *http.Request) bool {
	return r.Method == http.MethodPatch
}

// Get the alias from the URL
func getAlias(r *http.Request) string {
	return r.URL.Query().Get("response-alias")
}

// Find the response config with the given alias
func findResponseConfig(alias string, responsesCfg []config.Response) *config.Response {
	for _, cfg := range responsesCfg {
		if cfg.Alias == alias {
			return &cfg
		}
	}
	return nil
}

// Read the body of the request
func readRequestBody(r *http.Request) ([]byte, error) {
	return io.ReadAll(r.Body)
}

// Unmarshal the body into a new response
func unmarshalResponse(body []byte) (config.Response, error) {
	newResponse := config.Response{}
	err := yaml.Unmarshal(body, &newResponse)
	return newResponse, err
}

// Update the response with the new response
func updateResponseInCommunicator(alias string, newResponse *config.Response, communicator communicators.Communicator, responsesCfg *config.Response) {
	for _, resp := range communicator.GetResponses() {
		if resp.Alias == alias {
			updateResponse(responsesCfg, newResponse)
			builder := prototype.NewResponsePrototypeBuilder(nil)
			communicator.UpdateResponse(builder.BuildResponse(*responsesCfg))
			break
		}
	}
}

func updateResponse(current *config.Response, newResp *config.Response) {
	// Update the response
	if newResp.Body != nil {
		current.Body = newResp.Body
	}

	if newResp.Options != nil {
		current.Options = newResp.Options
		if newResp.Options.Quantity == 0 {
			current.Options.Quantity = 1
		}
	}
}
