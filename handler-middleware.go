package metadev

import (
	"github.com/MetaDiv-AI/metadev/types"
)

// Middleware creates a new middleware handler builder
func Middleware(module types.Module) *middlewareNameBuilder {
	return &middlewareNameBuilder{module: module}
}

type middlewareNameBuilder struct {
	module types.Module
}

// Name sets the name of the middleware handler
func (b *middlewareNameBuilder) Name(name string) *middlewareHandlerBuilder {
	return &middlewareHandlerBuilder{module: b.module, name: name}
}

type middlewareHandlerBuilder struct {
	module types.Module
	name   string
}

// Handler sets the handler for the middleware
func (b *middlewareHandlerBuilder) Handler(handler func(ctx types.MiddlewareContext)) types.MiddlewareHandler {
	return types.NewMiddlewareHandler(b.module, b.name, handler)
}
