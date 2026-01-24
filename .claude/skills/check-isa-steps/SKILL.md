---
name: check-isa-steps
description: Analyze a .isa.feature file to identify which step definitions are NOT YET IMPLEMENTED and create an implementation plan with TodoWrite. This skill runs in plan mode - after analysis, it creates todos for missing steps and waits for user approval (lgtm) before implementing.
allowed-tools: Read, Glob, Grep, TodoWrite, Edit, Write
---

# Check ISA Steps - Step Definition Analyzer & Implementer

分析 `.isa.feature` 檔案,找出缺少的 step definitions 並建立實作計畫（使用 TodoWrite）。

## 執行步驟

### 1. 讀取目標 feature 檔案

如果用戶提供了檔案路徑,讀取該檔案。如果沒有提供,檢查是否有打開的 `.isa.feature` 檔案。

### 2. 提取 feature 檔案中的步驟

提取所有的 Given/When/Then/And 步驟,建立待查找清單。

### 3. 建立步驟組索引 (讀取 1 個檔案)

**讀取 `tests/bdd/steps/register.go`**:

- 提取所有 `Register*Steps(sc)` 呼叫
- 建立包名到檔案路徑的映射

範例映射:

```
database.RegisterUserSteps    → tests/bdd/steps/database/user.go
http.RegisterRequestSteps     → tests/bdd/steps/http/request.go
```

### 4. 逐步匹配每個步驟 (Progressive Disclosure)

**對 feature 中的每個步驟**:

1. **智慧推斷最可能的檔案**

   讀取 `reference.md` 的「階段 2: 智慧匹配步驟到檔案」,使用以下策略:
   - **Entity 關鍵字匹配**: 步驟中的名詞 (user, journey, mission) → `database/{entity}.go`
   - **HTTP 動作關鍵字匹配**:
     - `I send`, `I set request` → `http/request.go`
     - `response status`, `response body` → `http/response.go`
     - `cookie` → `http/cookies.go`
     - `I store`, `{{variable}}` → `http/variables.go`

2. **只讀取候選檔案 (1-2 個)**

   使用 Grep 搜尋 `sc.Step(` 來提取該檔案中的所有步驟正則表達式。

3. **匹配步驟**

   將 feature 步驟與正則表達式匹配:
   - 找到 → 標記為已實作,記錄檔案位置和函式名稱
   - 找不到 → 標記為缺少

4. **如果有多個候選檔案**

   按優先順序逐一檢查,直到找到或確認缺少。

### 5. 提供實作建議 (針對缺少的步驟)

對於每個未匹配的步驟:

1. **讀取 `reference.md`** 了解命名慣例和正則表達式模式

2. **判斷實作位置**:
   - 是否為現有 entity 的新步驟? → 加入現有檔案
   - 是否為新 entity? → 建議創建新檔案

3. **建議內容**:
   - 函式名稱 (遵循 camelCase 慣例)
   - 正則表達式模式
   - 註冊方式 (在哪個 `Register*Steps` 函式中)
   - 是否需要在 `register.go` 中新增註冊呼叫

### 6. 建立實作計畫（使用 TodoWrite）

**不要輸出完整分析報告**，而是使用 TodoWrite 建立實作 todo list：

1. **為每個缺少的步驟建立一個 todo**:
   - content: "實作步驟: {步驟文字}" (例如: "實作步驟: Given the database has 5 journeys")
   - activeForm: "實作步驟: {步驟文字}"
   - status: "pending"

2. **Todo 描述中包含關鍵資訊**:
   - 建議實作位置 (檔案路徑)
   - 函式名稱建議
   - 是否需要新增註冊

3. **簡短摘要輸出**:
   - 找到 X 個缺少的步驟
   - 已建立實作計畫，等待用戶確認（說 "lgtm" 開始實作）

**範例 Todo**:
```
content: "實作步驟 'Given the database has 5 journeys' 於 tests/bdd/steps/database/journey.go (函式: theDatabaseHasJourneys)"
activeForm: "實作步驟 'Given the database has 5 journeys'"
status: "pending"
```

## Token 最佳化原則

1. ✅ **只讀 1 次 register.go** - 建立完整索引
2. ✅ **智慧推斷候選檔案** - 減少不必要的檔案讀取
3. ✅ **使用 Grep 提取步驟** - 不需要讀取整個檔案
4. ✅ **Progressive disclosure** - 只在需要時讀取 reference.md 和 examples.md
5. ✅ **動態適應** - 不寫死分類規則,從程式碼結構推斷

## 工作流程

這個 skill 在 **plan mode** 中運作:

1. **分析階段**: 找出缺少的步驟並判斷實作位置
2. **計畫階段**: 使用 TodoWrite 建立實作 todo list
3. **等待確認**: 等待用戶說 "lgtm" 才開始實作
4. **實作階段**: 逐一實作每個缺少的步驟

## 重要提醒

- **先建立計畫，再實作**: 使用 TodoWrite 建立 todo list，等待用戶確認
- **不要輸出冗長報告**: 只需簡短摘要 + todo list
- **提供清楚的檔案路徑**: 使用 markdown link 格式 `[file.go](path/to/file.go)`
- **建議實作時**: 在 todo 描述中包含函式名稱、正則表達式、註冊方式
