package types

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func NewWsHandler[InitRequestType any, RequestType any, ResponseType any](
	module Module, name string, route string,
	rateLimitDuration time.Duration, rateLimitLimit int,
	middlewares []MiddlewareHandler,
	messageHandlers map[string]func(ctx WsContext[InitRequestType, RequestType, ResponseType], action string, message *WsMessage[RequestType]),
	publicMessageHandlers map[string]func(ctx PublicWsContext[InitRequestType, RequestType, ResponseType], action string, message *WsMessage[RequestType]),
	handler func(ctx WsContext[InitRequestType, RequestType, ResponseType]),
	publicHandler func(ctx PublicWsContext[InitRequestType, RequestType, ResponseType]),
	periodicHandler func(ctx WsContext[InitRequestType, RequestType, ResponseType]),
	publicPeriodicHandler func(ctx PublicWsContext[InitRequestType, RequestType, ResponseType]),
	periodicInterval time.Duration,
) WsHandler {
	return &wsHandler[InitRequestType, RequestType, ResponseType]{module: module, name: name, route: route, rateLimitDuration: rateLimitDuration, rateLimitLimit: rateLimitLimit, middlewares: middlewares, messageHandlers: messageHandlers, publicMessageHandlers: publicMessageHandlers, handler: handler, publicHandler: publicHandler, periodicHandler: periodicHandler, publicPeriodicHandler: publicPeriodicHandler, periodicInterval: periodicInterval}
}

type WsHandler interface {
	// Name returns the name of the ws handler
	Name() string
	// Module returns the module of the ws handler
	Module() Module
	// Route returns the route of the ws handler
	Route() string
	// RateLimit returns the rate limit of the ws handler
	RateLimit() (duration time.Duration, limit int)
	// Middlewares returns the middlewares of the ws handler
	Middlewares() []MiddlewareHandler
	// GinHandler returns the gin handler of the ws handler
	GinHandler() gin.HandlerFunc
}

type wsHandler[InitRequestType any, RequestType any, ResponseType any] struct {
	module Module
	name   string
	route  string

	rateLimitDuration time.Duration
	rateLimitLimit    int
	middlewares       []MiddlewareHandler

	messageHandlers       map[string]func(ctx WsContext[InitRequestType, RequestType, ResponseType], action string, message *WsMessage[RequestType])
	publicMessageHandlers map[string]func(ctx PublicWsContext[InitRequestType, RequestType, ResponseType], action string, message *WsMessage[RequestType])

	handler       func(ctx WsContext[InitRequestType, RequestType, ResponseType])
	publicHandler func(ctx PublicWsContext[InitRequestType, RequestType, ResponseType])

	periodicHandler       func(ctx WsContext[InitRequestType, RequestType, ResponseType])
	publicPeriodicHandler func(ctx PublicWsContext[InitRequestType, RequestType, ResponseType])
	periodicInterval      time.Duration
}

func (h *wsHandler[InitRequestType, RequestType, ResponseType]) Name() string {
	return h.name
}

func (h *wsHandler[InitRequestType, RequestType, ResponseType]) Module() Module {
	return h.module
}

func (h *wsHandler[InitRequestType, RequestType, ResponseType]) Route() string {
	route := ""
	if h.publicHandler != nil {
		route += "/public"
	}
	route += "/ws"
	route += "/" + h.module.App().Name()
	route += "/" + strings.TrimPrefix(h.route, "/")
	return route
}

func (h *wsHandler[InitRequestType, RequestType, ResponseType]) RateLimit() (duration time.Duration, limit int) {
	return h.rateLimitDuration, h.rateLimitLimit
}

func (h *wsHandler[InitRequestType, RequestType, ResponseType]) Middlewares() []MiddlewareHandler {
	return h.middlewares
}

func (h *wsHandler[InitRequestType, RequestType, ResponseType]) GinHandler() gin.HandlerFunc {
	wsUpgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	return func(ctx *gin.Context) {

		ws, err := wsUpgrader.Upgrade(ctx.Writer, ctx.Request, nil)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error":   "WebSocket upgrade failed",
			})
			return
		}

		// Create a mutex to protect WebSocket writes
		wsWriteMutex := &sync.Mutex{}

		var wsCtx WsContext[InitRequestType, RequestType, ResponseType]
		var publicWsCtx PublicWsContext[InitRequestType, RequestType, ResponseType]
		if h.publicHandler != nil {
			publicWsCtx, err = NewPublicWsContext[InitRequestType, RequestType, ResponseType](ctx, ws, h, wsWriteMutex)
			if err != nil {
				// Create a proper error response for validation errors
				responseCtx := NewResponseContext[ResponseType]()
				responseCtx.Error(err)
				responseCtx.MakeGinResponse(ctx)
				ctx.Abort()
				return
			}
			h.publicHandler(publicWsCtx)
			if publicWsCtx.MakeGinResponse(ctx) {
				return
			}
		} else {
			wsCtx, err = NewWsContext[InitRequestType, RequestType, ResponseType](ctx, ws, h, wsWriteMutex)
			if err != nil {
				// Create a proper error response for validation errors
				responseCtx := NewResponseContext[ResponseType]()
				responseCtx.Error(err)
				responseCtx.MakeGinResponse(ctx)
				ctx.Abort()
				return
			}
			h.handler(wsCtx)
		}

		pingTimeout := 2 * time.Minute   // 2 minutes timeout
		pingInterval := 30 * time.Second // Send ping every 30 seconds
		lastPing := time.Now()

		ws.SetReadDeadline(time.Now().Add(pingTimeout))

		// Set up ping handler
		ws.SetPingHandler(func(appData string) error {
			lastPing = time.Now()
			ws.SetReadDeadline(time.Now().Add(pingTimeout))
			wsWriteMutex.Lock()
			defer wsWriteMutex.Unlock()
			return ws.WriteMessage(websocket.PongMessage, []byte(appData))
		})

		// Set up pong handler
		ws.SetPongHandler(func(appData string) error {
			lastPing = time.Now()
			ws.SetReadDeadline(time.Now().Add(pingTimeout))
			return nil
		})
		// Fix the goroutine leak by using a context for cleanup
		cancelCtx, cancel := context.WithCancel(context.Background())
		defer cancel() // This will stop the monitoring goroutine

		// Start a goroutine to send periodic pings
		go func() {
			ticker := time.NewTicker(pingInterval)
			defer ticker.Stop()

			for {
				select {
				case <-ticker.C:
					// Send ping message
					pingMessage := map[string]interface{}{
						"action": "ping",
						"time":   time.Now().Format(time.RFC3339),
					}
					if pingData, err := json.Marshal(pingMessage); err == nil {
						wsWriteMutex.Lock()
						ws.WriteMessage(websocket.TextMessage, pingData)
						wsWriteMutex.Unlock()
					}
				case <-cancelCtx.Done():
					return
				}
			}
		}()

		// Start a goroutine to monitor ping timeout
		go func() {
			ticker := time.NewTicker(1 * time.Second)
			defer ticker.Stop()

			for {
				select {
				case <-ticker.C:
					if time.Since(lastPing) > pingTimeout {
						ws.Close()
						return
					}
				case <-cancelCtx.Done():
					return
				}
			}
		}()

		// Start periodic handler if configured
		if h.periodicInterval > 0 {
			go func() {
				ticker := time.NewTicker(h.periodicInterval)
				defer ticker.Stop()

				for {
					select {
					case <-ticker.C:
						if h.publicPeriodicHandler != nil {
							h.publicPeriodicHandler(publicWsCtx)
						} else if h.periodicHandler != nil {
							h.periodicHandler(wsCtx)
						}
					case <-cancelCtx.Done():
						return
					}
				}
			}()
		}

		// Monitor incoming messages
		for {
			messageType, messageData, err := ws.ReadMessage()
			if err != nil {
				break
			}

			// Handle ping messages
			if messageType == websocket.PingMessage {
				lastPing = time.Now()
				ws.SetReadDeadline(time.Now().Add(pingTimeout))
				wsWriteMutex.Lock()
				ws.WriteMessage(websocket.PongMessage, messageData)
				wsWriteMutex.Unlock()
				continue
			}

			// Handle text messages
			if messageType == websocket.TextMessage {
				// Check if message is a ping action
				var quickCheck map[string]interface{}
				if err := json.Unmarshal(messageData, &quickCheck); err == nil {
					if action, exists := quickCheck["action"]; exists && action == "ping" {
						lastPing = time.Now()
						ws.SetReadDeadline(time.Now().Add(pingTimeout))

						// Send pong response
						pongResponse := map[string]interface{}{
							"action": "pong",
							"time":   time.Now().Format(time.RFC3339),
						}
						if pongData, err := json.Marshal(pongResponse); err == nil {
							wsWriteMutex.Lock()
							ws.WriteMessage(websocket.TextMessage, pongData)
							wsWriteMutex.Unlock()
						}
						continue
					}
				}

				// Parse as WsMessageDTO for non-ping messages
				var wsMessage WsMessage[RequestType]
				if err := json.Unmarshal(messageData, &wsMessage); err != nil {
					continue
				}

				// Route to message handler by action name
				if handler, exists := h.messageHandlers[wsMessage.Action]; exists {
					handler(wsCtx, wsMessage.Action, &wsMessage)
				}

				if handler, exists := h.publicMessageHandlers[wsMessage.Action]; exists {
					handler(publicWsCtx, wsMessage.Action, &wsMessage)
				}
			}
		}
		ws.Close()
	}
}
