package matchers

import (
	"strings"

	"github.com/nicksdlc/macaw/model"
)

// FieldMatcher is a matcher that matches request to response by field
type FieldContainsMatcher struct {
	Field string

	Value string
}

// Match matches request to response by field
func (fcm *FieldContainsMatcher) Match(request model.RequestMessage) bool {
	return strings.Contains(request.Headers[fcm.Field], fcm.Value)
}
