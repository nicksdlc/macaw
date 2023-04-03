package template

import (
	"encoding/json"
	"log"
)

// IncomingRequest is a map of incoming fields
type IncomingRequest struct {
	Fields map[string]interface{}

	Headers map[string]string
}

// Serialize creates input request from byte array
func Serialize(headers map[string]string, input []byte) IncomingRequest {
	var r IncomingRequest
	err := json.Unmarshal(input, &r.Fields)
	if err != nil {
		log.Println(err.Error())
	}
	r.Headers = headers

	return r
}
