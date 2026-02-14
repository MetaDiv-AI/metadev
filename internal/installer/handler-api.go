package installer

import (
	"net/http"
	"time"

	"github.com/MetaDiv-AI/metadev/types"

	"github.com/gin-contrib/cache"
	"github.com/gin-contrib/cache/persistence"
	"github.com/gin-gonic/gin"
)

func InstallApiHandler(e *gin.Engine, handler types.ApiHandler) {

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

	if handler.Cache() > 0 && handler.Method() == http.MethodGet {
		handlers = append(handlers, cache.CachePage(persistence.NewInMemoryStore(time.Second), handler.Cache(), handler.GinHandler()))
	} else {
		handlers = append(handlers, handler.GinHandler())
	}

	e.Handle(handler.Method(), handler.Route(), handlers...)
}
