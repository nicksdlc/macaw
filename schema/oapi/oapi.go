package oapi

import (
	"os"

	"github.com/hashicorp/go-multierror"
	"github.com/pb33f/libopenapi"

	v3 "github.com/pb33f/libopenapi/datamodel/high/v3"
)

func LoadModel(path string) (*libopenapi.DocumentModel[v3.Document], error) {

	schFile, _ := os.ReadFile(path)

	document, err := libopenapi.NewDocument(schFile)

	if err != nil {
		return nil, err
	}

	docModel, errors := document.BuildV3Model()
	if errors != nil {
		return nil, multierror.Append(nil, errors...)
	}

	return docModel, nil

}
