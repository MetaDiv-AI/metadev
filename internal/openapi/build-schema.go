package openapi

import (
	"fmt"
	"reflect"
	"strings"
)

type typeStack struct {
	processed map[string]bool
}

func newTypeStack() *typeStack {
	return &typeStack{
		processed: make(map[string]bool),
	}
}

func (s *typeStack) push(t reflect.Type) bool {
	key := t.PkgPath() + "." + t.Name()
	if s.processed[key] {
		return false
	}
	s.processed[key] = true
	return true
}

func (s *typeStack) pop(t reflect.Type) {
	key := t.PkgPath() + "." + t.Name()
	delete(s.processed, key)
}

func BuildSchemaForJson[Model any]() *Schema {
	t := reflect.TypeOf((*Model)(nil)).Elem()
	return buildSchemaFromType(t, "json", newTypeStack())
}

func BuildSchemaForParameters[Model any]() []Parameter {
	t := reflect.TypeOf((*Model)(nil)).Elem()
	return buildParametersFromType(t, newTypeStack())
}

func buildParametersFromType(t reflect.Type, stack *typeStack) []Parameter {
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if t.Kind() != reflect.Struct {
		return nil
	}

	var parameters []Parameter
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		if field.Anonymous {
			embeddedParams := buildParametersFromType(field.Type, stack)
			parameters = append(parameters, embeddedParams...)
			continue
		}

		if formTag := field.Tag.Get("form"); formTag != "" {
			name := strings.Split(formTag, ",")[0]
			if name != "-" {
				param := Parameter{
					Name:        name,
					In:          "query",
					Description: "",
					Required:    !isFieldOptional(field),
					Schema:      buildSchemaFromType(field.Type, "form", stack),
				}
				parameters = append(parameters, param)
			}
		}

		if uriTag := field.Tag.Get("uri"); uriTag != "" {
			name := strings.Split(uriTag, ",")[0]
			if name != "-" {
				param := Parameter{
					Name:        name,
					In:          "path",
					Description: "",
					Required:    true,
					Schema:      buildSchemaFromType(field.Type, "uri", stack),
				}
				parameters = append(parameters, param)
			}
		}
	}

	return parameters
}

func buildSchemaFromType(t reflect.Type, tagType string, stack *typeStack) *Schema {
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if t.Kind() == reflect.Slice {
		return &Schema{
			Type:  "array",
			Items: buildSchemaFromType(t.Elem(), tagType, stack),
		}
	}

	if schema := handlePrimitiveType(t); schema != nil {
		return schema
	}

	if t.Kind() == reflect.Struct {
		if !stack.push(t) {
			return &Schema{
				Type: "object",
				Ref:  fmt.Sprintf("#/components/schemas/%s", t.Name()),
			}
		}
		defer stack.pop(t)

		schema := &Schema{
			Type:       "object",
			Properties: make(map[string]*Schema),
			Required:   []string{},
		}

		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)

			if field.Anonymous {
				embeddedSchema := buildSchemaFromType(field.Type, tagType, stack)
				if embeddedSchema.Properties != nil {
					for k, v := range embeddedSchema.Properties {
						schema.Properties[k] = v
					}
					schema.Required = append(schema.Required, embeddedSchema.Required...)
				}
				continue
			}

			tag := field.Tag.Get(tagType)
			if tag == "" || tag == "-" {
				continue
			}

			name := strings.Split(tag, ",")[0]
			fieldSchema := buildSchemaFromType(field.Type, tagType, stack)
			schema.Properties[name] = fieldSchema

			if !isFieldOptional(field) {
				schema.Required = append(schema.Required, name)
			}
		}

		if len(schema.Properties) == 0 {
			schema.Properties = nil
		}
		if len(schema.Required) == 0 {
			schema.Required = nil
		}

		return schema
	}

	return &Schema{Type: "string"}
}

func handlePrimitiveType(t reflect.Type) *Schema {
	switch t.Kind() {
	case reflect.Bool:
		return &Schema{Type: "boolean"}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return &Schema{Type: "integer", Format: getIntegerFormat(t)}
	case reflect.Float32, reflect.Float64:
		return &Schema{Type: "number", Format: getFloatFormat(t)}
	case reflect.String:
		return &Schema{Type: "string"}
	default:
		return nil
	}
}

func isFieldOptional(field reflect.StructField) bool {
	bindingTag := field.Tag.Get("binding")
	return !strings.Contains(bindingTag, "required")
}

func getIntegerFormat(t reflect.Type) string {
	switch t.Kind() {
	case reflect.Int64, reflect.Uint64:
		return "int64"
	case reflect.Int32, reflect.Uint32:
		return "int32"
	default:
		return ""
	}
}

func getFloatFormat(t reflect.Type) string {
	switch t.Kind() {
	case reflect.Float64:
		return "float64"
	case reflect.Float32:
		return "float32"
	default:
		return ""
	}
}
