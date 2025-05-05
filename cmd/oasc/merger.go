package main

import (
	"encoding/json"

	"github.com/getkin/kin-openapi/openapi3"
)

func mergeOpenAPIPaths(paths1, paths2 map[string]any) map[string]any {
	result := make(map[string]any)

	for path, methods1 := range paths1 {
		methods2, exists := paths2[path]
		if !exists {
			result[path] = methods1
			continue
		}

		mergedMethods := make(map[string]any)
		m1 := methods1.(map[string]any)
		m2 := methods2.(map[string]any)

		for method, operation1 := range m1 {
			operation2, exists := m2[method]
			if !exists {
				mergedMethods[method] = operation1
				continue
			}
			mergedMethods[method] = mergeOperations(operation1.(map[string]any), operation2.(map[string]any))
		}

		for method, operation2 := range m2 {
			if _, exists := m1[method]; !exists {
				mergedMethods[method] = operation2
			}
		}

		result[path] = mergedMethods
	}

	for path, methods2 := range paths2 {
		if _, exists := paths1[path]; !exists {
			result[path] = methods2
		}
	}

	return result
}

func mergeOperations(op1, op2 map[string]any) map[string]any {
	result := make(map[string]any)

	for k, v1 := range op1 {
		v2, exists := op2[k]
		if !exists {
			result[k] = v1
			continue
		}

		switch k {
		case "parameters":
			result[k] = mergeParameters(v1.([]any), v2.([]any))
		case "responses":
			result[k] = mergeResponses(v1.(map[string]any), v2.(map[string]any))
		case "requestBody":
			result[k] = mergeRequestBody(v1.(map[string]any), v2.(map[string]any))
		default:
			result[k] = v2
		}
	}

	for k, v2 := range op2 {
		if _, exists := op1[k]; !exists {
			result[k] = v2
		}
	}

	return result
}

func mergeParameters(params1, params2 []any) []any {
	seen := make(map[string]bool)
	var result []any

	for _, param := range params1 {
		if p, ok := param.(map[string]any); ok {
			if name, ok := p["name"].(string); ok {
				seen[name] = true
			}
		}
		result = append(result, param)
	}

	for _, param := range params2 {
		if p, ok := param.(map[string]any); ok {
			if name, ok := p["name"].(string); ok {
				if !seen[name] {
					result = append(result, param)
				}
			}
		}
	}

	return result
}

func mergeResponses(resp1, resp2 map[string]any) map[string]any {
	result := make(map[string]any)

	for k, v1 := range resp1 {
		v2, exists := resp2[k]
		if !exists {
			result[k] = v1
			continue
		}
		result[k] = v2
	}

	for k, v2 := range resp2 {
		if _, exists := resp1[k]; !exists {
			result[k] = v2
		}
	}

	return result
}

func mergeRequestBody(body1, body2 map[string]any) map[string]any {
	result := make(map[string]any)

	for k, v1 := range body1 {
		v2, exists := body2[k]
		if !exists {
			result[k] = v1
			continue
		}
		result[k] = v2
	}

	for k, v2 := range body2 {
		if _, exists := body1[k]; !exists {
			result[k] = v2
		}
	}

	return result
}

func mergeArrays(a1, a2 []any) []any {
	seen := make(map[string]bool)
	var result []any

	// タグの場合は名前で重複を判定
	isTag := false
	if len(a1) > 0 {
		if tag, ok := a1[0].(map[string]any); ok {
			if _, ok := tag["name"]; ok {
				isTag = true
			}
		}
	}

	for _, item := range a1 {
		if isTag {
			if tag, ok := item.(map[string]any); ok {
				if name, ok := tag["name"].(string); ok {
					seen[name] = true
				}
			}
		} else {
			if str, ok := item.(string); ok {
				if !seen[str] {
					seen[str] = true
				}
			} else if m, ok := item.(map[string]any); ok {
				if jsonBytes, err := json.Marshal(m); err == nil {
					jsonStr := string(jsonBytes)
					if !seen[jsonStr] {
						seen[jsonStr] = true
					}
				}
			}
		}
		result = append(result, item)
	}

	for _, item := range a2 {
		if isTag {
			if tag, ok := item.(map[string]any); ok {
				if name, ok := tag["name"].(string); ok {
					if !seen[name] {
						seen[name] = true
						result = append(result, item)
					} else {
						// 既存のタグを更新
						for i, existing := range result {
							if existingTag, ok := existing.(map[string]any); ok {
								if existingName, ok := existingTag["name"].(string); ok && existingName == name {
									result[i] = item
									break
								}
							}
						}
					}
				}
			}
		} else {
			if str, ok := item.(string); ok {
				if !seen[str] {
					seen[str] = true
					result = append(result, str)
				}
			} else if m, ok := item.(map[string]any); ok {
				if jsonBytes, err := json.Marshal(m); err == nil {
					jsonStr := string(jsonBytes)
					if !seen[jsonStr] {
						seen[jsonStr] = true
						result = append(result, m)
					}
				}
			} else {
				result = append(result, item)
			}
		}
	}

	return result
}

func MergeSpecs(spec1, spec2 *openapi3.T) *openapi3.T {
	if spec1 == nil {
		return spec2
	}
	if spec2 == nil {
		return spec1
	}

	// Always use the latest spec's Info
	spec1.Info = spec2.Info

	// Merge paths
	if spec2.Paths != nil {
		if spec1.Paths == nil {
			spec1.Paths = spec2.Paths
		} else {
			for path, pathItem := range spec2.Paths.Map() {
				if existing, ok := spec1.Paths.Map()[path]; ok {
					mergePathItem(existing, pathItem)
				} else {
					spec1.Paths.Set(path, pathItem)
				}
			}
		}
	}

	// Merge components
	if spec2.Components != nil {
		if spec1.Components == nil {
			spec1.Components = spec2.Components
		} else {
			// Merge schemas
			if spec2.Components.Schemas != nil {
				if spec1.Components.Schemas == nil {
					spec1.Components.Schemas = spec2.Components.Schemas
				} else {
					for name, schema := range spec2.Components.Schemas {
						spec1.Components.Schemas[name] = schema
					}
				}
			}

			// Merge responses
			if spec2.Components.Responses != nil {
				if spec1.Components.Responses == nil {
					spec1.Components.Responses = spec2.Components.Responses
				} else {
					for name, response := range spec2.Components.Responses {
						spec1.Components.Responses[name] = response
					}
				}
			}

			// Merge parameters
			if spec2.Components.Parameters != nil {
				if spec1.Components.Parameters == nil {
					spec1.Components.Parameters = spec2.Components.Parameters
				} else {
					for name, parameter := range spec2.Components.Parameters {
						spec1.Components.Parameters[name] = parameter
					}
				}
			}

			// Merge request bodies
			if spec2.Components.RequestBodies != nil {
				if spec1.Components.RequestBodies == nil {
					spec1.Components.RequestBodies = spec2.Components.RequestBodies
				} else {
					for name, requestBody := range spec2.Components.RequestBodies {
						spec1.Components.RequestBodies[name] = requestBody
					}
				}
			}
		}
	}

	// Merge tags with update
	if spec2.Tags != nil {
		if spec1.Tags == nil {
			spec1.Tags = spec2.Tags
		} else {
			// Create a map for quick lookup
			tagMap := make(map[string]*openapi3.Tag)
			for _, tag := range spec1.Tags {
				tagMap[tag.Name] = tag
			}

			// Update or add tags from spec2
			for _, tag2 := range spec2.Tags {
				if tag1, exists := tagMap[tag2.Name]; exists {
					// Update existing tag
					tag1.Description = tag2.Description
					tag1.ExternalDocs = tag2.ExternalDocs
					tag1.Extensions = tag2.Extensions
				} else {
					// Add new tag
					spec1.Tags = append(spec1.Tags, tag2)
				}
			}
		}
	}

	// Merge security requirements
	if spec2.Security != nil {
		spec1.Security = spec2.Security
	}

	return spec1
}

func mergePathItem(dst, src *openapi3.PathItem) {
	if src.Get != nil {
		dst.Get = src.Get
	}
	if src.Put != nil {
		dst.Put = src.Put
	}
	if src.Post != nil {
		dst.Post = src.Post
	}
	if src.Delete != nil {
		dst.Delete = src.Delete
	}
	if src.Options != nil {
		dst.Options = src.Options
	}
	if src.Head != nil {
		dst.Head = src.Head
	}
	if src.Patch != nil {
		dst.Patch = src.Patch
	}
	if src.Trace != nil {
		dst.Trace = src.Trace
	}
	if src.Parameters != nil {
		dst.Parameters = src.Parameters
	}
}
