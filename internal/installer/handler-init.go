package installer

import (
	"github.com/MetaDiv-AI/metadev/types"

	"github.com/gin-gonic/gin"
)

func InstallInitHandler(e *gin.Engine, handler types.InitHandler) {
	logger := types.NewLogger(handler.Module().App().Name(), handler.Name())
	db := types.NewDbContext(handler.Module().App().SQL())
	mongoDb := types.NewMongoContext(handler.Module().App().Mongo())

	handler.Handler()(db.DB(), mongoDb.MongoDB(), logger)
}
