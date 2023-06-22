package matchers

import "github.com/nicksdlc/macaw/model"

// FieldMatcher is a matcher that matches request to response by field
type FieldMatcher struct {
	Field string

	Value string
}

// Match matches request to response by field
func (m *FieldMatcher) Match(request model.RequestMessage) bool {
	return request.Headers[m.Field] == m.Value
}

// FieldExcludingMatcher is a matcher that matches request to response by field
type FieldExcludingMatcher struct {
	Field string

	Value string
}

// Match matches request to response by field
func (m *FieldExcludingMatcher) Match(request model.RequestMessage) bool {
	return request.Headers[m.Field] != m.Value
}
