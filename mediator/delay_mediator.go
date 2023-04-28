package mediator

import (
	"time"

	"github.com/nicksdlc/macaw/model"
)

type DelayingMediator struct {
	Delay int
}

func NewDelayingMediator(delay int) *DelayingMediator {
	return &DelayingMediator{
		Delay: delay,
	}
}

func (dm *DelayingMediator) Mediate(message model.RequestMessage, responses <-chan model.ResponseMessage) <-chan model.ResponseMessage {
	out := make(chan model.ResponseMessage)
	go func() {
		defer close(out)
		for response := range responses {
			time.Sleep(time.Duration(dm.Delay) * time.Millisecond)

			out <- response
		}
	}()
	return out
}
