package types

import (
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func NewApiContext[RequestType any, ResponseType any](gin *gin.Context, handler ApiHandler) (ApiContext[RequestType, ResponseType], error) {
	req, err := NewRequestContext[RequestType](gin)
	mongoCtx := NewMongoContext(handler.Module().App().Mongo())
	return &ApiContextImpl[RequestType, ResponseType]{
		logger:          *NewLogger(handler.Module().App().Name(), handler.Name()),
		dbContext:       *NewDbContext(handler.Module().App().SQL()),
		mongoContext:    *mongoCtx,
		jwtContext:      *NewJwtContext(gin),
		requestContext:  *req,
		responseContext: *NewResponseContext[ResponseType](),
	}, err
}

type ApiContext[RequestType any, ResponseType any] interface {
	Logger
	DbContext
	MongoContext
	JwtContext
	RequestContext[RequestType]
	ResponseContext[ResponseType]
}

type ApiContextImpl[RequestType any, ResponseType any] struct {
	logger
	dbContext
	mongoContext
	jwtContext
	requestContext[RequestType]
	responseContext[ResponseType]
}

func NewPublicApiContext[RequestType any, ResponseType any](gin *gin.Context, handler ApiHandler) (PublicApiContext[RequestType, ResponseType], error) {
	req, err := NewRequestContext[RequestType](gin)
	mongoCtx := NewMongoContext(handler.Module().App().Mongo())
	return &PublicApiContextImpl[RequestType, ResponseType]{
		logger:          *NewLogger(handler.Module().App().Name(), handler.Name()),
		dbContext:       *NewDbContext(handler.Module().App().SQL()),
		mongoContext:    *mongoCtx,
		ginContext:      *NewGinContext(gin),
		requestContext:  *req,
		responseContext: *NewResponseContext[ResponseType](),
	}, err
}

type PublicApiContext[RequestType any, ResponseType any] interface {
	Logger
	DbContext
	MongoContext
	GinContext
	RequestContext[RequestType]
	ResponseContext[ResponseType]
}

type PublicApiContextImpl[RequestType any, ResponseType any] struct {
	logger
	dbContext
	mongoContext
	ginContext
	requestContext[RequestType]
	responseContext[ResponseType]
}

func NewWsContext[InitRequestType any, RequestType any, ResponseType any](gin *gin.Context, conn *websocket.Conn, handler WsHandler, writeMutex *sync.Mutex) (WsContext[InitRequestType, RequestType, ResponseType], error) {
	req, err := NewRequestContext[InitRequestType](gin)
	mongoCtx := NewMongoContext(handler.Module().App().Mongo())
	return &WsContextImpl[InitRequestType, RequestType, ResponseType]{
		logger:         *NewLogger(handler.Module().App().Name(), handler.Name()),
		dbContext:      *NewDbContext(handler.Module().App().SQL()),
		mongoContext:   *mongoCtx,
		jwtContext:     *NewJwtContext(gin),
		wsConnContext:  *NewWsConnContext[ResponseType](conn, writeMutex),
		requestContext: *req,
	}, err
}

type WsContext[InitRequestType any, RequestType any, ResponseType any] interface {
	Logger
	DbContext
	MongoContext
	JwtContext
	WsConnContext[ResponseType]
	RequestContext[InitRequestType]
}

type WsContextImpl[InitRequestType any, RequestType any, ResponseType any] struct {
	logger
	dbContext
	mongoContext
	jwtContext
	wsConnContext[ResponseType]
	requestContext[InitRequestType]
}

func NewPublicWsContext[InitRequestType any, RequestType any, ResponseType any](gin *gin.Context, conn *websocket.Conn, handler WsHandler, writeMutex *sync.Mutex) (PublicWsContext[InitRequestType, RequestType, ResponseType], error) {
	req, err := NewRequestContext[InitRequestType](gin)
	mongoCtx := NewMongoContext(handler.Module().App().Mongo())
	return &PublicWsContextImpl[InitRequestType, RequestType, ResponseType]{
		logger:          *NewLogger(handler.Module().App().Name(), handler.Name()),
		dbContext:       *NewDbContext(handler.Module().App().SQL()),
		mongoContext:    *mongoCtx,
		ginContext:      *NewGinContext(gin),
		wsConnContext:   *NewWsConnContext[ResponseType](conn, writeMutex),
		requestContext:  *req,
		responseContext: *NewResponseContext[ResponseType](),
	}, err
}

type PublicWsContext[InitRequestType any, RequestType any, ResponseType any] interface {
	Logger
	DbContext
	MongoContext
	GinContext
	WsConnContext[ResponseType]
	RequestContext[InitRequestType]
	ResponseContext[ResponseType]
}

type PublicWsContextImpl[InitRequestType any, RequestType any, ResponseType any] struct {
	logger
	dbContext
	mongoContext
	ginContext
	wsConnContext[ResponseType]
	requestContext[InitRequestType]
	responseContext[ResponseType]
}

func NewMiddlewareContext(gin *gin.Context, handler MiddlewareHandler) MiddlewareContext {
	logger := NewLogger(handler.Module().App().Name(), handler.Name())
	mongoCtx := NewMongoContext(handler.Module().App().Mongo())
	return &MiddlewareContextImpl{
		logger:       *logger,
		dbContext:    *NewDbContext(handler.Module().App().SQL()),
		mongoContext: *mongoCtx,
		jwtContext:   *NewJwtContext(gin),
	}
}

type MiddlewareContext interface {
	Logger
	DbContext
	MongoContext
	JwtContext
}

type MiddlewareContextImpl struct {
	logger
	dbContext
	mongoContext
	jwtContext
}
