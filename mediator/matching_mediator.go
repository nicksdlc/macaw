package mediator

import (
	"github.com/nicksdlc/macaw/model"
	"github.com/nicksdlc/macaw/prototype/matchers"
)

type MatchingMediator struct {
	pattern  matchers.Pattern
	matchers []matchers.Matcher
}

func NewMatchingMediator(pattern matchers.Pattern, matchers []matchers.Matcher) *MatchingMediator {
	return &MatchingMediator{
		pattern:  pattern,
		matchers: matchers,
	}
}

func (mm *MatchingMediator) Mediate(message model.RequestMessage, responses <-chan model.ResponseMessage) <-chan model.ResponseMessage {
	out := make(chan model.ResponseMessage)
	go func() {
		defer close(out)
		for response := range responses {
			if matchers.Match(mm.matchers, message, mm.pattern) {
				out <- response
			}
		}
	}()
	return out
}
