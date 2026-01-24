# Reference Guide: Feature Scaffolding

本參考指南補充 SKILL.md 的執行步驟,提供實戰範例與速查表。詳細的架構說明請參考 `/CLAUDE.md`。

## 專案結構查詢

**重要:** 專案結構會動態變化,使用以下指令獲取最新結構:

```bash
tree -L 3 -I 'vendor|node_modules|.git|tmp' --dirsfirst
```

## 常見實作模式

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
// Repository: 直接回傳資料庫層錯誤
if errors.Is(err, pgx.ErrNoRows) {
    return nil, err
}

// Service: 轉換為業務層錯誤
if errors.Is(err, pgx.ErrNoRows) {
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

## SQLc 命名慣例

- `GetXxxByYyy` - 單筆查詢
- `ListXxxByYyy` - 多筆查詢
- `CreateXxx` - 新增
- `UpdateXxx` - 更新
- `UpsertXxx` - 新增或更新
- `DeleteXxx` - 刪除

查詢後執行: `make sqlc`

## 路由註冊範例

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

## 依賴注入範例

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

## ISA Feature 測試範例

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

## Scaffold 檢查清單

實作順序 (依賴關係):

1. **Domain Models** (無依賴) - `internal/model/`
2. **DTOs** (依賴 models) - `internal/dto/`
3. **SQLc Queries** (無依賴) - `internal/db/queries/*.sql`
4. **錯誤碼** (可提前) - `internal/apperror/`
5. **Repository** (依賴 SQLc) - `internal/repository/`
6. **Service** (依賴 repository) - `internal/service/`
7. **Handler** (依賴 service) - `internal/handler/`
8. **路由註冊** (依賴 handler) - `internal/router/router.go`
9. **依賴注入** (依賴所有元件) - `internal/app/app.go`

必檢項目:

- [ ] 所有層級已建立
- [ ] 錯誤處理已加入
- [ ] 授權檢查已實作 (如需要)
- [ ] 路由已註冊
- [ ] 依賴注入已配置
- [ ] Migration 已建立 (如需要)
- [ ] SQLc 已執行 (`make sqlc`)
- [ ] 測試步驟已定義

## 常見陷阱

❌ **避免:**

- 跳過授權檢查
- 忘記註冊路由
- 遺漏依賴注入
- 使用縮寫命名 (repo, svc, cfg)
- 忘記執行 `make sqlc`
- 忽略檢查現有 migrations
- 忽略檢查現有 step definitions

✅ **務必:**

- 驗證使用者權限
- 透過建構子注入依賴
- 優雅處理 `pgx.ErrNoRows`
- 使用完整變數名稱 (repository, service, config)
- 新增 queries 後執行 `make sqlc`
- 確認資料表結構一致
- 重用現有測試步驟
