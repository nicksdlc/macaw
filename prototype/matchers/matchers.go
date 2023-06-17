package matchers

import (
	"github.com/nicksdlc/macaw/model"
)

// Matcher is a interface for matching request to response
type Matcher interface {
	Match(request model.RequestMessage) bool
}

// Match matches request to response by all matchers
func Match(matchers []Matcher, message model.RequestMessage, pattern Pattern) bool {
	if len(matchers) == 0 {
		return true
	}

	if pattern == All {
		return MatchAll(matchers, message)
	}

	return MatchAny(matchers, message)
}

// MatchAll matches request to response by all matchers
func MatchAll(matchers []Matcher, message model.RequestMessage) bool {
	if len(matchers) == 0 {
		return true
	}

	for _, matcher := range matchers {
		if !matcher.Match(message) {
			return false
		}
	}
	return true
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
