package matchers

import (
	"strings"

	"github.com/nicksdlc/macaw/model"
)

// BodyContainsMatcher is a matcher that matches request to response by body
type BodyContainsMatcher struct {
	Contains string
}

// Match matches request to response by body
func (m *BodyContainsMatcher) Match(request model.RequestMessage) bool {
	return strings.Contains(string(request.Body), m.Contains)
}
