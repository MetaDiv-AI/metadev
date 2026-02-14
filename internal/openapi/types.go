package openapi

// OpenAPI represents the root structure of an OpenAPI 3.0.1 document
type OpenAPI struct {
	OpenAPI    string     `json:"openapi"`
	Info       Info       `json:"info"`
	Tags       []Tag      `json:"tags"`
	Paths      Paths      `json:"paths"`
	Components Components `json:"components"`
	Servers    []Server   `json:"servers"`
}

// Info contains metadata about the API
type Info struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Version     string `json:"version"`
}

// Tag represents an API tag
type Tag struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// Paths holds the relative paths to the individual endpoints
type Paths map[string]*PathItem

// PathItem describes the operations available on a single path
type PathItem struct {
	Get     *Operation `json:"get,omitempty"`
	Post    *Operation `json:"post,omitempty"`
	Put     *Operation `json:"put,omitempty"`
	Delete  *Operation `json:"delete,omitempty"`
	Options *Operation `json:"options,omitempty"`
	Head    *Operation `json:"head,omitempty"`
	Patch   *Operation `json:"patch,omitempty"`
}

// Operation describes a single API operation on a path
type Operation struct {
	Summary     string                `json:"summary"`
	Description string                `json:"description"`
	Tags        []string              `json:"tags"`
	Parameters  []Parameter           `json:"parameters"`
	RequestBody *RequestBody          `json:"requestBody,omitempty"`
	Responses   map[string]Response   `json:"responses"`
	Deprecated  bool                  `json:"deprecated"`
	Security    []map[string][]string `json:"security"`
}

// Parameter describes a single operation parameter
type Parameter struct {
	Name        string  `json:"name"`
	In          string  `json:"in"` // query, header, path, cookie
	Description string  `json:"description"`
	Required    bool    `json:"required"`
	Schema      *Schema `json:"schema"`
}

// RequestBody describes a single request body
type RequestBody struct {
	Description string               `json:"description"`
	Content     map[string]MediaType `json:"content"`
	Required    bool                 `json:"required"`
}

// MediaType represents a media type object
type MediaType struct {
	Schema   *Schema            `json:"schema"`
	Example  interface{}        `json:"example,omitempty"`
	Examples map[string]Example `json:"examples,omitempty"`
}

// Example represents an example object
type Example struct {
	Summary     string      `json:"summary"`
	Description string      `json:"description"`
	Value       interface{} `json:"value"`
}

// Response describes a single response from an API Operation
type Response struct {
	Description string               `json:"description"`
	Headers     map[string]Header    `json:"headers"`
	Content     map[string]MediaType `json:"content"`
}

// Header represents a header parameter
type Header struct {
	Description string  `json:"description"`
	Schema      *Schema `json:"schema"`
}

// Schema represents a JSON Schema object
type Schema struct {
	Type       string             `json:"type,omitempty"`
	Format     string             `json:"format,omitempty"`
	Properties map[string]*Schema `json:"properties,omitempty"`
	Required   []string           `json:"required,omitempty"`
	Items      *Schema            `json:"items,omitempty"`
	Ref        string             `json:"$ref,omitempty"`
}

// Components holds various schemas for the specification
type Components struct {
	Schemas         map[string]Schema         `json:"schemas"`
	SecuritySchemes map[string]SecurityScheme `json:"securitySchemes"`
}

// SecurityScheme represents a security scheme
type SecurityScheme struct {
	Type         string `json:"type"`
	Description  string `json:"description"`
	Name         string `json:"name"`
	In           string `json:"in"`
	Scheme       string `json:"scheme"`
	BearerFormat string `json:"bearerFormat"`
}

// Server represents a server object
type Server struct {
	URL         string                    `json:"url"`
	Description string                    `json:"description"`
	Variables   map[string]ServerVariable `json:"variables"`
}

// ServerVariable represents a server variable for server URL template substitution
type ServerVariable struct {
	Enum        []string `json:"enum"`
	Default     string   `json:"default"`
	Description string   `json:"description"`
}
