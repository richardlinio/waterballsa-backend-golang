# Code Review - 使用者登入功能實作

**日期**: 2026-01-14
**審查範圍**: feat/login 分支
**審查者**: Claude Code Review

---

## 📋 概述

此次審查涵蓋使用者登入功能的完整實作，包括 JWT 認證機制、Refresh Token、安全性設計與相關配置。Login 功能作為認證系統的核心，安全性與正確性至關重要。

### 變更檔案清單

#### 配置相關

- [.env.example](.env.example) - JWT 環境變數範例
- [go.mod](go.mod) / [go.sum](go.sum) - 新增 JWT 相關依賴
- [internal/config/config.go](internal/config/config.go) - 整合 JWT 配置
- [internal/config/jwt.go](internal/config/jwt.go) - JWT 配置載入邏輯（新檔案）

#### 核心功能

- [internal/service/auth.go](internal/service/auth.go) - Login 業務邏輯
- [internal/middleware/jwt.go](internal/middleware/jwt.go) - JWT 中介層（新檔案）
- [internal/dto/auth.go](internal/dto/auth.go) - Login 請求/回應結構
- [internal/router/router.go](internal/router/router.go) - 路由註冊
- [internal/app/app.go](internal/app/app.go) - 應用程式初始化

#### 測試相關

- [tests/bdd/features/isa/auth/login.isa.feature](tests/bdd/features/isa/auth/login.isa.feature) - BDD 測試場景
- [tests/bdd/steps/common/http_steps.go](tests/bdd/steps/common/http_steps.go) - HTTP 測試步驟
- [tests/testutil/server.go](tests/testutil/server.go) - 測試伺服器設定

---

## ✅ 整體優點

### 1. 安全性設計優秀 🔒

- **密碼驗證安全**: 使用 bcrypt 進行密碼比對
- **時序攻擊防護**: Login 失敗使用統一錯誤訊息，不洩漏使用者是否存在
- **Secret 長度驗證**: 強制 JWT_SECRET 至少 32 字元
- **Cookie 安全配置**: HTTPOnly、SameSite、SecureCookie 選項完整

### 2. 架構清晰

- **職責分離明確**: Service 負責業務邏輯，Middleware 負責 JWT 處理
- **錯誤處理統一**: 使用 `apperror` 統一錯誤回應格式
- **配置管理完善**: JWT 配置集中管理，環境變數驗證完整

### 3. 使用者體驗良好

- **雙 Token 機制**: Access Token + Refresh Token
- **靈活的 Token 來源**: 支援 Authorization Header 和 Cookie
- **詳細的日誌記錄**: 認證過程完整記錄，便於追蹤問題

### 4. 程式碼品質高

- **註解完整**: 關鍵邏輯都有清晰註解
- **變數命名清楚**: 遵循 full name 命名規範
- **錯誤處理完整**: 所有錯誤路徑都有適當處理

---

## 🔴 高優先級問題（建議立即修正）

### 1. JWT Middleware 缺少 Refresh Token Store 實作

**檔案**: [internal/middleware/jwt.go](internal/middleware/jwt.go)

**問題嚴重性**: 🔴 嚴重 - 功能不完整

**問題描述**:

`gin-jwt` v3 版本的 Refresh Token 機制需要實作 `RefreshTokenStore` 介面來管理 refresh token 的生命週期。當前程式碼未實作此介面，雖然能通過編譯（因為 `RefreshTokenStore` 是可選的），但可能導致 refresh token 無法正確儲存和驗證。

從 [internal/app/app.go:77-80](internal/app/app.go#L77-L80) 的程式碼可以看出已經呼叫了 `MiddlewareInit()`：

```go
// Initialize middleware (must be called to enable refresh token store)
if err := jwtMiddleware.MiddlewareInit(); err != nil {
    return nil, fmt.Errorf("failed to initialize JWT middleware: %w", err)
}
```

但在 [internal/middleware/jwt.go](internal/middleware/jwt.go) 中，JWT middleware 的配置缺少 `RefreshTokenStore` 欄位。

**影響**:

- Refresh token 可能無法正確儲存
- Token 刷新機制可能失效
- 安全性風險：無法撤銷 refresh token

**建議方案**:

#### 方案 A: 使用記憶體儲存（開發/測試環境）

```go
// internal/middleware/jwt.go
import (
    jwt "github.com/appleboy/gin-jwt/v3"
    "github.com/appleboy/gin-jwt/v3/token/memory"
    // ... 其他 imports
)

func NewJWTMiddleware(
    cfg config.JWTConfig,
    authService service.AuthService,
    logger *slog.Logger,
) (*jwt.GinJWTMiddleware, error) {
    return jwt.New(&jwt.GinJWTMiddleware{
        // ... 現有配置 ...

        // 新增: Refresh Token Store (記憶體版本)
        RefreshTokenStore: memory.NewStore(),

        // ... 其他配置 ...
    })
}
```

#### 方案 B: 使用 Redis 儲存（推薦用於生產環境）

如果專案計畫使用 Redis，建議實作 Redis-based refresh token store：

**步驟 1**: 新增 Redis 配置到 [internal/config/redis.go](internal/config/redis.go)（新檔案）

```go
package config

import (
    "fmt"
    "os"
    "strconv"
)

type RedisConfig struct {
    Host     string
    Port     int
    Password string
    DB       int
}

func loadRedisConfig() (*RedisConfig, error) {
    host := os.Getenv("REDIS_HOST")
    if host == "" {
        host = "localhost" // default
    }

    port := 6379
    if portStr := os.Getenv("REDIS_PORT"); portStr != "" {
        var err error
        port, err = strconv.Atoi(portStr)
        if err != nil {
            return nil, fmt.Errorf("invalid REDIS_PORT: %w", err)
        }
    }

    password := os.Getenv("REDIS_PASSWORD")

    db := 0
    if dbStr := os.Getenv("REDIS_DB"); dbStr != "" {
        var err error
        db, err = strconv.Atoi(dbStr)
        if err != nil {
            return nil, fmt.Errorf("invalid REDIS_DB: %w", err)
        }
    }

    return &RedisConfig{
        Host:     host,
        Port:     port,
        Password: password,
        DB:       db,
    }, nil
}
```

**步驟 2**: 實作 Redis Refresh Token Store

```go
// internal/infrastructure/token/redis_store.go (新檔案)
package token

import (
    "context"
    "fmt"
    "time"

    "github.com/redis/rueidis"
)

type RedisRefreshTokenStore struct {
    client rueidis.Client
    prefix string
}

func NewRedisRefreshTokenStore(client rueidis.Client) *RedisRefreshTokenStore {
    return &RedisRefreshTokenStore{
        client: client,
        prefix: "refresh_token:",
    }
}

func (s *RedisRefreshTokenStore) Set(token string, expiration time.Duration) error {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    key := s.prefix + token
    cmd := s.client.B().Set().Key(key).Value("1").Ex(expiration).Build()

    return s.client.Do(ctx, cmd).Error()
}

func (s *RedisRefreshTokenStore) Get(token string) (bool, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    key := s.prefix + token
    cmd := s.client.B().Get().Key(key).Build()

    result := s.client.Do(ctx, cmd)
    if rueidis.IsRedisNil(result.Error()) {
        return false, nil // Token 不存在
    }
    if result.Error() != nil {
        return false, result.Error()
    }

    return true, nil // Token 存在
}

func (s *RedisRefreshTokenStore) Delete(token string) error {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    key := s.prefix + token
    cmd := s.client.B().Del().Key(key).Build()

    return s.client.Do(ctx, cmd).Error()
}
```

**步驟 3**: 更新 middleware 使用 Redis store

```go
// internal/middleware/jwt.go
func NewJWTMiddleware(
    cfg config.JWTConfig,
    authService service.AuthService,
    logger *slog.Logger,
    refreshTokenStore jwt.RefreshTokenStore, // 新增參數
) (*jwt.GinJWTMiddleware, error) {
    return jwt.New(&jwt.GinJWTMiddleware{
        // ... 現有配置 ...

        RefreshTokenStore: refreshTokenStore,

        // ... 其他配置 ...
    })
}
```

**建議採用方案**:

- 短期：使用方案 A（記憶體儲存）快速修復
- 長期：採用方案 B（Redis）以支援分散式環境和 Token 撤銷功能

---

### 2. Logout 缺少 Refresh Token 清理邏輯

**檔案**: [internal/middleware/jwt.go:182-190](internal/middleware/jwt.go#L182-L190)

**問題嚴重性**: 🔴 嚴重 - 安全性漏洞

**問題描述**:

當前的 `LogoutResponse` 只是返回成功訊息，但並未清理 refresh token。這意味著使用者登出後，攻擊者仍可使用舊的 refresh token 獲取新的 access token。

```go
// 當前實作
LogoutResponse: func(c *gin.Context) {
    claims := jwt.ExtractClaims(c)
    if userID, ok := claims["user_id"].(float64); ok {
        logger.Info("User logged out successfully", "user_id", int64(userID))
    }

    c.JSON(http.StatusOK, gin.H{
        "code":    http.StatusOK,
        "message": "Successfully logged out",
    })
},
```

**安全性影響**:

- 使用者登出後 refresh token 仍然有效
- 無法實作「強制登出」功能
- Token 洩漏風險無法降低

**建議方案**:

#### 完整的 Logout 實作

```go
// internal/middleware/jwt.go

// LogoutResponse customizes the logout response and cleans up tokens
LogoutResponse: func(c *gin.Context) {
    claims := jwt.ExtractClaims(c)
    userID := int64(0)

    if id, ok := claims["user_id"].(float64); ok {
        userID = int64(id)
    }

    // 清理 refresh token
    if err := cleanupRefreshToken(c, logger); err != nil {
        logger.Error("Failed to cleanup refresh token on logout",
            "user_id", userID,
            "error", err)
        // 繼續執行，不要因為清理失敗而阻止登出
    }

    // 清除 cookies
    clearAuthCookies(c, cfg)

    logger.Info("User logged out successfully", "user_id", userID)

    c.JSON(http.StatusOK, gin.H{
        "code":    http.StatusOK,
        "message": "Successfully logged out",
    })
},

// 在檔案中新增輔助函式

// cleanupRefreshToken removes the refresh token from storage
func cleanupRefreshToken(c *gin.Context, logger *slog.Logger) error {
    // 方式 1: 從 cookie 取得 refresh token
    refreshToken, err := c.Cookie("refresh_token")
    if err != nil {
        // 如果 cookie 中沒有，嘗試從 request body 取得
        var body struct {
            RefreshToken string `json:"refreshToken"`
        }
        if bindErr := c.ShouldBindJSON(&body); bindErr == nil && body.RefreshToken != "" {
            refreshToken = body.RefreshToken
        }
    }

    if refreshToken == "" {
        logger.Warn("No refresh token found during logout")
        return nil // 沒有 token 不視為錯誤
    }

    // TODO: 這裡需要透過 RefreshTokenStore 刪除 token
    // 當實作了 RefreshTokenStore 後，呼叫其 Delete 方法
    // 範例：middleware.RefreshTokenStore.Delete(refreshToken)

    logger.Debug("Refresh token cleaned up", "token_prefix", refreshToken[:10]+"...")
    return nil
}

// clearAuthCookies clears all authentication-related cookies
func clearAuthCookies(c *gin.Context, cfg config.JWTConfig) {
    cookieSettings := &http.Cookie{
        Path:     "/",
        Domain:   cfg.CookieDomain,
        MaxAge:   -1, // 刪除 cookie
        HttpOnly: true,
        Secure:   cfg.SecureCookie,
        SameSite: http.SameSiteLaxMode,
    }

    // 清除 access token cookie
    accessCookie := *cookieSettings
    accessCookie.Name = "jwt"
    http.SetCookie(c.Writer, &accessCookie)

    // 清除 refresh token cookie
    refreshCookie := *cookieSettings
    refreshCookie.Name = "refresh_token"
    http.SetCookie(c.Writer, &refreshCookie)
}
```

**重要**: 此方案需要配合問題 #1 的 RefreshTokenStore 實作才能完全發揮作用。

---

### 3. JWT Secret 驗證不夠嚴格

**檔案**: [internal/config/jwt.go:27-28](internal/config/jwt.go#L27-L28)

**問題嚴重性**: 🟠 中高 - 安全性風險

**問題描述**:

當前只檢查長度是否 >= 32 字元，但未驗證 secret 的熵（entropy）和複雜度。弱密碼即使長度足夠，仍可能被暴力破解。

```go
// 當前實作
if len(secret) < 32 {
    return nil, fmt.Errorf("JWT_SECRET must be at least 32 characters long for security (current: %d)", len(secret))
}
```

**風險範例**:

- `"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"` - 32 個 'a'（通過檢查但極弱）
- `"12345678901234567890123456789012"` - 重複數字（通過檢查但極弱）

**建議方案**:

```go
// internal/config/jwt.go

import (
    "math"
    "unicode"
    // ... 其他 imports
)

// loadJWTConfig loads JWT configuration from environment variables
func loadJWTConfig() (*JWTConfig, error) {
    secret := os.Getenv("JWT_SECRET")
    if secret == "" {
        return nil, fmt.Errorf("JWT_SECRET environment variable is required")
    }

    // 驗證 secret 安全性
    if err := validateJWTSecret(secret); err != nil {
        return nil, err
    }

    // ... 其餘程式碼保持不變
}

// validateJWTSecret checks if the JWT secret meets security requirements
func validateJWTSecret(secret string) error {
    const minLength = 32

    // 檢查最小長度
    if len(secret) < minLength {
        return fmt.Errorf("JWT_SECRET must be at least %d characters long (current: %d)",
            minLength, len(secret))
    }

    // 建議最小長度為 64 字元（提供更好的安全性）
    if len(secret) < 64 {
        // 只是警告，不阻止啟動（向後相容）
        fmt.Printf("WARNING: JWT_SECRET length is %d. Recommended minimum is 64 characters for better security.\n",
            len(secret))
    }

    // 檢查字元多樣性（至少需要 3 種類型）
    hasLower := false
    hasUpper := false
    hasDigit := false
    hasSpecial := false

    for _, char := range secret {
        switch {
        case unicode.IsLower(char):
            hasLower = true
        case unicode.IsUpper(char):
            hasUpper = true
        case unicode.IsDigit(char):
            hasDigit = true
        case unicode.IsPunct(char) || unicode.IsSymbol(char):
            hasSpecial = true
        }
    }

    varietyCount := 0
    if hasLower { varietyCount++ }
    if hasUpper { varietyCount++ }
    if hasDigit { varietyCount++ }
    if hasSpecial { varietyCount++ }

    if varietyCount < 3 {
        return fmt.Errorf("JWT_SECRET must contain at least 3 different character types (lowercase, uppercase, digits, special characters)")
    }

    // 檢查熵（簡單版本：計算唯一字元比例）
    entropy := calculateSimpleEntropy(secret)
    if entropy < 0.5 { // 唯一字元少於一半
        return fmt.Errorf("JWT_SECRET has low entropy (too many repeated characters). Use a more random secret.")
    }

    // 檢查常見弱密碼模式
    if isWeakPattern(secret) {
        return fmt.Errorf("JWT_SECRET contains weak patterns (e.g., repeated sequences). Generate a truly random secret.")
    }

    return nil
}

// calculateSimpleEntropy calculates a simple entropy measure (0.0 to 1.0)
func calculateSimpleEntropy(s string) float64 {
    if len(s) == 0 {
        return 0
    }

    // 計算唯一字元比例
    unique := make(map[rune]bool)
    for _, char := range s {
        unique[char] = true
    }

    return float64(len(unique)) / float64(len(s))
}

// isWeakPattern checks for common weak patterns
func isWeakPattern(s string) bool {
    // 檢查是否包含長重複序列
    for i := 0; i < len(s)-3; i++ {
        pattern := s[i:i+4]
        count := 0
        for j := 0; j < len(s)-3; j++ {
            if s[j:j+4] == pattern {
                count++
            }
        }
        if count > 2 { // 同樣的 4 字元序列出現超過 2 次
            return true
        }
    }

    // 檢查是否為連續字元
    consecutive := 0
    for i := 1; i < len(s); i++ {
        if s[i] == s[i-1] {
            consecutive++
            if consecutive > 5 { // 超過 5 個連續相同字元
                return true
            }
        } else {
            consecutive = 0
        }
    }

    return false
}
```

**同時更新 .env.example**:

```bash
# JWT Configuration (Required)
# CRITICAL SECURITY: Generate a strong random secret using:
#   openssl rand -base64 64
#   or
#   head -c 64 /dev/urandom | base64
#
# Requirements:
# - Minimum 32 characters (64+ recommended)
# - Mix of uppercase, lowercase, digits, and special characters
# - High entropy (avoid repeated patterns)
JWT_SECRET=your-super-secret-jwt-key-at-least-32-characters-long
```

---

### 4. Login 失敗缺少 Rate Limiting

**檔案**: [internal/service/auth.go:57-71](internal/service/auth.go#L57-L71)

**問題嚴重性**: 🔴 嚴重 - 安全性漏洞

**問題描述**:

當前 Login 實作沒有任何速率限制，攻擊者可以無限次嘗試暴力破解密碼。雖然使用了 bcrypt（本身有抗暴力破解的特性），但仍建議實作 rate limiting。

**攻擊場景**:

1. 針對特定使用者進行暴力破解
2. 針對常見使用者名稱進行撞庫攻擊
3. 分散式暴力破解攻擊

**建議方案**:

#### 方案 A: 使用現有的全域 Rate Limiter（短期方案）

在 [internal/router/router.go:86-98](internal/router/router.go#L86-L98) 為 login endpoint 加上更嚴格的 rate limit：

```go
// internal/router/router.go

func (r *Router) setupAuthRoutes() {
    auth := r.engine.Group("/auth")
    {
        // Register endpoint - 標準 rate limit
        auth.POST("/register", r.authHandler.Register)

        // Login endpoint - 嚴格 rate limit（防止暴力破解）
        loginRateLimit := middleware.NewRateLimiter(middleware.RateLimitConfig{
            RequestsPerMinute: 5,    // 每分鐘最多 5 次嘗試
            BurstSize:         2,    // 允許短時間內 2 次突發
            KeyPrefix:         "login_rate_limit:",
        })
        auth.POST("/login", loginRateLimit, r.jwtMiddleware.LoginHandler)

        // Refresh endpoint - 寬鬆 rate limit
        auth.POST("/refresh", r.jwtMiddleware.RefreshHandler)
    }

    // Protected routes
    authProtected := r.engine.Group("/auth")
    authProtected.Use(r.jwtMiddleware.MiddlewareFunc())
    {
        authProtected.POST("/logout", r.jwtMiddleware.LogoutHandler)
    }
}
```

#### 方案 B: 實作基於使用者的 Login Attempt Tracking（推薦方案）

在資料庫層面追蹤登入失敗次數，實作帳號鎖定機制。

**步驟 1**: 更新 User model 加入鎖定欄位

```sql
-- migrations/XXX_add_login_attempts.sql

-- +goose Up
ALTER TABLE users ADD COLUMN failed_login_attempts INT NOT NULL DEFAULT 0;
ALTER TABLE users ADD COLUMN locked_until TIMESTAMP NULL;
ALTER TABLE users ADD COLUMN last_failed_login TIMESTAMP NULL;

CREATE INDEX idx_users_locked_until ON users(locked_until) WHERE locked_until IS NOT NULL;

-- +goose Down
DROP INDEX IF EXISTS idx_users_locked_until;
ALTER TABLE users DROP COLUMN IF EXISTS last_failed_login;
ALTER TABLE users DROP COLUMN IF EXISTS locked_until;
ALTER TABLE users DROP COLUMN IF EXISTS failed_login_attempts;
```

**步驟 2**: 更新 Repository 新增鎖定相關方法

```go
// internal/repository/user.go

type UserRepository interface {
    GetByUsername(ctx context.Context, username string) (*model.User, error)
    Create(ctx context.Context, username, passwordHash string) (int64, error)
    ExistsByUsername(ctx context.Context, username string) (bool, error)

    // 新增: 登入失敗追蹤
    IncrementFailedLoginAttempts(ctx context.Context, username string) error
    ResetFailedLoginAttempts(ctx context.Context, username string) error
    IsAccountLocked(ctx context.Context, username string) (bool, error)
}
```

**步驟 3**: 更新 Service 實作帳號鎖定邏輯

```go
// internal/service/auth.go

const (
    maxLoginAttempts = 5                  // 最大失敗次數
    lockDuration     = 15 * time.Minute   // 鎖定時長
)

func (s *authService) Login(ctx context.Context, req dto.LoginRequest) (*model.User, error) {
    // 檢查帳號是否被鎖定
    isLocked, err := s.userRepository.IsAccountLocked(ctx, req.Username)
    if err != nil {
        logger.Error("Failed to check account lock status", "username", req.Username, "error", err)
        return nil, apperror.AuthFailed()
    }

    if isLocked {
        logger.Warn("Login attempt on locked account", "username", req.Username)
        return nil, apperror.AccountLocked() // 需要新增此錯誤類型
    }

    // Get user by username
    user, err := s.userRepository.GetByUsername(ctx, req.Username)
    if err != nil {
        // 紀錄失敗嘗試（即使使用者不存在，也記錄以防止使用者列舉）
        _ = s.userRepository.IncrementFailedLoginAttempts(ctx, req.Username)
        return nil, apperror.AuthFailed()
    }

    // Verify password
    if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
        // 記錄失敗嘗試
        if incErr := s.userRepository.IncrementFailedLoginAttempts(ctx, req.Username); incErr != nil {
            logger.Error("Failed to increment login attempts", "username", req.Username, "error", incErr)
        }

        return nil, apperror.AuthFailed()
    }

    // 登入成功：重置失敗計數
    if err := s.userRepository.ResetFailedLoginAttempts(ctx, req.Username); err != nil {
        logger.Error("Failed to reset login attempts", "username", req.Username, "error", err)
        // 不阻止登入，只記錄錯誤
    }

    return user, nil
}
```

**步驟 4**: 新增 AccountLocked 錯誤類型

```go
// internal/apperror/codes.go
const (
    // ... 現有錯誤碼
    CodeAccountLocked = "ERR_ACCOUNT_LOCKED"
)

// internal/apperror/messages.go
var errorMessages = map[string]string{
    // ... 現有訊息
    CodeAccountLocked: "帳號已被鎖定，請稍後再試",
}

// internal/apperror/status.go
var httpStatusMap = map[string]int{
    // ... 現有對應
    CodeAccountLocked: http.StatusForbidden,
}

// internal/apperror/error.go
func AccountLocked() *AppError {
    return New(CodeAccountLocked)
}
```

**安全性考量**:

- ✅ 即使使用者不存在也記錄失敗嘗試（防止使用者列舉）
- ✅ 鎖定時間有限（15 分鐘），避免永久鎖定
- ✅ 登入成功後重置計數
- ✅ 記錄所有可疑登入行為

---

## 🟡 中優先級問題（建議近期修正）

### 5. JWT Claims 缺少標準欄位

**檔案**: [internal/middleware/jwt.go:79-89](internal/middleware/jwt.go#L79-L89)

**問題嚴重性**: 🟡 中等 - 功能完整性

**問題描述**:

當前 JWT claims 只包含自訂欄位，缺少 JWT 標準的 registered claims（如 `iss`, `aud`, `jti` 等），這可能導致與其他服務整合時出現問題。

```go
// 當前實作
PayloadFunc: func(data interface{}) gojwt.MapClaims {
    if user, ok := data.(*model.User); ok {
        return gojwt.MapClaims{
            "user_id":    user.ID,
            "username":   user.Username,
            "role":       user.Role,
            "experience": user.ExperiencePoints,
        }
    }
    return gojwt.MapClaims{}
},
```

**建議方案**:

```go
// internal/middleware/jwt.go

PayloadFunc: func(data interface{}) gojwt.MapClaims {
    if user, ok := data.(*model.User); ok {
        now := time.Now()

        return gojwt.MapClaims{
            // 標準 JWT claims (Registered Claims)
            "iss": "waterballsa",                    // Issuer: 發行者
            "sub": fmt.Sprintf("%d", user.ID),       // Subject: 主體（使用者 ID）
            "aud": []string{"waterballsa-api"},      // Audience: 受眾
            "exp": now.Add(cfg.AccessTokenTimeout).Unix(), // Expiration: 過期時間
            "iat": now.Unix(),                       // Issued At: 發行時間
            "jti": generateJTI(user.ID, now),        // JWT ID: 唯一識別碼

            // 自訂 claims (Private Claims)
            "user_id":    user.ID,
            "username":   user.Username,
            "role":       user.Role,
            "experience": user.ExperiencePoints,
        }
    }
    return gojwt.MapClaims{}
},

// 新增輔助函式生成唯一 JWT ID
func generateJTI(userID int64, t time.Time) string {
    // 使用 user ID + timestamp + random 生成唯一 ID
    return fmt.Sprintf("%d-%d-%s", userID, t.Unix(), randomString(8))
}

func randomString(n int) string {
    const letters = "abcdefghijklmnopqrstuvwxyz0123456789"
    b := make([]byte, n)
    for i := range b {
        b[i] = letters[rand.Intn(len(letters))]
    }
    return string(b)
}
```

**優點**:

- 符合 JWT RFC 7519 標準
- 便於與第三方服務整合
- `jti` 欄位可用於實作 token 撤銷黑名單

---

### 6. IdentityHandler 缺少錯誤處理

**檔案**: [internal/middleware/jwt.go:93-120](internal/middleware/jwt.go#L93-L120)

**問題嚴重性**: 🟡 中等 - 錯誤處理

**問題描述**:

當 JWT claims 解析失敗時，`IdentityHandler` 只是記錄 warning 並返回 `nil`，但沒有明確的錯誤回應機制。這可能導致後續的 `Authorizer` 或業務邏輯出現 panic。

```go
// 當前實作
IdentityHandler: func(c *gin.Context) interface{} {
    claims := jwt.ExtractClaims(c)

    userID, ok := claims["user_id"].(float64)
    if !ok {
        logger.Warn("Invalid user_id in JWT claims")
        return nil // 返回 nil 可能導致後續處理出錯
    }
    // ...
},
```

**建議方案**:

```go
// internal/middleware/jwt.go

IdentityHandler: func(c *gin.Context) interface{} {
    claims := jwt.ExtractClaims(c)

    // 提取並驗證 user_id
    userID, ok := claims["user_id"].(float64)
    if !ok {
        logger.Error("Invalid user_id type in JWT claims",
            "type", fmt.Sprintf("%T", claims["user_id"]))
        c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
            "code":  "ERR_INVALID_TOKEN",
            "error": "Token 格式錯誤",
        })
        return nil
    }

    // 提取並驗證 username
    username, ok := claims["username"].(string)
    if !ok || username == "" {
        logger.Error("Invalid or missing username in JWT claims")
        c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
            "code":  "ERR_INVALID_TOKEN",
            "error": "Token 格式錯誤",
        })
        return nil
    }

    // 提取並驗證 role
    role, ok := claims["role"].(string)
    if !ok || role == "" {
        logger.Error("Invalid or missing role in JWT claims")
        c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
            "code":  "ERR_INVALID_TOKEN",
            "error": "Token 格式錯誤",
        })
        return nil
    }

    // 驗證 user_id 合法性
    if userID <= 0 {
        logger.Error("Invalid user_id value in JWT claims", "user_id", userID)
        c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
            "code":  "ERR_INVALID_TOKEN",
            "error": "Token 格式錯誤",
        })
        return nil
    }

    user := &model.User{
        ID:       int64(userID),
        Username: username,
        Role:     role,
    }

    logger.Debug("Identity extracted from JWT",
        "user_id", user.ID,
        "username", user.Username,
        "role", user.Role)

    return user
},
```

**改進重點**:

- 明確驗證所有必要欄位
- 提供清晰的錯誤訊息
- 使用 `c.AbortWithStatusJSON` 立即中斷請求
- 增加 debug 日誌便於追蹤

---

### 7. LoginResponse 錯誤處理可以改進

**檔案**: [internal/middleware/jwt.go:136-167](internal/middleware/jwt.go#L136-L167)

**問題嚴重性**: 🟡 中等 - 使用者體驗

**問題描述**:

當 `authenticated_user` 從 context 取出失敗時，使用的錯誤碼是通用的 `ERR_INTERNAL_ERROR`，但實際上這是一個認證流程的邏輯錯誤，應該有更明確的錯誤碼。

**建議方案**:

```go
// internal/apperror/codes.go
const (
    // ... 現有錯誤碼
    CodeAuthStateError = "ERR_AUTH_STATE_ERROR" // 認證狀態錯誤
)

// internal/apperror/messages.go
var errorMessages = map[string]string{
    // ... 現有訊息
    CodeAuthStateError: "認證狀態異常，請重新登入",
}

// internal/apperror/status.go
var httpStatusMap = map[string]int{
    // ... 現有對應
    CodeAuthStateError: http.StatusInternalServerError,
}

// internal/apperror/error.go
func AuthStateError(err error) *AppError {
    return NewWithError(CodeAuthStateError, err)
}
```

```go
// internal/middleware/jwt.go

LoginResponse: func(c *gin.Context, token *core.Token) {
    // Get user data from context
    userData, exists := c.Get("authenticated_user")
    if !exists {
        logger.Error("Authenticated user not found in context during login response")
        _ = c.Error(apperror.AuthStateError(
            fmt.Errorf("user data missing from context after authentication")))
        return
    }

    user, ok := userData.(*model.User)
    if !ok {
        logger.Error("Invalid user data type in context",
            "expected", "*model.User",
            "got", fmt.Sprintf("%T", userData))
        _ = c.Error(apperror.AuthStateError(
            fmt.Errorf("invalid user data type: %T", userData)))
        return
    }

    logger.Info("Login successful",
        "user_id", user.ID,
        "username", user.Username,
        "token_expires_at", time.Unix(token.ExpiresAt, 0).Format(time.RFC3339))

    c.JSON(http.StatusOK, dto.LoginResponse{
        AccessToken: token.AccessToken,
        User: dto.UserInfo{
            ID:         user.ID,
            Username:   user.Username,
            Experience: user.ExperiencePoints,
        },
    })
},
```

---

### 8. 缺少 Token 過期時間在回應中

**檔案**: [internal/dto/auth.go:18-22](internal/dto/auth.go#L18-L22)

**問題嚴重性**: 🟡 中等 - 使用者體驗

**問題描述**:

前端無法得知 access token 的過期時間，無法提前刷新 token，只能等到 401 錯誤後才知道需要刷新。

```go
// 當前結構
type LoginResponse struct {
    AccessToken  string   `json:"accessToken"`
    RefreshToken string   `json:"refreshToken,omitempty"`
    User         UserInfo `json:"user"`
}
```

**建議方案**:

```go
// internal/dto/auth.go

type LoginResponse struct {
    AccessToken  string   `json:"accessToken"`
    RefreshToken string   `json:"refreshToken,omitempty"`
    ExpiresAt    int64    `json:"expiresAt"`    // Unix timestamp
    ExpiresIn    int      `json:"expiresIn"`    // Seconds until expiration
    TokenType    string   `json:"tokenType"`    // "Bearer"
    User         UserInfo `json:"user"`
}

// 可選：新增輔助方法
func NewLoginResponse(accessToken string, expiresAt int64, user *model.User) LoginResponse {
    now := time.Now().Unix()
    expiresIn := int(expiresAt - now)
    if expiresIn < 0 {
        expiresIn = 0
    }

    return LoginResponse{
        AccessToken: accessToken,
        ExpiresAt:   expiresAt,
        ExpiresIn:   expiresIn,
        TokenType:   "Bearer",
        User: UserInfo{
            ID:         user.ID,
            Username:   user.Username,
            Experience: user.ExperiencePoints,
        },
    }
}
```

更新 middleware 的 LoginResponse：

```go
// internal/middleware/jwt.go

LoginResponse: func(c *gin.Context, token *core.Token) {
    userData, exists := c.Get("authenticated_user")
    if !exists {
        // ... 錯誤處理
        return
    }

    user, ok := userData.(*model.User)
    if !ok {
        // ... 錯誤處理
        return
    }

    logger.Info("Login successful", "user_id", user.ID, "username", user.Username)

    response := dto.NewLoginResponse(token.AccessToken, token.ExpiresAt, user)
    c.JSON(http.StatusOK, response)
},
```

**前端使用範例**:

```typescript
// 前端可以提前刷新 token
const response = await fetch('/auth/login', { ... });
const data = await response.json();

// 在 token 過期前 1 分鐘自動刷新
const refreshTime = (data.expiresIn - 60) * 1000; // 轉換為毫秒
setTimeout(() => refreshToken(), refreshTime);
```

---

### 9. 測試覆蓋不足

**檔案**: [tests/bdd/features/isa/auth/login.isa.feature](tests/bdd/features/isa/auth/login.isa.feature)

**問題嚴重性**: 🟡 中等 - 測試完整性

**問題描述**:

BDD 測試檔案已建立但內容不完整。Login 功能作為核心功能，應該有完整的測試覆蓋。

**建議的測試場景**:

```gherkin
# tests/bdd/features/isa/auth/login.isa.feature

Feature: User Login
  As a registered user
  I want to login to the system
  So that I can access protected resources

  Background:
    Given the database has a user
      | username | testuser       |
      | password | SecurePass123! |

  Scenario: Successful login with valid credentials
    When I send "POST" request to "/auth/login" with JSON body:
      """
      {
        "username": "testuser",
        "password": "SecurePass123!"
      }
      """
    Then the response code should be 200
    And the response should have JSON key "accessToken"
    And the response should have JSON key "user"
    And the response JSON at "user.username" should be "testuser"
    And the response should set cookie "jwt"
    And the response should set cookie "refresh_token"

  Scenario: Login fails with wrong password
    When I send "POST" request to "/auth/login" with JSON body:
      """
      {
        "username": "testuser",
        "password": "WrongPassword"
      }
      """
    Then the response code should be 401
    And the response JSON at "code" should be "ERR_AUTH_FAILED"
    And the response should not set cookie "jwt"

  Scenario: Login fails with non-existent username
    When I send "POST" request to "/auth/login" with JSON body:
      """
      {
        "username": "nonexistent",
        "password": "SomePassword123!"
      }
      """
    Then the response code should be 401
    And the response JSON at "code" should be "ERR_AUTH_FAILED"

  Scenario: Login fails with missing username
    When I send "POST" request to "/auth/login" with JSON body:
      """
      {
        "password": "SecurePass123!"
      }
      """
    Then the response code should be 400
    And the response JSON at "code" should be "ERR_VALIDATION_FAILED"

  Scenario: Login fails with missing password
    When I send "POST" request to "/auth/login" with JSON body:
      """
      {
        "username": "testuser"
      }
      """
    Then the response code should be 400
    And the response JSON at "code" should be "ERR_VALIDATION_FAILED"

  Scenario: Login fails with empty request body
    When I send "POST" request to "/auth/login" with JSON body:
      """
      {}
      """
    Then the response code should be 400

  Scenario: Can access protected endpoint with valid token
    Given I send "POST" request to "/auth/login" with JSON body:
      """
      {
        "username": "testuser",
        "password": "SecurePass123!"
      }
      """
    And I store the response JSON at "accessToken" as "token"
    When I send "POST" request to "/auth/logout" with header "Authorization" value "Bearer {{token}}"
    Then the response code should be 200

  Scenario: Cannot access protected endpoint without token
    When I send "POST" request to "/auth/logout"
    Then the response code should be 401

  Scenario: Cannot access protected endpoint with invalid token
    When I send "POST" request to "/auth/logout" with header "Authorization" value "Bearer invalid.token.here"
    Then the response code should be 401

  # 如果實作了 Rate Limiting（問題 #4）
  Scenario: Login is rate limited after multiple failures
    When I send "POST" request to "/auth/login" with JSON body 6 times:
      """
      {
        "username": "testuser",
        "password": "WrongPassword"
      }
      """
    Then the response code should be 429
    And the response JSON at "code" should be "ERR_RATE_LIMIT_EXCEEDED"

  # 如果實作了 Account Locking（問題 #4）
  Scenario: Account is locked after multiple failed login attempts
    Given I send "POST" request to "/auth/login" with JSON body 5 times:
      """
      {
        "username": "testuser",
        "password": "WrongPassword"
      }
      """
    When I send "POST" request to "/auth/login" with JSON body:
      """
      {
        "username": "testuser",
        "password": "SecurePass123!"
      }
      """
    Then the response code should be 403
    And the response JSON at "code" should be "ERR_ACCOUNT_LOCKED"
```

**還需要新增的 Step Definitions**:

```go
// tests/bdd/steps/common/http_steps.go

// 新增支援儲存回應資料的步驟
func iStoreTheResponseJSONAtAs(ctx context.Context, jsonPath, varName string) (context.Context, error) {
    // 實作 JSON path 解析並儲存到 context
}

// 新增支援使用變數的步驟
func iSendRequestToWithHeaderValue(ctx context.Context, method, path, headerName, headerValue string) (context.Context, error) {
    // 支援 {{varName}} 語法替換變數
}

// 新增檢查 cookie 的步驟
func theResponseShouldSetCookie(ctx context.Context, cookieName string) error {
    // 檢查回應是否設定了指定的 cookie
}

func theResponseShouldNotSetCookie(ctx context.Context, cookieName string) error {
    // 檢查回應是否沒有設定指定的 cookie
}
```

---

## 🟢 低優先級建議（可選改進）

### 10. JWT 配置可以更靈活

**檔案**: [internal/config/jwt.go](internal/config/jwt.go)

**建議**: 新增更多可選配置項目

```go
// internal/config/jwt.go

type JWTConfig struct {
    Secret              []byte
    AccessTokenTimeout  time.Duration
    RefreshTokenTimeout time.Duration
    SecureCookie        bool
    CookieDomain        string

    // 新增：可選配置
    Issuer              string        // JWT issuer
    Audience            []string      // JWT audience
    EnableRefreshToken  bool          // 是否啟用 refresh token
    TokenLookupSources  []string      // Token 來源優先順序
}

func loadJWTConfig() (*JWTConfig, error) {
    // ... 現有邏輯 ...

    // Load optional issuer (default: "waterballsa")
    issuer := os.Getenv("JWT_ISSUER")
    if issuer == "" {
        issuer = "waterballsa"
    }

    // Load optional audience (default: ["waterballsa-api"])
    audience := []string{"waterballsa-api"}
    if aud := os.Getenv("JWT_AUDIENCE"); aud != "" {
        audience = strings.Split(aud, ",")
    }

    // Load optional enable refresh token (default: true)
    enableRefreshToken := true
    if val := os.Getenv("JWT_ENABLE_REFRESH_TOKEN"); val == "false" {
        enableRefreshToken = false
    }

    return &JWTConfig{
        // ... 現有欄位 ...
        Issuer:             issuer,
        Audience:           audience,
        EnableRefreshToken: enableRefreshToken,
    }, nil
}
```

---

### 11. 日誌記錄可以更結構化

**檔案**: [internal/middleware/jwt.go](internal/middleware/jwt.go)

**建議**: 統一日誌格式，增加追蹤資訊

```go
// internal/middleware/jwt.go

// 為所有認證相關日誌加上統一的 prefix
const logPrefix = "jwt_middleware"

// 在每個日誌中加入 request_id
Authenticator: func(c *gin.Context) (interface{}, error) {
    var req dto.LoginRequest
    requestID := c.GetString("request_id") // 假設已有 request ID middleware

    if err := c.ShouldBindJSON(&req); err != nil {
        logger.Warn("Login request validation failed",
            "component", logPrefix,
            "request_id", requestID,
            "error", err)
        return nil, jwt.ErrMissingLoginValues
    }

    // ... 其餘邏輯

    logger.Info("User authenticated successfully",
        "component", logPrefix,
        "request_id", requestID,
        "user_id", user.ID,
        "username", user.Username,
        "ip", c.ClientIP())

    return user, nil
},
```

---

### 12. 環境變數範例可以更詳細

**檔案**: [.env.example](.env.example)

**建議**: 增加更多註解和範例

```bash
# JWT Configuration (Required)
# ============================================================================
#
# JWT_SECRET - 用於簽名 JWT 的密鑰
#
# 安全要求:
#   - 最少 32 字元（強烈建議 64 字元以上）
#   - 必須包含大小寫字母、數字、特殊符號（至少 3 種）
#   - 避免使用重複字元或簡單模式
#   - 不要提交到版本控制系統
#
# 產生強密碼的方法:
#   Linux/macOS:  openssl rand -base64 64
#   Linux:        head -c 64 /dev/urandom | base64
#   PowerShell:   [Convert]::ToBase64String((1..64 | ForEach-Object { Get-Random -Minimum 0 -Maximum 255 }))
#   Online:       https://generate-secret.vercel.app/64
#
# ⚠️  WARNING: 更換此密鑰會使所有現有 token 失效
#
JWT_SECRET=REPLACE_WITH_GENERATED_SECRET_AT_LEAST_64_CHARACTERS

# JWT_ACCESS_TOKEN_TIMEOUT - Access Token 有效時長
# 建議值: 15m (15分鐘) 到 1h (1小時)
# 格式: 使用 Go duration 格式 (例如: 15m, 1h, 30s)
# 預設: 15m
JWT_ACCESS_TOKEN_TIMEOUT=15m

# JWT_REFRESH_TOKEN_TIMEOUT - Refresh Token 有效時長
# 建議值: 168h (7天) 到 720h (30天)
# 格式: 使用 Go duration 格式 (例如: 168h, 720h)
# 預設: 168h (7天)
JWT_REFRESH_TOKEN_TIMEOUT=168h

# JWT_SECURE_COOKIE - 是否只允許 HTTPS 傳輸 cookie
# 開發環境: false
# 生產環境: true (強烈建議)
# 預設: false
JWT_SECURE_COOKIE=false

# JWT_COOKIE_DOMAIN - Cookie 的 domain 屬性
# 範例: .example.com (允許所有子網域)
# 留空表示使用當前 domain
# 預設: (空白)
JWT_COOKIE_DOMAIN=

# JWT_ISSUER - JWT issuer claim (可選)
# 用於識別 token 的發行者
# 預設: waterballsa
# JWT_ISSUER=waterballsa

# JWT_AUDIENCE - JWT audience claim (可選)
# 用於指定 token 的預期接收者（可用逗號分隔多個）
# 預設: waterballsa-api
# JWT_AUDIENCE=waterballsa-api

# JWT_ENABLE_REFRESH_TOKEN - 是否啟用 refresh token 功能
# true: 啟用 refresh token (推薦)
# false: 停用 refresh token
# 預設: true
# JWT_ENABLE_REFRESH_TOKEN=true
```

---

## 📊 問題優先級總結

### 🔴 高優先級（建議立即處理）- 安全性關鍵

1. **Refresh Token Store 實作** - 功能完整性（必須修正）
2. **Logout 清理 Refresh Token** - 安全性漏洞
3. **JWT Secret 驗證強化** - 安全性風險
4. **Login Rate Limiting** - 防暴力破解（安全性關鍵）

### 🟡 中優先級（建議近期處理）- 功能完善

5. **JWT Claims 標準化** - 相容性與標準化
6. **IdentityHandler 錯誤處理** - 穩定性
7. **LoginResponse 錯誤處理改進** - 錯誤明確性
8. **Token 過期時間在回應中** - 使用者體驗
9. **測試覆蓋不足** - 品質保證

### 🟢 低優先級（可選改進）- 錦上添花

10. **JWT 配置彈性** - 可配置性
11. **日誌結構化** - 可觀測性
12. **環境變數文件** - 文件完整性

---

## 🎯 建議的實作順序

### 第一階段：安全性修復（1-2 天）

**必須完成的項目**:

1. **實作 Refresh Token Store**（4-6 小時）

   - 使用 memory store 快速實作
   - 或直接實作 Redis store（推薦）

2. **實作 Logout Token 清理**（2-3 小時）

   - 配合 Refresh Token Store
   - 清除 cookies

3. **強化 JWT Secret 驗證**（1-2 小時）

   - 實作 `validateJWTSecret` 函式
   - 更新 `.env.example` 文件

4. **實作 Login Rate Limiting**（3-5 小時）
   - 短期：使用全域 rate limiter
   - 長期：實作帳號鎖定機制（可分階段）

### 第二階段：功能完善（1 天）

5. **標準化 JWT Claims**（1-2 小時）
6. **改進 IdentityHandler**（1 小時）
7. **改進 LoginResponse 錯誤處理**（1 小時）
8. **新增 Token 過期資訊**（1-2 小時）
9. **補充 BDD 測試**（3-4 小時）

### 第三階段：優化與文件（半天）

10. **JWT 配置優化**（1 小時，可選）
11. **日誌優化**（1 小時，可選）
12. **文件完善**（1 小時）

---

## 📝 其他建議

### 安全性最佳實踐

1. **Secret 管理**

   - 使用環境變數或 secret 管理服務（如 AWS Secrets Manager、HashiCorp Vault）
   - 定期輪換 JWT secret（需要有遷移策略）
   - 不同環境使用不同的 secret

2. **Token 撤銷機制**

   - 考慮實作 token blacklist（使用 Redis）
   - 在關鍵操作後（如修改密碼）強制重新登入

3. **HTTPS**

   - 生產環境必須使用 HTTPS
   - 設定 `JWT_SECURE_COOKIE=true`

4. **監控與告警**
   - 監控異常登入行為（異地登入、短時間多次失敗等）
   - 記錄所有認證事件用於安全審計

### 文件建議

建議在 [CLAUDE.md](CLAUDE.md) 中新增 Authentication 章節：

```markdown
## Authentication & Authorization

### JWT Authentication

The application uses JWT (JSON Web Token) for authentication with a dual-token system:

- **Access Token**: Short-lived token (default: 15 minutes) for API access
- **Refresh Token**: Long-lived token (default: 7 days) for obtaining new access tokens

### Login Flow

1. User sends credentials to `POST /auth/login`
2. Server validates credentials and generates tokens
3. Tokens are sent both in response body and secure HTTP-only cookies
4. Client can use either Authorization header or cookie for subsequent requests

### Token Refresh

When access token expires, use `POST /auth/refresh` with refresh token to obtain a new access token.

### Logout

Call `POST /auth/logout` (requires authentication) to invalidate tokens and clear cookies.

### Security Features

- bcrypt password hashing
- Rate limiting on login endpoint
- Account locking after failed attempts
- Secure cookie configuration (HTTPOnly, SameSite)
- Token storage in Redis for revocation support

### Configuration

See [.env.example](.env.example) for JWT configuration options.
```

### 效能考量

1. **Redis 連線池**

   - 如果使用 Redis 作為 refresh token store，需要配置適當的連線池大小
   - 監控 Redis 的記憶體使用

2. **Bcrypt Cost**
   - 當前預設 cost 為 10（在 register 中），對於大多數場景已足夠
   - 如果效能成為瓶頸，考慮使用 goroutine 非同步處理部分邏輯

### 相容性考量

1. **跨域 (CORS)**

   - 確保 CORS 設定允許 credentials（cookie）傳遞
   - 檢查 `Access-Control-Allow-Credentials` header

2. **移動端**
   - 移動應用可能無法使用 cookie，確保 Authorization header 方式可用
   - 考慮提供 refresh token 在 response body 中（當前已有 `omitempty`）

---

## ✅ 結論

整體而言，Login 功能的實作展現了良好的架構設計和安全意識。主要優點包括：

- ✅ 使用業界標準的 JWT 認證
- ✅ 實作了 Refresh Token 機制
- ✅ 密碼使用 bcrypt 安全儲存
- ✅ 統一的錯誤處理機制
- ✅ 清晰的程式碼結構

**關鍵改進點**:

1. **必須修正**: Refresh Token Store 實作（當前功能不完整）
2. **安全性**: 實作 Rate Limiting 和帳號鎖定機制
3. **安全性**: 強化 JWT Secret 驗證
4. **安全性**: Logout 時清理 Refresh Token
5. **功能性**: 補充完整的測試覆蓋

建議優先處理第一階段的安全性修復（預估 1-2 天），這些是影響系統安全的關鍵問題。其他改進可以根據專案時程逐步實施。

完成這些改進後，Login 功能將達到生產環境的品質標準。
