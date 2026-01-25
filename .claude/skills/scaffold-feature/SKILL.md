---
name: scaffold-feature
description: 根據 ISA feature 檔案分析需求,建立完整的後端功能實作計畫,包含分層架構的所有元件 (DTO, Model, Repository, Service, Handler, Router, DI)。
allowed-tools: Read Glob Grep Bash Write
---

此 skill 旨在分析 ISA feature 檔案的需求，並為新的後端功能生成一個完整、逐步的實作計畫，嚴格遵循專案的既定分層架構。

有關詳細的實作模式、程式碼範例和檢查清單，請參閱[參考指南](references/REFERENCE.md)。

## 執行步驟

### 階段 1: 分析需求

首先，透過閱讀相關文件來收集所有必要的背景資訊。

1.  **讀取 ISA feature 檔案** 以了解：
    *   API 端點 (方法、路徑、參數)。
    *   請求/回應的 JSON 結構。
    *   業務邏輯與規則。
    *   測試場景。
2.  **讀取 API 規格文件** (`/docs/api-docs/swagger.yaml`, `/docs/api-docs/openapi/paths/*.yaml`) 以提取：
    *   完整的端點規格。
    *   請求/回應的 schemas。
    *   HTTP 狀態碼。
    *   認證需求。
3.  **讀取資料庫綱要** (`/docs/db-schema.dbml`) 以確認：
    *   相關的資料表結構。
    *   欄位定義與約束。
    *   關聯與索引。
4.  **讀取專案慣例** (`/CLAUDE.md`) 以理解：
    *   分層架構模式。
    *   錯誤處理機制。
    *   命名與設定慣例。

### 階段 2: 探索現有程式碼

使用程式碼搜尋工具來識別可重用的元件，並確定需要建立哪些新元件。

1.  **使用 `Grep` 和 `Glob` 搜尋程式碼庫**，尋找：
    *   與該功能相關的現有 handlers、services 和 repositories。
    *   相關的 domain models 和 DTOs。
    *   現有的錯誤處理與認證模式。
2.  **檢查 `migrations/` 目錄**，確認資料表是否已存在或需要修改。
3.  **檢查 `tests/bdd/steps/` 目錄**，尋找可重用的 BDD step definitions。
4.  **總結你的發現**，將元件分類為「✅ 已存在」和「❌ 缺少」。

### 階段 3: 產生 Scaffold 計畫

根據你的分析，在一個 markdown 檔案中建立詳細的實作計畫。請使用以下結構。

```markdown
# Feature Scaffold 計畫: [功能名稱]

## 1. 現況分析

### 缺少元件
- ❌ [需要建立的元件列表]

## 2. API 端點需求

### [端點名稱: 例如，建立使用者]
- **方法與路徑**: `POST /users`
- **請求範例**: [JSON]
- **回應範例**: [JSON]
- **業務邏輯**: [關鍵規則]

## 3. 實作步驟

The implementation should follow the dependency order outlined in the checklist below. For detailed code examples and anti-patterns, consult the [Implementation Guide](references/implementation_guide.md).
## Scaffold 檢查清單

實作順序 (依賴關係):

1.  **Domain Models** (無依賴) - `internal/model/`
2.  **DTOs** (依賴 models) - `internal/dto/`
3.  **SQLc Queries** (無依賴) - `internal/db/queries/*.sql`
4.  **錯誤碼** (可提前) - `internal/apperror/`
5.  **Repository** (依賴 SQLc) - `internal/repository/`
6.  **Service** (依賴 repository) - `internal/service/`
7.  **Handler** (依賴 service) - `internal/handler/`
8.  **路由註冊** (依賴 handler) - `internal/router/router.go`
9.  **依賴注入** (依賴所有元件)
   - `internal/app/app.go` - 生產環境
   - `tests/testutil/server.go` - 測試環境 ⚠️ **必須同步更新**

必檢項目:

- [ ] 所有層級已建立
- [ ] 錯誤處理已加入 (參考 `/add-error-handling`)
- [ ] 授權檢查已實作 (如需要)
- [ ] 路由已註冊
- [ ] 依賴注入已配置 (`internal/app/app.go`)
- [ ] **測試伺服器已更新** (`tests/testutil/server.go`) ⚠️
- [ ] Migration 已建立 (如需要)
- [ ] SQLc 已執行 (`make sqlc`)
- [ ] 測試步驟已定義
- [ ] 遵循 Go idiomatic patterns (參考 `/go-idiomatic`)
## 4. 檔案總結

**需建立的檔案**:
- `[file/path_1.go]`
- `[file/path_2.go]`

**需修改的檔案**:
- `[file/path_3.go]`
- `[file/path_4.go]`

## 5. 驗證計畫
1.  **建置與檢查**: `make sqlc && make fmt && make lint`
2.  **測試**: `make test`
```

### 階段 4: 審查計畫

在執行前，審查產生的計畫以確保其完整與正確。
- ✅ 涵蓋 ISA 檔案中的所有場景。
- ✅ 符合 OpenAPI/Swagger 規格。
- ✅ 與資料庫綱要匹配。
- ✅ 遵循 `CLAUDE.md` 中的專案慣例。
- ✅ 包含所有必要的檔案與修改。
- ✅ 考慮到錯誤處理與授權。
- ✅ 定義了清晰的驗證策略。
