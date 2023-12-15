package matchers

import (
	"testing"

	"github.com/nicksdlc/macaw/model"

	"github.com/stretchr/testify/assert"
)

func TestIfNoMatcherExistShouldReturnTrue(t *testing.T) {
	request := compexRequest()

	assert.True(t, Match(nil, request, Any))
}

func TestFieldMatcher_Match(t *testing.T) {
	request := compexRequest()

	matcher := FieldMatcher{
		Field: "host",
		Value: "localhost:8080",
	}

	assert.True(t, Match([]Matcher{&matcher}, request, Any))
}

func TestFieldMatcher_MatchShouldReturnFalse(t *testing.T) {
	request := compexRequest()

	matcher := FieldMatcher{
		Field: "host",
		Value: "localhost:8081",
	}

	assert.False(t, Match([]Matcher{&matcher}, request, Any))
}

func TestFieldExcludingMatcher_Match(t *testing.T) {
	request := compexRequest()

	matcher := FieldExcludingMatcher{
		Field: "host",
		Value: "localhost:8081",
	}

	assert.True(t, Match([]Matcher{&matcher}, request, Any))
}

func TestFieldExcludingMatcher_MatchShouldReturnFalse(t *testing.T) {
	request := compexRequest()

	matcher := FieldExcludingMatcher{
		Field: "host",
		Value: "localhost:8080",
	}

	assert.False(t, Match([]Matcher{&matcher}, request, Any))
}

func TestBodyContainsMatcher_Match(t *testing.T) {
	request := compexRequest()

	matcher := BodyContainsMatcher{
		Contains: "test",
	}

	assert.True(t, Match([]Matcher{&matcher}, request, Any))
}

func TestBodyContainsMatcher_MatchShouldReturnFalse(t *testing.T) {
	request := compexRequest()

	matcher := BodyContainsMatcher{
		Contains: "test1",
	}

	assert.False(t, Match([]Matcher{&matcher}, request, Any))
}

func TestExcludesMatcher_Match(t *testing.T) {
	request := compexRequest()

	matcher := ExcludesMatcher{
		Value: "test1",
	}

	assert.True(t, Match([]Matcher{&matcher}, request, Any))
}

func TestExcludesMatcher_MatchShouldReturnFalse(t *testing.T) {
	request := compexRequest()

	matcher := ExcludesMatcher{
		Value: "test",
	}

	assert.False(t, Match([]Matcher{&matcher}, request, Any))
}

func TestFieldContainsMatcher_Match(t *testing.T) {
	request := compexRequest()

	matcher := FieldContainsMatcher{
		Field: "host",
		Value: "localhost",
	}

	assert.True(t, Match([]Matcher{&matcher}, request, Any))
}

func TestFieldContainsMatcher_MatchShouldReturnFalse(t *testing.T) {
	request := compexRequest()

	matcher := FieldContainsMatcher{
		Field: "host",
		Value: "localhost1",
	}

	assert.False(t, Match([]Matcher{&matcher}, request, Any))
}

func TestAnyMatchWithMultipleMatchers_ShouldReturnTrueIfOneMatcherMatches(t *testing.T) {
	request := compexRequest()

	matcher1 := FieldContainsMatcher{
		Field: "host",
		Value: "localhost",
	}

	matcher2 := FieldContainsMatcher{
		Field: "host",
		Value: "localhost1",
	}

	assert.True(t, Match([]Matcher{&matcher1, &matcher2}, request, Any))
}

func TestAnyMatchWithMultipleMatchers_ShouldReturnFalseIfNoMatcherMatches(t *testing.T) {
	request := compexRequest()

	matcher1 := FieldContainsMatcher{
		Field: "host",
		Value: "localhost1",
	}

	matcher2 := BodyContainsMatcher{
		Contains: "test1",
	}

	assert.False(t, Match([]Matcher{&matcher1, &matcher2}, request, Any))
}

func TestAllMatchWithMultipleMatchers_ShouldReturnFalseIfOneMatcherDoesNotMatch(t *testing.T) {
	request := compexRequest()

	matcher1 := FieldContainsMatcher{
		Field: "host",
		Value: "localhost",
	}

	matcher2 := FieldContainsMatcher{
		Field: "host",
		Value: "localhost1",
	}

	assert.False(t, Match([]Matcher{&matcher1, &matcher2}, request, All))
}

func TestAllMatchWithMultipleMatchers_ShouldReturnTrueIfAllMatchersMatch(t *testing.T) {
	request := compexRequest()

	matcher1 := FieldContainsMatcher{
		Field: "host",
		Value: "localhost",
	}

	matcher2 := BodyContainsMatcher{
		Contains: "test",
	}

	assert.True(t, Match([]Matcher{&matcher1, &matcher2}, request, All))
}

func compexRequest() model.RequestMessage {
	return model.RequestMessage{
		Headers: map[string]string{
			"container":  "buildpack-deps:stretch-curl",
			"host":       "localhost:8080",
			"user-agent": "curl/7.64.0",
		},
		Body: []byte("{\"name\":\"test\",\"tag\":\"latest\"}"),
	}
}
