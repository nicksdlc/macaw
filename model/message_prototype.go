package model

import "strings"

// MessagePrototype is a template for message with mediators, bodyTemplate and headers
type MessagePrototype struct {
	Headers map[string]string

	BodyTemplate string

	Mediators []Mediator

	From string

	To string

	Matcher []Matcher
}

// Matcher is a interface for matching request to response
type Matcher interface {
	Match(request RequestMessage) bool
}

// FieldMatcher is a matcher that matches request to response by field
type FieldMatcher struct {
	Field string

	Value string
}

// Match matches request to response by field
func (m *FieldMatcher) Match(request RequestMessage) bool {
	return request.Headers[m.Field] == m.Value
}

// BodyContainsMatcher is a matcher that matches request to response by body
type BodyContainsMatcher struct {
	Contains string
}

// Match matches request to response by body
func (m *BodyContainsMatcher) Match(request RequestMessage) bool {
	return strings.Contains(string(request.Body), m.Contains)
}
