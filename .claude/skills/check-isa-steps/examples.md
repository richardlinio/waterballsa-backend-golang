# Report Examples

本文件包含分析報告的輸出格式範例。

## 完整報告格式

````markdown
## 📋 ISA Feature 步驟分析報告

**檔案**: [檔案路徑](檔案路徑)

### ✅ 已實作的步驟 (X 個)

| 步驟                                                   | 位置                                | 函式                            |
| ------------------------------------------------------ | ----------------------------------- | ------------------------------- |
| Given the database has a user with username "testuser" | tests/bdd/steps/database/user.go:54 | theDatabaseHasAUserWithUsername |
| When I send "GET" request to "/api/users"              | tests/bdd/steps/http/request.go:23  | iSendRequestTo                  |
| Then the response status code should be 200            | tests/bdd/steps/http/response.go:15 | theResponseStatusCodeShouldBe   |

### ❌ 缺少的步驟 (Y 個)

#### 步驟 1: "Given the database has 5 journeys"

**建議實作位置**: tests/bdd/steps/database/journey.go

**原因**: 這是 database setup 步驟,用於設定多筆 journey 測試資料,應該放在 database/journey.go

**建議實作**:

```go
// theDatabaseHasJourneys creates multiple journeys in the database
func theDatabaseHasJourneys(ctx context.Context, count int) (context.Context, error) {
    // 實作邏輯:
    // 1. 從 context 取得 database pool
    // 2. 使用迴圈插入 count 筆 journeys
    // 3. 儲存 journey IDs 到 context 供後續步驟使用
    return ctx, nil
}

// 在 RegisterJourneySteps 中註冊:
func RegisterJourneySteps(sc *godog.ScenarioContext) {
    // ... 其他步驟註冊
    sc.Step(`^the database has (\d+) journeys$`, theDatabaseHasJourneys)
}
```
````

**註冊狀態**: `database.RegisterJourneySteps(sc)` 已在 register.go 中註冊 ✓

---

#### 步驟 2: "And I store the response field "id" as "journeyId""

**建議實作位置**: tests/bdd/steps/http/variables.go

**原因**: 這是變數儲存步驟,用於從 response 中提取欄位並存為變數,應該放在 http/variables.go

**建議實作**:

```go
// iStoreTheResponseFieldAs extracts a field from response and stores it as a variable
func iStoreTheResponseFieldAs(ctx context.Context, fieldPath, variableName string) (context.Context, error) {
    // 實作邏輯:
    // 1. 從 context 取得最後的 response
    // 2. 使用 fieldPath 提取 JSON 欄位值
    // 3. 儲存到 context 的變數 map 中
    return ctx, nil
}

// 在 RegisterVariableSteps 中註冊:
func RegisterVariableSteps(sc *godog.ScenarioContext) {
    // ... 其他步驟註冊
    sc.Step(`^I store the response field "([^"]*)" as "([^"]*)"$`, iStoreTheResponseFieldAs)
}
```

**需要在 register.go 中加入**:

```go
http.RegisterVariableSteps(sc)
```

---

### 📊 統計

- 總步驟數: 15
- 已實作: 10 (67%)
- 需要實作: 5 (33%)

### 💡 下一步行動

1. 實作缺少的步驟 (依照上述建議)
2. 在對應的 `Register*Steps` 函式中註冊新步驟
3. 確保 `register.go` 中已呼叫所有必要的 `Register*Steps` 函式
4. 執行測試驗證步驟是否正確實作

````

## 簡化版報告 (當所有步驟都已實作時)

```markdown
## 📋 ISA Feature 步驟分析報告

**檔案**: [tests/bdd/features/isa/auth/login.isa.feature](tests/bdd/features/isa/auth/login.isa.feature)

### ✅ 所有步驟都已實作! (8 個)

| 步驟 | 位置 | 函式 |
|------|------|------|
| Given the database has a user with username "johndoe" | tests/bdd/steps/database/user.go:54 | theDatabaseHasAUserWithUsername |
| When I send "POST" request to "/api/auth/login" | tests/bdd/steps/http/request.go:23 | iSendRequestTo |
| Then the response status code should be 200 | tests/bdd/steps/http/response.go:15 | theResponseStatusCodeShouldBe |
| And the response body should contain "accessToken" | tests/bdd/steps/http/response.go:42 | theResponseBodyShouldContain |
| And cookie "refreshToken" should be set | tests/bdd/steps/http/cookies.go:18 | cookieShouldBeSet |

### 📊 統計

- 總步驟數: 8
- 已實作: 8 (100%)
- 需要實作: 0 (0%)

### 🎉 測試已準備就緒

此 feature 檔案的所有步驟都已實作,可以執行測試。
````
