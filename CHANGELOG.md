# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.1] - 2025-02-15

### Added

- AI warning header on generated TypeScript files (general.ts, models.ts, api.ts) instructing AI to never edit them

## [1.0.0] - 2025-02-14

### Added

- App and Module architecture for organizing applications
- Type-safe REST API handlers (GET, POST, PUT, PATCH, DELETE) with generic request/response types
- WebSocket handlers with init, message, and periodic handlers
- Cron job support with cron expressions
- Init handlers for startup logic (migrations, seeding)
- Middleware support (JWT auth, rate limiting, caching)
- SQL database support via metaorm (MySQL, PostgreSQL, SQLite)
- MongoDB support via metamongo
- TypeScript type generation from Go structs (`--ts` flag)
- OpenAPI spec generation (`--openapi` flag)
- Graceful shutdown on SIGINT/SIGTERM
- Request ID, CORS, and structured logging middleware
