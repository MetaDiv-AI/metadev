package metadev

import (
	"time"

	"github.com/MetaDiv-AI/metadev/types"
)

// GET creates a new API route builder for GET requests
func GET[RequestType any, ResponseType any](module types.Module) *apiNameBuilder[RequestType, ResponseType] {
	return &apiNameBuilder[RequestType, ResponseType]{method: "GET", module: module}
}

// POST creates a new API route builder for POST requests
func POST[RequestType any, ResponseType any](module types.Module) *apiNameBuilder[RequestType, ResponseType] {
	return &apiNameBuilder[RequestType, ResponseType]{method: "POST", module: module}
}

// PUT creates a new API route builder for PUT requests
func PUT[RequestType any, ResponseType any](module types.Module) *apiNameBuilder[RequestType, ResponseType] {
	return &apiNameBuilder[RequestType, ResponseType]{method: "PUT", module: module}
}

// DELETE creates a new API route builder for DELETE requests
func DELETE[RequestType any, ResponseType any](module types.Module) *apiNameBuilder[RequestType, ResponseType] {
	return &apiNameBuilder[RequestType, ResponseType]{method: "DELETE", module: module}
}

// PATCH creates a new API route builder for PATCH requests
func PATCH[RequestType any, ResponseType any](module types.Module) *apiNameBuilder[RequestType, ResponseType] {
	return &apiNameBuilder[RequestType, ResponseType]{method: "PATCH", module: module}
}

type apiNameBuilder[RequestType any, ResponseType any] struct {
	method string
	module types.Module
}

func (a *apiNameBuilder[RequestType, ResponseType]) Name(name string) *apiRouteBuilder[RequestType, ResponseType] {
	return &apiRouteBuilder[RequestType, ResponseType]{method: a.method, name: name, module: a.module}
}

// apiRouteBuilder is a builder for API routes
type apiRouteBuilder[RequestType any, ResponseType any] struct {
	method string
	name   string
	module types.Module
}

// Route sets the route path for the API route
func (a *apiRouteBuilder[RequestType, ResponseType]) Route(route string) *apiOptionals[RequestType, ResponseType] {
	return &apiOptionals[RequestType, ResponseType]{
		method: a.method,
		module: a.module,
		name:   a.name,
		route:  route,
	}
}

// apiOptionals is a builder for API route options
type apiOptionals[RequestType any, ResponseType any] struct {
	method            string
	module            types.Module
	route             string
	name              string
	rateLimitDuration time.Duration
	rateLimitLimit    int
	cacheDuration     time.Duration
	skipTypescript    bool
	skipOpenApi       bool
	middlewares       []types.MiddlewareHandler
}

// SkipTypescript sets the skip typescript flag for the API route
func (a *apiOptionals[RequestType, ResponseType]) SkipTypescript() *apiOptionals[RequestType, ResponseType] {
	a.skipTypescript = true
	return a
}

// SkipOpenApi sets the skip openapi flag for the API route
func (a *apiOptionals[RequestType, ResponseType]) SkipOpenApi() *apiOptionals[RequestType, ResponseType] {
	a.skipOpenApi = true
	return a
}

// RateLimit sets the rate limit for the API route
func (a *apiOptionals[RequestType, ResponseType]) RateLimit(duration time.Duration, limit int) *apiOptionals[RequestType, ResponseType] {
	a.rateLimitDuration = duration
	a.rateLimitLimit = limit
	return a
}

// Cache sets the cache duration for the API route
func (a *apiOptionals[RequestType, ResponseType]) Cache(duration time.Duration) *apiOptionals[RequestType, ResponseType] {
	a.cacheDuration = duration
	return a
}

// Middleware sets the middleware for the API route
func (a *apiOptionals[RequestType, ResponseType]) Middleware(middleware types.MiddlewareHandler) *apiOptionals[RequestType, ResponseType] {
	a.middlewares = append(a.middlewares, middleware)
	return a
}

// Handler sets the handler for the API route
func (a *apiOptionals[RequestType, ResponseType]) Handler(handler func(ctx types.ApiContext[RequestType, ResponseType])) types.Module {
	h := types.NewApiHandler(a.module, a.method, a.route, a.name, a.rateLimitDuration, a.rateLimitLimit, a.cacheDuration, a.middlewares, handler, nil, a.skipTypescript, a.skipOpenApi)
	a.module.RegisterApiHandler(h)
	return a.module
}

// PublicHandler sets the public handler for the API route (no authentication required)
func (a *apiOptionals[RequestType, ResponseType]) PublicHandler(publicHandler func(ctx types.PublicApiContext[RequestType, ResponseType])) types.Module {
	h := types.NewApiHandler(a.module, a.method, a.route, a.name, a.rateLimitDuration, a.rateLimitLimit, a.cacheDuration, a.middlewares, nil, publicHandler, a.skipTypescript, a.skipOpenApi)
	a.module.RegisterApiHandler(h)
	return a.module
}
