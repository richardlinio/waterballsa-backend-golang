# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Critical Constraints

**DO NOT USE SUBAGENT** - Direct tool usage is required for all operations.

## Project Overview

WaterBall SA Backend - A Golang REST API backend for an educational platform using Gin framework, PostgreSQL, and layer-based architecture.

- **Go Version**: 1.25.5
- **Main Framework**: Gin (HTTP server)
- **Database**: PostgreSQL with pgx/v5 driver
- **Migration Tool**: Goose
- **Development**: Docker-based with Air for hot reload

## Common Commands

All commands below run inside Docker containers. Ensure Docker is running and containers are up before executing.

### Development

```bash
# Format code
make fmt

# Run linter
make lint

# Build the application (inside Docker)
make build

# Run all tests (inside Docker)
make test
```

### Database Migrations

```bash
# Check migration status
make migrate-status

# Run pending migrations
make migrate-up

# Rollback last migration
make migrate-down

# Reset database (drop all + re-run migrations)
make migrate-reset
```

## Architecture

### Layer-Based Structure

The codebase follows a clean layer-based architecture with clear separation of concerns:

```
cmd/server/          → Application entry point (main.go)
internal/
├── app/             → Application lifecycle management (initialization, graceful shutdown)
├── config/          → Configuration loading from environment variables
├── handler/         → HTTP request handlers (presentation layer)
├── service/         → Business logic layer
├── repository/      → Data access layer (interacts with database)
├── dto/             → Data Transfer Objects (API request/response structures)
├── model/           → Domain models (database entities)
├── middleware/      → HTTP middleware (e.g., authentication, logging)
├── router/          → Route definitions and setup
├── validator/       → Custom validation logic
├── util/            → Utility functions and helpers
└── infrastructure/  → Infrastructure concerns
    ├── database/    → Database connection pooling (pgxpool)
    ├── logger/      → Structured logging (slog)
    └── server/      → HTTP server configuration
```

### Configuration Management

Configuration is loaded exclusively from environment variables (see [.env.example](.env.example)):

- **Server Config**: Port, host, timeouts (read, write, shutdown, request)
- **Database Config**: Host, port, user, password, dbname, connection pool settings
- **Logger Config**: Level (debug/info/warn/error), format (json/text)
- **CORS Config**: Allowed origins, credentials, max age
- **JWT Config**: Secret key (planned feature)

All config structs are in `internal/config/` with dedicated files:

- [config.go](internal/config/config.go) - Main config struct and loader
- [server.go](internal/config/server.go) - Server config
- [database.go](internal/config/database.go) - Database config
- [logger.go](internal/config/logger.go) - Logger config
- [cors.go](internal/config/cors.go) - CORS config

#### Migration Files

Migrations use Goose format with `-- +goose Up` and `-- +goose Down` directives:

- Filenames: `NNN_description.sql` (e.g., `001_create_users_table.sql`)
- Location: `migrations/` directory
- Schema includes: users, access_tokens (JWT blacklist), courses module tables

### Router and Handlers

Routes are registered in [internal/router/router.go](internal/router/router.go):

Handler pattern:

- Each handler is a struct with dependencies (pool, logger, timeout)
- Constructor function `NewXxxHandler()` for dependency injection
- Methods are Gin handler functions with signature `func(c *gin.Context)`

### Middleware

The application uses middleware in a specific order to ensure correct behavior:

**Adding New Middleware:**

1. Create middleware function in `internal/middleware/` returning `gin.HandlerFunc`
2. Add middleware configuration to `internal/config/` if needed
3. Register middleware in `router.Setup()` in the appropriate order
4. Update `.env.example` with any new environment variables

### Error Handling

- Application-level errors use custom error types (e.g., `ErrShutdownSignal`)
- Handlers should return appropriate HTTP status codes
- Structured logging for error context

### Testing

- Unit tests use `_test.go` suffix
- CI runs tests with race detector: `go test -v -race ./...`
- Config package has test coverage for environment variable parsing
- BDD test structure exists in `tests/features/` and `tests/steps/` (currently empty)

## Development Workflow

### Adding a New Feature

1. **Define DTOs** in `internal/dto/` for request/response structures
2. **Add domain models** in `internal/model/` if needed
3. **Create repository** in `internal/repository/` for data access
4. **Implement service** in `internal/service/` for business logic
5. **Create handler** in `internal/handler/` for HTTP layer
6. **Register routes** in `internal/router/router.go`
7. **Add middleware** in `internal/middleware/` if cross-cutting concerns needed

### Database Changes

1. **Create migration** in `migrations/NNN_description.sql`
2. **Run migration**: `make migrate-up`
3. **Verify status**: `make migrate-status`
4. Using sqlc:
   - Write SQL queries in `internal/db/queries/*.sql`
   - Run `make sqlc` to create type-safe Go code

### Code Quality

- **Formatting**: Use `make fmt` (gofumpt) before commits
- **Linting**: Use `make lint` (golangci-lint) - CI enforces this
- **Testing**: Add tests for new code - CI runs with race detector
- **Config**: `.golangci.yml` defines linting rules

### Code Style

**Naming Conventions:**

- **Use full names instead of abbreviations** for better clarity and consistency
- Variable names should be descriptive and avoid abbreviations

Examples:

```go
// Preferred ✓
authService
userRepository
courseRepository

// Avoid ✗
authSvc
userRepo
courseRepo
```

When implementing these features, follow the layer-based architecture and maintain consistency with existing patterns.
