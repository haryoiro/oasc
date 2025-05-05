package main

import (
	"context"
	"os"

	"github.com/getkin/kin-openapi/openapi3"
	"gopkg.in/yaml.v3"
)

func ReadFile(path string) (*openapi3.T, error) {
	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	loader.Context = context.Background()

	doc, err := loader.LoadFromFile(path)
	if err != nil {
		return nil, err
	}

	// Skip validation for now
	return doc, nil
}

func WriteFile(doc *openapi3.T, path string, isJSON bool) error {
	var data []byte
	var err error

	if isJSON {
		data, err = doc.MarshalJSON()
	} else {
		yamlData, err := doc.MarshalYAML()
		if err != nil {
			return err
		}
		data, err = yaml.Marshal(yamlData)
	}
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}
