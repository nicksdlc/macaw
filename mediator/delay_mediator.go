package mediator

import (
	"log"
	"math/rand"
	"strconv"
	"time"

	"github.com/nicksdlc/macaw/model"
)

type DelayingMediator struct {
	Delay time.Duration
}

func NewDelayingMediator(delay string) *DelayingMediator {
	duration, err := parseDurationWithDefaultUnit(delay, "ms")
	if err != nil {
		log.Fatal(err)
	}
	return &DelayingMediator{
		Delay: duration,
	}
}

func (dm *DelayingMediator) Mediate(message model.RequestMessage, responses <-chan model.ResponseMessage) <-chan model.ResponseMessage {
	out := make(chan model.ResponseMessage)
	go func() {
		defer close(out)
		for response := range responses {
			log.Println("Delaying response for ", dm.Delay)
			time.Sleep(dm.Delay)

			out <- response
		}
	}()
	return out
}

type RandomDelayingMediator struct {
	MinDelay time.Duration
	MaxDelay time.Duration
}

func NewRandomDelayingMediator(minDelay, maxDelay string) *RandomDelayingMediator {
	minDuration, err := parseDurationWithDefaultUnit(minDelay, "ms")
	if err != nil {
		log.Fatal(err)
	}
	maxDuration, err := parseDurationWithDefaultUnit(maxDelay, "ms")
	if err != nil {
		log.Fatal(err)
	}
	return &RandomDelayingMediator{
		MinDelay: minDuration,
		MaxDelay: maxDuration,
	}
}

func (rdm *RandomDelayingMediator) Mediate(message model.RequestMessage, responses <-chan model.ResponseMessage) <-chan model.ResponseMessage {
	out := make(chan model.ResponseMessage)
	go func() {
		defer close(out)
		for response := range responses {
			if rdm.MinDelay < 0 || rdm.MaxDelay < 0 {
				log.Println("Error: MinDelay and MaxDelay must be non-negative")
				return
			}
			if rdm.MinDelay > rdm.MaxDelay {
				log.Println("Error: MinDelay cannot be greater than MaxDelay")
				return
			}
			if rdm.MaxDelay > 0 && rdm.MinDelay < rdm.MaxDelay {
				time.Sleep(rdm.MinDelay + time.Duration(rand.Intn(int(rdm.MaxDelay-rdm.MinDelay))))
			} else {
				time.Sleep(rdm.MinDelay)
			}
			out <- response
		}
	}()
	return out
}

func parseDurationWithDefaultUnit(s, defaultUnit string) (time.Duration, error) {
	if s == "" {
		return 0, nil
	}

	// Check if s is a number without a unit
	if _, err := strconv.Atoi(s); err == nil {
		s += defaultUnit
	}
	return time.ParseDuration(s)
}
