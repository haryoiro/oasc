package main

import (
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/stretchr/testify/assert"
)

func TestMergeSpecs(t *testing.T) {
	tests := []struct {
		name     string
		spec1    *openapi3.T
		spec2    *openapi3.T
		expected *openapi3.T
	}{
		{
			name: "basic merge",
			spec1: &openapi3.T{
				OpenAPI: "3.0.0",
				Info: &openapi3.Info{
					Title:   "API 1",
					Version: "1.0.0",
				},
				Paths: &openapi3.Paths{},
			},
			spec2: &openapi3.T{
				OpenAPI: "3.0.0",
				Info: &openapi3.Info{
					Title:   "API 2",
					Version: "1.0.0",
				},
				Paths: &openapi3.Paths{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Add test paths
			tt.spec1.Paths.Set("/test1", &openapi3.PathItem{
				Get: &openapi3.Operation{
					OperationID: "getTest1",
				},
			})
			tt.spec2.Paths.Set("/test2", &openapi3.PathItem{
				Get: &openapi3.Operation{
					OperationID: "getTest2",
				},
			})

			merged := MergeSpecs(tt.spec1, tt.spec2)
			assert.NotNil(t, merged)
			assert.Equal(t, tt.spec1.OpenAPI, merged.OpenAPI)
			assert.Equal(t, tt.spec1.Info.Title, merged.Info.Title)
			assert.Equal(t, tt.spec1.Info.Version, merged.Info.Version)
			assert.Equal(t, 2, len(merged.Paths.Map()))
		})
	}
}

func TestMergePathItem(t *testing.T) {
	tests := []struct {
		name     string
		dst      *openapi3.PathItem
		src      *openapi3.PathItem
		expected *openapi3.PathItem
	}{
		{
			name: "merge path items",
			dst:  &openapi3.PathItem{},
			src: &openapi3.PathItem{
				Get: &openapi3.Operation{
					OperationID: "getTest",
				},
				Post: &openapi3.Operation{
					OperationID: "postTest",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mergePathItem(tt.dst, tt.src)
			assert.NotNil(t, tt.dst.Get)
			assert.NotNil(t, tt.dst.Post)
			assert.Equal(t, "getTest", tt.dst.Get.OperationID)
			assert.Equal(t, "postTest", tt.dst.Post.OperationID)
		})
	}
}
