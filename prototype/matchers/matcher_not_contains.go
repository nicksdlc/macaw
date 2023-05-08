package matchers

import (
	"strings"

	"github.com/nicksdlc/macaw/model"
)

// ExcludesMatcher is a matcher that checkes if body does not contain a string
type ExcludesMatcher struct {
	Value string
}

// Match matches request to response by body
func (m *ExcludesMatcher) Match(request model.RequestMessage) bool {
	return !strings.Contains(string(request.Body), m.Value)
}
