package main

import (
	"encoding/json"
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

func MergeSpecs(spec1, spec2 map[string]any) map[string]any {
	result := make(map[string]any)

	for k, v := range spec1 {
		result[k] = v
	}

	for k, v2 := range spec2 {
		v1, exists := result[k]
		if !exists {
			result[k] = v2
			continue
		}

		switch k {
		case "paths":
			result[k] = mergeOpenAPIPaths(v1.(map[string]any), v2.(map[string]any))
		case "components":
			if m1, ok1 := v1.(map[string]any); ok1 {
				if m2, ok2 := v2.(map[string]any); ok2 {
					result[k] = MergeSpecs(m1, m2)
				}
			}
		case "info":
			if m1, ok1 := v1.(map[string]any); ok1 {
				if m2, ok2 := v2.(map[string]any); ok2 {
					merged := make(map[string]any)
					for k, v := range m1 {
						merged[k] = v
					}
					for k, v := range m2 {
						if _, exists := merged[k]; !exists {
							merged[k] = v
						}
					}
					result[k] = merged
				}
			}
		case "tags":
			result[k] = mergeArrays(v1.([]any), v2.([]any))
		default:
			result[k] = v2
		}
	}

	return result
}
