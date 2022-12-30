package sender

// Sender sends message to the external interface
type Sender interface {
	Send() error
}
