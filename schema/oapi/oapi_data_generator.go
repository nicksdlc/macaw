package oapi

import (
	gen "github.com/nicksdlc/macaw/generator"
	"github.com/pb33f/libopenapi/datamodel/high/base"
)

func NewGeneratorFromSchema(model *base.Schema) gen.ObjGenerator {
	prop := model.Properties.First()
	var g gen.Generator = nil
	for prop != nil {
		g = gen.Compose(g, gen.GenerateUint(prop.Key()))
		prop = prop.Next()
	}

	return gen.NewObjGenerator(g)
}
