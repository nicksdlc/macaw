package template

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"strings"

	"github.com/nicksdlc/macaw/template/types"
	"github.com/nicksdlc/macaw/template/types/number"
)

// Response is a response application will send to the consumer
//
//	Response is generated based on the incomint request and on the go template
//	All the names starting from "FromRequest" are the names of the functions in the template
//	If go template engine meets such a function in the text of the BaseTemplate it will
//	use the code of the function to fill it with relevant value
type Response struct {
	BaseTemplate string
	Quantity     int

	tmpl         *template.Template
	placeholders map[int]types.Type
	index        int
	request      *IncomingRequest
}

// NewResponse ctor
func NewResponse(base string, quantity int, req *IncomingRequest) Response {
	return Response{
		BaseTemplate: base,
		Quantity:     quantity,
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
		log.Print(err.Error())
	}

	r.index = 0
	return result.String()
}

// FromRequest gets the field from request
func (r *Response) FromRequest(field string) string {
	path := strings.Split(field, "/")

	intermidiate := r.request.Fields
	for _, p := range path {
		if val, ok := intermidiate[p]; ok {
			if value, ok := val.(map[string]interface{}); ok {
				intermidiate = value
				continue
			}
			return fmt.Sprint(val)
		}

	}

	return "UNDEFINED"
}

// FromRequestHeaders gets the field values from request headers if those present
func (r *Response) FromRequestHeaders(field string) string {
	if val, ok := r.request.Headers[field]; ok {
		return val
	}

	return "UNDEFINED"
}

// Number Represents number in template
func (r *Response) Number(parameters ...string) string {
	currentIndex := r.updatePlaceholders(parameters, number.Number)
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
	listResponse := NewResponse(input, r.Quantity, r.request)
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
