package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

func validateOpenAPIVersion(spec map[string]any) error {
	openapi, ok := spec["openapi"].(string)
	if !ok {
		return fmt.Errorf("invalid OpenAPI specification: missing 'openapi' field")
	}

	if !strings.HasPrefix(openapi, "3.") {
		return fmt.Errorf("unsupported OpenAPI version: %s (only OpenAPI 3.x is supported)", openapi)
	}

	return nil
}

func ReadFile(filePath string) (map[string]any, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("error reading file: %v", err)
	}

	var result map[string]any
	ext := strings.ToLower(filepath.Ext(filePath))

	if ext == ".json" {
		if err := json.Unmarshal(data, &result); err != nil {
			return nil, fmt.Errorf("error parsing JSON: %v", err)
		}
	} else {
		if err := yaml.Unmarshal(data, &result); err != nil {
			return nil, fmt.Errorf("error parsing YAML: %v", err)
		}
	}

	if err := validateOpenAPIVersion(result); err != nil {
		return nil, err
	}

	return result, nil
}

func WriteFile(data map[string]any, outputPath string, isJSON bool) error {
	var output []byte
	var err error

	if isJSON {
		output, err = json.MarshalIndent(data, "", "  ")
	} else {
		output, err = yaml.Marshal(data)
	}
	if err != nil {
		return fmt.Errorf("error marshaling result: %v", err)
	}

	if err := os.WriteFile(outputPath, output, 0644); err != nil {
		return fmt.Errorf("error writing output: %v", err)
	}

	return nil
}
