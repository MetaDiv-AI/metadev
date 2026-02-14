# metadev

A Go framework for building web applications with Gin, SQL (via metaorm), MongoDB (via metamongo), and type-safe handlers. Part of the MetaDiv-AI ecosystem.

## Features

- **App & Module architecture** — Organize your application into apps and modules with clear boundaries
- **Type-safe REST API handlers** — Generic request/response types with automatic validation
- **WebSocket support** — Full-duplex connections with typed message handlers
- **Cron jobs** — Scheduled tasks with cron expressions
- **Init handlers** — Run logic at startup (e.g., migrations, seeding)
- **Middleware** — JWT auth, rate limiting, caching, and custom middleware
- **Code generation** — TypeScript types and OpenAPI specs from your Go structs

## Requirements

- Go 1.24+
- MySQL, PostgreSQL, or SQLite (for SQL apps)
- MongoDB (optional, for apps using Mongo)

## Installation

```bash
go get github.com/MetaDiv-AI/metadev
```

## Quick Start

```go
package main

import (
	"github.com/MetaDiv-AI/metadev"
	"github.com/MetaDiv-AI/metadev/types"
)

func main() {
	app := metadev.NewApp("myapp").
		SQLMigrates(/* your GORM models */).
		Build()

	module := metadev.NewModule(app).Name("users")

	metadev.GET[GetUserRequest, GetUserResponse](module).
		Name("get_user").
		Route("/users/:id").
		Handler(func(ctx types.ApiContext[GetUserRequest, GetUserResponse]) {
			// ctx.DB(), ctx.Mongo(), ctx.Jwt(), ctx.Logger() available
			ctx.OK(GetUserResponse{}) // return your response
		})

	engine := metadev.NewEngine()
	engine.Run()
}
```

## Configuration

Create a `.env` file (see `.env.example`):

| Variable | Description | Default |
|----------|-------------|---------|
| `GIN_MODE` | `debug` or `release` | `release` |
| `GIN_HOST` | Server host | `127.0.0.1` |
| `GIN_PORT` | Server port | `5000` |
| `SQL_HOST` | MySQL/Postgres host | `localhost` |
| `SQL_PORT` | SQL port | `3306` |
| `SQL_USER` | SQL username | `root` |
| `SQL_PASSWORD` | SQL password | - |
| `SQL_DATABASE` | SQL database name | `operagent` |
| `MONGO_URI` | MongoDB connection URI | `mongodb://localhost:27017` |
| `MONGO_DATABASE` | MongoDB database name | `operagent` |

## Core Concepts

### App

An app is the top-level container. It owns database connections and modules.

```go
app := metadev.NewApp("myapp").
	SQLMigrates(User{}, Post{}).  // GORM models for auto-migration
	Build()
```

### Module

Modules group related handlers (API, WebSocket, cron, init) under a logical unit.

```go
module := metadev.NewModule(app).Name("users")
```

### API Handlers

Define REST endpoints with typed request and response:

```go
// Authenticated handler (requires JWT)
metadev.POST[CreateUserRequest, CreateUserResponse](module).
	Name("create_user").
	Route("/users").
	RateLimit(time.Minute, 10).
	Cache(5 * time.Minute).  // GET only
	Middleware(myMiddleware).
	Handler(func(ctx types.ApiContext[CreateUserRequest, CreateUserResponse]) {
		user := ctx.Request()
		ctx.DB().Create(&user)
		ctx.OK(user)
	})

// Public handler (no auth)
metadev.GET[HealthRequest, HealthResponse](module).
	Name("health").
	Route("/health").
	PublicHandler(func(ctx types.PublicApiContext[HealthRequest, HealthResponse]) {
		ctx.OK(HealthResponse{Status: "ok"})
	})
```

Supported methods: `GET`, `POST`, `PUT`, `PATCH`, `DELETE`.

### WebSocket Handlers

```go
metadev.NewWsHandler[InitRequest, MessageRequest, MessageResponse](module).
	Name("chat").
	Route("/ws/chat").
	RateLimit(time.Minute, 60).
	InitHandler(func(ctx types.WsContext[InitRequest, MessageRequest, MessageResponse]) {
		// Called when connection opens
	}).
	MessageHandler("send", func(ctx types.WsContext[InitRequest, MessageRequest, MessageResponse], action string, msg *types.WsMessage[MessageRequest]) {
		// Handle incoming messages with action "send"
		ctx.Send(action, response)
	}).
	PeriodicHandler(30*time.Second, func(ctx types.WsContext[InitRequest, MessageRequest, MessageResponse]) {
		// Periodic tick
	}).
	Build()
```

### Cron Handlers

```go
metadev.NewCronHandler(module).
	Name("cleanup").
	Spec("@daily").  // Cron expression
	Handler(func(db metaorm.Database, mongo metamongo.Database, logger types.Logger) {
		// Run cleanup logic
	})
```

### Init Handlers

Run once at startup, before the server accepts requests:

```go
metadev.InitFunc(module).
	Name("seed").
	Handler(func(db metaorm.Database, mongo metamongo.Database, logger types.Logger) {
		// Seed database, run migrations, etc.
	})
```

### Middleware

```go
authMiddleware := metadev.Middleware(module).
	Name("auth").
	Handler(func(ctx types.MiddlewareContext) {
		if ctx.Jwt() == nil {
			ctx.Gin().AbortWithStatus(401)
			return
		}
		ctx.Gin().Next()
	})

// Use on specific routes
metadev.GET[Request, Response](module).
	Name("protected").
	Route("/protected").
	Middleware(authMiddleware).
	Handler(handler)
```

## Code Generation

Generate TypeScript types and API client from your Go structs:

```bash
go run . --ts ./frontend/types
```

Generate OpenAPI specification:

```bash
go run . --openapi ./docs
```

## Context Types

Handlers receive a context that combines:

- **Logger** — Structured logging
- **DbContext** — SQL via GORM (`ctx.DB()`)
- **MongoContext** — MongoDB (`ctx.Mongo()`)
- **JwtContext** — JWT parsing and validation via `ctx.Jwt()` (API/WS with auth)
- **RequestContext** — Parsed request body/query
- **ResponseContext** — `ctx.OK()`, `ctx.Error()`, etc.
- **WsConnContext** — WebSocket send (`ctx.Send()`)

## Graceful Shutdown

The engine handles `SIGINT` and `SIGTERM`, stopping the cron scheduler and closing database connections before exit.

## License

See repository for license information.
