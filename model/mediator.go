package model

// Mediator is a function that should be done with message
type Mediator func(message RequestMessage, response ResponseMessage) (ResponseMessage, error)
