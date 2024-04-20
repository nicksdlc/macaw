package prototype

import (
	"path"

	"github.com/nicksdlc/macaw/config"
	"github.com/nicksdlc/macaw/mediator"
	"github.com/nicksdlc/macaw/prototype/matchers"
	"github.com/nicksdlc/macaw/schema/oapi"
	"github.com/pb33f/libopenapi/datamodel/high/base"
	v3 "github.com/pb33f/libopenapi/datamodel/high/v3"
)

// BuildFromSpec builds the response prototype from OpenAPI spec
func (rtb *ResponsePrototypeBuilder) BuildFromOAPISpec(responseConfig config.Response) ([]MessagePrototype, error) {

	model, err := oapi.LoadModel(responseConfig.FromOpenAPI)

	if err != nil {
		return nil, err
	}

	result := []MessagePrototype{}
	path := model.Model.Paths.PathItems.First()

	for path != nil {
		//TODO: Supports only GET, remove when adding support for other methods
		if path.Value().Get == nil {
			path = path.Next()
			continue
		}
		sch := findSchema(path.Value())
		if sch == nil {
			path = path.Next()
			continue
		}
		from, err := combineTo(responseConfig.ResponseRequest.To, path.Key())

		if err != nil {
			return nil, err
		}

		prototype := MessagePrototype{
			// TODO: Change when support of other methods is added
			Type:      "GET",
			Alias:     responseConfig.Alias,
			Mediators: buildOAPIMediators(sch, responseConfig.Options),
			From:      from,
			Matcher:   buildMatcher(responseConfig.ResponseRequest.Matchers),
		}
		// prepend matchingMediator to the mediators list
		// so that it is the first mediator to be executed
		// and all others can be skipped if it doesn't match
		prototype.Mediators.Prepend(
			mediator.NewMatchingMediator(
				matchers.ParsePattern(responseConfig.ResponseRequest.Match),
				prototype.Matcher))

		result = append(result, prototype)
		path = path.Next()

	}

	return result, nil
}

func combineTo(base, addr string) (string, error) {
	return path.Join(base, addr), nil
}

func buildOAPIMediators(sch *base.Schema, opt *config.Options) mediator.MediatorChain {
	var chain mediator.MediatorChain

	chain.Append(mediator.NewOAPIGeneratingMediator(opt.Quantity, sch))
	chain.Append(mediator.NewDelayingMediator(opt.Delay))
	chain.Append(mediator.NewRandomDelayingMediator(opt.RandomDelay.Min, opt.RandomDelay.Max))

	return chain
}

// findSchema finds the schema under path. In first version supports only 200 and json
func findSchema(pi *v3.PathItem) *base.Schema {
	code := pi.Get.Responses.Codes.First()
	for code != nil {
		if code.Key() != "200" {
			code = code.Next()
			continue
		}
		content := code.Value().Content.First()

		for content != nil {
			if content.Key() == "application/json" {
				return content.Value().Schema.Schema()
			}
			content = content.Next()
		}
	}
	return nil
}
