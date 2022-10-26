package template

import (
	"encoding/json"
	"log"
)

type Request struct {
	Fields map[string]interface{}
}

func Serialize(input []byte) Request {
	var r Request
	err := json.Unmarshal(input, &r.Fields)
	if err != nil {
		log.Println(err.Error())
	}

	return r
}
