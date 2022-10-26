package template

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"macaw/template/types"
)

// Response is a response application will send to the consumer
type Response struct {
	BaseTemplate string
	Amount       int

	tmpl         *template.Template
	placeholders map[int]types.Type
	index        int
	request      *Request
}

// NewResponse ctor
func NewResponse(base string, amount int, req *Request) Response {
	return Response{
		BaseTemplate: base,
		Amount:       amount,
		tmpl:         template.New(""),
		placeholders: make(map[int]types.Type),
		index:        0,
		request:      req,
	}
}

// Create creates new response
func (r *Response) Create() string {
	r.tmpl.Parse(r.BaseTemplate)

	var result bytes.Buffer
	err := r.tmpl.Execute(&result, r)
	if err != nil {
		log.Printf(err.Error())
	}

	r.index = 0
	return result.String()
}

// FromRequest gets the field from request
func (r *Response) FromRequest(field string) string {
	return fmt.Sprint(r.request.Fields[field])
}

// Number Represents number in template
func (r *Response) Number(parameters ...string) string {
	currentIndex := r.updatePlaceholders(parameters, types.Number)
	return r.placeholders[currentIndex].Value()
}

// String Represents string in template
func (r *Response) String(parameters ...string) string {
	currentIndex := r.updatePlaceholders(parameters, types.String)
	return r.placeholders[currentIndex].Value()
}

// Date represents date in template
func (r *Response) Date(parameters ...string) string {
	currentIndex := r.updatePlaceholders(parameters, types.Date)
	return r.placeholders[currentIndex].Value()
}

// List represents list in template
func (r *Response) List(input string, times int) string {
	listResponse := NewResponse(input, r.Amount, r.request)
	var result string
	for i := 0; i < times; i++ {
		result += listResponse.Create() + "\n"
	}
	return result
}

func (r *Response) updatePlaceholders(parameters []string, t func(params []string) types.Type) int {
	current := r.index
	if _, ok := r.placeholders[current]; !ok {
		r.placeholders[current] = t(parameters)
	}
	r.index++
	return current
}
