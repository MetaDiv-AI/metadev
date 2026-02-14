package metadev

import (
	"fmt"
	"time"

	"github.com/MetaDiv-AI/metadev/types"
)

func NewWsHandler[InitRequestType any, RequestType any, ResponseType any](module types.Module) *wsNameBuilder[InitRequestType, RequestType, ResponseType] {
	return &wsNameBuilder[InitRequestType, RequestType, ResponseType]{module: module}
}

type wsNameBuilder[InitRequestType any, RequestType any, ResponseType any] struct {
	module types.Module
}

func (b *wsNameBuilder[InitRequestType, RequestType, ResponseType]) Name(name string) *wsRouteBuilder[InitRequestType, RequestType, ResponseType] {
	return &wsRouteBuilder[InitRequestType, RequestType, ResponseType]{module: b.module, name: name}
}

type wsRouteBuilder[InitRequestType any, RequestType any, ResponseType any] struct {
	module types.Module
	name   string
}

// Route sets the route path for the WebSocket handler
func (b *wsRouteBuilder[InitRequestType, RequestType, ResponseType]) Route(route string) *wsOptionals[InitRequestType, RequestType, ResponseType] {
	return &wsOptionals[InitRequestType, RequestType, ResponseType]{module: b.module, name: b.name, route: route}
}

type wsOptionals[InitRequestType any, RequestType any, ResponseType any] struct {
	module types.Module
	name   string
	route  string

	rateLimitDuration time.Duration
	rateLimitLimit    int
	middlewares       []types.MiddlewareHandler
}

func (b *wsOptionals[InitRequestType, RequestType, ResponseType]) RateLimit(duration time.Duration, limit int) *wsOptionals[InitRequestType, RequestType, ResponseType] {
	b.rateLimitDuration = duration
	b.rateLimitLimit = limit
	return b
}

// Middleware sets the middleware for the WebSocket handler
func (b *wsOptionals[InitRequestType, RequestType, ResponseType]) Middleware(middleware types.MiddlewareHandler) *wsOptionals[InitRequestType, RequestType, ResponseType] {
	b.middlewares = append(b.middlewares, middleware)
	return b
}

func (b *wsOptionals[InitRequestType, RequestType, ResponseType]) InitHandler(handler func(ctx types.WsContext[InitRequestType, RequestType, ResponseType])) *wsHandlerBuilder[InitRequestType, RequestType, ResponseType] {
	return &wsHandlerBuilder[InitRequestType, RequestType, ResponseType]{module: b.module, name: b.name, route: b.route, rateLimitDuration: b.rateLimitDuration, rateLimitLimit: b.rateLimitLimit, middlewares: b.middlewares, messageHandlers: make(map[string]func(ctx types.WsContext[InitRequestType, RequestType, ResponseType], action string, message *types.WsMessage[RequestType])), handler: handler}
}

type wsHandlerBuilder[InitRequestType any, RequestType any, ResponseType any] struct {
	module types.Module
	name   string
	route  string

	rateLimitDuration time.Duration
	rateLimitLimit    int
	middlewares       []types.MiddlewareHandler

	messageHandlers map[string]func(ctx types.WsContext[InitRequestType, RequestType, ResponseType], action string, message *types.WsMessage[RequestType])
	handler         func(ctx types.WsContext[InitRequestType, RequestType, ResponseType])

	periodicHandler  func(ctx types.WsContext[InitRequestType, RequestType, ResponseType])
	periodicInterval time.Duration
}

func (b *wsHandlerBuilder[InitRequestType, RequestType, ResponseType]) MessageHandler(action string, handler func(ctx types.WsContext[InitRequestType, RequestType, ResponseType], action string, message *types.WsMessage[RequestType])) *wsHandlerBuilder[InitRequestType, RequestType, ResponseType] {
	_, ok := b.messageHandlers[action]
	if ok {
		panic(fmt.Sprintf("message handler for action %s already registered for ws handler %s", action, b.name))
	}
	b.messageHandlers[action] = handler
	return b
}

func (b *wsHandlerBuilder[InitRequestType, RequestType, ResponseType]) PeriodicHandler(interval time.Duration, handler func(ctx types.WsContext[InitRequestType, RequestType, ResponseType])) *wsHandlerBuilder[InitRequestType, RequestType, ResponseType] {
	b.periodicHandler = handler
	b.periodicInterval = interval
	return b
}

func (b *wsHandlerBuilder[InitRequestType, RequestType, ResponseType]) Build() types.WsHandler {
	h := types.NewWsHandler(b.module, b.name, b.route, b.rateLimitDuration, b.rateLimitLimit, b.middlewares, b.messageHandlers, nil, b.handler, nil, b.periodicHandler, nil, b.periodicInterval)
	b.module.RegisterWsHandler(h)
	return h
}
