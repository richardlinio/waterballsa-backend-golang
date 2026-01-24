# Reference Guide: ISA Feature Analysis & Implementation

This document provides quick-reference information for the analyze-isa-feature skill.

## Project Structure

```
waterballsa-backend-golang/
├── cmd/server/              # Application entry point
├── internal/
│   ├── app/                 # App lifecycle (init, shutdown)
│   ├── config/              # Environment config loading
│   ├── handler/             # HTTP handlers (presentation)
│   ├── service/             # Business logic
│   ├── repository/          # Data access layer
│   ├── dto/                 # Request/Response DTOs
│   ├── model/               # Domain models
│   ├── middleware/          # HTTP middleware
│   ├── router/              # Route registration
│   ├── validator/           # Custom validators
│   ├── util/                # Helper functions
│   ├── infrastructure/      # Infrastructure setup
│   │   ├── database/        # pgxpool connection
│   │   ├── logger/          # slog logger
│   │   └── server/          # HTTP server config
│   └── db/                  # SQLc generated code
│       ├── queries/         # SQL query definitions
│       └── *.sql.go         # Generated Go code
├── migrations/              # Goose database migrations
├── tests/
│   ├── bdd/                 # BDD test infrastructure
│   │   ├── features/        # Gherkin .feature files
│   │   │   └── isa/         # ISA layer tests
│   │   └── steps/           # Step definitions
│   │       ├── database/    # DB setup steps
│   │       └── http/        # HTTP request/response steps
│   └── testutil/            # Test utilities
└── docs/
    ├── api-docs/            # OpenAPI specifications
    │   ├── swagger.yaml     # Main API doc
    │   └── openapi/
    │       ├── paths/       # Endpoint definitions
    │       └── schemas/     # Request/Response schemas
    └── db-schema.dbml       # Database schema documentation
```

## Layer Architecture Pattern

### 1. Handler Layer (Presentation)

**Location:** `internal/handler/*_handler.go`

**Responsibilities:**

- Extract path/query parameters
- Bind JSON request bodies
- Validate user authorization
- Call service layer
- Convert domain models to DTOs
- Return HTTP responses
- Pass errors to middleware

**Pattern:**

```go
type FooHandler struct {
    service *service.FooService
    logger  *slog.Logger
    timeout time.Duration
}

func NewFooHandler(service *service.FooService, logger *slog.Logger, timeout time.Duration) *FooHandler {
    return &FooHandler{
        service: service,
        logger:  logger,
        timeout: timeout,
    }
}

func (h *FooHandler) GetFoo(c *gin.Context) {
    // 1. Extract params
    id := c.Param("id")

    // 2. Authorization
    user := c.MustGet("JWT_PAYLOAD").(*model.User)
    if user.ID != requestedUserID {
        _ = c.Error(apperror.UnauthorizedAccess())
        return
    }

    // 3. Create context with timeout
    ctx, cancel := context.WithTimeout(c.Request.Context(), h.timeout)
    defer cancel()

    // 4. Call service
    result, err := h.service.GetFoo(ctx, id)
    if err != nil {
        _ = c.Error(err)
        return
    }

    // 5. Convert to DTO and return
    response := dto.ToFooResponse(result)
    c.JSON(http.StatusOK, response)
}
```

### 2. Service Layer (Business Logic)

**Location:** `internal/service/*_service.go`

**Responsibilities:**

- Validate business rules
- Coordinate multiple repository calls
- Calculate derived values
- Enforce state transitions
- Wrap repository errors with AppError

**Pattern:**

```go
type FooService struct {
    repository *repository.FooRepository
    logger     *slog.Logger
    timeout    time.Duration
}

func NewFooService(repository *repository.FooRepository, logger *slog.Logger, timeout time.Duration) *FooService {
    return &FooService{
        repository: repository,
        logger:     logger,
        timeout:    timeout,
    }
}

func (s *FooService) GetFoo(ctx context.Context, id int64) (*model.Foo, error) {
    // 1. Validate input
    if id <= 0 {
        return nil, apperror.InvalidInput()
    }

    // 2. Call repository
    foo, err := s.repository.GetByID(ctx, id)
    if err != nil {
        return nil, apperror.DatabaseError(err)
    }

    // 3. Business logic / calculations
    foo.CalculatedField = computeSomething(foo)

    return foo, nil
}
```

### 3. Repository Layer (Data Access)

**Location:** `internal/repository/*_repository.go`

**Responsibilities:**

- Execute SQLc queries
- Convert sqlc models to domain models
- Handle database errors
- Manage transactions

**Pattern:**

```go
type FooRepository struct {
    pool *pgxpool.Pool
}

func NewFooRepository(pool *pgxpool.Pool) *FooRepository {
    return &FooRepository{pool: pool}
}

func (r *FooRepository) GetByID(ctx context.Context, id int64) (*model.Foo, error) {
    queries := db.New(r.pool)

    dbFoo, err := queries.GetFooByID(ctx, id)
    if err != nil {
        if errors.Is(err, pgx.ErrNoRows) {
            return nil, nil // No error, just no record
        }
        return nil, err
    }

    // Convert to domain model
    return &model.Foo{
        ID:   dbFoo.ID,
        Name: dbFoo.Name,
        // ... other fields
    }, nil
}
```

### 4. Data Layer (SQLc Queries)

**Location:** `internal/db/queries/*.sql`

**Pattern:**

```sql
-- name: GetFooByID :one
SELECT id, name, created_at, updated_at
FROM foos
WHERE id = $1 AND deleted_at IS NULL;

-- name: CreateFoo :one
INSERT INTO foos (name, created_at, updated_at)
VALUES ($1, NOW(), NOW())
RETURNING id, name, created_at, updated_at;

-- name: UpdateFoo :one
UPDATE foos
SET name = $1, updated_at = NOW()
WHERE id = $2 AND deleted_at IS NULL
RETURNING id, name, created_at, updated_at;
```

**After changes:** Run `make sqlc` to generate Go code.

## Common Patterns

### Authorization Check

```go
// In handler
user := c.MustGet("JWT_PAYLOAD").(*model.User)
requestedUserID, _ := strconv.ParseInt(c.Param("userId"), 10, 64)

if user.ID != requestedUserID {
    _ = c.Error(apperror.UnauthorizedAccess())
    return
}
```

### Request Binding with Validation

```go
var req dto.FooRequest
if err := c.ShouldBindJSON(&req); err != nil {
    _ = c.Error(apperror.NewWithError(apperror.CodeValidationFailed, err))
    return
}
```

### Error Wrapping

```go
// In service layer
result, err := s.repository.GetSomething(ctx, id)
if err != nil {
    return nil, apperror.DatabaseError(err) // Wrap with AppError
}
```

### Handling Not Found

```go
// In repository - return nil, nil (not an error)
if errors.Is(err, pgx.ErrNoRows) {
    return nil, nil
}

// In service - decide how to handle
result, err := s.repository.GetByID(ctx, id)
if err != nil {
    return nil, apperror.DatabaseError(err)
}
if result == nil {
    return nil, apperror.FooNotFound() // Or return default
}
```

### UPSERT Pattern

```sql
-- name: UpsertFoo :one
INSERT INTO foos (user_id, status, value)
VALUES ($1, $2, $3)
ON CONFLICT (user_id)
DO UPDATE SET
  status = EXCLUDED.status,
  value = EXCLUDED.value,
  updated_at = NOW()
RETURNING id, user_id, status, value, created_at, updated_at;
```

## Error Handling System

### Adding New Error Codes

**1. Define code constant** (`internal/apperror/codes.go`):

```go
const (
    CodeFooNotFound = "ERR_FOO_NOT_FOUND"
    CodeInvalidFoo  = "ERR_INVALID_FOO"
)
```

**2. Add message** (`internal/apperror/messages.go`):

```go
var errorMessages = map[string]string{
    CodeFooNotFound: "找不到指定的資源",
    CodeInvalidFoo:  "無效的輸入資料",
}
```

**3. Map HTTP status** (`internal/apperror/status.go`):

```go
var httpStatusMap = map[string]int{
    CodeFooNotFound: http.StatusNotFound,      // 404
    CodeInvalidFoo:  http.StatusBadRequest,    // 400
}
```

**4. Create constructor** (`internal/apperror/error.go`):

```go
func FooNotFound() *AppError {
    return New(CodeFooNotFound)
}

func InvalidFoo() *AppError {
    return New(CodeInvalidFoo)
}
```

**5. Use in code**:

```go
if foo == nil {
    return apperror.FooNotFound()
}

if input < 0 {
    return apperror.InvalidFoo()
}
```

## DTO Patterns

### Request DTO

```go
type CreateFooRequest struct {
    Name  string `json:"name" binding:"required,min=3,max=100"`
    Value int    `json:"value" binding:"required,min=0"`
}
```

### Response DTO

```go
type FooResponse struct {
    ID        int64  `json:"id"`
    Name      string `json:"name"`
    Value     int    `json:"value"`
    CreatedAt int64  `json:"createdAt"` // Unix timestamp in milliseconds
}
```

### Conversion Function

```go
func ToFooResponse(foo *model.Foo) *FooResponse {
    return &FooResponse{
        ID:        foo.ID,
        Name:      foo.Name,
        Value:     foo.Value,
        CreatedAt: foo.CreatedAt.UnixMilli(),
    }
}
```

## Route Registration

**Location:** `internal/router/router.go`

```go
func Setup(
    engine *gin.Engine,
    cfg config.Config,
    authHandler *handler.AuthHandler,
    fooHandler *handler.FooHandler, // Add handler parameter
    pool *pgxpool.Pool,
    logger *slog.Logger,
) *gin.Engine {
    // ... middleware setup ...

    // Public routes
    public := engine.Group("/")
    {
        public.GET("/healthz", healthCheck(pool, logger))
        public.POST("/auth/login", authHandler.Login)
    }

    // Protected routes
    protected := engine.Group("/")
    protected.Use(
        middleware.JWTAuth(cfg.JWT),
        middleware.BlacklistChecker(pool),
        middleware.Authorize(pool),
    )
    {
        protected.GET("/foos/:id", fooHandler.GetFoo)
        protected.POST("/foos", fooHandler.CreateFoo)
    }

    return engine
}
```

## Dependency Wiring

**Location:** `internal/app/app.go`

```go
func (a *App) Run(ctx context.Context) error {
    // ... config and pool setup ...

    // Repositories
    fooRepository := repository.NewFooRepository(pool)

    // Services
    fooService := service.NewFooService(
        fooRepository,
        logger,
        cfg.Server.RequestTimeout,
    )

    // Handlers
    fooHandler := handler.NewFooHandler(
        fooService,
        logger,
        cfg.Server.RequestTimeout,
    )

    // Router
    engine := router.Setup(
        gin.New(),
        cfg,
        authHandler,
        fooHandler, // Pass handler
        pool,
        logger,
    )

    // ... server start ...
}
```

## ISA Feature File Structure

```gherkin
@isa
Feature: Feature Name

  Scenario: Scenario description
    # Setup: Database records
    Given the database has a journey:
      | title   | Java 基礎課程 |
      | slug    | java-basics  |

    And the database has a user:
      | username | Alice     |
      | password | Test1234! |

    # Authentication
    And I set request body to:
      """
      {
        "username": "Alice",
        "password": "Test1234!"
      }
      """
    When I send "POST" request to "/auth/login"
    And I store the response field "accessToken" as "accessToken"
    Given I set Authorization header to "{{accessToken}}"

    # Action: API call
    When I send "GET" request to "/users/1/foo"

    # Verification: Response
    Then the response status code should be 200
    And the response body should contain field "id"
    And the response body field "name" should equal string "Expected Value"
```

## Naming Conventions

### ✅ Correct (Full Names)

```go
progressRepository := repository.NewProgressRepository(pool)
authenticationService := service.NewAuthenticationService(repo)
userMissionProgress := &model.UserMissionProgress{}
```

### ❌ Incorrect (Abbreviations)

```go
progressRepo := repository.NewProgressRepository(pool)
authSvc := service.NewAuthenticationService(repo)
ump := &model.UserMissionProgress{}
```

## Configuration Pattern

### ✅ Correct (Pass-by-Value)

```go
func NewHandler(cfg config.ServerConfig) *Handler {
    return &Handler{config: cfg}
}

func loadServerConfig() (ServerConfig, error) {
    return ServerConfig{Port: 8080}, nil
}

type Config struct {
    Server ServerConfig // Value, not pointer
}
```

### ❌ Incorrect (Pass-by-Pointer)

```go
func NewHandler(cfg *config.ServerConfig) *Handler {
    return &Handler{config: cfg}
}

func loadServerConfig() (*ServerConfig, error) {
    return &ServerConfig{Port: 8080}, nil
}

type Config struct {
    Server *ServerConfig // Pointer
}
```

## Common Make Commands

```bash
# Format code
make fmt

# Run linter
make lint

# Build application
make build

# Run tests
make test

# Generate SQLc code
make sqlc

# Database migrations
make migrate-status
make migrate-up
make migrate-down
make migrate-reset
```

## Documentation References

- **API Spec**: `/docs/api-docs/swagger.yaml` and `/docs/api-docs/openapi/paths/*.yaml`
- **DB Schema**: `/docs/db-schema.dbml`
- **Project Guide**: `/CLAUDE.md`
- **ISA Features**: `/tests/bdd/features/isa/**/*.isa.feature`

## Quick Checklist for Implementation

- [ ] Read ISA feature file
- [ ] Read swagger.yaml for endpoint spec
- [ ] Read db-schema.dbml for tables
- [ ] Launch Explore agent to find existing code
- [ ] Create domain model (internal/model/)
- [ ] Create DTOs (internal/dto/)
- [ ] Write SQLc queries (internal/db/queries/)
- [ ] Run `make sqlc`
- [ ] Create repository (internal/repository/)
- [ ] Create service with business logic (internal/service/)
- [ ] Create handler (internal/handler/)
- [ ] Add error codes (internal/apperror/)
- [ ] Register routes (internal/router/router.go)
- [ ] Wire dependencies (internal/app/app.go)
- [ ] Run `make fmt && make lint && make build`
- [ ] Run BDD tests: `cd tests/bdd && go test -v -tags=isa ./...`
