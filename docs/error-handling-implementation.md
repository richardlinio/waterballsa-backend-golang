# Gin Error Handling Middleware 實作計畫

## 目標

實作統一的錯誤處理機制，使用 middleware + 自訂錯誤類型，支援錯誤代碼 + 中文訊息，覆蓋全層級（service、handler、middleware）。

## 設計概述

### 自訂錯誤類型架構

```
internal/apperror/
├── error.go          - 定義 AppError 結構和工廠函數
├── codes.go          - 定義所有錯誤代碼常數
├── messages.go       - 定義錯誤代碼對應的中文訊息
└── status.go         - 定義錯誤代碼對應的 HTTP 狀態碼
```

### 錯誤流程

1. **Service 層**：拋出 `AppError`（只需指定 Code，HTTP 狀態碼由 apperror 包自動對應）
2. **Handler 層**：調用 service 後，使用 `c.Error(err)` 將錯誤附加到 context，然後 `return`
3. **Middleware**：檢查 `c.Errors`，根據錯誤類型回傳統一格式的 JSON

### 設計原則

- **Service 層與 HTTP 完全解耦**：Service 層不需要 `import "net/http"`，只需要知道業務錯誤代碼
- **集中管理對應關係**：錯誤代碼到 HTTP 狀態碼的對應關係在 `apperror` 包中統一管理
- **降低耦合度**：Service 層只關心業務邏輯錯誤，不關心 HTTP 協議細節

## 實作步驟

### Step 1: 建立錯誤代碼定義

**檔案**: `internal/apperror/codes.go`

定義所有錯誤代碼常數：

```go
const (
    // Validation errors (4000x)
    CodeValidationFailed    = "ERR_VALIDATION_FAILED"
    CodeInvalidRequest      = "ERR_INVALID_REQUEST"

    // Auth errors (4010x)
    CodeUsernameExists      = "ERR_USERNAME_EXISTS"
    CodePasswordTooLong     = "ERR_PASSWORD_TOO_LONG"
    CodeUnauthorized        = "ERR_UNAUTHORIZED"
    CodeInvalidCredentials  = "ERR_INVALID_CREDENTIALS"

    // Resource errors (4040x)
    CodeResourceNotFound    = "ERR_RESOURCE_NOT_FOUND"

    // Server errors (5000x)
    CodeInternalError       = "ERR_INTERNAL_ERROR"
    CodeDatabaseError       = "ERR_DATABASE_ERROR"
    CodeServiceUnavailable  = "ERR_SERVICE_UNAVAILABLE"
)
```

### Step 2: 建立錯誤訊息對應

**檔案**: `internal/apperror/messages.go`

建立錯誤代碼到中文訊息的對應表：

```go
package apperror

var errorMessages = map[string]string{
    CodeValidationFailed:    "輸入資料驗證失敗",
    CodeInvalidRequest:      "請求格式無效",
    CodeUsernameExists:      "使用者名稱已存在",
    CodePasswordTooLong:     "使用者名稱或密碼格式無效",
    CodeUnauthorized:        "未授權存取",
    CodeInvalidCredentials:  "使用者名稱或密碼錯誤",
    CodeResourceNotFound:    "找不到指定的資源",
    CodeInternalError:       "伺服器內部錯誤",
    CodeDatabaseError:       "資料庫操作失敗",
    CodeServiceUnavailable:  "服務暫時無法使用",
}

func GetMessage(code string) string {
    if msg, ok := errorMessages[code]; ok {
        return msg
    }
    return "未知錯誤"
}
```

### Step 3: 建立 HTTP 狀態碼對應

**檔案**: `internal/apperror/status.go`

建立錯誤代碼到 HTTP 狀態碼的對應表：

```go
package apperror

var httpStatusMap = map[string]int{
    // Validation errors (400)
    CodeValidationFailed: 400,
    CodeInvalidRequest:   400,
    CodePasswordTooLong:  400,

    // Auth errors (401/409)
    CodeUnauthorized:       401,
    CodeInvalidCredentials: 401,
    CodeUsernameExists:     409,

    // Resource errors (404)
    CodeResourceNotFound: 404,

    // Server errors (500/503)
    CodeInternalError:      500,
    CodeDatabaseError:      500,
    CodeServiceUnavailable: 503,
}

func GetHTTPStatus(code string) int {
    if status, ok := httpStatusMap[code]; ok {
        return status
    }
    return 500 // default to Internal Server Error
}
```

### Step 4: 建立 AppError 結構

**檔案**: `internal/apperror/error.go`

```go
package apperror

type AppError struct {
    Code       string // 錯誤代碼 (e.g., "ERR_USERNAME_EXISTS")
    Message    string // 中文錯誤訊息
    HTTPStatus int    // HTTP 狀態碼
    Err        error  // 原始錯誤 (用於日誌記錄)
}

func (e *AppError) Error() string {
    return e.Message
}

func (e *AppError) Unwrap() error {
    return e.Err
}

// 工廠函數 - HTTP 狀態碼由 code 自動對應
func New(code string) *AppError {
    return &AppError{
        Code:       code,
        Message:    GetMessage(code),
        HTTPStatus: GetHTTPStatus(code), // 自動對應 HTTP 狀態碼
    }
}

func NewWithError(code string, err error) *AppError {
    return &AppError{
        Code:       code,
        Message:    GetMessage(code),
        HTTPStatus: GetHTTPStatus(code), // 自動對應 HTTP 狀態碼
        Err:        err,
    }
}

// 預定義的錯誤建立函數 - 不需要指定 HTTP 狀態碼
func ValidationFailed() *AppError {
    return New(CodeValidationFailed)
}

func UsernameExists() *AppError {
    return New(CodeUsernameExists)
}

func PasswordTooLong() *AppError {
    return New(CodePasswordTooLong)
}

func InternalError(err error) *AppError {
    return NewWithError(CodeInternalError, err)
}

func DatabaseError(err error) *AppError {
    return NewWithError(CodeDatabaseError, err)
}

func ServiceUnavailable() *AppError {
    return New(CodeServiceUnavailable)
}
```

### Step 5: 建立 Error Handling Middleware

**檔案**: `internal/middleware/error_handler.go`

```go
package middleware

import (
    "errors"
    "fmt"
    "log/slog"
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/go-playground/validator/v10"
    "github.com/richardlinio/waterballsa-backend-golang/internal/apperror"
    "github.com/richardlinio/waterballsa-backend-golang/internal/dto"
)

func ErrorHandler(logger *slog.Logger) gin.HandlerFunc {
    return func(c *gin.Context) {
        // Step 1: 先處理請求
        c.Next()

        // Step 2: 檢查是否有錯誤
        if len(c.Errors) == 0 {
            return
        }

        // Step 3: 檢查是否已經回應（避免重複寫入）
        if c.Writer.Written() {
            return
        }

        // Step 4: 取得最後一個錯誤
        err := c.Errors.Last().Err

        // Step 5: 判斷錯誤類型並回傳
        var appErr *apperror.AppError
        if errors.As(err, &appErr) {
            // 自訂錯誤 - 記錄日誌並回傳結構化錯誤
            logger.Warn("Application error",
                "code", appErr.Code,
                "message", appErr.Message,
                "status", appErr.HTTPStatus,
                "path", c.Request.URL.Path,
                "method", c.Request.Method,
            )

            // 如果有原始錯誤,記錄詳細錯誤
            if appErr.Err != nil {
                logger.Error("Underlying error", "error", appErr.Err)
            }

            response := dto.ErrorResponse{
                Code:  appErr.Code,
                Error: appErr.Message,
            }

            // 解析驗證錯誤細節
            if appErr.Code == apperror.CodeValidationFailed && appErr.Err != nil {
                response.Details = parseValidationErrors(appErr.Err)
            }

            c.JSON(appErr.HTTPStatus, response)
        } else {
            // 未預期的錯誤 - 記錄完整錯誤並回傳通用錯誤
            logger.Error("Unexpected error",
                "error", err,
                "path", c.Request.URL.Path,
                "method", c.Request.Method,
            )

            c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
                Code:  apperror.CodeInternalError,
                Error: apperror.GetMessage(apperror.CodeInternalError),
            })
        }
    }
}

// parseValidationErrors 解析 Gin 驗證錯誤
func parseValidationErrors(err error) map[string]string {
    var ve validator.ValidationErrors
    if errors.As(err, &ve) {
        details := make(map[string]string)
        for _, fe := range ve {
            details[fe.Field()] = getValidationMessage(fe)
        }
        return details
    }
    return nil
}

// getValidationMessage 根據驗證規則返回中文錯誤訊息
func getValidationMessage(fe validator.FieldError) string {
    switch fe.Tag() {
    case "required":
        return "此欄位為必填"
    case "min":
        return fmt.Sprintf("長度至少需要 %s 個字元", fe.Param())
    case "max":
        return fmt.Sprintf("長度不可超過 %s 個字元", fe.Param())
    default:
        return "格式無效"
    }
}
```

### Step 6: 修改 Service 層錯誤

**檔案**: `internal/service/auth.go`

**重要原則**：

- Service 層**不應記錄錯誤日誌**，統一由 middleware 處理
- Service 層**只需返回包裝好的 AppError**
- 移除所有 `logger.Error(...)` 呼叫

將現有的 `errors.New` 改為 `apperror`：

**Before**:

```go
var (
    ErrUsernameExists  = errors.New("username already exists")
    ErrPasswordTooLong = errors.New("password too long")
)
```

**After**:

```go
// 移除 var 區塊中的錯誤定義,直接使用 apperror
```

在 `Register` 方法中：

```go
// Before
if len(req.Password) > 72 {
    return 0, ErrPasswordTooLong
}

// After
if len(req.Password) > 72 {
    return 0, apperror.PasswordTooLong()
}

// Before
if exists {
    return 0, ErrUsernameExists
}

// After
if exists {
    return 0, apperror.UsernameExists()
}

// Before (資料庫錯誤 + 日誌記錄)
s.logger.Error("Failed to check username existence", "error", err, "username", req.Username)
return 0, fmt.Errorf("failed to check username existence [username=%s]: %w", req.Username, err)

// After (只返回錯誤，不記錄日誌)
return 0, apperror.DatabaseError(err)
```

**完整的 Service 層範例**（移除所有日誌記錄）：

```go
func (s *authService) Register(ctx context.Context, req dto.RegisterRequest) (int64, error) {
    if len(req.Password) > 72 {
        return 0, apperror.PasswordTooLong()
    }

    exists, err := s.userRepository.ExistsByUsername(ctx, req.Username)
    if err != nil {
        return 0, apperror.DatabaseError(err)
    }

    if exists {
        return 0, apperror.UsernameExists()
    }

    passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
    if err != nil {
        return 0, apperror.InternalError(err)
    }

    userID, err := s.userRepository.Create(ctx, req.Username, string(passwordHash))
    if err != nil {
        return 0, apperror.DatabaseError(err)
    }

    return userID, nil
}
```

### Step 7: 修改 Handler 層

**檔案**: `internal/handler/auth.go`

簡化錯誤處理邏輯，使用 `c.Error()` 代替直接回傳 JSON：

**Before**:

```go
if err := c.ShouldBindJSON(&req); err != nil {
    h.logger.Warn("Invalid registration request", "error", err)
    c.JSON(http.StatusBadRequest, dto.ErrorResponse{
        Error: "使用者名稱或密碼格式無效",
    })
    return
}

// ...

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
            Error: "使用者名稱或密碼格式無效",
        })
        return
    }

    h.logger.Error("Registration failed", "error", err)
    c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
        Error: "註冊失敗,請稍後再試",
    })
    return
}
```

**After**:

```go
if err := c.ShouldBindJSON(&req); err != nil {
    // 使用 NewWithError 保留原始驗證錯誤，middleware 會解析欄位細節
    _ = c.Error(apperror.NewWithError(apperror.CodeValidationFailed, err))
    return
}

// ...

userID, err := h.authService.Register(ctx, req)
if err != nil {
    _ = c.Error(err)
    return
}
```

**檔案**: `internal/handler/health.go`

簡化 health check 錯誤處理：

**Before**:

```go
if err := h.pool.Ping(ctx); err != nil {
    h.logger.Error("Health check failed: database ping error", "error", err)
    c.JSON(http.StatusServiceUnavailable, dto.HealthCheckResponse{
        Status:   "DOWN",
        Database: "DOWN",
    })
    return
}
```

**After**:

```go
if err := h.pool.Ping(ctx); err != nil {
    _ = c.Error(apperror.ServiceUnavailable())
    return
}
```

### Step 8: 註冊 Middleware

**檔案**: `internal/router/router.go`

在 `Setup()` 方法中註冊 error handling middleware：

**Before**:

```go
func (r *Router) Setup() {
    // 1. Setup middlewares
    r.engine.Use(middleware.Security())
    r.engine.Use(middleware.CORS(r.corsConfig))

    // 2. Setup routes
    r.setupHealthRoutes()
    r.setupAuthRoutes()
}
```

**After**:

```go
func (r *Router) Setup() {
    // 1. Setup middlewares
    r.engine.Use(middleware.Security())
    r.engine.Use(middleware.CORS(r.corsConfig))
    r.engine.Use(middleware.ErrorHandler(r.logger)) // 新增 error handler

    // 2. Setup routes
    r.setupHealthRoutes()
    r.setupAuthRoutes()
}
```

需要在 `Router` 結構中新增 logger 欄位：

```go
type Router struct {
    engine        *gin.Engine
    corsConfig    config.CORSConfig
    logger        *slog.Logger  // 新增
    healthHandler *handler.HealthHandler
    authHandler   *handler.AuthHandler
}

func NewRouter(
    engine *gin.Engine,
    corsConfig config.CORSConfig,
    logger *slog.Logger,  // 新增參數
    healthHandler *handler.HealthHandler,
    authHandler *handler.AuthHandler,
) *Router {
    return &Router{
        engine:        engine,
        corsConfig:    corsConfig,
        logger:        logger,  // 新增
        healthHandler: healthHandler,
        authHandler:   authHandler,
    }
}
```

### Step 9: 更新 App 初始化

**檔案**: `internal/app/app.go`

在建立 Router 時傳入 logger：

找到 `NewRouter` 的調用位置，新增 logger 參數：

```go
router := router.NewRouter(
    engine,
    cfg.CORS,
    logger,  // 新增 logger 參數
    healthHandler,
    authHandler,
)
```

### Step 10: 更新 DTO

**檔案**: `internal/dto/common.go`

確保 `ErrorResponse` 結構支援錯誤代碼（已經有 `Code` 欄位，無需修改）：

```go
type ErrorResponse struct {
    Code    string `json:"code,omitempty"`
    Error   string `json:"error"`
    Details any    `json:"details,omitempty"`
}
```

## 關鍵檔案清單

### 新建檔案

- `internal/apperror/codes.go` - 錯誤代碼常數
- `internal/apperror/messages.go` - 錯誤訊息對應
- `internal/apperror/status.go` - HTTP 狀態碼對應
- `internal/apperror/error.go` - AppError 結構定義
- `internal/middleware/error_handler.go` - Error handling middleware

### 修改檔案

- `internal/service/auth.go` - 改用 apperror
- `internal/handler/auth.go` - 簡化錯誤處理
- `internal/handler/health.go` - 簡化錯誤處理
- `internal/router/router.go` - 註冊 middleware，新增 logger
- `internal/app/app.go` - 傳遞 logger 給 Router

## 測試計畫

1. **單元測試**：測試 AppError 的建立和方法
2. **整合測試**：
   - 測試 validation 錯誤回傳 400 + 正確錯誤代碼
   - 測試使用者名稱重複回傳 409 + ERR_USERNAME_EXISTS
   - 測試未預期錯誤回傳 500 + ERR_INTERNAL_ERROR
3. **手動測試**：使用 curl/Postman 測試各種錯誤情境

## 優勢

1. **統一錯誤格式**：所有錯誤都透過 middleware 統一處理
2. **錯誤代碼標準化**：前端可根據 code 做客製化處理
3. **中文訊息**：保留使用者友善的中文錯誤訊息
4. **Handler 簡化**：不再需要大量的 if-else 判斷和 JSON 回傳
5. **集中式日誌**：所有錯誤都在 middleware 記錄，便於監控
6. **可擴展性**：未來新增錯誤類型只需在 `apperror` 新增定義
7. **完全解耦 Service 層與 HTTP**：Service 層不需要 `import "net/http"`，只需要知道業務錯誤代碼
8. **集中管理對應關係**：HTTP 狀態碼對應關係統一在 `status.go` 中管理，易於維護

## 注意事項

1. **Middleware 順序**：ErrorHandler 應該在 CORS 之後註冊，確保錯誤回應也包含 CORS headers

2. **統一錯誤處理**：不要在 handler 中直接使用 `c.JSON` 回傳錯誤，統一使用 `c.Error()`

3. **完整覆蓋**：確保所有 service 層的錯誤都使用 `apperror`，避免出現未處理的錯誤類型

4. **日誌記錄責任分離**：
   - **Service 層**：不記錄錯誤日誌，只返回包裝好的 `AppError`
   - **Middleware**：統一記錄所有錯誤日誌（Warn 級別記錄 AppError，Error 級別記錄底層錯誤）
   - **Handler 層**：不記錄錯誤日誌，使用 `c.Error()` 傳遞錯誤給 middleware

5. **驗證錯誤處理**：
   - Handler 使用 `apperror.NewWithError(apperror.CodeValidationFailed, err)` 保留原始驗證錯誤
   - Middleware 會自動解析 Gin 驗證錯誤並填入 `Details` 欄位
   - 前端可根據 `Details` 顯示欄位級別的錯誤訊息

6. **Service 層不需要知道 HTTP**：Service 層只需要 import `internal/apperror`，不需要 import `net/http`

7. **對應關係集中管理**：所有錯誤代碼到 HTTP 狀態碼的對應關係都在 `status.go` 中定義

8. **避免重複寫入**：Middleware 會檢查 `c.Writer.Written()`，防止重複寫入 response
