package mediator

import (
	"log"

	"github.com/nicksdlc/macaw/model"
	"github.com/nicksdlc/macaw/template"
)

type GeneratingMediator struct {
	bodyTempalte string
	amount       int
}

func NewGeneratingMediator(amount int, bodyTempalte string) *GeneratingMediator {
	return &GeneratingMediator{
		bodyTempalte: bodyTempalte,
		amount:       amount,
	}
}

func (gm *GeneratingMediator) Mediate(message model.RequestMessage, responses <-chan model.ResponseMessage) <-chan model.ResponseMessage {
	out := make(chan model.ResponseMessage)
	go func() {
		defer close(out)
		for response := range responses {
			base := gm.bodyTempalte
			req := template.Serialize(message.Headers, message.Body)

			bodyGenerator := template.NewResponse(string(base), gm.amount, &req)

			log.Printf("Generating %d responses", gm.amount)
			for i := 0; i < gm.amount; i++ {
				out <- model.ResponseMessage{
					Response: bodyGenerator.Create(),
					Metadata: response.Metadata,
				}
			}
		}
	}()
	return out
}
