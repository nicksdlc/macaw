package mediator

import (
	"github.com/nicksdlc/macaw/model"
	"github.com/nicksdlc/macaw/template"
)

type GeneratingMediator struct {
	bodyTempalte string
	quantity     int
}

func NewGeneratingMediator(quantity int, bodyTempalte string) *GeneratingMediator {
	return &GeneratingMediator{
		bodyTempalte: bodyTempalte,
		quantity:     quantity,
	}
}

func (gm *GeneratingMediator) Mediate(message model.RequestMessage, responses <-chan model.ResponseMessage) <-chan model.ResponseMessage {
	out := make(chan model.ResponseMessage)
	go func() {
		defer close(out)
		for response := range responses {
			base := gm.bodyTempalte
			req := template.Serialize(message.Headers, message.Body)

			bodyGenerator := template.NewResponse(base, gm.quantity, &req)

			for i := 0; i < gm.quantity; i++ {
				out <- model.ResponseMessage{
					Body:     bodyGenerator.Create(),
					Metadata: response.Metadata,
				}
			}
		}
	}()
	return out
}
