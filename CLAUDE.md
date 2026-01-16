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

**Adding New Config:**

When adding new configuration, maintain consistency with pass-by-value pattern:
- Loader functions return values, not pointers: `func loadXxxConfig() (XxxConfig, error)`
- Config struct stores values: `type Config struct { Xxx XxxConfig }`
- Functions receive values: `func NewHandler(cfg config.XxxConfig)`

This expresses config immutability and ensures consistent behavior across the codebase.

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

#### Architecture

**Error Flow:**

```
Handler → Service → Repository
         ↓ returns AppError
      Middleware intercepts and responds with JSON
```

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

#### Error Response Format

All error responses follow this JSON structure:

```json
{
	"code": "ERR_USERNAME_EXISTS",
	"error": "使用者名稱已存在",
	"details": {
		"Username": "長度至少需要 3 個字元"
	}
}
```

- `code`: Error code constant (e.g., `ERR_USERNAME_EXISTS`)
- `error`: User-friendly Chinese error message
- `details`: Optional field for validation error details (field-level errors)

#### How to Use in Code

**In Service Layer:**

Services should return `*apperror.AppError` for all business logic errors:

```go
// Return predefined errors
if exists {
    return 0, apperror.UsernameExists()
}

// Wrap underlying errors (database, internal errors)
if err != nil {
    return 0, apperror.DatabaseError(err)
}
```

**In Handler Layer:**

Handlers should **NOT** manually create JSON responses. Instead, pass errors to middleware using `c.Error()`:

```go
// For validation errors
if err := c.ShouldBindJSON(&req); err != nil {
    _ = c.Error(apperror.NewWithError(apperror.CodeValidationFailed, err))
    return
}

// For service errors (service already returns AppError)
if err := h.service.SomeMethod(ctx, req); err != nil {
    _ = c.Error(err)
    return
}
```

**Key Pattern:**

- ✅ Service returns `*apperror.AppError`
- ✅ Handler calls `c.Error(err)` and returns early
- ✅ Middleware automatically creates JSON response
- ❌ DO NOT manually call `c.JSON()` for errors in handlers

#### Adding New Error Types

Follow these 5 steps to add a new error type:

**Step 1:** Add error code constant in [internal/apperror/codes.go](internal/apperror/codes.go):

```go
const (
    CodeCourseNotFound = "ERR_COURSE_NOT_FOUND"
)
```

**Step 2:** Add error message in [internal/apperror/messages.go](internal/apperror/messages.go):

```go
var errorMessages = map[string]string{
    CodeCourseNotFound: "課程不存在",
}
```

**Step 3:** Map to HTTP status in [internal/apperror/status.go](internal/apperror/status.go):

```go
var httpStatusMap = map[string]int{
    CodeCourseNotFound: http.StatusNotFound,
}
```

**Step 4:** Create constructor function in [internal/apperror/error.go](internal/apperror/error.go):

```go
func CourseNotFound() *AppError {
    return New(CodeCourseNotFound)
}
```

**Step 5:** Use in service layer:

```go
if !exists {
    return apperror.CourseNotFound()
}
```

#### Logging

Error logging is **automatic** via the error handler middleware:

- **AppError with underlying error** (`Err` field set): Logs at ERROR level with underlying error details
- **AppError without underlying error**: Logs at WARN level with code, message, status, path, method
- **Unexpected errors** (non-AppError): Logs at ERROR level with full error details

**No manual logging needed in handlers** - the middleware handles all error logging automatically.

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
