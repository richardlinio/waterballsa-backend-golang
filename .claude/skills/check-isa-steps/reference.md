# Step Definition Reference

本文件說明如何動態分析步驟分類和命名慣例。

## 🎯 Token 最佳化查找策略

### 階段 1: 建立步驟組索引 (只讀 1 個檔案)

1. **讀取 `tests/bdd/steps/register.go`**
   - 提取所有 `Register*Steps(sc)` 呼叫
   - 解析出包名 (database/http) 和函式名稱

   範例:

   ```go
   database.RegisterUserSteps(sc)      → database/user.go
   database.RegisterJourneySteps(sc)   → database/journey.go
   http.RegisterRequestSteps(sc)       → http/request.go
   http.RegisterResponseSteps(sc)      → http/response.go
   ```

2. **建立檔案路徑映射**
   ```
   database.RegisterUserSteps      → tests/bdd/steps/database/user.go
   database.RegisterJourneySteps   → tests/bdd/steps/database/journey.go
   database.RegisterChapterSteps   → tests/bdd/steps/database/chapter.go
   database.RegisterMissionSteps   → tests/bdd/steps/database/mission.go
   database.RegisterRewardSteps    → tests/bdd/steps/database/reward.go
   database.RegisterMissionResourceSteps → tests/bdd/steps/database/mission_resource.go
   http.RegisterRequestSteps       → tests/bdd/steps/http/request.go
   http.RegisterResponseSteps      → tests/bdd/steps/http/response.go
   http.RegisterCookieSteps        → tests/bdd/steps/http/cookies.go
   http.RegisterVariableSteps      → tests/bdd/steps/http/variables.go
   ```

### 階段 2: 智慧匹配步驟到檔案 (Progressive Disclosure)

針對 feature 中的每個步驟,使用以下規則判斷最可能的檔案:

#### 規則 1: Entity 關鍵字匹配 (Database 層)

步驟中包含 entity 名稱 → 對應的 database/{entity}.go

```
"Given the database has a user"        → database/user.go
"Given there are 5 journeys"           → database/journey.go
"And the chapter has 3 missions"       → database/mission.go 或 database/chapter.go
"Given the mission has a reward"       → database/reward.go 或 database/mission.go
```

**判斷邏輯**:

- 提取步驟中的名詞 (user, journey, chapter, mission, reward, mission_resource)
- 匹配到 `database.Register{Entity}Steps` 中的 entity 名稱
- 優先選擇最明確的 entity (例如 "user" → user.go)

#### 規則 2: HTTP 動作關鍵字匹配 (HTTP 層)

```
"When I send"                → http/request.go
"Given I set request body"   → http/request.go
"Then the response status"   → http/response.go
"And the response body"      → http/response.go
"And cookie ... should"      → http/cookies.go
"And I store"                → http/variables.go
"{{variableName}}" 語法      → http/variables.go
```

**判斷邏輯**:

- 偵測關鍵動詞和名詞組合
- Request 層: `I send`, `I set request`, `I set Authorization`
- Response 層: `response status`, `response body`, `response should`
- Cookie 層: `cookie`
- Variable 層: `I store`, `I extract`, 雙大括號語法

#### 規則 3: 當有多個候選檔案時

如果步驟可能對應多個檔案 (例如 "Given the chapter has missions"):

1. 優先檢查主要 entity 的檔案 (chapter.go)
2. 如果找不到,再檢查關聯 entity 的檔案 (mission.go)
3. 都找不到才標記為缺少

### 階段 3: 精準讀取並匹配 (只讀必要的檔案)

1. **只讀取候選檔案** (從階段 2 推斷出的 1-2 個檔案)
2. **使用 Grep 搜尋 `sc.Step(` 模式**
3. **提取正則表達式並匹配**
4. **記錄結果** (已實作/缺少)

## 命名慣例 (從現有程式碼推斷)

### 函式命名模式

- 使用 camelCase,開頭小寫
- 直接描述步驟動作
- 常見模式:
  - `theDatabaseHas{Entity}` - database setup
  - `iSend{Action}` - HTTP request
  - `theResponseShouldBe{Condition}` - HTTP response validation
  - `iStore{What}As{VariableName}` - variable operations

### 正則表達式模式

- 字串參數: `"([^"]*)"`
- 數字參數: `(\d+)` (整數) 或 `(.+)` (浮點數)
- 開頭: `^`
- 結尾: `$`
- 完整範例: `^I send "([^"]*)" request to "([^"]*)"$`

## 避免寫死的規則

**不要**在分析時假設固定的分類規則,而是:

1. ✅ 從 `register.go` 動態讀取註冊的步驟組
2. ✅ 從檔案名稱推斷 entity 和功能
3. ✅ 使用關鍵字匹配推斷最可能的檔案
4. ✅ Progressive disclosure - 只讀取必要的檔案

這樣即使之後新增了 `database.RegisterOrderSteps(sc)`,也能自動適應。
