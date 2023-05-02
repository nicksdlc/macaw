package model

var StatusHeader = "status"

// ResponseMetadata is the metadata of the response
type ResponseMetadata map[string]string

// Status returns the status of the response
func (rm ResponseMetadata) Status() string {
	return rm[StatusHeader]
}

// NewResponseMetadata creates a new response metadata
func NewResponseMetadata(status string) ResponseMetadata {
	return ResponseMetadata{
		StatusHeader: status,
	}
}

// ResponseMessage is used to build response before replying
type ResponseMessage struct {
	Metadata ResponseMetadata

	Body string
}
