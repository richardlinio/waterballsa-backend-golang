---
name: check-isa-steps
description: Analyze a .isa.feature file to identify which step definitions are NOT YET IMPLEMENTED and generate a step implementation plan. Use ONLY when user explicitly asks to "check missing steps", "find unimplemented steps", or "analyze step coverage" - NOT when implementing application code to pass tests.
allowed-tools: Read, Glob, Grep
---

# Check ISA Steps - Step Definition Analyzer

分析 `.isa.feature` 檔案,找出缺少的 step definitions 並建議實作位置。

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

### 6. 輸出分析報告

讀取 `examples.md` 查看報告格式,然後生成包含以下內容的報告:

- ✅ 已實作的步驟 (檔案位置、函式名稱)
- ❌ 缺少的步驟 (建議實作位置、程式碼範例)
- 📊 統計資訊
- 💡 下一步行動

## Token 最佳化原則

1. ✅ **只讀 1 次 register.go** - 建立完整索引
2. ✅ **智慧推斷候選檔案** - 減少不必要的檔案讀取
3. ✅ **使用 Grep 提取步驟** - 不需要讀取整個檔案
4. ✅ **Progressive disclosure** - 只在需要時讀取 reference.md 和 examples.md
5. ✅ **動態適應** - 不寫死分類規則,從程式碼結構推斷

## 重要提醒

- **只分析並報告**,不要自動實作步驟
- 提供清楚的檔案路徑連結,方便使用者查看
- 建議實作時,提供完整的程式碼範例和註冊方式
