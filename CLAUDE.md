!important: DO NOT USE SUBAGENT

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

# Run tests
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

### Build

```bash
# Build the application (inside Docker)
docker exec backend go build -o ./tmp/main ./cmd/server

# Build outside Docker (CI environment)
go build -o /dev/null ./...
```

### Running Tests

```bash
# Run all tests (inside Docker)
make test
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

### Application Lifecycle

1. **Initialization** ([cmd/server/main.go](cmd/server/main.go)):

   - Delegates to `app.New()` for setup
   - Calls `application.Run()` to start

2. **Setup** ([internal/app/app.go](internal/app/app.go)):

   - Loads configuration from environment variables
   - Initializes logger (slog)
   - Creates database connection pool (pgxpool)
   - Sets up Gin router with routes
   - Creates HTTP server

3. **Running**:

   - HTTP server runs in errgroup goroutine
   - Signal handler runs in errgroup goroutine
   - Waits for shutdown signal or error

4. **Graceful Shutdown**:
   - Triggered by SIGINT/SIGTERM or application error
   - Stops HTTP server with configurable timeout
   - Closes database connection pool
   - Uses `errgroup` for coordinated shutdown

### Configuration Management

Configuration is loaded exclusively from environment variables (see [.env.example](.env.example)):

- **Server Config**: Port, host, timeouts (read, write, shutdown, request)
- **Database Config**: Host, port, user, password, dbname, connection pool settings
- **Logger Config**: Level (debug/info/warn/error), format (json/text)
- **JWT Config**: Secret key (planned feature)

All config structs are in `internal/config/` with dedicated files:

- [config.go](internal/config/config.go) - Main config struct and loader
- [server.go](internal/config/server.go) - Server config
- [database.go](internal/config/database.go) - Database config
- [logger.go](internal/config/logger.go) - Logger config

### Database Layer

- **Driver**: pgx/v5 (not using an ORM)
- **Connection Pool**: pgxpool for efficient connection management
- **Migrations**: Goose with SQL files in `migrations/`
- **Planned**: sqlc for type-safe query generation (see [docs/auth-module-implementation.md](docs/auth-module-implementation.md))

#### Migration Files

Migrations use Goose format with `-- +goose Up` and `-- +goose Down` directives:

- Filenames: `NNN_description.sql` (e.g., `001_create_users_table.sql`)
- Location: `migrations/` directory
- Schema includes: users, access_tokens (JWT blacklist), courses module tables

### Router and Handlers

Routes are registered in [internal/router/router.go](internal/router/router.go):

- Takes Gin engine, database pool, logger, and timeout as parameters
- Handlers are initialized with dependencies
- Currently implements: `/healthz` endpoint

Handler pattern:

- Each handler is a struct with dependencies (pool, logger, timeout)
- Constructor function `NewXxxHandler()` for dependency injection
- Methods are Gin handler functions with signature `func(c *gin.Context)`

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
4. If using sqlc (planned):
   - Write SQL queries in `internal/db/queries/*.sql`
   - Run `sqlc generate` to create type-safe Go code

### Code Quality

- **Formatting**: Use `make fmt` (gofumpt) before commits
- **Linting**: Use `make lint` (golangci-lint) - CI enforces this
- **Testing**: Add tests for new code - CI runs with race detector
- **Config**: `.golangci.yml` defines linting rules

When implementing these features, follow the layer-based architecture and maintain consistency with existing patterns.
