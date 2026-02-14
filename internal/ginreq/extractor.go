package ginreq

import (
	"fmt"
	"io"
	"reflect"

	"github.com/gin-gonic/gin"
)

func NewExtractor[RequestType any](ginContext *gin.Context) Extractor[RequestType] {
	return &ExtractorImpl[RequestType]{
		GinContext: ginContext,
	}
}

type Extractor[RequestType any] interface {
	ExtractFile() (file []byte, filename string, fileSize int64, err error)
	ExtractRequest() (request *RequestType, err error)
}

type ExtractorImpl[RequestType any] struct {
	GinContext *gin.Context
}

func (e *ExtractorImpl[RequestType]) ExtractFile() (file []byte, filename string, fileSize int64, err error) {
	fileHeader, err := e.GinContext.FormFile("file")
	if err != nil {
		return nil, "", 0, err
	}

	fileReader, err := fileHeader.Open()
	if err != nil {
		return nil, "", 0, err
	}
	defer fileReader.Close()

	content, err := io.ReadAll(fileReader)
	if err != nil {
		return nil, "", 0, err
	}

	return content, fileHeader.Filename, fileHeader.Size, nil
}

func (e *ExtractorImpl[RequestType]) ExtractRequest() (request *RequestType, err error) {
	// Check if RequestType is Empty using type assertion instead of reflection
	request = new(RequestType)
	if e.isEmptyRequest(request) {
		return request, nil
	}

	objects := make([]RequestType, 0)
	tags := parseTags(request)

	// If no binding tags are present, return empty struct
	if len(tags) == 0 {
		return request, nil
	}

	for _, tag := range tags {
		switch tag {
		case "json":
			jsonRequest := new(RequestType)
			if err := e.GinContext.ShouldBindJSON(jsonRequest); err != nil {
				return nil, fmt.Errorf("failed to bind JSON: %w", err)
			}
			objects = append(objects, *jsonRequest)
		case "form":
			formRequest := new(RequestType)
			if err := e.GinContext.ShouldBindQuery(formRequest); err != nil {
				return nil, fmt.Errorf("failed to bind query parameters: %w", err)
			}
			objects = append(objects, *formRequest)
		case "uri":
			uriRequest := new(RequestType)
			if err := e.GinContext.ShouldBindUri(uriRequest); err != nil {
				return nil, fmt.Errorf("failed to bind URI parameters: %w", err)
			}
			objects = append(objects, *uriRequest)
		}
	}

	result := e.syncAllRequestsToOneRequest(objects)
	if result == nil {
		return new(RequestType), nil
	}

	// Auto-validate if the request implements ValidatedRequest interface
	if err := e.validateRequest(result); err != nil {
		return nil, err
	}

	return result, nil
}

func (e *ExtractorImpl[RequestType]) isEmptyRequest(request *RequestType) bool {
	typeName := reflect.TypeOf(request).Elem().String()
	return typeName == "metadev.Empty"
}

func (e *ExtractorImpl[RequestType]) syncAllRequestsToOneRequest(objects []RequestType) *RequestType {
	if len(objects) == 0 {
		return nil
	}

	objectVal := reflect.ValueOf(new(RequestType)).Elem()
	for _, o := range objects {
		val := reflect.ValueOf(o)
		for i := 0; i < val.NumField(); i++ {
			field := val.Field(i)
			if !field.IsZero() && objectVal.Field(i).CanSet() {
				if objectVal.Field(i).IsZero() && !objectVal.Field(i).CanSet() {
					continue
				}
				switch field.Kind() {
				case reflect.String:
					objectVal.Field(i).SetString(field.String())
				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
					objectVal.Field(i).SetInt(field.Int())
				case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
					objectVal.Field(i).SetUint(field.Uint())
				case reflect.Float32, reflect.Float64:
					objectVal.Field(i).SetFloat(field.Float())
				case reflect.Bool:
					objectVal.Field(i).SetBool(field.Bool())
				case reflect.Slice, reflect.Array, reflect.Map, reflect.Struct, reflect.Ptr:
					objectVal.Field(i).Set(field)
				}
			}
		}
	}

	output := objectVal.Interface().(RequestType)
	return &output
}

// validateRequest checks if the request implements ValidatedRequest and validates it
func (e *ExtractorImpl[RequestType]) validateRequest(request *RequestType) error {
	// Try type assertion to ValidatedRequest interface
	if validatedReq, ok := any(request).(interface{ Validate() error }); ok {
		return validatedReq.Validate()
	}
	// If not implemented, validation passes (optional validation)
	return nil
}
