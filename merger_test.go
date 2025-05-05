package main

import (
	"testing"
)

func TestMergeOpenAPIPaths(t *testing.T) {
	paths1 := map[string]any{
		"/pets": map[string]any{
			"get": map[string]any{
				"summary": "List pets",
				"parameters": []any{
					map[string]any{"name": "limit", "in": "query"},
				},
			},
		},
	}

	paths2 := map[string]any{
		"/pets": map[string]any{
			"post": map[string]any{
				"summary": "Create pet",
			},
		},
		"/pets/{id}": map[string]any{
			"get": map[string]any{
				"summary": "Get pet by ID",
			},
		},
	}

	merged := mergeOpenAPIPaths(paths1, paths2)

	if len(merged) != 2 {
		t.Errorf("Expected 2 paths, got %d", len(merged))
	}

	pets, ok := merged["/pets"].(map[string]any)
	if !ok {
		t.Error("Expected /pets path to exist")
	}

	if len(pets) != 2 {
		t.Errorf("Expected 2 methods for /pets, got %d", len(pets))
	}
}

func TestMergeOperations(t *testing.T) {
	op1 := map[string]any{
		"summary": "List pets",
		"parameters": []any{
			map[string]any{"name": "limit", "in": "query"},
		},
		"responses": map[string]any{
			"200": map[string]any{
				"description": "OK",
			},
		},
	}

	op2 := map[string]any{
		"summary": "List pets with tags",
		"parameters": []any{
			map[string]any{"name": "tags", "in": "query"},
		},
		"responses": map[string]any{
			"400": map[string]any{
				"description": "Bad Request",
			},
		},
	}

	merged := mergeOperations(op1, op2)

	if merged["summary"] != "List pets with tags" {
		t.Error("Expected summary to be from op2")
	}

	params := merged["parameters"].([]any)
	if len(params) != 2 {
		t.Errorf("Expected 2 parameters, got %d", len(params))
	}
}

func TestMergeParameters(t *testing.T) {
	params1 := []any{
		map[string]any{"name": "limit", "in": "query"},
	}

	params2 := []any{
		map[string]any{"name": "limit", "in": "query"},
		map[string]any{"name": "offset", "in": "query"},
	}

	merged := mergeParameters(params1, params2)

	if len(merged) != 2 {
		t.Errorf("Expected 2 parameters, got %d", len(merged))
	}
}

func TestMergeSpecs(t *testing.T) {
	spec1 := map[string]any{
		"openapi": "3.0.0",
		"info": map[string]any{
			"title": "Petstore API",
		},
		"paths": map[string]any{
			"/pets": map[string]any{
				"get": map[string]any{
					"summary": "List pets",
				},
			},
		},
	}

	spec2 := map[string]any{
		"openapi": "3.0.0",
		"info": map[string]any{
			"version": "1.0.0",
		},
		"paths": map[string]any{
			"/pets/{id}": map[string]any{
				"get": map[string]any{
					"summary": "Get pet by ID",
				},
			},
		},
	}

	merged := MergeSpecs(spec1, spec2)

	if merged["openapi"] != "3.0.0" {
		t.Error("Expected openapi version to be 3.0.0")
	}

	info := merged["info"].(map[string]any)
	if info["title"] != "Petstore API" {
		t.Error("Expected title to be from spec1")
	}
	if info["version"] != "1.0.0" {
		t.Error("Expected version to be from spec2")
	}

	paths := merged["paths"].(map[string]any)
	if len(paths) != 2 {
		t.Errorf("Expected 2 paths, got %d", len(paths))
	}
}

func TestIntegrationWithRealFiles(t *testing.T) {
	spec1, err := ReadFile("test/test1.openapi.yaml")
	if err != nil {
		t.Fatalf("Failed to read test1.openapi.yaml: %v", err)
	}

	spec2, err := ReadFile("test/test2.openapi.yaml")
	if err != nil {
		t.Fatalf("Failed to read test2.openapi.yaml: %v", err)
	}

	merged := MergeSpecs(spec1, spec2)

	if merged["openapi"] != "3.0.0" {
		t.Error("Expected openapi version to be 3.0.0")
	}

	paths := merged["paths"].(map[string]any)
	if len(paths) != 3 {
		t.Errorf("Expected 3 paths, got %d", len(paths))
	}

	users, ok := paths["/users"].(map[string]any)
	if !ok {
		t.Error("Expected /users path to exist")
	}
	if len(users) != 2 {
		t.Errorf("Expected 2 methods for /users, got %d", len(users))
	}
	if _, exists := users["get"]; !exists {
		t.Error("Expected GET method to exist for /users")
	}
	if _, exists := users["post"]; !exists {
		t.Error("Expected POST method to exist for /users")
	}

	userById, ok := paths["/users/{id}"].(map[string]any)
	if !ok {
		t.Error("Expected /users/{id} path to exist")
	}
	if _, exists := userById["get"]; !exists {
		t.Error("Expected GET method to exist for /users/{id}")
	}

	products, ok := paths["/products"].(map[string]any)
	if !ok {
		t.Error("Expected /products path to exist")
	}
	if _, exists := products["get"]; !exists {
		t.Error("Expected GET method to exist for /products")
	}

	components := merged["components"].(map[string]any)
	schemas := components["schemas"].(map[string]any)
	expectedSchemas := []string{"User", "Product", "UserInput"}
	for _, schema := range expectedSchemas {
		if _, exists := schemas[schema]; !exists {
			t.Errorf("Expected schema %s to exist", schema)
		}
	}

	tags := merged["tags"].([]any)
	if len(tags) != 2 {
		t.Errorf("Expected 2 tags, got %d", len(tags))
	}

	tagNames := make(map[string]bool)
	for _, tag := range tags {
		tagMap := tag.(map[string]any)
		tagNames[tagMap["name"].(string)] = true
	}
	expectedTags := []string{"users", "products"}
	for _, expectedTag := range expectedTags {
		if !tagNames[expectedTag] {
			t.Errorf("Expected tag %s to exist", expectedTag)
		}
	}
}
