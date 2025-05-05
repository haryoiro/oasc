package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/stretchr/testify/assert"
)

type testCase struct {
	name           string
	inputFiles     []string
	outputFormat   string
	expectedSchema string
	expectedTags   []string
	expectedPaths  []string
}

func assertOpenAPIEqual(t *testing.T, expected, actual *openapi3.T) {
	assert.Equal(t, expected.OpenAPI, actual.OpenAPI)
	assert.Equal(t, expected.Info.Title, actual.Info.Title)
	assert.Equal(t, expected.Info.Version, actual.Info.Version)
	assert.Equal(t, expected.Info.Description, actual.Info.Description)
}

func assertComponentsEqual(t *testing.T, expected, actual openapi3.Components) {
	assert.Equal(t, len(expected.Schemas), len(actual.Schemas))
	for schema := range expected.Schemas {
		assert.Contains(t, actual.Schemas, schema)
	}
}

func assertPathsEqual(t *testing.T, expected, actual openapi3.Paths) {
	assert.Equal(t, len(expected.Map()), len(actual.Map()))
	for path, expectedPath := range expected.Map() {
		actualPath, ok := actual.Map()[path]
		assert.True(t, ok, "Path %s not found in actual spec", path)
		assert.Equal(t, expectedPath.Get != nil, actualPath.Get != nil)
		assert.Equal(t, expectedPath.Post != nil, actualPath.Post != nil)
		assert.Equal(t, expectedPath.Put != nil, actualPath.Put != nil)
		assert.Equal(t, expectedPath.Delete != nil, actualPath.Delete != nil)
	}
}

func TestEndToEnd(t *testing.T) {
	t.Parallel()

	tests := []testCase{
		{
			name: "merge YAML files",
			inputFiles: []string{
				"./testdata/valid/test1.openapi.yaml",
				"./testdata/valid/test2.openapi.yaml",
			},
			outputFormat:   "yaml",
			expectedSchema: "Product",
			expectedTags:   []string{"products", "users"},
			expectedPaths:  []string{"/products", "/users", "/users/{id}"},
		},
		{
			name: "merge JSON files",
			inputFiles: []string{
				"./testdata/valid/test1.openapi.yaml",
				"./testdata/valid/test3.openapi.json",
			},
			outputFormat:   "json",
			expectedSchema: "Order",
			expectedTags:   []string{"orders", "products", "users"},
			expectedPaths:  []string{"/orders", "/products/{id}", "/users", "/users/{id}"},
		},
		{
			name: "merge empty file",
			inputFiles: []string{
				"./testdata/edge/empty.openapi.yaml",
				"./testdata/valid/test2.openapi.yaml",
			},
			outputFormat:   "yaml",
			expectedSchema: "Product",
			expectedTags:   []string{"products", "users"},
			expectedPaths:  []string{"/products", "/users"},
		},
	}

	for _, tt := range tests {
		tt := tt // capture range variable
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Create temp directory
			tempDir := t.TempDir()
			outputPath := filepath.Join(tempDir, "merged."+tt.outputFormat)

			// Read first file
			merged, err := ReadFile(tt.inputFiles[0])
			assert.NoError(t, err)
			assert.NotNil(t, merged)

			// Merge remaining files
			for i := 1; i < len(tt.inputFiles); i++ {
				spec, err := ReadFile(tt.inputFiles[i])
				assert.NoError(t, err)
				merged = MergeSpecs(merged, spec)
			}

			// Verify merge result
			assert.NotNil(t, merged)
			assert.NotNil(t, merged.Components)
			assert.NotNil(t, merged.Components.Schemas)
			assert.Contains(t, merged.Components.Schemas, tt.expectedSchema)

			// Verify expected tags
			if tt.expectedTags != nil {
				for _, tag := range tt.expectedTags {
					found := false
					for _, mergedTag := range merged.Tags {
						if mergedTag.Name == tag {
							found = true
							break
						}
					}
					assert.True(t, found, "Expected tag %s not found", tag)
				}
			}

			// Verify expected paths
			if tt.expectedPaths != nil {
				for _, path := range tt.expectedPaths {
					assert.Contains(t, merged.Paths.Map(), path)
				}
			}

			// Write output file
			isJSON := tt.outputFormat == "json"
			err = WriteFile(merged, outputPath, isJSON)
			assert.NoError(t, err)

			// Verify file exists
			_, err = os.Stat(outputPath)
			assert.NoError(t, err)

			// Verify written file
			writtenSpec, err := ReadFile(outputPath)
			assert.NoError(t, err)
			assert.NotNil(t, writtenSpec)
			assertOpenAPIEqual(t, merged, writtenSpec)
			assertComponentsEqual(t, *merged.Components, *writtenSpec.Components)
			assertPathsEqual(t, *merged.Paths, *writtenSpec.Paths)
		})
	}
}

func TestCompleteEndToEnd(t *testing.T) {
	// Test files
	inputFiles := []string{
		"./testdata/valid/test1.openapi.yaml",
		"./testdata/valid/test2.openapi.yaml",
		"./testdata/valid/test3.openapi.json",
	}
	expectedFile := "./testdata/valid/expect.openapi.yaml"

	// Create temp directory
	tempDir := t.TempDir()
	outputPath := filepath.Join(tempDir, "merged.yaml")

	// Read and merge input files
	var merged *openapi3.T
	for i, file := range inputFiles {
		spec, err := ReadFile(file)
		assert.NoError(t, err)
		if err != nil {
			t.Fatalf("Failed to read file %s: %v", file, err)
			return
		}
		assert.NotNil(t, spec)
		if spec == nil {
			t.Fatalf("Spec is nil for file %s", file)
			return
		}

		if i == 0 {
			merged = spec
		} else {
			merged = MergeSpecs(merged, spec)
		}
	}

	if merged == nil {
		t.Fatal("Merged spec is nil")
		return
	}

	// Write merged result
	err := WriteFile(merged, outputPath, false)
	assert.NoError(t, err)
	if err != nil {
		t.Fatalf("Failed to write merged file: %v", err)
		return
	}

	// Read expected result
	expected, err := ReadFile(expectedFile)
	assert.NoError(t, err)
	if err != nil {
		t.Fatalf("Failed to read expected file: %v", err)
		return
	}
	assert.NotNil(t, expected)
	if expected == nil {
		t.Fatal("Expected spec is nil")
		return
	}

	// Compare OpenAPI version
	assert.Equal(t, expected.OpenAPI, merged.OpenAPI)

	// Compare Info
	assert.Equal(t, expected.Info.Title, merged.Info.Title)
	assert.Equal(t, expected.Info.Version, merged.Info.Version)
	assert.Equal(t, expected.Info.Description, merged.Info.Description)

	// Compare Paths
	if expected.Paths != nil && merged.Paths != nil {
		assert.Equal(t, len(expected.Paths.Map()), len(merged.Paths.Map()))
		for path, expectedPath := range expected.Paths.Map() {
			mergedPath, ok := merged.Paths.Map()[path]
			assert.True(t, ok, "Path %s not found in merged spec", path)
			if ok {
				assert.Equal(t, expectedPath.Get != nil, mergedPath.Get != nil)
				assert.Equal(t, expectedPath.Post != nil, mergedPath.Post != nil)
				assert.Equal(t, expectedPath.Put != nil, mergedPath.Put != nil)
				assert.Equal(t, expectedPath.Delete != nil, mergedPath.Delete != nil)
			}
		}
	}

	// Compare Components
	if expected.Components != nil && merged.Components != nil {
		assert.Equal(t, len(expected.Components.Schemas), len(merged.Components.Schemas))
		for schema := range expected.Components.Schemas {
			assert.Contains(t, merged.Components.Schemas, schema)
		}
	}

	// Compare Tags
	if expected.Tags != nil && merged.Tags != nil {
		assert.Equal(t, len(expected.Tags), len(merged.Tags))
		for i, tag := range expected.Tags {
			assert.Equal(t, tag.Name, merged.Tags[i].Name)
			assert.Equal(t, tag.Description, merged.Tags[i].Description)
		}
	}
}
