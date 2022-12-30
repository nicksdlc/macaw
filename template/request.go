package template

import (
	"bytes"
	"html/template"
	"log"
	"macaw/template/types"
)

// Request is a request macaw will send to the consumer
type Request struct {
	BaseTemplate string

	tmpl         *template.Template
	placeholders map[int]types.Type
	index        int
}

// NewRequest ctor
func NewRequest(base string) Request {
	return Request{
		BaseTemplate: base,
		tmpl:         template.New(""),
		placeholders: make(map[int]types.Type),
		index:        0,
	}
}

// Create creates new Request
func (r *Request) Create() string {
	r.tmpl.Parse(r.BaseTemplate)

	var result bytes.Buffer
	err := r.tmpl.Execute(&result, r)
	if err != nil {
		log.Printf(err.Error())
	}

	r.index = 0
	return result.String()
}
