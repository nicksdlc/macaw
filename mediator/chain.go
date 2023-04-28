package mediator

import (
	"log"

	"github.com/nicksdlc/macaw/model"
)

type MediatorChain struct {
	linkedMediator *chainedMediator
}

type chainedMediator struct {
	mediator Mediator
	next     *chainedMediator
}

func (mc *MediatorChain) Append(mediator Mediator) {
	if mc.linkedMediator == nil {
		mc.linkedMediator = &chainedMediator{
			mediator: mediator,
		}
		return
	}

	mc.linkedMediator.next = &chainedMediator{
		mediator: mediator,
	}

}

func (mc *MediatorChain) Run(request model.RequestMessage, base model.ResponseMessage) <-chan model.ResponseMessage {
	m := mc.linkedMediator
	channel := m.mediator.Mediate(request, mc.generateChan(base))
	log.Printf("Mediator %T", m.mediator)

	for m.next != nil {
		log.Printf("Mediator %T", m.next.mediator)
		m = m.next
		channel = m.mediator.Mediate(request, channel)
	}

	return channel
}

func (mc *MediatorChain) generateChan(base model.ResponseMessage) <-chan model.ResponseMessage {
	in := make(chan model.ResponseMessage)
	go func() {
		in <- base
		close(in)
	}()

	return in
}
