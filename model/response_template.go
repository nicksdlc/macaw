package model

type ResponseTemplate struct {
	Body string

	Headers map[string]string

	Mediators []Mediator
}
