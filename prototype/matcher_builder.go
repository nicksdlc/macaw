package prototype

import (
	"github.com/nicksdlc/macaw/config"
	"github.com/nicksdlc/macaw/prototype/matchers"
)

var matcherTypes = make(map[string]matcherBuilder)

func init() {
	matcherTypes["contains"] = func(matcherConfig config.Matcher) matchers.Matcher {
		return &matchers.BodyContainsMatcher{
			Contains: matcherConfig.Value,
		}
	}

	matcherTypes["field"] = func(matcherConfig config.Matcher) matchers.Matcher {
		return &matchers.FieldMatcher{
			Field: matcherConfig.Name,
			Value: matcherConfig.Value,
		}
	}

	matcherTypes["excludesfield"] = func(matcherConfig config.Matcher) matchers.Matcher {
		return &matchers.FieldExcludingMatcher{
			Field: matcherConfig.Name,
			Value: matcherConfig.Value,
		}
	}

	matcherTypes["excludes"] = func(matcherConfig config.Matcher) matchers.Matcher {
		return &matchers.ExcludesMatcher{
			Value: matcherConfig.Value,
		}
	}

	matcherTypes["fieldcontains"] = func(matcherConfig config.Matcher) matchers.Matcher {
		return &matchers.FieldContainsMatcher{
			Field: matcherConfig.Name,
			Value: matcherConfig.Value,
		}
	}
}

type matcherBuilder func(config.Matcher) matchers.Matcher
