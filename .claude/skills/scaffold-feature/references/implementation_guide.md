# Implementation Guide: Patterns & Pitfalls

This guide provides common Go implementation patterns for this project, along with anti-patterns to avoid. Following these conventions is crucial for maintaining code consistency and quality.

Each section presents the recommended approach (✅ **Do**) and contrasts it with common mistakes (❌ **Don't**).

## Slice Initialization Style

✅ **Do:** Use the `make` and `append` pattern for efficiency and clarity, especially when the final size is known.

```go
// Pre-allocate capacity when source length is known
items := make([]T, 0, len(source))
for _, item := range source {
    items = append(items, T{...})
}
```

Further details on idiomatic Go can be found in the `/go-idiomatic` skill.

## Return Value Design

✅ **Do:** If a Service method needs to return more than two data fields, create a dedicated `domain model` struct. This improves readability and makes the function signature stable.

```go
// Domain Model (internal/model/)
type MissionDeliveryResult struct {
    Message          string
    ExperienceGained int
    TotalExperience  int32
    CurrentLevel     int32
}

// Service Method
func (s *Service) DeliverMission(ctx context.Context, userID, missionID int64) (*model.MissionDeliveryResult, error) {
    // ... business logic
    return &model.MissionDeliveryResult{...}, nil
}
```

❌ **Don't:** Return numerous values from a single function. This is a code smell and makes the code harder to refactor and read.

```go
// ANTI-PATTERN: Avoid this in service layers
func (s *Service) DeliverMission(...) (string, int, int32, int32, error) {
    // ...
}
```

## Handling "Not Found" Errors

✅ **Do:** Maintain separation of concerns. The Repository layer should convert database-specific errors (like `pgx.ErrNoRows`) into a generic domain error. The Service layer should only check for this domain error.

```go
// Repository: Convert to a domain error, hiding pgx details.
if errors.Is(err, pgx.ErrNoRows) {
    return nil, repository.ErrResourceNotFound // A custom error defined in the repository package
}

// Service: Check the repository's custom error, not the pgx error.
if errors.Is(err, repository.ErrResourceNotFound) {
    return nil, apperror.ResourceNotFound() // Return a standard application error
}
```

❌ **Don't:** Let service layers depend on `pgx` or any other database-specific implementation details.

```go
// ANTI-PATTERN: Service checking pgx.ErrNoRows breaks architectural boundaries.
if errors.Is(err, pgx.ErrNoRows) {
    // ...
}
```

## Authorization Checks

✅ **Do:** Perform authorization checks in the handler to protect endpoints. Retrieve the user from the context and verify their permissions against the requested resource.

```go
// Handler
user := c.MustGet("JWT_PAYLOAD").(*model.User)
requestedUserID, _ := strconv.ParseInt(c.Param("userId"), 10, 64)

if user.ID != requestedUserID {
    _ = c.Error(apperror.UnauthorizedAccess())
    return
}
```

❌ **Don't:** Skip authorization checks. This is a critical security vulnerability.

## Variable Naming

✅ **Do:** Use descriptive, full-word variable names. This improves code readability.

```go
// Correct
var repository repository.UserRepository
var service service.UserService
var config config.Config
```

❌ **Don't:** Use abbreviations. While it saves a few characters, it makes the code harder to understand for new developers.

```go
// ANTI-PATTERN
var repo repository.UserRepository // -> repository
var svc service.UserService     // -> service
var cfg config.Config         // -> config is acceptable, but be consistent
```

## UPSERT Pattern in SQL

✅ **Do:** Use the `INSERT ... ON CONFLICT DO UPDATE` statement for creating or updating records in a single database round-trip.

```sql
-- name: UpsertProgress :one
INSERT INTO user_mission_progress (user_id, mission_id, status, watch_position_seconds)
VALUES ($1, $2, $3, $4)
ON CONFLICT (user_id, mission_id)
DO UPDATE SET
  status = EXCLUDED.status,
  watch_position_seconds = EXCLUDED.watch_position_seconds,
  updated_at = NOW()
RETURNING *;
```

## SQLc Naming Conventions

✅ **Do:** Follow consistent naming conventions for SQLc queries.

- `GetXxxByYyy`: Single record query
- `ListXxxByYyy`: Multiple record query
- `CreateXxx`: Create a new record
- `UpdateXxx`: Update an existing record
- `UpsertXxx`: Create or update a record
- `DeleteXxx`: Delete a record

Remember to run `make sqlc` after adding or modifying queries.

## Router Registration

✅ **Do:** Register new handlers in `internal/router/router.go`.

```go
func Setup(
    engine *gin.Engine,
    cfg config.Config,
    fooHandler *handler.FooHandler, // Add new handler as a parameter
    pool *pgxpool.Pool,
    logger *slog.Logger,
) *gin.Engine {
    // ...
    protected := engine.Group("/")
    protected.Use(middleware.JWTAuth(cfg.JWT))
    {
        protected.GET("/foo/:id", fooHandler.GetFoo) // Register route
    }
    return engine
}
```

❌ **Don't:** Forget to register the route. The endpoint will not be available otherwise.
