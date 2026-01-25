---
name: go-idiomatic
description: Go 語言慣用寫法與設計模式指南，包含回傳值設計、Aggregate Structures、Slice 初始化、命名慣例等。
allowed-tools: Read
---

# Go Idiomatic 寫法指南

本 skill 提供 Go 語言慣用寫法與專案特定的設計模式。

## 回傳值設計 (Go Idiomatic)

### ❌ 避免多個回傳值

```go
// 不佳: 太多回傳值，難以維護
func DeliverMission(ctx context.Context, userID, missionID int64) (message string, expGained int, totalExp, currentLevel int32, err error)
```

### ✅ 使用 struct 封裝

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

## Aggregate Structures 設計原則

專案區分不同類型的 aggregate structures，依據其用途與所在層級。

### Model Layer

包含兩種結構：

1. **Database Entity**: 單一資料表映射
   - 範例: `Journey`, `Mission`, `Order`, `User`
   - 代表資料庫中的單一 row/document
   - 可包含實體特定的業務方法 (例: `User.CalculateLevel()`)

2. **Query Result Aggregate**: 多表查詢結果容器
   - 範例: `JourneyDetail`, `MissionDetail`
   - 純資料容器，無業務邏輯
   - 由 Service 層組裝，但定義在 Model 層
   - 目的: 避免 Service 方法回傳太多參數
   - 識別: 組合多個資料表的資料用於讀取操作

### Service Layer

包含 **Business Result Structs**: 業務操作的結果

- 範例: `OrderResult`, `LoginResult`
- 包含業務邏輯產生的資料 (反正規化資料、計算欄位)
- 命名慣例: `*Result` suffix
- 識別: 包含業務邏輯產生的資料，而非純 DB 查詢
- **絕對不可** 依賴 DTO layer (避免 `service → dto` 依賴)

### DTO Layer

包含 **Request/Response DTOs**: API 序列化/反序列化

- 範例: `OrderResponse`, `LoginRequest`, `UserInfo`
- 僅處理 JSON tags、validation tags、API 格式化
- **絕對不可** 被 Service layer 引用
- 可以依賴 Model layer 進行轉換

### Dependency Rules

```
Handler → DTO ← (converts from) ← Service → Model
   ↓                                   ↓
(uses)                              (uses)
   ↓                                   ↓
Service                            Repository
```

- ✅ Handler 可依賴 Service 和 DTO
- ✅ Service 可依賴 Model 和 Repository
- ✅ DTO 可依賴 Model (用於轉換函數)
- ❌ Service **絕對不可** 依賴 DTO (防止循環依賴與層級混淆)

### 正確模式範例

```go
// service/auth.go
type LoginResult struct {
    AccessToken string
    UserID      int64      // ✅ 獨立欄位
    Username    string     // ✅ 無 dto.UserInfo 依賴
    Experience  int32
}

// handler/auth.go
result, _ := authService.Login(ctx, req)
c.JSON(http.StatusOK, dto.LoginResponse{
    AccessToken: result.AccessToken,
    User: dto.UserInfo{  // ✅ Handler 從 Result 建構 DTO
        ID:         result.UserID,
        Username:   result.Username,
        Experience: result.Experience,
    },
})
```

### 反模式 (Anti-pattern)

```go
// ❌ 錯誤: Service 依賴 DTO
type LoginResult struct {
    AccessToken string
    UserInfo    dto.UserInfo  // ❌ 建立 service → dto 依賴
}
```

## Slice 初始化風格

```go
// 使用 append 模式 (推薦)
items := make([]T, 0, len(source))
for _, item := range source {
    items = append(items, T{...})
}
```

**原則**:
- 使用 `make([]T, 0, capacity)` 預分配容量
- 透過 `append` 逐一加入元素
- 避免使用索引賦值 `items[i] = ...`

## 命名慣例

**使用完整名稱而非縮寫** 以提升可讀性與一致性。

### 變數命名

```go
// 推薦 ✓
authService
userRepository
courseRepository
serverConfig
databaseConfig

// 避免 ✗
authSvc
userRepo
courseRepo
srvCfg
dbCfg
```

### 原則

- 變數名稱應具描述性
- 避免縮寫，除非是業界標準 (如 ID, URL, HTTP)
- 保持整個專案的一致性

## 設定管理模式

**Pass-by-value pattern** 表達設定的不可變性：

```go
// Loader functions 回傳值而非指標
func loadXxxConfig() (XxxConfig, error)

// Config struct 儲存值
type Config struct {
    Xxx XxxConfig
}

// Functions 接收值
func NewHandler(cfg config.XxxConfig)
```

**原則**:
- 設定載入後不應被修改
- 透過值傳遞表達不可變性
- 整個專案保持一致的模式

## Repository 層錯誤處理

**重要**: Service 層不應檢查資料庫特定錯誤 (如 `pgx.ErrNoRows`)。

### 正確模式

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

### 反模式

```go
// ❌ Service 層不應檢查 pgx.ErrNoRows
if errors.Is(err, pgx.ErrNoRows) {
    return nil, apperror.ResourceNotFound()
}
```

**原則**:
- Repository 負責將資料庫特定錯誤轉換為 domain errors
- Service 僅檢查 repository 定義的錯誤
- 維持分層架構原則，Service 不應知道底層資料庫實作

## 總結檢查清單

- [ ] Service 回傳多個資料時使用 domain model struct
- [ ] Service 不依賴 DTO layer
- [ ] Service 檢查 repository 自訂錯誤，而非 pgx 錯誤
- [ ] 使用完整變數名稱 (repository, service, config)
- [ ] Slice 使用 append 模式初始化
- [ ] Config 使用 pass-by-value pattern
- [ ] Repository 將資料庫錯誤轉為 domain error
