package oapi

import (
	gen "github.com/nicksdlc/macaw/generator"
	"github.com/pb33f/libopenapi/datamodel/high/base"
	"github.com/pb33f/libopenapi/orderedmap"
)

// Each type in schema is represented in the of pair, where key is the name and Schema is the schema
// Different types may have different formats(e.g. number can be float) and additional information (e.g. array of strings)
type typeHandler func(orderedmap.Pair[string, *base.SchemaProxy]) gen.Generator

var typeHandlers = map[string]typeHandler{}

func init() {
	//TODO: bug, number can be float or double and now handled as integer
	typeHandlers["integer"] = handleInteger
	typeHandlers["number"] = handleNumber
	typeHandlers["string"] = handleString
}

func NewGeneratorFromSchema(model *base.Schema) gen.ObjGenerator {
	prop := model.Properties.First()
	var g gen.Generator = nil
	for prop != nil {
		// It is possible that Type contains multiple types but until oneOf, anyOf, allOf are not supported
		// it is not relevant
		sch := prop.Value().Schema()
		h, exists := typeHandlers[sch.Type[0]]
		if exists && h != nil {
			g = gen.Compose(g, h(prop))
		}
		prop = prop.Next()
	}

	return gen.NewObjGenerator(g)
}

func handleInteger(p orderedmap.Pair[string, *base.SchemaProxy]) gen.Generator {
	return gen.GenerateInt(p.Key())
}

func handleString(p orderedmap.Pair[string, *base.SchemaProxy]) gen.Generator {
	return gen.GenerateString(p.Key())
}

func handleNumber(p orderedmap.Pair[string, *base.SchemaProxy]) gen.Generator {
	return gen.GenerateFloat(p.Key())
}
