package types

import (
	"strings"
	"time"

	"github.com/MetaDiv-AI/metadev/internal/openapi"
	"github.com/MetaDiv-AI/metadev/internal/typescript"

	"github.com/gin-gonic/gin"
)

func NewApiHandler[RequestType any, ResponseType any](
	module Module, method string, route string, name string,
	rateLimitDuration time.Duration, rateLimitLimit int,
	cacheDuration time.Duration,
	middlewares []MiddlewareHandler,
	handler func(ctx ApiContext[RequestType, ResponseType]),
	publicHandler func(ctx PublicApiContext[RequestType, ResponseType]),
	skipTypescript bool,
	skipOpenApi bool,
) ApiHandler {
	return &apiHandler[RequestType, ResponseType]{
		module:            module,
		method:            method,
		name:              name,
		route:             route,
		rateLimitDuration: rateLimitDuration,
		rateLimitLimit:    rateLimitLimit,
		cacheDuration:     cacheDuration,
		middlewares:       middlewares,
		handler:           handler,
		publicHandler:     publicHandler,
		skipTypescript:    skipTypescript,
		skipOpenApi:       skipOpenApi,
	}
}

type ApiHandler interface {
	Name() string
	Module() Module
	Method() string
	Route() string
	RateLimit() (duration time.Duration, limit int)
	Cache() (duration time.Duration)
	Middlewares() []MiddlewareHandler
	GinHandler() gin.HandlerFunc
	SkipTypescript() bool
	SkipOpenApi() bool
	TypescriptInfo() typescript.ApiInfo
	RegisterOpenApi(builder *openapi.Builder)
}

type apiHandler[RequestType any, ResponseType any] struct {
	module Module
	name   string
	method string
	route  string

	middlewares []MiddlewareHandler

	rateLimitDuration time.Duration
	rateLimitLimit    int

	cacheDuration time.Duration

	skipTypescript bool
	skipOpenApi    bool

	handler       func(ctx ApiContext[RequestType, ResponseType])
	publicHandler func(ctx PublicApiContext[RequestType, ResponseType])
}

func (h *apiHandler[RequestType, ResponseType]) Module() Module {
	return h.module
}

func (h *apiHandler[RequestType, ResponseType]) Name() string {
	return h.name
}

func (h *apiHandler[RequestType, ResponseType]) Method() string {
	return h.method
}

func (h *apiHandler[RequestType, ResponseType]) Route() string {
	route := ""
	if h.publicHandler != nil {
		route += "/public"
	}
	route += "/" + h.module.App().Name()
	route += "/" + strings.TrimPrefix(h.route, "/")
	return route
}

func (h *apiHandler[RequestType, ResponseType]) SkipTypescript() bool {
	return h.skipTypescript
}

func (h *apiHandler[RequestType, ResponseType]) SkipOpenApi() bool {
	return h.skipOpenApi
}

func (h *apiHandler[RequestType, ResponseType]) RateLimit() (duration time.Duration, limit int) {
	return h.rateLimitDuration, h.rateLimitLimit
}

func (h *apiHandler[RequestType, ResponseType]) Cache() (duration time.Duration) {
	return h.cacheDuration
}

func (h *apiHandler[RequestType, ResponseType]) Middlewares() []MiddlewareHandler {
	return h.middlewares
}

func (h *apiHandler[RequestType, ResponseType]) GinHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if h.publicHandler != nil {
			c, err := NewPublicApiContext[RequestType, ResponseType](ctx, h)
			if err != nil {
				// Create a proper error response for validation errors
				responseCtx := NewResponseContext[ResponseType]()
				responseCtx.Error(err)
				responseCtx.MakeGinResponse(ctx)
				ctx.Abort()
				return
			}
			h.publicHandler(c)
			c.MakeGinResponse(ctx)
		} else {
			c, err := NewApiContext[RequestType, ResponseType](ctx, h)
			if err != nil {
				// Create a proper error response for validation errors
				responseCtx := NewResponseContext[ResponseType]()
				responseCtx.Error(err)
				responseCtx.MakeGinResponse(ctx)
				ctx.Abort()
				return
			}
			h.handler(c)
			c.MakeGinResponse(ctx)
		}
	}
}

func (h *apiHandler[RequestType, ResponseType]) TypescriptInfo() typescript.ApiInfo {
	apiInfo := typescript.ApiInfo{
		Name:   h.Name(),
		Route:  h.Route(),
		Method: h.Method(),
		Uris:   typescript.GetUris[RequestType](),
		Forms:  typescript.GetForms[RequestType](),
	}

	if typescript.CheckTypeIsJson[RequestType]() {
		apiInfo.Request = typescript.GetName[RequestType]()
		apiInfo.RequestType = typescript.GetType[RequestType]()
	}

	response := typescript.GetName[ResponseType]()
	if response != "Empty" && typescript.CheckTypeIsJson[ResponseType]() {
		apiInfo.Response = response
		apiInfo.ResponseType = typescript.GetType[ResponseType]()
	}
	return apiInfo
}

func (h *apiHandler[RequestType, ResponseType]) RegisterOpenApi(builder *openapi.Builder) {
	operation := new(openapi.Operation)
	operation.Summary = h.Name()

	parameters := openapi.BuildSchemaForParameters[RequestType]()
	operation.Parameters = parameters

	requestBodySchema := openapi.BuildSchemaForJson[RequestType]()
	requestBody := new(openapi.RequestBody)
	requestBody.Content = map[string]openapi.MediaType{
		"application/json": {
			Schema: requestBodySchema,
		},
	}
	operation.RequestBody = requestBody

	responseSchema := openapi.BuildSchemaForJson[Response[ResponseType]]()
	response := new(openapi.Response)
	response.Description = "Success"
	response.Content = map[string]openapi.MediaType{
		"application/json": {
			Schema: responseSchema,
		},
	}
	operation.Responses = map[string]openapi.Response{
		"200": *response,
	}

	// Format the route before adding it to the OpenAPI spec
	formattedRoute := formatRouteForOpenAPI(h.Route())
	builder.AddPath(formattedRoute, h.Method(), operation)
}

// formatRouteForOpenAPI converts gin route format to OpenAPI format
// e.g.: /users/:id -> /users/{id}
func formatRouteForOpenAPI(route string) string {
	parts := strings.Split(route, "/")
	for i, part := range parts {
		if strings.HasPrefix(part, ":") {
			parts[i] = "{" + strings.TrimPrefix(part, ":") + "}"
		}
	}
	return strings.Join(parts, "/")
}
