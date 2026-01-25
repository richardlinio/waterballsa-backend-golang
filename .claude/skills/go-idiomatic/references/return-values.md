# 回傳值設計 (Go Idiomatic)

當 Service 方法需要回傳多個資料欄位時，應使用 struct 進行封裝，以提高程式碼的可讀性與可維護性。

### ❌ 避免多個回傳值 (Anti-Pattern)

```go
// 不佳: 太多回傳值，難以維護
func DeliverMission(ctx context.Context, userID, missionID int64) (message string, expGained int, totalExp, currentLevel int32, err error)
```

### ✅ 使用 struct 封裝 (Correct Pattern)

將回傳的資料封裝在一個業務結果 (Business Result) 的 struct 中。

**1. 定義 Domain Model (internal/model/)**

```go
type MissionDeliveryResult struct {
    Message          string
    ExperienceGained int
    TotalExperience  int32
    CurrentLevel     int32
}
```

**2. Service 方法回傳 Model**

```go
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
```

**3. DTO Converter 轉換為 API 回應 (internal/dto/)**

```go
func ToDeliverResponse(result *model.MissionDeliveryResult) DeliverResponse {
    return DeliverResponse{
        Message:          result.Message,
        ExperienceGained: result.ExperienceGained,
        TotalExperience:  result.TotalExperience,
        CurrentLevel:     result.CurrentLevel,
    }
}
```

#### 原則總結

- Service 的職責是回傳代表業務結果的 domain model。
- DTO converter 的職責是將 domain model 轉換為 API 的回應格式。
- 當回傳的資料欄位超過 2 個時，優先考慮使用 struct 封裝。
