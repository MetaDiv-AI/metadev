package ginreq

import "reflect"

func parseTags[T any](request T) []string {
	m := make(map[string]bool)

	for i := 0; i < reflect.TypeOf(request).Elem().NumField(); i++ {
		f := reflect.TypeOf(request).Elem().Field(i)

		if f.Anonymous {
			tags := parseTags(reflect.New(f.Type).Interface())
			for _, tag := range tags {
				m[tag] = true
			}
			continue
		}

		tag := reflect.TypeOf(request).Elem().Field(i).Tag

		// Check which tags exist on this field
		hasJSON := len(tag.Get("json")) > 0
		hasForm := len(tag.Get("form")) > 0
		hasURI := len(tag.Get("uri")) > 0

		// Apply priority rules: ignore "json" if "form" or "uri" is present
		if hasForm {
			m["form"] = true
		}
		if hasURI {
			m["uri"] = true
		}
		if hasJSON && !hasForm && !hasURI {
			m["json"] = true
		}
	}

	result := make([]string, 0)
	for key := range m {
		result = append(result, key)
	}
	return result
}
