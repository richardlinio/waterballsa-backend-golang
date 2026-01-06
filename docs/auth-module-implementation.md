# Auth Module 實作計畫

> **目標**：實作基於 gin-jwt 的認證模組，支援登入/refresh token 機制
> **時程**：預計 3-4 小時
> **策略**：使用 sqlc 生成 type-safe repository code，採用 gin-jwt 標準雙 token 設計

---

## 📋 需求確認

### 功能範圍

- ✅ **登入功能**：帳號密碼認證
- ✅ **Token 刷新**：使用 refresh token 換取新的 access token
- ❌ **登出功能**：不在此次實作範圍
- ❌ **註冊功能**：不在此次實作範圍

### 技術決策

- ✅ 使用 **sqlc** 生成資料存取層程式碼（不手寫 repository）
- ✅ 使用 **gin-jwt v3** 標準雙 token 機制
- ✅ 採用 **cookie-based** 認證（httpOnly cookies）
- ✅ Role 使用**單一 enum**（STUDENT/TEACHER/ADMIN）
- ✅ **單機部署**（refresh token 存 memory store）

---

## 🗂️ 專案結構

### 採用 Layer-based 架構

保持現有的分層架構模式：

```
internal/
├── db/                      # sqlc 生成的資料存取層
│   ├── models.go           # 自動生成：User struct
│   ├── queries.sql.go      # 自動生成：type-safe query functions
│   └── db.go               # 自動生成：DBTX interface
│
├── db/queries/             # SQL query 定義（手寫）
│   └── users.sql           # User 相關的 SQL queries
│
├── dto/                    # Data Transfer Objects（手寫）
│   └── auth.go             # LoginRequest, LoginResponse
│
├── service/                # 業務邏輯層（手寫）
│   └── auth_service.go     # 認證邏輯（密碼驗證）
│
├── handler/                # HTTP 處理層（手寫）
│   └── auth_handler.go     # 認證 HTTP handlers
│
└── middleware/             # 中間件（手寫）
    └── auth.go             # JWT 中間件設定
```

### 為什麼使用 sqlc？

- ✅ **Type-safe**：編譯時檢查 SQL 語法錯誤
- ✅ **程式碼生成**：自動生成 Go structs 和 query functions
- ✅ **維護性高**：修改 SQL 後重新生成即可
- ✅ **效能好**：直接使用 `pgx` driver，無 ORM overhead
- ✅ **Go 社群推薦**：適合中小型專案，學習曲線低

---

## 📝 實作步驟

### Phase 1: 整合 sqlc

**1.1 安裝 sqlc**

```bash
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
```

**1.2 建立 sqlc 配置檔**

在專案根目錄建立 `sqlc.yaml`：

```yaml
version: '2'
sql:
  - engine: 'postgresql'
    queries: 'internal/db/queries'
    schema: 'migrations'
    gen:
      go:
        package: 'db'
        out: 'internal/db'
        sql_package: 'pgx/v5'
        emit_json_tags: true
        emit_interface: true
        emit_empty_slices: true
```

**1.3 建立 SQL queries 檔案**

建立 `internal/db/queries/users.sql`：

```sql
-- name: GetUserByID :one
SELECT * FROM users
WHERE id = $1 AND deleted_at IS NULL;

-- name: GetUserByUsername :one
SELECT * FROM users
WHERE username = $1 AND deleted_at IS NULL;
```

**1.4 執行 sqlc 生成程式碼**

```bash
sqlc generate
```

這會自動生成：

- `internal/db/models.go` - User struct（對應 DB schema）
- `internal/db/queries.sql.go` - GetUserByID, GetUserByUsername functions
- `internal/db/db.go` - DBTX interface, Queries struct

---

### Phase 2: 安裝依賴與基礎設定

**2.1 安裝 Go 套件**

```bash
go get github.com/appleboy/gin-jwt/v3
go get golang.org/x/crypto/bcrypt
```

**2.2 更新環境變數**

在 `.env.example` 新增：

```env
# JWT Configuration
JWT_SECRET=<使用 openssl rand -base64 32 生成>
JWT_REALM=WaterBall SA
JWT_TIMEOUT=1h
JWT_REFRESH_TIMEOUT=168h  # 7 days
```

**2.3 更新 Config**

在 `internal/config/config.go` 新增 JWT 配置結構：

```go
type JWTConfig struct {
    Secret         string
    Realm          string
    Timeout        time.Duration
    RefreshTimeout time.Duration
}
```

---

### Phase 3: 建立 DTO

建立 `internal/dto/auth.go`，定義 API 請求/回應格式：

- `LoginRequest` - 登入請求（username, password）
- `LoginResponse` - 登入回應（access_token, token_type, expires_in, refresh_token, user）
- `UserInfo` - 使用者資訊（id, username, experience）

---

### Phase 4: 建立認證服務

建立 `internal/service/auth_service.go`，實作業務邏輯：

- `VerifyPassword(hashedPassword, plainPassword string) error` - 驗證密碼
- `HashPassword(plainPassword string) (string, error)` - bcrypt hash（未來註冊用）

---

### Phase 5: 建立 JWT 中間件 ⭐

**這是最複雜的部分**

建立 `internal/middleware/auth.go`，配置 gin-jwt：

**5.1 核心組件**

- **Authenticator**: 登入驗證邏輯

  - 呼叫 `sqlc` 生成的 `GetUserByUsername` 查詢使用者
  - 使用 `auth_service.VerifyPassword` 驗證密碼
  - 返回 User 物件

- **PayloadFunc**: JWT Claims 設定

  - 只存 `user_id` 到 JWT（其他資訊會變動）

- **IdentityHandler**: 從 JWT 還原使用者

  - 從 claims 取得 `user_id`
  - 呼叫 `GetUserByID` 查詢完整 User 資訊

- **Authorizer**: 權限控制
  - 目前允許所有認證使用者
  - 未來可依 role 或其他條件控制

**5.2 自訂回應格式**

- **LoginResponse**: 登入成功回應
- **RefreshResponse**: Token 刷新成功回應

**5.3 Cookie 設定**

```go
SendCookie:            true,
SecureCookie:          false,  // dev 環境，production 設為 true
CookieHTTPOnly:        true,   // 防 XSS
CookieSameSite:        http.SameSiteStrictMode,  // 防 CSRF
CookieName:            "jwt",
RefreshTokenCookieName: "refresh_token",
```

---

### Phase 6: 建立 HTTP Handlers

建立 `internal/handler/auth_handler.go`：

封裝 gin-jwt 提供的 handlers：

- `LoginHandler` - 處理 POST /auth/login
- `RefreshHandler` - 處理 POST /auth/refresh

提供統一的錯誤處理和日誌記錄。

---

### Phase 7: 註冊路由

更新 `internal/router/router.go`，註冊認證相關路由：

```go
// Public routes
r.POST("/auth/login", authHandler.LoginHandler)
r.POST("/auth/refresh", authHandler.RefreshHandler)
```

---

### Phase 8: 更新 API 文檔

**重要**：API 文檔需要更新以反映新的認證機制

**8.1 新增 `/auth/refresh` 端點**

更新 `docs/api-docs/swagger.yaml`：

```yaml
/auth/refresh:
  $ref: './openapi/paths/auth.yaml#/refresh'
```

**8.2 更新 LoginResponse schema**

更新 `docs/api-docs/openapi/schemas/auth.yaml`：

- 新增 `refresh_token` 欄位
- 新增 `token_type` 欄位（固定為 "Bearer"）
- 新增 `expires_in` 欄位（秒數）
- 修改 `accessToken` → `access_token`（符合 OAuth 2.0 標準）

**8.3 新增 RefreshResponse schema**

定義 refresh endpoint 的回應格式。

**8.4 更新 JWT 說明**

修改 `securitySchemes.bearerAuth` 描述：

- Access Token 有效期：1 天 → **1 小時**
- 新增 Refresh Token 說明（7 天有效期）
- 移除 token 黑名單描述

---

### Phase 9: 驗證 Migration

確認 `migrations/001_create_users_table.sql` 正確：

- ✅ user_role ENUM (STUDENT, TEACHER, ADMIN)
- ✅ users 表欄位完整
- ❌ **不需要** access_tokens 表（移除或忽略 `002_create_access_tokens_table.sql`）

---

## 🔑 關鍵設計決策

### JWT Token 策略

採用 **gin-jwt 標準雙 token 設計**：

| 項目         | Access Token                           | Refresh Token             |
| ------------ | -------------------------------------- | ------------------------- |
| **有效期**   | 1 小時                                 | 7 天                      |
| **傳遞方式** | httpOnly cookie + Authorization header | httpOnly cookie（僅）     |
| **Claims**   | user_id, exp, iat                      | -                         |
| **類型**     | JWT (有簽名)                           | Opaque token (不透明字串) |
| **儲存**     | 客戶端（cookie/localStorage）          | Server-side memory store  |
| **用途**     | API 認證                               | 換取新的 access token     |

### JWT Claims 設計

**只存 `user_id`**，其他資訊需要時再查 DB：

```json
{
	"user_id": 123,
	"exp": 1234567890,
	"iat": 1234567890
}
```

**為什麼不存 role、level、experience_points？**

- `level`, `experience_points` **經常變動**
- 如果存在 JWT，會導致 token 中的資料與 DB 不同步
- 需要時用 `user_id` 查 DB 即可（`IdentityHandler` 會處理）


### Role 設計

- **單一 role enum**：每個 user 只有一個 role
- **三種角色**：STUDENT, TEACHER, ADMIN
- **DB 欄位**：`role user_role NOT NULL DEFAULT 'STUDENT'`
- **Go type**：sqlc 會自動生成對應的 enum type

**未來擴充**：如需支援多個 roles，需要：

1. 建立 `user_roles` 關聯表
2. 修改 JWT claims 結構
3. 修改 Authorizer 邏輯

---

## 📦 檔案清單

### 新增檔案（8 個）

**配置檔**

1. `sqlc.yaml` - sqlc 配置檔

**SQL Queries（手寫）** 2. `internal/db/queries/users.sql` - User SQL queries

**自動生成（by sqlc）** 3. `internal/db/models.go` - User struct 4. `internal/db/queries.sql.go` - Query functions 5. `internal/db/db.go` - DBTX interface

**手寫程式碼** 6. `internal/dto/auth.go` - 認證 DTO 7. `internal/service/auth_service.go` - 認證業務邏輯 8. `internal/handler/auth_handler.go` - 認證 HTTP handlers 9. `internal/middleware/auth.go` - JWT 中間件

### 修改檔案（7 個）

1. `go.mod` - 新增依賴（gin-jwt, bcrypt）
2. `internal/config/config.go` - 新增 JWT 配置
3. `.env.example` - 新增環境變數範例
4. `internal/router/router.go` - 註冊認證路由
5. `docs/api-docs/openapi/paths/auth.yaml` - 新增 refresh 端點
6. `docs/api-docs/openapi/schemas/auth.yaml` - 更新回傳格式
7. `docs/api-docs/swagger.yaml` - 更新 JWT 說明

### 確認檔案（1 個）

1. `migrations/001_create_users_table.sql` - 確認 users 表正確

---

## 🧪 測試計畫

### 手動測試流程

**測試前準備**：先用 SQL 建立測試帳號

```sql
INSERT INTO users (username, password_hash, role)
VALUES ('testuser', '<bcrypt_hash>', 'STUDENT');
```

**1. 登入測試**

```bash
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","password":"password123"}' \
  -c cookies.txt -v
```

驗證：

- ✅ HTTP 200 OK
- ✅ 回傳包含 `access_token`, `refresh_token`, `token_type`, `expires_in`, `user`
- ✅ Cookie 中有 `jwt` 和 `refresh_token`（httpOnly）

**2. Token 刷新測試**

```bash
curl -X POST http://localhost:8080/auth/refresh \
  -b cookies.txt \
  -c cookies.txt -v
```

驗證：

- ✅ HTTP 200 OK
- ✅ 取得新的 access_token 和 refresh_token
- ✅ 舊的 refresh_token 失效

**3. 訪問保護路由測試**（未來業務 API 實作後）

```bash
curl http://localhost:8080/journeys \
  -b cookies.txt -v
```

驗證：

- ✅ 能正確驗證 JWT token
- ✅ 能取得使用者資訊（透過 IdentityHandler）

---

## ⏱️ 估計工作量

| Phase   | 任務                                          | 預估時間  |
| ------- | --------------------------------------------- | --------- |
| Phase 1 | 整合 sqlc（安裝、配置、撰寫 SQL、生成程式碼） | ~40 分鐘  |
| Phase 2 | 安裝依賴與基礎設定                            | ~15 分鐘  |
| Phase 3 | 建立 DTO                                      | ~15 分鐘  |
| Phase 4 | 建立認證服務                                  | ~20 分鐘  |
| Phase 5 | 建立 JWT 中間件 ⭐                            | ~1.5 小時 |
| Phase 6 | 建立 HTTP Handlers                            | ~20 分鐘  |
| Phase 7 | 註冊路由                                      | ~10 分鐘  |
| Phase 8 | 更新 API 文檔                                 | ~30 分鐘  |
| Phase 9 | 驗證 Migration                                | ~10 分鐘  |

**總計**：~4 小時

---

## 📚 參考資料

### 官方文檔

- [sqlc Documentation](https://docs.sqlc.dev/)
- [sqlc Getting Started](https://docs.sqlc.dev/en/latest/tutorials/getting-started-postgresql.html)
- [gin-jwt v3 Documentation](https://github.com/appleboy/gin-jwt)
- [gin-jwt Basic Example](https://github.com/appleboy/gin-jwt/tree/master/_example/basic)

### 技術標準

- [RFC 6749 - OAuth 2.0](https://tools.ietf.org/html/rfc6749) - Refresh Token 標準
- [RFC 6750 - OAuth 2.0 Bearer Token Usage](https://tools.ietf.org/html/rfc6750)

### Go 套件

- [Go bcrypt Package](https://pkg.go.dev/golang.org/x/crypto/bcrypt)
- [pgx v5](https://github.com/jackc/pgx) - PostgreSQL driver

---

## 💡 常見問題

### Q1: 為什麼不用 GORM？

sqlc 適合中小型專案：

- ✅ 更輕量、學習曲線低
- ✅ 完全控制 SQL（效能更好）
- ✅ Type-safe，編譯時檢查錯誤
- ❌ 不支援複雜的 relation mapping（但我們專案用不到）

### Q2: Access Token 1 小時會不會太短？

- ✅ 這是業界標準做法（OAuth 2.0 推薦）
- ✅ 配合 refresh token，使用者體驗不受影響
- ✅ 降低 token 洩漏風險
- ✅ 可以用 refresh token 自動續期

### Q3: 為什麼不在 JWT 存 role？

- ❌ Role 可能會變動（例如 STUDENT → TEACHER）
- ❌ JWT 在過期前無法修改
- ✅ 只存 user_id，需要時查 DB（確保資料是最新的）

### Q4: Migration 中有 access_tokens 表怎麼辦？

- 可以保留（未來可能用到）
- 或移除（目前用不到）
- 不影響此次實作

---

## ✅ 完成檢查清單

實作完成後確認：

- [ ] sqlc 成功生成程式碼（`internal/db/` 目錄有檔案）
- [ ] 環境變數設定正確（JWT_SECRET 已生成）
- [ ] API 文檔已更新（新增 `/auth/refresh`）
- [ ] 所有 3 個手動測試通過
- [ ] Cookie 設定正確（httpOnly, SameSite）
- [ ] Migration 確認正確
- [ ] 程式碼可以編譯（`go build ./...`）
- [ ] 使用 bcrypt 加密密碼（cost = 10-12）
