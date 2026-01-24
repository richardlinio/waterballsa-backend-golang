# Reference Guide: ISA Feature Analysis

本參考指南補充 SKILL.md 的執行步驟，提供額外的模式範例與注意事項。

## 專案結構

```
waterballsa-backend-golang/
├── cmd/server/              # 程式入口
├── internal/
│   ├── handler/             # HTTP handlers (presentation)
│   ├── service/             # Business logic
│   ├── repository/          # Data access
│   ├── dto/                 # Request/Response DTOs
│   ├── model/               # Domain models
│   ├── db/queries/          # SQLc 查詢定義
│   ├── apperror/            # 錯誤定義
│   ├── middleware/          # HTTP middleware
│   ├── router/              # 路由註冊
│   └── app/                 # 依賴注入
├── migrations/              # Goose migrations
├── tests/bdd/
│   ├── features/isa/        # ISA feature 檔案
│   └── steps/               # Step definitions
└── docs/
    ├── api-docs/swagger.yaml
    └── db-schema.dbml
```

## 分層架構職責

### Handler (Presentation)
- 提取路徑/查詢參數
- 綁定 JSON 請求
- 驗證使用者授權
- 呼叫 service
- 轉換為 DTO
- 回傳 HTTP 回應

### Service (Business Logic)
- 驗證業務規則
- 協調多個 repository 呼叫
- 計算衍生值
- 狀態轉換邏輯
- 用 AppError 包裝錯誤

### Repository (Data Access)
- 執行 SQLc queries
- 轉換 sqlc models 為 domain models
- 處理資料庫錯誤
- 管理 transactions

## 常見模式

### 授權檢查

```go
user := c.MustGet("JWT_PAYLOAD").(*model.User)
requestedUserID, _ := strconv.ParseInt(c.Param("userId"), 10, 64)

if user.ID != requestedUserID {
    _ = c.Error(apperror.UnauthorizedAccess())
    return
}
```

### 處理 Not Found

```go
// Repository: 回傳 nil, nil (不是錯誤)
if errors.Is(err, pgx.ErrNoRows) {
    return nil, nil
}

// Service: 決定如何處理
if result == nil {
    return nil, apperror.ResourceNotFound()
}
```

### UPSERT 模式

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

## 錯誤處理流程

1. **定義錯誤碼** (`internal/apperror/codes.go`):
   ```go
   const CodeResourceNotFound = "ERR_RESOURCE_NOT_FOUND"
   ```

2. **定義訊息** (`internal/apperror/messages.go`):
   ```go
   var errorMessages = map[string]string{
       CodeResourceNotFound: "找不到指定的資源",
   }
   ```

3. **對應 HTTP 狀態** (`internal/apperror/status.go`):
   ```go
   var httpStatusMap = map[string]int{
       CodeResourceNotFound: http.StatusNotFound,
   }
   ```

4. **建立建構子** (`internal/apperror/error.go`):
   ```go
   func ResourceNotFound() *AppError {
       return New(CodeResourceNotFound)
   }
   ```

## SQLc 命名慣例

- `GetXxxByYyy` - 單筆查詢
- `ListXxxByYyy` - 多筆查詢
- `CreateXxx` - 新增
- `UpdateXxx` - 更新
- `UpsertXxx` - 新增或更新
- `DeleteXxx` - 刪除

查詢後執行: `make sqlc`

## 路由註冊

在 `internal/router/router.go` 新增:

```go
func Setup(
    engine *gin.Engine,
    cfg config.Config,
    fooHandler *handler.FooHandler, // 新增參數
    pool *pgxpool.Pool,
    logger *slog.Logger,
) *gin.Engine {
    // Protected routes
    protected := engine.Group("/")
    protected.Use(middleware.JWTAuth(cfg.JWT))
    {
        protected.GET("/foo/:id", fooHandler.GetFoo)
    }
    return engine
}
```

## 依賴注入

在 `internal/app/app.go` 配置:

```go
// Repositories
fooRepository := repository.NewFooRepository(pool)

// Services
fooService := service.NewFooService(fooRepository, logger, cfg.Server.RequestTimeout)

// Handlers
fooHandler := handler.NewFooHandler(fooService, logger, cfg.Server.RequestTimeout)

// Router
engine := router.Setup(engine, cfg, fooHandler, pool, logger)
```

## ISA Feature 範例結構

```gherkin
@isa
Feature: 功能名稱

  Scenario: 場景描述
    # 資料準備
    Given the database has a user:
      | username | Alice |
      | password | Test1234! |

    # 認證
    When I send "POST" request to "/auth/login"
    And I store the response field "accessToken" as "accessToken"
    Given I set Authorization header to "{{accessToken}}"

    # 執行動作
    When I send "GET" request to "/users/1/progress"

    # 驗證回應
    Then the response status code should be 200
    And the response body field "status" should equal string "UNCOMPLETED"
```

## 命名慣例

### ✅ 正確 (完整名稱)

```go
progressRepository := repository.NewProgressRepository(pool)
authenticationService := service.NewAuthenticationService(repo)
userMissionProgress := &model.UserMissionProgress{}
```

### ❌ 錯誤 (縮寫)

```go
progressRepo := repository.NewProgressRepository(pool)
authSvc := service.NewAuthenticationService(repo)
ump := &model.UserMissionProgress{}
```

## 設定模式

### ✅ 正確 (Pass-by-Value)

```go
func NewHandler(cfg config.ServerConfig) *Handler {
    return &Handler{config: cfg}
}

type Config struct {
    Server ServerConfig // Value
}
```

### ❌ 錯誤 (Pass-by-Pointer)

```go
func NewHandler(cfg *config.ServerConfig) *Handler {
    return &Handler{config: cfg}
}

type Config struct {
    Server *ServerConfig // Pointer
}
```

## 常見指令

```bash
make fmt           # 格式化程式碼
make lint          # 執行 linter
make build         # 編譯
make test          # 執行測試
make sqlc          # 生成 SQLc 程式碼
make migrate-up    # 執行 migration
make migrate-status # 檢查 migration 狀態
```

## 常見陷阱

❌ **避免:**
- 跳過授權檢查
- 忘記註冊路由
- 遺漏依賴注入
- 使用縮寫命名
- 建立 pointer configs
- 忘記執行 `make sqlc`
- 忽略檢查現有 migrations
- 忽略檢查現有 step definitions

✅ **務必:**
- 驗證使用者權限
- 透過建構子注入依賴
- 優雅處理 `pgx.ErrNoRows`
- 使用完整變數名稱
- Config 使用 pass-by-value
- 新增 queries 後執行 `make sqlc`
- 確認資料表結構一致
- 重用現有測試步驟

## 文件參考

- **API 規格**: `/docs/api-docs/swagger.yaml`
- **資料庫結構**: `/docs/db-schema.dbml`
- **專案指南**: `/CLAUDE.md`
- **ISA Features**: `/tests/bdd/features/isa/`
