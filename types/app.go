package types

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/MetaDiv-AI/metadev/internal/openapi"
	"github.com/MetaDiv-AI/metadev/internal/typescript"
	"github.com/MetaDiv-AI/metamongo"
	"github.com/MetaDiv-AI/metaorm"
)

func NewApp(name string, database string, mongoDatabase string, migrations []Model) App {
	return &app{
		name:          name,
		database:      database,
		mongoDatabase: mongoDatabase,
		migrations:    migrations,
		modules:       make(map[string]Module),
	}
}

type App interface {
	// Name returns the name of the app
	Name() string
	// SQL returns the SQL database of the app
	SQL() metaorm.Database
	// Mongo returns the MongoDB database of the app
	Mongo() metamongo.Database
	// Migration returns the migrations of the app
	Migrations() []Model
	// Modules returns the modules of the app
	Modules() []Module
	// ApiHandlers returns the api handlers of the app
	ApiHandlers() []ApiHandler
	// CronHandlers returns the cron handlers of the app
	CronHandlers() []CronHandler
	// InitHandlers returns the init handlers of the app
	InitHandlers() []InitHandler
	// WsHandlers returns the ws handlers of the app
	WsHandlers() []WsHandler
	// RegisterModule registers a module to the app
	RegisterModule(module Module)
	// GenerateTypescript generates the typescript files for the app
	GenerateTypescript(folderPath string)
	// GenerateOpenAPISpec generates the openapi spec for the app
	GenerateOpenAPISpec(folderPath string)
	// RequiresDatabase returns true if the app requires a database connection
	RequiresDatabase() bool
}

type Model interface {
	// TableName returns the table name of the model
	TableName() string
}

type app struct {
	name          string
	database      string
	mongoDatabase string
	migrations    []Model
	modules       map[string]Module
}

func (a *app) Name() string {
	return a.name
}

func (a *app) SQL() metaorm.Database {
	if a.database == "" {
		return nil
	}
	return dbPool.GetConnection(a.database)
}

func (a *app) Mongo() metamongo.Database {
	if a.mongoDatabase == "" {
		return nil
	}
	return mongoPool.GetConnection(a.mongoDatabase)
}

func (a *app) Migrations() []Model {
	return a.migrations
}

func (a *app) Modules() []Module {
	keys := make([]string, 0, len(a.modules))
	for key := range a.modules {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	modules := make([]Module, 0, len(a.modules))
	for _, key := range keys {
		modules = append(modules, a.modules[key])
	}
	return modules
}

func (a *app) ApiHandlers() []ApiHandler {
	modules := a.Modules()
	handlers := make([]ApiHandler, 0)
	for _, module := range modules {
		handlers = append(handlers, module.ApiHandlers()...)
	}
	return handlers
}

func (a *app) CronHandlers() []CronHandler {
	modules := a.Modules()
	handlers := make([]CronHandler, 0)
	for _, module := range modules {
		handlers = append(handlers, module.CronHandlers()...)
	}
	return handlers
}

func (a *app) InitHandlers() []InitHandler {
	modules := a.Modules()
	handlers := make([]InitHandler, 0)
	for _, module := range modules {
		handlers = append(handlers, module.InitHandlers()...)
	}
	return handlers
}

func (a *app) WsHandlers() []WsHandler {
	modules := a.Modules()
	handlers := make([]WsHandler, 0)
	for _, module := range modules {
		handlers = append(handlers, module.WsHandlers()...)
	}
	return handlers
}

func (a *app) RequiresDatabase() bool {
	return a.database != ""
}

func (a *app) RegisterModule(module Module) {
	_, ok := a.modules[module.Name()]
	if ok {
		panic(fmt.Sprintf("module %s already registered for app %s", module.Name(), a.Name()))
	}
	a.modules[module.Name()] = module
}

func (a *app) GenerateTypescript(folderPath string) {

	builder := typescript.NewBuilder()

	var skipBuild = true
	for _, handler := range a.ApiHandlers() {
		if handler.SkipTypescript() {
			continue
		}
		skipBuild = false
		info := handler.TypescriptInfo()
		if info.Request != "" {
			builder.Typescriptify.AddType(info.RequestType)
		}
		if info.Response != "" {
			builder.Typescriptify.AddType(info.ResponseType)
		}
		builder.ApiInfos = append(builder.ApiInfos, info)
	}

	if skipBuild {
		return
	}

	os.MkdirAll(folderPath, os.ModePerm)
	builder.Build(a.Name(), folderPath)
}

func (a *app) GenerateOpenAPISpec(folderPath string) {
	builder := openapi.NewBuilder()

	var skipBuild = true
	for _, handler := range a.ApiHandlers() {
		if !handler.SkipOpenApi() {
			handler.RegisterOpenApi(builder)
			skipBuild = false
		}
	}

	if skipBuild {
		return
	}

	j := builder.ToJSON()
	os.MkdirAll(folderPath, os.ModePerm)
	os.WriteFile(filepath.Join(folderPath, "openapi-"+a.Name()+".json"), []byte(j), os.ModePerm)
}
