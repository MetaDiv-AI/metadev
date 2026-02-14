package openapi

import "encoding/json"

func NewBuilder() *Builder {
	return &Builder{
		OpenAPI: &OpenAPI{
			OpenAPI: "3.0.1",
			Info: Info{
				Title:       "MetaDev APIs",
				Description: "Powered by Metadiv Studio & Lab",
				Version:     "1.0.0",
			},
			Paths: make(Paths),
		},
	}
}

type Builder struct {
	OpenAPI *OpenAPI
}

func (b *Builder) AddPath(path string, method string, operation *Operation) {
	if _, ok := b.OpenAPI.Paths[path]; !ok {
		b.OpenAPI.Paths[path] = new(PathItem)
	}
	switch method {
	case "GET":
		b.OpenAPI.Paths[path].Get = operation
	case "POST":
		b.OpenAPI.Paths[path].Post = operation
	case "PUT":
		b.OpenAPI.Paths[path].Put = operation
	case "DELETE":
		b.OpenAPI.Paths[path].Delete = operation
	}
}

func (b *Builder) ToJSON() string {
	j, _ := json.MarshalIndent(b.OpenAPI, "", "  ")
	return string(j)
}
