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

**Configuration Patterns:** For detailed config design patterns and best practices, refer to the `/go-idiomatic` skill.

### Layer Responsibilities

The codebase follows clean separation of concerns across different layers:

- **Handler Layer**: HTTP request/response handling, calls services
- **Service Layer**: Business logic, orchestrates repositories
- **Repository Layer**: Data access, database operations
- **DTO Layer**: API request/response structures
- **Model Layer**: Domain entities and business data structures

For detailed information about:

- Aggregate structures design (Model, Service, DTO layers)
- Layer dependency rules and anti-patterns
- Return value patterns
- Service/Repository error handling patterns

Refer to the `/go-idiomatic` skill.

### Store Pattern (Transactions)

When business operations require **atomic multi-step database operations**, use the **Store pattern**:

- **Location**: `internal/store/`
- **Purpose**: Encapsulate transaction logic separate from business logic
- **Pattern**: Each transaction function has `{Action}TxParams`, `{Action}TxResult`, and uses `execTx()`
- **Error Handling**: Domain-specific errors in `internal/store/errors.go`
- **Usage**: Service receives `*store.Store` via DI, calls transaction methods

**Examples:**

- [store/auth.go](internal/store/auth.go) - Atomic token rotation
- [store/order_create.go](internal/store/order_create.go) - Order creation with race protection
- [store/progress.go](internal/store/progress.go) - Reward claiming with TOCTOU prevention

**When to use:**

- Store: TOCTOU prevention, atomic multi-table updates, complex validation requiring locks
- Repository: Single-table operations, read-only queries, simple CRUD

For detailed implementation patterns, refer to the `/scaffold-feature` skill's [store-transactions reference](/.claude/skills/scaffold-feature/references/store-transactions.md).

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

The application uses a **unified error handling system** with centralized middleware to ensure consistent error responses and logging.

**Package Structure:**

```
internal/apperror/
├── error.go      → AppError struct and constructor functions
├── codes.go      → Error code constants (ERR_*)
├── messages.go   → Chinese error messages
└── status.go     → HTTP status code mappings

internal/middleware/
└── error_handler.go → Centralized error handling and logging
```

**Key Usage Patterns:**

- Service returns `*apperror.AppError` for business logic errors
- Handler calls `c.Error(err)` and returns early
- Middleware automatically creates JSON response and logs errors
- DO NOT manually call `c.JSON()` for errors in handlers

**Adding New Error Types:** Use the `/add-error-handling` skill for detailed 5-step implementation guide.

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

- Use full names instead of abbreviations for better clarity
- Variable names should be descriptive (e.g., `authService` not `authSvc`)

For comprehensive Go idiomatic patterns, refer to the `/go-idiomatic` skill.

When implementing features, follow the layer-based architecture and maintain consistency with existing patterns.
