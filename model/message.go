package model

// RequestMessage message that will be used and processed in the system
type RequestMessage struct {
	Headers map[string]string

	Body []byte
}
