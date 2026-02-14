package metadev

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/MetaDiv-AI/logger"
	"github.com/MetaDiv-AI/metadev/internal/installer"
	"github.com/MetaDiv-AI/metadev/types"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/requestid"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Apps = make(map[string]types.App)

func NewEngine() Engine {
	return &engine{
		ginEngine: gin.Default(),
		cron:      cron.New(),
		zapLogger: zap.NewNop(),
	}
}

type Engine interface {
	Run()
	GinEngine() *gin.Engine
}

type engine struct {
	ginEngine *gin.Engine
	cron      *cron.Cron
	zapLogger *zap.Logger
}

func (e *engine) GinEngine() *gin.Engine {
	return e.ginEngine
}

func (e *engine) Run() {
	// Create zap logger for middleware
	var zapLogger *zap.Logger
	var err error
	mode := gin.Mode()
	if mode == gin.DebugMode {
		config := zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		zapLogger, err = config.Build()
	} else {
		config := zap.NewProductionConfig()
		zapLogger, err = config.Build()
	}
	if err != nil {
		zapLogger, _ = zap.NewProduction()
	}
	e.zapLogger = zapLogger
	defer e.zapLogger.Sync()

	// Add request ID middleware
	e.ginEngine.Use(requestid.New())

	// Add zap logging middleware
	e.ginEngine.Use(ginzap.Ginzap(e.zapLogger, "2006-01-02T15:04:05Z07:00", true))

	// Add CORS middleware
	e.ginEngine.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowHeaders: []string{"Origin", "Authorization", "Content-Type", "X-Locale"},
	}))

	fmt.Println("found apps: ", len(Apps))

	skipRun := false
	// If there is args --ts folder_path, build typescript
	if hasFlag, folderPath := hasFlag("--ts"); hasFlag {
		buildTypescript(folderPath)
		skipRun = true
	}

	// If there is args --openapi folder_path, build openapi
	if hasFlag, folderPath := hasFlag("--openapi"); hasFlag {
		buildOpenApi(folderPath)
		skipRun = true
	}

	if skipRun {
		fmt.Println("skipped run")
		return
	}

	for _, app := range Apps {
		func() {
			defer func() {
				if r := recover(); r != nil {
					fmt.Printf("Warning: Failed to initialize database for app %s: %v\n", app.Name(), r)
				}
			}()
			database := app.SQL()
			if database == nil {
				return
			}
			db := types.NewDbContext(database)
			migrates := app.Migrations()
			for _, migrate := range migrates {
				db.DB().AutoMigrate(migrate)
			}
		}()
	}

	for _, app := range Apps {
		// Check if app requires database but connection is nil
		if app.SQL() == nil && app.Name() != "system" {
			log := logger.New().Build()
			log.Warn("Skipping handler registration for app: database connection is nil",
				zap.String("app", app.Name()),
			)
			continue
		}

		for _, handler := range app.CronHandlers() {
			installer.InstallCronHandler(e.cron, handler)
		}
		for _, handler := range app.ApiHandlers() {
			installer.InstallApiHandler(e.ginEngine, handler)
		}
		for _, handler := range app.WsHandlers() {
			installer.InstallWsHandler(e.ginEngine, handler)
		}
	}

	for _, app := range Apps {
		// Skip non-system apps if database is not configured
		if app.Name() != "system" {
			// Check if database is configured using the same logic as main.go
			if app.SQL() == nil {
				log := logger.New().Build()
				log.Warn("Skipping init handler registration for app: database connection is nil",
					zap.String("app", app.Name()),
				)
				continue
			}
			// Check if app requires database but connection is nil
			if app.RequiresDatabase() && app.SQL() == nil {
				log := logger.New().Build()
				log.Warn("Skipping init handler registration for app: database connection is nil",
					zap.String("app", app.Name()),
				)
				continue
			}
		}

		for _, handler := range app.InitHandlers() {
			installer.InstallInitHandler(e.ginEngine, handler)
		}
	}

	go e.cron.Start()

	// Set up graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("Shutting down gracefully...")
		e.cron.Stop()
		types.CloseAllConnections()
		types.CloseAllMongoConnections()
		os.Exit(0)
	}()

	host := os.Getenv("GIN_HOST")
	if host == "" {
		host = "127.0.0.1"
	}
	port := os.Getenv("GIN_PORT")
	if port == "" {
		port = "5000"
	}

	e.ginEngine.Run(host + ":" + port)
}

func buildTypescript(folderPath string) {
	if folderPath == "" {
		folderPath = "./_typescript"
	}
	for _, app := range Apps {
		app.GenerateTypescript(folderPath)
	}
}

func buildOpenApi(folderPath string) {
	if folderPath == "" {
		folderPath = "./_openapi"
	}
	for _, app := range Apps {
		app.GenerateOpenAPISpec(folderPath)
	}
}

func hasFlag(flag string) (bool, string) {
	for i, arg := range os.Args {
		if arg == flag {
			// Check if there's a value after the flag
			if i+1 < len(os.Args) && !strings.HasPrefix(os.Args[i+1], "-") {
				return true, os.Args[i+1]
			}
			return true, ""
		}
	}
	return false, ""
}
