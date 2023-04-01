package model

// ResponseMessage is used to build response before replying
type ResponseMessage struct {
	Metadata map[string]string

	Responses []string
}
