package matchers

import (
	"github.com/nicksdlc/macaw/model"
)

// Matcher is a interface for matching request to response
type Matcher interface {
	Match(request model.RequestMessage) bool
}

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

// Should be moved to a mediator maybe
func MatchAny(matchers []Matcher, message model.RequestMessage) bool {
	if len(matchers) == 0 {
		return true
	}

	for _, matcher := range matchers {
		if matcher.Match(message) {
			return true
		}
	}
	return false
}
