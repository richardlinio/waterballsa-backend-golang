# Step Definition Reference

快速查詢：匹配規則、命名慣例

## 匹配規則

### Entity 關鍵字 → Database 檔案

| 步驟包含關鍵字                     | 檔案路徑                                     |
| ---------------------------------- | -------------------------------------------- |
| user                               | tests/bdd/steps/database/user.go             |
| journey                            | tests/bdd/steps/database/journey.go          |
| chapter                            | tests/bdd/steps/database/chapter.go          |
| mission                            | tests/bdd/steps/database/mission.go          |
| reward                             | tests/bdd/steps/database/reward.go           |
| mission_resource, mission resource | tests/bdd/steps/database/mission_resource.go |

### HTTP 動作關鍵字 → HTTP 檔案

| 步驟包含關鍵字                                  | 檔案路徑                          |
| ----------------------------------------------- | --------------------------------- |
| I send, I set request, I set Authorization      | tests/bdd/steps/http/request.go   |
| response status, response body, response should | tests/bdd/steps/http/response.go  |
| cookie                                          | tests/bdd/steps/http/cookies.go   |
| I store, I extract, {{variable}}                | tests/bdd/steps/http/variables.go |

### 多候選檔案優先順序

1. 主要 entity 檔案優先（例如 "chapter has missions" → `chapter.go` 優先）
2. 找不到再檢查關聯 entity（`mission.go`）
3. 都找不到標記為缺少

## 命名慣例

### 函式命名

| 步驟類型              | 命名模式                      | 範例                              |
| --------------------- | ----------------------------- | --------------------------------- |
| Database setup        | `theDatabaseHas{Entity}`      | `theDatabaseHasJourneys`          |
| Database setup (複數) | `theDatabaseHas{Entities}`    | `theDatabaseHasJourneys`          |
| HTTP request          | `iSend{Action}`, `iSet{What}` | `iSendRequest`, `iSetRequestBody` |
| HTTP response         | `theResponse{Check}`          | `theResponseStatusCodeShouldBe`   |
| Variable              | `iStore{What}As{Name}`        | `iStoreResponseFieldAs`           |

**規則**: camelCase，開頭小寫

### 正則表達式模式

| 參數類型  | 正則模式    | 範例步驟                                  | 匹配範例 |
| --------- | ----------- | ----------------------------------------- | -------- |
| 字串      | `"([^"]*)"` | `I send "GET" request`                    | `"GET"`  |
| 數字      | `(\d+)`     | `database has (\d+) journeys`             | `5`      |
| 開頭/結尾 | `^...$`     | `^I send "([^"]*)" request to "([^"]*)"$` | 完整匹配 |

### 步驟註冊範例

```go
// 在對應的 Register*Steps 函式中
sc.Step(`^the database has (\d+) journeys$`, theDatabaseHasJourneys)
sc.Step(`^I send "([^"]*)" request to "([^"]*)"$`, iSendRequestTo)
sc.Step(`^the response status code should be (\d+)$`, theResponseStatusCodeShouldBe)
```

### register.go 註冊範例

如果建立了新的 `Register*Steps` 函式，需在 `tests/bdd/steps/register.go` 中呼叫：

```go
func RegisterSteps(sc *godog.ScenarioContext) {
    database.RegisterNewEntitySteps(sc)  // 新增這行
    // ... 其他註冊
}
```

## Database Schema 參考

實作步驟時可參考：`/Users/linporu/Documents/world-of-code/waterballsa-project/docs/db-schema.dbml`
