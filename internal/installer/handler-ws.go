package installer

import (
	"net/http"

	"github.com/MetaDiv-AI/metadev/types"

	"github.com/gin-gonic/gin"
)

func InstallWsHandler(e *gin.Engine, handler types.WsHandler) {
	handlers := make([]gin.HandlerFunc, 0)

	duration, limit := handler.RateLimit()
	if duration > 0 {
		handlers = append(handlers, newRateLimitHandler(duration, int64(limit)))
	}

	if len(handler.Middlewares()) > 0 {
		for _, middleware := range handler.Middlewares() {
			handlers = append(handlers, middlewareToGin(middleware))
		}
	}

	handlers = append(handlers, handler.GinHandler())
	e.Handle(http.MethodGet, handler.Route(), handlers...)
}
