package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/stretchr/testify/assert"
)

func TestReadFile(t *testing.T) {
	t.Run("read YAML file", func(t *testing.T) {
		spec, err := ReadFile("./testdata/valid/test1.openapi.yaml")
		assert.NoError(t, err)
		if err != nil {
			t.Fatalf("Failed to read YAML file: %v", err)
			return
		}
		assert.NotNil(t, spec)
		if spec == nil {
			t.Fatal("Spec is nil")
			return
		}
		assert.Equal(t, "3.0.0", spec.OpenAPI)
	})

	t.Run("read JSON file", func(t *testing.T) {
		spec, err := ReadFile("./testdata/valid/test3.openapi.json")
		assert.NoError(t, err)
		if err != nil {
			t.Fatalf("Failed to read JSON file: %v", err)
			return
		}
		assert.NotNil(t, spec)
		if spec == nil {
			t.Fatal("Spec is nil")
			return
		}
		assert.Equal(t, "3.0.0", spec.OpenAPI)
	})

	t.Run("read empty file", func(t *testing.T) {
		spec, err := ReadFile("./testdata/edge/empty.openapi.yaml")
		assert.NoError(t, err)
		if err != nil {
			t.Fatalf("Failed to read empty file: %v", err)
			return
		}
		assert.NotNil(t, spec)
		if spec == nil {
			t.Fatal("Spec is nil")
			return
		}
		assert.Equal(t, "3.0.0", spec.OpenAPI)
	})

	t.Run("read invalid file", func(t *testing.T) {
		spec, err := ReadFile("./testdata/invalid/invalid.openapi.yaml")
		assert.NoError(t, err)
		if err != nil {
			t.Fatalf("Failed to read invalid file: %v", err)
			return
		}
		assert.NotNil(t, spec)
		if spec == nil {
			t.Fatal("Spec is nil")
			return
		}
		assert.Equal(t, "3.0.0", spec.OpenAPI)
	})

	t.Run("read non-existent file", func(t *testing.T) {
		_, err := ReadFile("./testdata/non-existent.yaml")
		assert.Error(t, err)
	})
}

func TestWriteFile(t *testing.T) {
	tempDir := t.TempDir()
	tests := []struct {
		name      string
		spec      *openapi3.T
		filePath  string
		isJSON    bool
		wantErr   bool
	}{
		{
			name: "write YAML file",
			spec: &openapi3.T{
				OpenAPI: "3.0.0",
				Info: &openapi3.Info{
					Title:   "Test API",
					Version: "1.0.0",
				},
			},
			filePath: filepath.Join(tempDir, "test.yaml"),
			isJSON:   false,
			wantErr:  false,
		},
		{
			name: "write JSON file",
			spec: &openapi3.T{
				OpenAPI: "3.0.0",
				Info: &openapi3.Info{
					Title:   "Test API",
					Version: "1.0.0",
				},
			},
			filePath: filepath.Join(tempDir, "test.json"),
			isJSON:   true,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := WriteFile(tt.spec, tt.filePath, tt.isJSON)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				_, err := os.Stat(tt.filePath)
				assert.NoError(t, err)
			}
		})
	}
}
