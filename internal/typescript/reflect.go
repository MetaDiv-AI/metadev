package typescript

import "reflect"

func GetType[Model any]() reflect.Type {
	typ := reflect.TypeOf((*Model)(nil))
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	if typ.Kind() == reflect.Slice {
		typ = typ.Elem()
	}
	return typ
}

func GetName[Model any]() string {
	name := GetType[Model]().Name()
	typ := reflect.TypeOf((*Model)(nil))
	if typ.Kind() == reflect.Ptr && typ.Elem().Kind() == reflect.Slice {
		name = name + "[]"
	}
	return name
}

func GetForms[Model any]() []string {
	return getForms(GetType[Model]())
}

func getForms(typ reflect.Type) []string {
	forms := make([]string, 0)

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)

		// Handle embedded struct
		if field.Anonymous {
			forms = append(forms, getForms(field.Type)...)
			continue
		}

		if formTag := field.Tag.Get("form"); formTag != "" && formTag != "-" {
			forms = append(forms, formTag)
		}
	}

	return forms
}

func GetUris[Model any]() []string {
	return getUris(GetType[Model]())
}

func getUris(typ reflect.Type) []string {
	uris := make([]string, 0)

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)

		// Handle embedded struct
		if field.Anonymous {
			uris = append(uris, getUris(field.Type)...)
			continue
		}

		if uriTag := field.Tag.Get("uri"); uriTag != "" && uriTag != "-" {
			uris = append(uris, uriTag)
		}
	}

	return uris
}

func CheckTypeIsJson[Model any]() bool {
	return checkTypeIsJson(GetType[Model]())
}

func checkTypeIsJson(typ reflect.Type) bool {
	if typ.Name() == "RequestListing" || typ.Name() == "Empty" || typ.Name() == "RequestPathId" || typ.Name() == "RequestPathUUID" {
		return false
	}
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)

		// Handle embedded struct
		if field.Anonymous {
			if field.Type.Name() == "RequestListing" || field.Type.Name() == "Empty" || field.Type.Name() == "RequestPathId" || typ.Name() == "RequestPathUUID" {
				continue
			}
			if checkTypeIsJson(field.Type) {
				return true
			}
			continue
		}

		// Check if field has json tag
		jsonTag := field.Tag.Get("json")
		if jsonTag != "" && jsonTag != "-" {
			return true
		}
	}

	return false
}
