package main

import (
	"os"
	"testing"
)

func TestValidateOpenAPIVersion(t *testing.T) {
	tests := []struct {
		name    string
		spec    map[string]any
		wantErr bool
	}{
		{
			name: "valid OpenAPI 3.0.0",
			spec: map[string]any{
				"openapi": "3.0.0",
			},
			wantErr: false,
		},
		{
			name: "valid OpenAPI 3.1.0",
			spec: map[string]any{
				"openapi": "3.1.0",
			},
			wantErr: false,
		},
		{
			name: "invalid OpenAPI 2.0",
			spec: map[string]any{
				"openapi": "2.0",
			},
			wantErr: true,
		},
		{
			name: "missing openapi field",
			spec: map[string]any{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateOpenAPIVersion(tt.spec)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateOpenAPIVersion() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestReadFile(t *testing.T) {
	// テスト用の一時ファイルを作成
	yamlContent := `openapi: 3.0.0
info:
  title: Test API
paths:
  /test:
    get:
      summary: Test endpoint`

	tmpFile, err := os.CreateTemp("", "test-*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.WriteString(yamlContent); err != nil {
		t.Fatal(err)
	}
	tmpFile.Close()

	spec, err := ReadFile(tmpFile.Name())
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}

	if spec["openapi"] != "3.0.0" {
		t.Error("Expected openapi version to be 3.0.0")
	}

	info := spec["info"].(map[string]any)
	if info["title"] != "Test API" {
		t.Error("Expected title to be 'Test API'")
	}
}
