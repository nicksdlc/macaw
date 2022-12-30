package template

import (
	"encoding/json"
	"log"
)

// IncomingRequest is a map of incoming fields
type IncomingRequest struct {
	Fields map[string]interface{}
}

// Serialize creates input request from byte array
func Serialize(input []byte) IncomingRequest {
	var r IncomingRequest
	err := json.Unmarshal(input, &r.Fields)
	if err != nil {
		log.Println(err.Error())
	}

	return r
}
