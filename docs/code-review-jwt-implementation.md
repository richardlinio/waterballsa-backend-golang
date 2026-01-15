# JWT 實作程式碼審查報告

**審查日期**: 2026-01-15
**審查範圍**: JWT 認證系統實作 (feat/login 分支)
**審查者**: Claude Code

---

## 概述

本次審查針對 JWT 認證系統的實作進行完整檢視，發現多項需要修正的問題。這些問題涵蓋安全性、一致性、可測試性與可維護性等面向。

---

## 高優先級問題


---

### 4. Cookie 名稱硬編碼重複

**位置**: 多個檔案
- [internal/handler/auth.go:106](../internal/handler/auth.go#L106) - `setCookie`
- [internal/handler/auth.go:117](../internal/handler/auth.go#L117) - `clearCookie`
- [internal/handler/auth.go:141](../internal/handler/auth.go#L141) - `extractToken`
- [internal/infrastructure/auth/jwt.go:36](../internal/infrastructure/auth/jwt.go#L36) - middleware 配置

**問題描述**:
```go
// handler/auth.go
c.SetCookie("jwt", ...)      // 位置 1
c.Cookie("jwt")              // 位置 2
c.SetCookie("jwt", ...)      // 位置 3

// auth/jwt.go
CookieName: "jwt",           // 位置 4
```

**問題分析**:
- 字串 `"jwt"` 在 4 個地方重複出現
- 如果需要修改 cookie 名稱，必須更新所有位置
- 容易遺漏導致不一致

**建議修正**:
```go
// 在 config/jwt.go 中新增
type JWTConfig struct {
    Secret              []byte
    AccessTokenTimeout  time.Duration
    RefreshTokenTimeout time.Duration
    SecureCookie        bool
    CookieDomain        string
    CookieName          string        // 新增：可配置的 cookie 名稱
}

func loadJWTConfig() (*JWTConfig, error) {
    // ... 其他配置 ...

    // 載入可選的 cookie 名稱 (預設: "jwt")
    cookieName := os.Getenv("JWT_COOKIE_NAME")
    if cookieName == "" {
        cookieName = "jwt"
    }

    return &JWTConfig{
        // ... 其他欄位 ...
        CookieName: cookieName,
    }, nil
}

// 在 handler 中使用
func (h *AuthHandler) setCookie(c *gin.Context, token string, expire time.Time) {
    c.SetCookie(
        h.jwtConfig.CookieName,  // 使用配置中的名稱
        token,
        // ... 其他參數 ...
    )
}

// 在 jwt.go 中使用
CookieName: cfg.CookieName,
```

**影響等級**: 🟠 中 - 影響可維護性，容易出錯

---

### 5. Refresh 端點為公開但可能需要 Token

**位置**: [internal/router/router.go:90](../internal/router/router.go#L90)

**問題描述**:
```go
// Public routes
auth := r.engine.Group("/auth")
{
    auth.POST("/register", r.authHandler.Register)
    auth.POST("/login", r.authHandler.Login)
    auth.POST("/refresh", r.authHandler.Refresh)  // 公開路由
}
```

**問題分析**:
- `/auth/refresh` 位於公開路由組
- 但 refresh token 機制通常需要驗證現有的 token (可能已過期但未超過 MaxRefresh)
- 需要確認 `gin-jwt` 的 `RefreshHandler` 行為是否符合預期

**建議修正**:

**方案 A: 確認 gin-jwt 行為後保持公開**
```go
// 如果 gin-jwt 的 RefreshHandler 內部會驗證 token，則保持公開
// 但需要在文件中說明：
// - refresh 端點需要在 header 或 cookie 中提供有效的 refresh token
// - refresh token 在 MaxRefresh 時間內有效 (預設 7 天)

auth.POST("/refresh", r.authHandler.Refresh)  // 公開但需要 refresh token
```

**方案 B: 使用自定義 middleware 保護**
```go
// 如果需要更明確的保護，可以使用 gin-jwt 的 middleware
auth.POST("/refresh",
    middleware.JWTAuth(r.jwtMiddleware),  // 驗證 token (可過期)
    r.authHandler.Refresh,
)
```

**方案 C: 檢查並記錄**
```go
// 在 handler 中明確記錄 refresh 的行為
func (h *AuthHandler) Refresh(c *gin.Context) {
    // gin-jwt RefreshHandler 會：
    // 1. 驗證 token (允許過期但未超過 MaxRefresh)
    // 2. 生成新的 access token
    // 3. 返回新 token
    h.jwtMiddleware.RefreshHandler(c)
}
```

**影響等級**: 🟠 中 - 可能影響安全性，需要確認行為

---

## 低優先級問題

### 6. RefreshToken 欄位從未使用

**位置**: [internal/dto/auth.go:17](../internal/dto/auth.go#L17)

**問題描述**:
```go
type LoginResponse struct {
    AccessToken  string   `json:"accessToken"`
    RefreshToken string   `json:"refreshToken,omitempty"` // 宣告但從未設定
    User         UserInfo `json:"user"`
}
```

**問題分析**:
- DTO 中定義了 `RefreshToken` 欄位
- 但在 handler 中從未設定此欄位
- 客戶端永遠收不到 refresh token (只能從 cookie 取得)

**建議修正方案**:

**方案 A: 移除未使用的欄位**
```go
type LoginResponse struct {
    AccessToken string   `json:"accessToken"`
    User        UserInfo `json:"user"`
    // 移除 RefreshToken 欄位
}
```

**方案 B: 實際填充此欄位**
```go
// 在 handler 中同時返回 refresh token (給不支援 cookie 的客戶端)
func (h *AuthHandler) Login(c *gin.Context) {
    // ... 登入邏輯 ...

    // 生成 refresh token
    refreshToken, refreshExpire, err := h.tokenGenerator.GenerateRefresh(result.UserInfo)
    if err != nil {
        _ = c.Error(apperror.InternalError(err))
        return
    }

    h.setCookie(c, result.Token, result.Expire)

    c.JSON(http.StatusOK, dto.LoginResponse{
        AccessToken:  result.Token,
        RefreshToken: refreshToken,  // 填充此欄位
        User:         result.UserInfo,
    })
}
```

**影響等級**: 🟡 低 - 僅影響程式碼清晰度

---

### 7. Magic Number 在 Token 提取中

**位置**: [internal/handler/auth.go:136-138](../internal/handler/auth.go#L136-L138)

**問題描述**:
```go
if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
    return authHeader[7:], nil
}
```

**問題分析**:
- 硬編碼數字 `7` (Bearer 加空格的長度)
- 不易理解且容易出錯

**建議修正**:
```go
const bearerPrefix = "Bearer "

func (h *AuthHandler) extractToken(c *gin.Context) (string, error) {
    authHeader := c.GetHeader("Authorization")
    if authHeader != "" {
        // 使用 strings.TrimPrefix 更清晰
        if strings.HasPrefix(authHeader, bearerPrefix) {
            return strings.TrimPrefix(authHeader, bearerPrefix), nil
        }
    }

    // Try cookie
    token, err := c.Cookie(h.jwtConfig.CookieName)
    if err != nil {
        return "", fmt.Errorf("token not found in header or cookie: %w", err)
    }

    return token, nil
}
```

**影響等級**: 🟡 低 - 影響程式碼可讀性

---

### 8. 身份提取失敗時無日誌記錄

**位置**: [internal/middleware/auth.go:18-38](../internal/middleware/auth.go#L18-L38)

**問題描述**:
```go
func ExtractIdentity(c *gin.Context) interface{} {
    claims := jwt.ExtractClaims(c)

    userID, ok := claims["user_id"].(float64)
    if !ok {
        return nil  // 靜默失敗，無日誌
    }

    username, ok := claims["username"].(string)
    if !ok {
        return nil  // 靜默失敗，無日誌
    }

    // ...
}
```

**問題分析**:
- 當 claims 格式不正確時靜默失敗
- 難以除錯，不知道為什麼驗證失敗

**建議修正**:
```go
func ExtractIdentity(c *gin.Context) interface{} {
    claims := jwt.ExtractClaims(c)

    userID, ok := claims["user_id"].(float64)
    if !ok {
        // 記錄警告，包含實際收到的類型
        log := slog.With("claims", claims)
        log.Warn("failed to extract user_id from claims",
            "type", fmt.Sprintf("%T", claims["user_id"]),
        )
        return nil
    }

    username, ok := claims["username"].(string)
    if !ok {
        log := slog.With("claims", claims)
        log.Warn("failed to extract username from claims",
            "type", fmt.Sprintf("%T", claims["username"]),
        )
        return nil
    }

    // ... 其他欄位類似處理
}
```

**注意**: 需要在 middleware 中注入 logger

**影響等級**: 🟡 低 - 影響除錯效率

---

### 9. Time 未注入導致測試困難

**位置**: [internal/infrastructure/auth/token_generator.go:27](../internal/infrastructure/auth/token_generator.go#L27)

**問題描述**:
```go
func (tg *jwtTokenGenerator) Generate(user *model.User) (string, time.Time, error) {
    expireTime := time.Now().Add(tg.config.AccessTokenTimeout)  // 直接使用 time.Now()
    // ...
}
```

**問題分析**:
- 直接使用 `time.Now()` 使單元測試難以驗證過期時間
- 無法在測試中控制時間

**建議修正**:
```go
// 定義時間介面
type TimeProvider interface {
    Now() time.Time
}

type realTimeProvider struct{}

func (rtp *realTimeProvider) Now() time.Time {
    return time.Now()
}

// 修改 TokenGenerator
type jwtTokenGenerator struct {
    config       config.JWTConfig
    timeProvider TimeProvider  // 新增
}

func NewTokenGenerator(config config.JWTConfig) TokenGenerator {
    return &jwtTokenGenerator{
        config:       config,
        timeProvider: &realTimeProvider{},  // 預設使用真實時間
    }
}

// 提供測試用建構函式
func NewTokenGeneratorWithTime(config config.JWTConfig, tp TimeProvider) TokenGenerator {
    return &jwtTokenGenerator{
        config:       config,
        timeProvider: tp,
    }
}

func (tg *jwtTokenGenerator) Generate(user *model.User) (string, time.Time, error) {
    expireTime := tg.timeProvider.Now().Add(tg.config.AccessTokenTimeout)
    // ...
}

// 測試中使用
type mockTimeProvider struct {
    fixedTime time.Time
}

func (mtp *mockTimeProvider) Now() time.Time {
    return mtp.fixedTime
}

func TestTokenGenerator_Generate(t *testing.T) {
    mockTime := mockTimeProvider{fixedTime: time.Date(2026, 1, 15, 12, 0, 0, 0, time.UTC)}
    generator := NewTokenGeneratorWithTime(cfg, &mockTime)
    // ... 測試邏輯
}
```

**影響等級**: 🟡 低 - 影響可測試性

---

### 10. SameSite Cookie 模式不可配置

**位置**: [internal/infrastructure/auth/jwt.go:34](../internal/infrastructure/auth/jwt.go#L34)

**問題描述**:
```go
CookieSameSite: http.SameSiteLaxMode,  // 硬編碼
```

**問題分析**:
- SameSite 屬性影響 CSRF 防護與跨域行為
- `Lax` 模式在某些場景下可能不適用：
  - `Strict`: 更嚴格的 CSRF 防護
  - `None`: 跨域請求需要 (需配合 Secure=true)

**建議修正**:
```go
// 在 config/jwt.go 中新增
type JWTConfig struct {
    // ... 現有欄位 ...
    CookieSameSite string  // "lax", "strict", "none"
}

func loadJWTConfig() (*JWTConfig, error) {
    // ... 其他配置 ...

    // 載入 SameSite 設定 (預設: lax)
    sameSite := os.Getenv("JWT_COOKIE_SAMESITE")
    if sameSite == "" {
        sameSite = "lax"
    }

    return &JWTConfig{
        // ... 其他欄位 ...
        CookieSameSite: sameSite,
    }, nil
}

// 在 jwt.go 中使用
func NewJWTMiddleware(cfg config.JWTConfig, ...) (*jwt.GinJWTMiddleware, error) {
    sameSiteMode := parseSameSiteMode(cfg.CookieSameSite)

    return jwt.New(&jwt.GinJWTMiddleware{
        // ... 其他配置 ...
        CookieSameSite: sameSiteMode,
        // ... 其他配置 ...
    })
}

func parseSameSiteMode(mode string) http.SameSite {
    switch strings.ToLower(mode) {
    case "strict":
        return http.SameSiteStrictMode
    case "none":
        return http.SameSiteNoneMode
    case "lax":
        fallthrough
    default:
        return http.SameSiteLaxMode
    }
}
```

**影響等級**: 🟡 低 - 影響靈活性，特定場景下可能需要

---

### 11. Authorize 函式未使用 Context 參數

**位置**: [internal/middleware/auth.go:43](../internal/middleware/auth.go#L43)

**問題描述**:
```go
func Authorize(_ *gin.Context, data interface{}) bool {
    if _, ok := data.(*model.User); ok {
        return true
    }
    return false
}
```

**問題分析**:
- Context 參數被忽略 (`_`)
- 未來如果需要基於路由、方法或其他 context 資訊進行權限檢查，需要重構

**建議修正**:

**方案 A: 保留參數以便未來擴展**
```go
func Authorize(c *gin.Context, data interface{}) bool {
    user, ok := data.(*model.User)
    if !ok {
        return false
    }

    // 未來可以在這裡新增基於路由的權限檢查
    // 例如：
    // requiredRole := c.GetString("required_role")
    // return user.Role == requiredRole

    return true
}
```

**方案 B: 實作基本的角色檢查**
```go
// 在路由中設定需要的角色
func (r *Router) setupProtectedRoutes() {
    admin := r.engine.Group("/admin")
    admin.Use(middleware.JWTAuth(r.jwtMiddleware))
    admin.Use(middleware.RequireRole("admin"))  // 新增角色檢查
    {
        admin.GET("/users", r.adminHandler.ListUsers)
    }
}

// 在 Authorize 中檢查角色
func Authorize(c *gin.Context, data interface{}) bool {
    user, ok := data.(*model.User)
    if !ok {
        return false
    }

    // 檢查是否需要特定角色
    requiredRole, exists := c.Get("required_role")
    if exists {
        return user.Role == requiredRole.(string)
    }

    // 預設只要是有效使用者就通過
    return true
}
```

**影響等級**: 🟡 低 - 目前功能正常，但限制未來擴展性

---

## 修正優先級總結

| 優先級 | 問題 | 建議工時 | 影響範圍 |
|--------|------|----------|----------|
| 🔴 高 | 1. 登出功能為空實作 | 4-8 小時 | 安全性、使用者體驗 |
| 🟠 中 | 2. 雙重 Token 生成來源不一致 | 3-4 小時 | 資料一致性 |
| 🟠 中 | 3. ExtractIdentity 未提取 Experience | 0.5 小時 | 資料完整性 |
| 🟠 中 | 4. Cookie 名稱硬編碼重複 | 1 小時 | 可維護性 |
| 🟠 中 | 5. Refresh 端點權限確認 | 1 小時 | 安全性 |
| 🟡 低 | 6. RefreshToken 欄位未使用 | 0.25 小時 | 程式碼清晰度 |
| 🟡 低 | 7. Magic Number | 0.25 小時 | 可讀性 |
| 🟡 低 | 8. 身份提取失敗無日誌 | 0.5 小時 | 可除錯性 |
| 🟡 低 | 9. Time 未注入 | 1.5 小時 | 可測試性 |
| 🟡 低 | 10. SameSite 不可配置 | 0.5 小時 | 靈活性 |
| 🟡 低 | 11. Authorize 未使用 Context | 0.5 小時 | 未來擴展性 |

---

## 建議修正順序

### 第一階段 (必須修正)
1. **修正登出功能** - 選擇實作方案並執行
2. **修正 ExtractIdentity** - 快速修正，避免資料不一致

### 第二階段 (強烈建議)
3. **統一 Token 生成機制** - 避免未來出現難以追蹤的 bug
4. **修正 Cookie 名稱硬編碼** - 提升可維護性
5. **確認 Refresh 端點行為** - 確保安全性

### 第三階段 (改進品質)
6. **移除未使用的 RefreshToken 欄位**
7. **修正 Magic Number**
8. **新增身份提取日誌**

### 第四階段 (可選，提升可測試性)
9. **注入 Time Provider**
10. **使 SameSite 可配置**
11. **改善 Authorize 函式**

---

## 額外建議

### 1. 新增整合測試
針對完整的認證流程撰寫測試：
```go
func TestAuthFlow(t *testing.T) {
    // 1. 註冊
    // 2. 登入 -> 取得 token
    // 3. 使用 token 存取受保護端點
    // 4. Refresh token
    // 5. 登出
    // 6. 確認 token 失效
}
```

### 2. 更新 CLAUDE.md 文件
在專案文件中新增 JWT 認證系統的說明：
- Token 生命週期
- Cookie vs Header 使用時機
- Refresh 機制說明
- 登出行為說明

### 3. 新增環境變數範例
在 `.env.example` 中補充新的配置選項：
```bash
# JWT Cookie Configuration (Advanced)
# JWT_COOKIE_NAME=jwt                  # Optional: Cookie name (Default: jwt)
# JWT_COOKIE_SAMESITE=lax              # Optional: lax/strict/none (Default: lax)
```

---

## 結論

此次 JWT 實作整體架構良好，遵循了層次化架構原則。主要問題集中在：
1. **安全性認知** - 登出功能需要明確實作或文件說明
2. **一致性** - Token 生成機制需統一
3. **可維護性** - 減少硬編碼，提升配置彈性

建議優先處理高優先級與中優先級問題，低優先級問題可在後續重構時逐步改進。
