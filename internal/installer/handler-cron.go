package installer

import (
	"github.com/MetaDiv-AI/metadev/types"

	"github.com/robfig/cron/v3"
)

func InstallCronHandler(e *cron.Cron, handler types.CronHandler) {
	logger := types.NewLogger(handler.Module().App().Name(), handler.Name())
	db := types.NewDbContext(handler.Module().App().SQL())
	mongoDb := types.NewMongoContext(handler.Module().App().Mongo())

	e.AddFunc(handler.Spec(), func() {
		handler.Handler()(db.DB(), mongoDb.MongoDB(), logger)
	})
}
