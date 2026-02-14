package types

func NewMiddlewareHandler(module Module, name string, handler func(ctx MiddlewareContext)) MiddlewareHandler {
	return &middlewareHandler{module: module, name: name, handler: handler}
}

type MiddlewareHandler interface {
	// Module returns the module of the middleware handler
	Module() Module
	// Name returns the name of the middleware handler
	Name() string
	// Handler
	Handler() func(ctx MiddlewareContext)
}

type middlewareHandler struct {
	module  Module
	name    string
	handler func(ctx MiddlewareContext)
}

func (h *middlewareHandler) Module() Module {
	return h.module
}

func (h *middlewareHandler) Name() string {
	return h.name
}

func (h *middlewareHandler) Handler() func(ctx MiddlewareContext) {
	return h.handler
}
