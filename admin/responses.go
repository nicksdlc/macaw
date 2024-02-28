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

func ResponsesEndpoint(communicator communicators.Communicator, responsesCfg *[]config.Response) Endpoint {
	return Endpoint{
		Path: "/responses",
		Function: func(w http.ResponseWriter, r *http.Request) {
			if isGetRequest(r) {
				handleGetResponses(w, r, communicator, responsesCfg)
				return
			}

			if isPostRequest(r) {
				handleUpdateAllRequest(w, r, communicator, responsesCfg)
				return
			}

			if isPatchRequest(r) {
				handleUpdateRequest(w, r, communicator, responsesCfg)
				return
			}
		},
	}
}

func isPostRequest(r *http.Request) bool {
	return r.Method == http.MethodPost
}

func isGetRequest(r *http.Request) bool {
	return r.Method == http.MethodGet
}

// Check if the request is a Patch
func isPatchRequest(r *http.Request) bool {
	return r.Method == http.MethodPatch
}

func handleGetResponses(w http.ResponseWriter, r *http.Request, communicator communicators.Communicator, responsesCfg *[]config.Response) {
	w.Header().Set("Content-Type", "application/yaml")
	w.WriteHeader(http.StatusOK)

	body, err := yaml.Marshal(responsesCfg)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = w.Write(body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func handleUpdateAllRequest(w http.ResponseWriter, r *http.Request, communicator communicators.Communicator, responsesCfg *[]config.Response) {
	body, err := readRequestBody(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("Error reading request body: %s", err)
		return
	}

	newResponses, err := unmarshalResponses(body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("Error unmarshalling response: %s", err.Error())
		return
	}

	// update responses configuration
	*responsesCfg = newResponses

	updateAllResponsesInCommunicator(newResponses, communicator, responsesCfg)

	w.WriteHeader(http.StatusOK)

}

func updateAllResponsesInCommunicator(newResponses []config.Response, communicator communicators.Communicator, responsesCfg *[]config.Response) {
	communicator.RespondWith(prototype.NewResponsePrototypeBuilder(newResponses).Build())
}

func unmarshalResponses(body []byte) ([]config.Response, error) {
	newResponses := []config.Response{}
	err := yaml.Unmarshal(body, &newResponses)
	return newResponses, err
}

func handleUpdateRequest(w http.ResponseWriter, r *http.Request, communicator communicators.Communicator, responsesCfg *[]config.Response) {
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

// Get the alias from the URL
func getAlias(r *http.Request) string {
	return r.URL.Query().Get("response-alias")
}

// Find the response config with the given alias
func findResponseConfig(alias string, responsesCfg *[]config.Response) *config.Response {
	for _, cfg := range *responsesCfg {
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
