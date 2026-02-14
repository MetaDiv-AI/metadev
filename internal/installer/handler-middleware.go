package installer

import (
	"github.com/MetaDiv-AI/metadev/types"

	"github.com/gin-gonic/gin"
)

func middlewareToGin(middleware types.MiddlewareHandler) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		middleware.Handler()(types.NewMiddlewareContext(ctx, middleware))
	}
}
