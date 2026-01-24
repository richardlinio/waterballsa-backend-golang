# Reference Guide: Feature Scaffolding

本參考指南補充 SKILL.md 的執行步驟,提供實戰範例與速查表。詳細的架構說明請參考 `/CLAUDE.md`。

## 專案結構查詢

**策略:** 使用分層查詢避免 token 浪費

### 快速導覽 (優先使用)

```bash
# 僅顯示主要目錄結構 (1-2 層)
tree -L 2 -d -I 'vendor|node_modules|.git|tmp' --dirsfirst

# 查看特定模組的檔案
tree internal/handler -I '*_test.go'
tree internal/service -I '*_test.go'
tree internal/repository -I '*_test.go'
```

### 常見路徑速查

| 需求       | 路徑                        | 說明            |
| ---------- | --------------------------- | --------------- |
| HTTP 處理  | `internal/handler/`         | Gin 路由處理器  |
| 業務邏輯   | `internal/service/`         | 服務層實作      |
| 資料存取   | `internal/repository/`      | Repository 層   |
| 資料結構   | `internal/dto/`             | API 請求/回應   |
| 領域模型   | `internal/model/`           | 資料庫實體      |
| SQL 查詢   | `internal/db/queries/`      | SQLc 查詢定義   |
| 錯誤定義   | `internal/apperror/`        | 統一錯誤處理    |
| 路由註冊   | `internal/router/router.go` | 路由設定        |
| DI 配置    | `internal/app/app.go`       | 依賴注入        |
| 資料庫遷移 | `migrations/`               | Goose migration |
| BDD 測試   | `tests/features/`           | Gherkin 功能檔  |

### 完整結構 (僅必要時使用)

```bash
# ⚠️ 會產生大量輸出，僅在需要全貌時使用
tree -L 3 -I 'vendor|node_modules|.git|tmp' --dirsfirst
```

## 常見實作模式

### 回傳值設計 (Go Idiomatic)

❌ **避免多個回傳值**:

```go
// 不佳: 太多回傳值，難以維護
func DeliverMission(ctx context.Context, userID, missionID int64) (message string, expGained int, totalExp, currentLevel int32, err error)
```

✅ **使用 struct 封裝**:

```go
// Domain Model (internal/model/)
type MissionDeliveryResult struct {
    Message          string
    ExperienceGained int
    TotalExperience  int32
    CurrentLevel     int32
}

// Service 方法
func (s *Service) DeliverMission(ctx context.Context, userID, missionID int64) (*model.MissionDeliveryResult, error) {
    // ...
    return &model.MissionDeliveryResult{
        Message:          "成功",
        ExperienceGained: reward.Value,
        TotalExperience:  newExp,
        CurrentLevel:     newLevel,
    }, nil
}

// DTO Converter (internal/dto/)
func ToDeliverResponse(result *model.MissionDeliveryResult) DeliverResponse {
    return DeliverResponse{
        Message:          result.Message,
        ExperienceGained: result.ExperienceGained,
        TotalExperience:  result.TotalExperience,
        CurrentLevel:     result.CurrentLevel,
    }
}
```

**原則**:

- Service 回傳 domain model (業務結果)
- DTO converter 負責轉換為 API 回應格式
- 超過 2 個資料欄位時，優先考慮 struct

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
// Repository: 轉換為 domain error (隱藏 pgx 實作細節)
if errors.Is(err, pgx.ErrNoRows) {
    return nil, repository.ErrResourceNotFound
}

// Service: 檢查 repository 的自訂錯誤 (不依賴 pgx)
if errors.Is(err, repository.ErrResourceNotFound) {
    return nil, apperror.ResourceNotFound()
}
```

**重要**: Service 層絕對不應該檢查 `pgx.ErrNoRows`，必須檢查 repository 定義的錯誤，維持分層架構原則。

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

### 生產環境 (`internal/app/app.go`)

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

### 測試環境 (`tests/testutil/server.go`)

**重要:** 新增功能時必須同步更新此檔案，否則測試編譯會失敗！

```go
// 1. Repository 初始化 (第 93-99 行附近)
fooRepository := repository.NewFooRepository(queries)

// 2. Service 初始化 (第 104-108 行附近)
fooService := service.NewFooService(fooRepository)

// 3. Handler 初始化 (第 122-126 行附近)
fooHandler := handler.NewFooHandler(fooService, log, cfg.Server.RequestTimeout)

// 4. Router Setup (第 135 行附近)
r := router.NewRouter(
    ginEngine,
    cfg.CORS,
    cfg.RateLimit,
    cfg.JWT,
    log,
    healthHandler,
    authHandler,
    journeyHandler,
    missionHandler,
    progressHandler,
    fooHandler, // 新增此參數
    jwtMiddleware,
    blacklistChecker,
)
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
9. **依賴注入** (依賴所有元件)
   - `internal/app/app.go` - 生產環境
   - `tests/testutil/server.go` - 測試環境 ⚠️ **必須同步更新**

必檢項目:

- [ ] 所有層級已建立
- [ ] 錯誤處理已加入
- [ ] 授權檢查已實作 (如需要)
- [ ] 路由已註冊
- [ ] 依賴注入已配置 (`internal/app/app.go`)
- [ ] **測試伺服器已更新** (`tests/testutil/server.go`) ⚠️
- [ ] Migration 已建立 (如需要)
- [ ] SQLc 已執行 (`make sqlc`)
- [ ] 測試步驟已定義

## 常見陷阱

❌ **避免:**

- Service 方法回傳太多值 (超過 2 個資料欄位應使用 struct)
- **Service 層檢查 `pgx.ErrNoRows`** (違反分層架構)
- 跳過授權檢查
- 忘記註冊路由
- 遺漏依賴注入
- **忘記更新測試伺服器** (`tests/testutil/server.go`)
- 使用縮寫命名 (repo, svc, cfg)
- 忘記執行 `make sqlc`
- 忽略檢查現有 migrations
- 忽略檢查現有 step definitions

✅ **務必:**

- Service 回傳多個資料時使用 domain model struct
- **Service 檢查 repository 自訂錯誤，而非 pgx 錯誤**
- 驗證使用者權限
- 透過建構子注入依賴
- **同步更新生產與測試環境的 DI 配置**
- Repository 將 `pgx.ErrNoRows` 轉為 domain error
- 使用完整變數名稱 (repository, service, config)
- 新增 queries 後執行 `make sqlc`
- 確認資料表結構一致
- 重用現有測試步驟
