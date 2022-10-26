package receiver

// Listener listens to the incoming message and updates the responder
type Listener interface {
	Listen()

	Notify(message []byte)
}
