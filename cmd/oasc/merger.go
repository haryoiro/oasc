package main

import (
	"github.com/getkin/kin-openapi/openapi3"
)

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
