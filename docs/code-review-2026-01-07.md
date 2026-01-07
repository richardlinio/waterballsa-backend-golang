# Code Review Report - 2026-01-07

## 專案概述

WaterBall SA Backend - 使用 Gin framework、PostgreSQL 和分層架構的線上課程平台後端系統。

**審查範圍：**

- 專案架構設計
- Go/Gin 社群最佳實踐符合度
- 分層邏輯與實作
- 目前已實作的 Registry API

**審查時間：** 2026-01-07
**審查版本：** 初始架構（Registry API 已實作）

---

## ✅ 優點

### 1. 分層架構清晰且符合 Go 慣例

專案採用標準的分層架構：

```
Handler (HTTP層) → Service (業務邏輯) → Repository (資料存取) → DB (sqlc生成的type-safe queries)
```

這符合 Go 社群的主流做法，也是 Gin 官方推薦的架構模式。

**相關檔案：**

- `internal/handler/auth.go` - HTTP 層處理請求和回應
- `internal/service/auth.go` - 業務邏輯層
- `internal/repository/user.go` - 資料存取層
- `internal/db/users.sql.go` - sqlc 生成的 type-safe queries

### 2. 依賴注入實作得很好

`internal/app/app.go` 的 `New()` 函數展現了優秀的依賴注入模式：

```go
// 由外而內建構依賴關係
queries := db.New(pool)
userRepo := repository.NewUserRepository(queries)
authService := service.NewAuthService(userRepo, log)
authHandler := handler.NewAuthHandler(authService, log, cfg.Server.RequestTimeout)
```

**優點：**

- 沒有使用複雜的 DI framework（符合 Go 簡潔哲學）
- 依賴流向清晰且單向
- 易於理解和維護

### 3. 使用 sqlc 產生 type-safe 程式碼

選擇 sqlc 而非 ORM 是非常好的決定：

- 避免了 ORM 的複雜性和隱藏邏輯
- 保證了類型安全
- SQL 查詢可視且可優化
- 符合 Go 社群的主流趨勢

### 4. Graceful Shutdown 實作完善

`internal/app/app.go` 的 `Run()` 方法使用 `errgroup` 管理並發 goroutines：

```go
g, ctx := errgroup.WithContext(context.Background())

// Start HTTP server
g.Go(func() error {
    return a.server.Start()
})

// Handle OS signals
g.Go(func() error {
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    // ...
})
```

這是 production-ready 的標準做法，確保服務可以優雅地關閉。

### 5. 錯誤處理良好

**自定義錯誤類型：**

- `service.ErrUsernameExists` - 業務邏輯層錯誤
- `repository.ErrUserNotFound` - 資料層錯誤
- `app.ErrShutdownSignal` - 應用程式層錯誤

**錯誤包裝和日誌：**

```go
if err != nil {
    s.logger.Error("Failed to create user", "error", err, "username", req.Username)
    return 0, fmt.Errorf("failed to create user [username=%s]: %w", req.Username, err)
}
```

提供了足夠的 context 用於除錯。

### 6. Context Timeout 使用正確

每個 HTTP handler 都正確使用 context timeout：

```go
ctx, cancel := context.WithTimeout(c.Request.Context(), h.requestTimeout)
defer cancel()

userID, err := h.authService.Register(ctx, req)
```

這確保了請求不會無限期地佔用資源。

### 7. 配置管理完善

- 所有配置從環境變數載入
- 配置結構清晰分類（Server、Database、Logger）
- 有完整的 `.env.example` 範例
- 包含單元測試驗證配置載入邏輯

---

## ⚠️ 需要改進的地方

### 1. ❗ 關鍵：應該使用 Interface 而非具體型別

#### 問題描述

目前 Handler 直接依賴 Service 的具體實作，違反了 Go 的 "Accept interfaces, return structs" 原則。

**目前的做法：**

```go
// internal/handler/auth.go:15-18
type AuthHandler struct {
    authService    *service.AuthService  // ❌ 依賴具體實作
    logger         *slog.Logger
    requestTimeout time.Duration
}
```

**影響：**

- 難以進行單元測試（無法輕易 mock）
- 增加了層與層之間的耦合度
- 不符合 Go 社群的強烈共識

#### 建議改法

**步驟 1：在 service package 定義 interface**

```go
// internal/service/auth.go

// AuthService 定義認證服務的介面
type AuthService interface {
    Register(ctx context.Context, req dto.RegisterRequest) (int64, error)
    Login(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, error)
    // 未來的方法...
}

// authService 實作 AuthService interface（小寫開頭變為私有）
type authService struct {
    userRepo repository.UserRepository
    logger   *slog.Logger
}

// NewAuthService 回傳 interface 型別
func NewAuthService(userRepo repository.UserRepository, logger *slog.Logger) AuthService {
    return &authService{
        userRepo: userRepo,
        logger:   logger,
    }
}

func (s *authService) Register(ctx context.Context, req dto.RegisterRequest) (int64, error) {
    // 實作保持不變
}
```

**步驟 2：Handler 依賴 interface**

```go
// internal/handler/auth.go
type AuthHandler struct {
    authService    service.AuthService  // ✅ 依賴 interface
    logger         *slog.Logger
    requestTimeout time.Duration
}

func NewAuthHandler(authService service.AuthService, logger *slog.Logger, requestTimeout time.Duration) *AuthHandler {
    return &AuthHandler{
        authService:    authService,
        logger:         logger,
        requestTimeout: requestTimeout,
    }
}
```

**優點：**

- 符合 Go 語言的設計哲學
- 可以輕鬆建立 mock 用於單元測試
- 降低層與層之間的耦合度
- 提高程式碼的可維護性和可測試性

---

### 2. Router 的依賴注入方式可改進

#### 問題描述

目前 `internal/router/router.go` 使用函數式方法註冊路由：

```go
func SetupRoutes(
    r *gin.Engine,
    healthHandler *handler.HealthHandler,
    authHandler *handler.AuthHandler,
) {
    r.GET("/healthz", healthHandler.HealthCheck)
    r.POST("/auth/register", authHandler.Register)
}
```

當專案成長時，這個函數會變得很長，且難以組織路由群組和 middleware。

#### 建議改法

**使用 Router 結構體：**

```go
// internal/router/router.go
type Router struct {
    engine         *gin.Engine
    healthHandler  *handler.HealthHandler
    authHandler    *handler.AuthHandler
    // 未來的 handlers...
}

func NewRouter(
    engine *gin.Engine,
    healthHandler *handler.HealthHandler,
    authHandler *handler.AuthHandler,
) *Router {
    return &Router{
        engine:        engine,
        healthHandler: healthHandler,
        authHandler:   authHandler,
    }
}

func (rt *Router) Setup() {
    rt.setupHealthRoutes()
    rt.setupAuthRoutes()
    // rt.setupCourseRoutes()  // 未來擴展
}

func (rt *Router) setupHealthRoutes() {
    rt.engine.GET("/healthz", rt.healthHandler.HealthCheck)
}

func (rt *Router) setupAuthRoutes() {
    auth := rt.engine.Group("/auth")
    {
        auth.POST("/register", rt.authHandler.Register)
        auth.POST("/login", rt.authHandler.Login)
    }
}

// 未來可以輕鬆加入 middleware
func (rt *Router) setupCourseRoutes() {
    courses := rt.engine.Group("/courses")
    courses.Use(middleware.RequireAuth())  // 需要認證的路由
    {
        courses.GET("", rt.courseHandler.List)
        courses.POST("", rt.courseHandler.Create)
    }
}
```

**更新 app.go：**

```go
// internal/app/app.go
router := router.NewRouter(ginRouter, healthHandler, authHandler)
router.Setup()
```

**優點：**

- 更好的組織性（路由分組清晰）
- 未來加 middleware 更方便
- 符合 Gin 官方推薦的路由分組模式
- 每個路由群組可以獨立管理

---

### 3. 密碼驗證的邏輯問題

#### 問題描述

目前在 DTO 層限制密碼最大長度為 72：

```go
// internal/dto/auth.go:5
Password string `json:"password" binding:"required,min=8,max=72,password_charset"`
```

這個 72 字元的限制來自 bcrypt 的實作細節，但這是**業務規則**而非**格式驗證**，不應該在 DTO 層處理。

#### 問題分析

- bcrypt 只能處理最多 72 bytes（不是字元）
- 這個限制與加密實作綁定，屬於業務邏輯層的責任
- 未來如果更換加密演算法，需要修改 DTO 定義（違反單一職責原則）

#### 建議改法

**DTO 層只做基本格式驗證：**

```go
// internal/dto/auth.go
type RegisterRequest struct {
    Username string `json:"username" binding:"required,min=3,max=50,alphanum_underscore"`
    Password string `json:"password" binding:"required,min=8,max=128,password_charset"`
}
```

**Service 層檢查業務規則：**

```go
// internal/service/auth.go

var (
    ErrUsernameExists  = errors.New("username already exists")
    ErrPasswordTooLong = errors.New("password too long")  // 新增
)

func (s *authService) Register(ctx context.Context, req dto.RegisterRequest) (int64, error) {
    // bcrypt 只能處理最多 72 bytes
    if len(req.Password) > 72 {
        return 0, ErrPasswordTooLong
    }

    // 檢查使用者名稱是否存在
    exists, err := s.userRepo.ExistsByUsername(ctx, req.Username)
    // ...
}
```

**Handler 層處理新錯誤：**

```go
// internal/handler/auth.go
func (h *AuthHandler) Register(c *gin.Context) {
    // ... binding 邏輯 ...

    userID, err := h.authService.Register(ctx, req)
    if err != nil {
        if errors.Is(err, service.ErrUsernameExists) {
            c.JSON(http.StatusConflict, dto.ErrorResponse{
                Error: "使用者名稱已存在",
            })
            return
        }
        if errors.Is(err, service.ErrPasswordTooLong) {
            c.JSON(http.StatusBadRequest, dto.ErrorResponse{
                Error: "密碼長度超過限制",
            })
            return
        }
        // ...
    }
}
```

**優點：**

- 符合單一職責原則
- DTO 層專注於格式驗證
- Service 層處理業務規則
- 未來更換加密演算法時更靈活

---

### 4. Error Response 結構太簡單

#### 問題描述

目前的錯誤回應結構：

```go
// internal/dto/common.go
type ErrorResponse struct {
    Error string `json:"error"`
}
```

只有人類可讀的錯誤訊息，缺少機器可讀的錯誤碼。

#### 影響

- 前端難以精確判斷錯誤類型（只能靠字串比對，容易出錯）
- 國際化困難（錯誤訊息寫死中文）
- 不符合 RESTful API 最佳實踐（如 Google API Design Guide）

#### 建議改法

```go
// internal/dto/common.go

// ErrorResponse 標準錯誤回應格式
type ErrorResponse struct {
    Code    string `json:"code"`              // 錯誤碼（機器可讀）
    Message string `json:"message"`           // 錯誤訊息（人類可讀）
    Details any    `json:"details,omitempty"` // 選填：額外細節（如驗證錯誤詳情）
}

// 定義錯誤碼常數
const (
    ErrCodeInvalidInput    = "INVALID_INPUT"
    ErrCodeUsernameExists  = "USERNAME_EXISTS"
    ErrCodeInternalError   = "INTERNAL_ERROR"
    ErrCodeUnauthorized    = "UNAUTHORIZED"
    ErrCodeNotFound        = "NOT_FOUND"
)
```

**使用範例：**

```go
// internal/handler/auth.go
func (h *AuthHandler) Register(c *gin.Context) {
    var req dto.RegisterRequest

    if err := c.ShouldBindJSON(&req); err != nil {
        h.logger.Warn("Invalid registration request", "error", err)
        c.JSON(http.StatusBadRequest, dto.ErrorResponse{
            Code:    dto.ErrCodeInvalidInput,
            Message: "使用者名稱或密碼格式無效",
            Details: err.Error(),  // 開發環境可以提供詳細資訊
        })
        return
    }

    userID, err := h.authService.Register(ctx, req)
    if err != nil {
        if errors.Is(err, service.ErrUsernameExists) {
            c.JSON(http.StatusConflict, dto.ErrorResponse{
                Code:    dto.ErrCodeUsernameExists,
                Message: "使用者名稱已存在",
            })
            return
        }

        h.logger.Error("Registration failed", "error", err)
        c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
            Code:    dto.ErrCodeInternalError,
            Message: "註冊失敗，請稍後再試",
        })
        return
    }

    c.JSON(http.StatusCreated, dto.RegisterResponse{
        Message: "註冊成功",
        UserID:  userID,
    })
}
```

**前端使用範例：**

```javascript
// 前端可以根據 code 做精確處理
fetch('/auth/register', {
	method: 'POST',
	body: JSON.stringify({ username, password })
})
	.then((res) => res.json())
	.then((data) => {
		if (data.code === 'USERNAME_EXISTS') {
			showError('這個使用者名稱已經有人使用囉！')
		} else if (data.code === 'INVALID_INPUT') {
			showError('輸入格式不正確')
		}
	})
```

**優點：**

- 前端可以根據 `code` 做更精確的錯誤處理
- 支援國際化（前端根據 code 顯示對應語言的訊息）
- 符合業界標準（Google API Design Guide、Stripe API 等都採用類似結構）
- 機器可讀，方便自動化測試

---

### 5. 缺少 Request ID / Trace ID

#### 問題描述

目前沒有實作 request tracing 機制，當發生問題時：

- 難以在日誌中追蹤單一請求的完整流程
- 無法關聯前端錯誤與後端日誌
- 在分散式系統中無法追蹤請求鏈路

#### 建議改法

**建立 Request ID Middleware：**

```go
// internal/middleware/request_id.go
package middleware

import (
    "github.com/gin-gonic/gin"
    "github.com/google/uuid"
)

const RequestIDKey = "request_id"

// RequestID middleware 為每個請求生成或提取 Request ID
func RequestID() gin.HandlerFunc {
    return func(c *gin.Context) {
        // 優先使用客戶端提供的 Request ID（如果有的話）
        requestID := c.GetHeader("X-Request-ID")
        if requestID == "" {
            // 沒有的話就生成新的
            requestID = uuid.New().String()
        }

        // 存入 context
        c.Set(RequestIDKey, requestID)

        // 回傳給客戶端（方便前端追蹤）
        c.Header("X-Request-ID", requestID)

        c.Next()
    }
}
```

**在 Router 中啟用：**

```go
// internal/router/router.go
func (rt *Router) Setup() {
    // 在所有路由之前啟用 Request ID middleware
    rt.engine.Use(middleware.RequestID())

    rt.setupHealthRoutes()
    rt.setupAuthRoutes()
}
```

**在 Logger 中使用：**

```go
// internal/handler/auth.go
func (h *AuthHandler) Register(c *gin.Context) {
    requestID := c.GetString(middleware.RequestIDKey)

    // 在 logger 中帶入 request_id
    logger := h.logger.With("request_id", requestID)

    var req dto.RegisterRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        logger.Warn("Invalid registration request", "error", err)
        // ...
    }

    // ...
}
```

**優點：**

- 可以在日誌中追蹤單一請求的完整流程
- 前端可以回報錯誤時附上 Request ID，方便除錯
- 為未來的分散式追蹤（如 OpenTelemetry）打下基礎

---

### 6. Logger 缺乏 Context 整合

#### 問題描述

目前每次記錄日誌時，都需要手動傳入 request_id、user_id 等欄位：

```go
h.logger.Warn("Invalid request", "error", err, "request_id", requestID)
```

這種做法容易遺漏重要資訊，且重複性高。

#### 建議改法

**建立 Logger Context 工具函數：**

```go
// internal/infrastructure/logger/context.go
package logger

import (
    "context"
    "log/slog"
)

type contextKey string

const loggerKey contextKey = "logger"

// WithContext 將 logger 存入 context
func WithContext(ctx context.Context, logger *slog.Logger) context.Context {
    return context.WithValue(ctx, loggerKey, logger)
}

// FromContext 從 context 取得 logger（帶有 request_id 等資訊）
func FromContext(ctx context.Context) *slog.Logger {
    if logger, ok := ctx.Value(loggerKey).(*slog.Logger); ok {
        return logger
    }
    return slog.Default()
}
```

**建立 Logger Middleware：**

```go
// internal/middleware/logger.go
package middleware

import (
    "log/slog"
    "time"

    "github.com/gin-gonic/gin"
    "github.com/linporu/waterballsa-backend-golang/internal/infrastructure/logger"
)

// Logger middleware 為每個請求建立帶有 context 的 logger
func Logger(baseLogger *slog.Logger) gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()
        path := c.Request.URL.Path

        // 從 context 取得 request_id
        requestID := c.GetString(RequestIDKey)

        // 建立帶有 request_id 的 logger
        reqLogger := baseLogger.With(
            "request_id", requestID,
            "method", c.Request.Method,
            "path", path,
        )

        // 存入 gin.Context（方便 handler 取用）
        c.Set("logger", reqLogger)

        // 也存入 Request.Context（方便 service/repository 層取用）
        ctx := logger.WithContext(c.Request.Context(), reqLogger)
        c.Request = c.Request.WithContext(ctx)

        c.Next()

        // 請求結束後記錄
        latency := time.Since(start)
        status := c.Writer.Status()

        reqLogger.Info("Request completed",
            "status", status,
            "latency_ms", latency.Milliseconds(),
        )
    }
}
```

**在 Handler 中使用：**

```go
// internal/handler/auth.go
func (h *AuthHandler) Register(c *gin.Context) {
    // 從 context 取得帶有 request_id 的 logger
    logger := logger.FromContext(c.Request.Context())

    var req dto.RegisterRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        logger.Warn("Invalid registration request", "error", err)  // 自動帶 request_id
        c.JSON(http.StatusBadRequest, dto.ErrorResponse{
            Error: "使用者名稱或密碼格式無效",
        })
        return
    }

    // ...
}
```

**在 Service 層使用：**

```go
// internal/service/auth.go
func (s *authService) Register(ctx context.Context, req dto.RegisterRequest) (int64, error) {
    // Service 層也可以從 context 取得 logger
    logger := logger.FromContext(ctx)

    exists, err := s.userRepo.ExistsByUsername(ctx, req.Username)
    if err != nil {
        logger.Error("Failed to check username existence", "error", err, "username", req.Username)
        return 0, fmt.Errorf("failed to check username existence: %w", err)
    }

    // ...
}
```

**優點：**

- 所有日誌自動包含 request_id，無需手動傳入
- 減少重複程式碼
- 未來可以輕鬆加入 user_id、tenant_id 等資訊
- Service 和 Repository 層也能享受到相同的 logger context

---

### 7. Repository 的錯誤轉換可以更細緻

#### 問題描述

目前 `Create` 方法直接回傳資料庫錯誤：

```go
// internal/repository/user.go:60-69
func (r *userRepository) Create(ctx context.Context, username, passwordHash string) (int64, error) {
    userID, err := r.queries.CreateUser(ctx, db.CreateUserParams{
        Username:     username,
        PasswordHash: passwordHash,
    })
    if err != nil {
        return 0, err  // ❌ 直接返回資料庫錯誤
    }
    return userID, nil
}
```

雖然在 Service 層有檢查使用者名稱是否存在，但作為防禦性編程，Repository 層也應該處理可能的 unique constraint 錯誤。

#### 建議改法

```go
// internal/repository/user.go
import (
    "errors"
    "fmt"

    "github.com/jackc/pgx/v5"
    "github.com/jackc/pgx/v5/pgconn"
    "github.com/linporu/waterballsa-backend-golang/internal/db"
    "github.com/linporu/waterballsa-backend-golang/internal/model"
)

var (
    ErrUserNotFound    = errors.New("user not found")
    ErrUsernameExists  = errors.New("username already exists")  // 新增
)

func (r *userRepository) Create(ctx context.Context, username, passwordHash string) (int64, error) {
    userID, err := r.queries.CreateUser(ctx, db.CreateUserParams{
        Username:     username,
        PasswordHash: passwordHash,
    })
    if err != nil {
        // 檢查是否為 PostgreSQL 錯誤
        var pgErr *pgconn.PgError
        if errors.As(err, &pgErr) {
            // 23505 = unique_violation
            if pgErr.Code == "23505" {
                return 0, ErrUsernameExists
            }
        }

        return 0, fmt.Errorf("failed to create user in database: %w", err)
    }

    return userID, nil
}
```

**Service 層也需要處理這個錯誤：**

```go
// internal/service/auth.go
func (s *authService) Register(ctx context.Context, req dto.RegisterRequest) (int64, error) {
    // ... 密碼長度檢查 ...

    // 檢查使用者名稱是否存在
    exists, err := s.userRepo.ExistsByUsername(ctx, req.Username)
    if err != nil {
        logger.Error("Failed to check username existence", "error", err)
        return 0, fmt.Errorf("failed to check username existence: %w", err)
    }
    if exists {
        return 0, ErrUsernameExists
    }

    // Hash password
    passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
    if err != nil {
        logger.Error("Failed to hash password", "error", err)
        return 0, fmt.Errorf("failed to hash password: %w", err)
    }

    // Create user
    userID, err := s.userRepo.Create(ctx, req.Username, string(passwordHash))
    if err != nil {
        // 雖然已經檢查過，但還是處理可能的 race condition
        if errors.Is(err, repository.ErrUsernameExists) {
            return 0, ErrUsernameExists
        }

        logger.Error("Failed to create user", "error", err)
        return 0, fmt.Errorf("failed to create user: %w", err)
    }

    return userID, nil
}
```

**優點：**

- 防禦性編程，處理可能的 race condition
- Repository 層回傳的錯誤更具語意
- 不會洩漏資料庫實作細節到上層

---

### 8. Validator 註冊位置需要加註解

#### 問題描述

Validator 在 `internal/app/app.go:50-52` 註冊：

```go
// Register custom validators (application-level setup)
if err := validator.RegisterAuthValidators(); err != nil {
    return nil, fmt.Errorf("failed to register validators: %w", err)
}
```

註解說明了是 "application-level setup"，但沒有解釋**為什麼**要在這裡註冊。

#### 建議改善

```go
// Register custom validators (must be done before any request handling starts)
// These validators extend Gin's default validation rules for auth-specific fields
// and are registered globally for the application lifetime
if err := validator.RegisterAuthValidators(); err != nil {
    return nil, fmt.Errorf("failed to register validators: %w", err)
}
```

**說明：**

- 為什麼在這裡：必須在處理任何請求之前註冊
- 做什麼：擴展 Gin 的預設驗證規則
- 作用範圍：全域註冊，整個應用程式生命週期有效

---

### 9. 建議的 Middleware 優先順序

#### 未來實作 Middleware 時的建議順序

當你實作 CORS、Auth、Error Handling、Rate Limiting、Log middleware 時，註冊順序非常重要。以下是 Go/Gin 社群的最佳實踐：

```go
// internal/router/router.go

func (rt *Router) Setup() {
    // 1. Recovery（最外層，捕捉 panic）
    // 必須第一個註冊，確保能捕捉到所有 middleware 的 panic
    rt.engine.Use(gin.Recovery())

    // 2. CORS（需要在其他 middleware 之前處理 preflight requests）
    rt.engine.Use(middleware.CORS())

    // 3. Request ID（為後續 log 提供追蹤 ID）
    rt.engine.Use(middleware.RequestID())

    // 4. Logger（需要 request ID，所以在 RequestID 之後）
    rt.engine.Use(middleware.Logger(rt.baseLogger))

    // 5. Rate Limiting（阻擋過多請求，避免浪費後續資源）
    rt.engine.Use(middleware.RateLimit())

    // 6. 設定不需要認證的路由
    rt.setupHealthRoutes()
    rt.setupAuthRoutes()  // login, register 不需要 auth

    // 7. Auth 相關路由（需要 auth middleware 的路由群組）
    rt.setupProtectedRoutes()
}

func (rt *Router) setupProtectedRoutes() {
    // 需要認證的路由群組
    protected := rt.engine.Group("")
    protected.Use(middleware.RequireAuth())  // Auth middleware 只套用在這個群組
    {
        rt.setupCourseRoutes(protected)
        rt.setupUserRoutes(protected)
    }
}

func (rt *Router) setupCourseRoutes(group *gin.RouterGroup) {
    courses := group.Group("/courses")
    {
        courses.GET("", rt.courseHandler.List)
        courses.POST("", rt.courseHandler.Create)
        courses.GET("/:id", rt.courseHandler.Get)
    }
}
```

**順序邏輯：**

1. **Recovery** - 最外層，確保 panic 不會讓服務崩潰
2. **CORS** - 處理跨域請求，必須早於其他邏輯
3. **Request ID** - 為每個請求分配 ID，方便追蹤
4. **Logger** - 記錄請求資訊（依賴 Request ID）
5. **Rate Limiting** - 限流，避免惡意請求消耗資源
6. **Auth** - 僅套用在需要認證的路由群組（不是全域）

**注意事項：**

- **不要**在全域套用 Auth middleware（login/register 會失敗）
- Error Handling 可以透過 gin.Recovery() + 自定義 error handler 實作
- 每個 middleware 的順序都有其邏輯，不要隨意調整

---

### 10. 小建議：提升程式碼一致性

#### a) 命名一致性

目前的命名混合使用縮寫和完整名稱：

```go
authService  // 完整名稱
userRepo     // 縮寫
```

**建議統一：**

**選項 1：都使用縮寫（推薦，更簡潔）**

```go
authSvc
userRepo
courseRepo
```

**選項 2：都使用完整名稱（更明確）**

```go
authService
userRepository
courseRepository
```

選擇其中一種並在整個專案中保持一致。

#### b) 錯誤訊息常數化

目前錯誤訊息硬編碼在 handler 層：

```go
c.JSON(http.StatusBadRequest, dto.ErrorResponse{
    Error: "使用者名稱或密碼格式無效",
})
```

**建議定義錯誤訊息常數：**

```go
// internal/handler/errors.go
package handler

const (
    ErrMsgInvalidCredentials = "使用者名稱或密碼格式無效"
    ErrMsgUsernameExists     = "使用者名稱已存在"
    ErrMsgRegistrationFailed = "註冊失敗，請稍後再試"
    ErrMsgPasswordTooLong    = "密碼長度超過限制"
)
```

**使用：**

```go
c.JSON(http.StatusBadRequest, dto.ErrorResponse{
    Code:    dto.ErrCodeInvalidInput,
    Message: ErrMsgInvalidCredentials,
})
```

**優點：**

- 統一管理錯誤訊息
- 未來國際化時更容易（只需修改常數檔案）
- 避免同一訊息在不同地方有不同寫法

#### c) HTTP Status Code 常數

可以考慮定義常用的 status code 組合：

```go
// internal/handler/response.go
package handler

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "github.com/linporu/waterballsa-backend-golang/internal/dto"
)

// RespondError 統一的錯誤回應函數
func RespondError(c *gin.Context, statusCode int, errCode, message string) {
    c.JSON(statusCode, dto.ErrorResponse{
        Code:    errCode,
        Message: message,
    })
}

// RespondSuccess 統一的成功回應函數
func RespondSuccess(c *gin.Context, statusCode int, data any) {
    c.JSON(statusCode, data)
}
```

**使用：**

```go
func (h *AuthHandler) Register(c *gin.Context) {
    var req dto.RegisterRequest

    if err := c.ShouldBindJSON(&req); err != nil {
        h.logger.Warn("Invalid registration request", "error", err)
        RespondError(c, http.StatusBadRequest, dto.ErrCodeInvalidInput, ErrMsgInvalidCredentials)
        return
    }

    userID, err := h.authService.Register(ctx, req)
    if err != nil {
        if errors.Is(err, service.ErrUsernameExists) {
            RespondError(c, http.StatusConflict, dto.ErrCodeUsernameExists, ErrMsgUsernameExists)
            return
        }

        h.logger.Error("Registration failed", "error", err)
        RespondError(c, http.StatusInternalServerError, dto.ErrCodeInternalError, ErrMsgRegistrationFailed)
        return
    }

    RespondSuccess(c, http.StatusCreated, dto.RegisterResponse{
        Message: "註冊成功",
        UserID:  userID,
    })
}
```

**優點：**

- 減少重複程式碼
- 統一回應格式
- 未來加入額外邏輯（如 metrics、logging）更方便

---

## 📊 整體評價

### 架構評分：8.5/10

**優點：**

- ✅ 分層清晰，符合 Go/Gin 社群標準
- ✅ 使用 sqlc 而非 ORM（優秀選擇）
- ✅ 依賴注入實作良好
- ✅ 配置管理完善

**改進空間：**

- ⚠️ 缺少 interface 抽象（最需要改進）
- ⚠️ Router 組織方式可以更模組化

### 程式碼品質：8/10

**優點：**

- ✅ 錯誤處理完整且有語意
- ✅ Context timeout 使用正確
- ✅ Graceful shutdown 實作完善
- ✅ 日誌記錄充足

**改進空間：**

- ⚠️ 缺少 request tracing 機制
- ⚠️ Logger 沒有與 context 深度整合
- ⚠️ 錯誤回應結構過於簡單

### Go 慣用性：7.5/10

**優點：**

- ✅ 大部分符合 Go idioms
- ✅ 錯誤處理遵循 Go 慣例
- ✅ 沒有過度設計

**改進空間：**

- ⚠️ 應遵循 "Accept interfaces, return structs" 原則
- ⚠️ 業務邏輯與格式驗證混淆（密碼長度限制）

### 可維護性：8/10

**優點：**

- ✅ 程式碼結構清晰
- ✅ 檔案組織合理
- ✅ 註解適量

**改進空間：**

- ⚠️ 缺少統一的回應處理機制
- ⚠️ 錯誤訊息硬編碼

### 可測試性：6.5/10

**改進空間：**

- ⚠️ 缺少 interface 抽象，難以 mock
- ⚠️ 沒有單元測試範例
- ⚠️ Handler 層難以獨立測試

---

## 🎯 優先改進建議

### 高優先（架構性問題）

1. **引入 Service 和 Repository 的 interface 抽象**

   - 影響：可測試性、可維護性
   - 難度：中
   - 預估工時：2-3 小時

2. **重構 Router 為結構體 + 路由分組**
   - 影響：可擴展性、middleware 管理
   - 難度：低
   - 預估工時：1-2 小時

### 中優先（可維護性）

3. **改進 Error Response 結構**

   - 影響：API 品質、前端體驗
   - 難度：低
   - 預估工時：1 小時

4. **加入 Request ID middleware**

   - 影響：可除錯性
   - 難度：低
   - 預估工時：1 小時

5. **整合 Logger 與 Context**
   - 影響：日誌品質、開發體驗
   - 難度：中
   - 預估工時：2 小時

### 低優先（錦上添花）

6. **錯誤訊息常數化**

   - 影響：國際化準備
   - 難度：低
   - 預估工時：0.5 小時

7. **命名一致性調整**

   - 影響：程式碼可讀性
   - 難度：低
   - 預估工時：0.5 小時

8. **Repository 錯誤轉換**
   - 影響：健壯性
   - 難度：低
   - 預估工時：0.5 小時

---

## 總結

你的專案架構基礎非常紮實，分層清晰且符合 Go 社群的主流做法。主要的改進點集中在：

1. **Interface 抽象** - 這是最關鍵的改進，會大幅提升可測試性和可維護性
2. **Request Tracing** - 加入 Request ID 會讓除錯變得更容易
3. **錯誤處理優化** - 改進錯誤回應結構，提供更好的 API 體驗

其他建議都是基於 production 環境的經驗，但現階段你的架構已經可以支撐線上課程平台的開發了。

### 架構優勢

- 清晰的分層架構（Handler → Service → Repository → DB）
- 使用 sqlc 避免 ORM 複雜性
- 良好的依賴注入模式
- 完善的 graceful shutdown 機制
- 正確的 context timeout 使用

### 值得保持的實踐

- 不過度設計，保持簡潔
- 適當的錯誤包裝和日誌記錄
- 配置管理完善
- 使用 errgroup 協調並發

繼續保持這種清晰的分層和良好的錯誤處理習慣，相信這個專案會越來越好！

---

## 附錄：參考資源

### Go 最佳實踐

- [Effective Go](https://go.dev/doc/effective_go)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- [Uber Go Style Guide](https://github.com/uber-go/guide/blob/master/style.md)

### Gin Framework

- [Gin 官方文檔](https://gin-gonic.com/docs/)
- [Gin Examples](https://github.com/gin-gonic/examples)

### Architecture

- [Standard Go Project Layout](https://github.com/golang-standards/project-layout)
- [Go Clean Architecture](https://github.com/bxcodec/go-clean-arch)

### API Design

- [Google API Design Guide](https://cloud.google.com/apis/design)
- [RESTful API Best Practices](https://stackoverflow.blog/2020/03/02/best-practices-for-rest-api-design/)
