# 常見實作模式 (Implementation Patterns)

本文件提供在專案中常見的 Go 程式碼實作模式。

## Slice 初始化風格

```go
// 使用 append 模式 (推薦)
items := make([]T, 0, len(source))
for _, item := range source {
    items = append(items, T{...})
}
```

**詳細說明**: 參考 `/go-idiomatic` skill

## 回傳值設計

**原則**: Service 方法若需回傳 2 個以上資料欄位，應建立 domain model struct

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
    // ... 業務邏輯
    return &model.MissionDeliveryResult{...}, nil
}
```

**詳細說明與反模式**: 參考 `/go-idiomatic` skill

## 授權檢查

```go
user := c.MustGet("JWT_PAYLOAD").(*model.User)
requestedUserID, _ := strconv.ParseInt(c.Param("userId"), 10, 64)

if user.ID != requestedUserID {
    _ = c.Error(apperror.UnauthorizedAccess())
    return
}
```

## 處理 Not Found

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

## UPSERT 模式

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
