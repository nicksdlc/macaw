package model

// Mediator is a function that should be done with message
// response is a pointer to response message that should be modified
// by each successive mediator
type Mediator func(message RequestMessage, response *ResponseMessage) error
