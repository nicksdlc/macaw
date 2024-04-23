package mediator

import (
	"encoding/json"

	gen "github.com/nicksdlc/macaw/generator"
	"github.com/nicksdlc/macaw/model"
	"github.com/nicksdlc/macaw/schema/oapi"
	"github.com/pb33f/libopenapi/datamodel/high/base"
)

type oapiGeneratingMediator struct {
	generator gen.ObjGenerator
	quantity  int
}

func NewOAPIGeneratingMediator(quantity int, sch *base.Schema) Mediator {
	return &oapiGeneratingMediator{
		generator: oapi.NewGeneratorFromSchema(sch),
		quantity:  quantity,
	}
}

func (gm *oapiGeneratingMediator) Mediate(message model.RequestMessage, responses <-chan model.ResponseMessage) <-chan model.ResponseMessage {
	out := make(chan model.ResponseMessage)
	go func() {
		defer close(out)
		for response := range responses {
			for i := 0; i < gm.quantity; i++ {
				obj := gm.generator()
				json, _ := json.Marshal(obj)
				if response.Metadata == nil {
					response.Metadata = map[string]string{}
				}
				response.Metadata["Content-Type"] = "application/json"
				out <- model.ResponseMessage{
					Body:     string(json),
					Metadata: response.Metadata,
				}
			}
		}
	}()
	return out
}
