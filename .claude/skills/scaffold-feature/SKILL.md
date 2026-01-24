---
name: scaffold-feature
description: 根據 ISA feature 檔案分析需求,建立完整的後端功能實作計畫,包含分層架構的所有元件 (DTO, Model, Repository, Service, Handler, Router, DI)。
allowed-tools: Read, Glob, Grep, Bash, Write
---

# 執行步驟

## 階段 1: 讀取文件

1. 讀取 ISA feature 檔案，識別：
   - API 端點 (方法、路徑、參數)
   - 請求/回應 JSON 結構
   - 業務邏輯規則
   - 測試場景

2. 讀取 `/docs/api-docs/swagger.yaml` 和 `/docs/api-docs/openapi/paths/*.yaml`，提取：
   - 完整端點規格
   - 請求/回應 schema
   - HTTP 狀態碼
   - 認證需求

3. 讀取 `/docs/db-schema.dbml`，確認：
   - 相關資料表結構
   - 欄位定義與限制
   - 關聯與索引

4. 讀取 `/CLAUDE.md`，理解：
   - 分層架構模式
   - 錯誤處理機制
   - 命名與設定慣例

## 階段 2: 探索程式碼

1. 使用 **Grep** 和 **Glob** 搜尋：
   - 現有的相關 handlers, services, repositories
   - Domain models 和 DTOs
   - 錯誤處理與認證模式

2. 檢查 `migrations/` 目錄：
   - 確認所需資料表是否存在
   - 驗證欄位定義是否一致
   - 標記是否需要新增 migration

3. 檢查 `tests/bdd/steps/` 目錄：
   - 確認已實作的 step definitions
   - 識別可重用的測試步驟

4. 分類結果：
   - ✅ 已存在: [列出可用的基礎設施]
   - ❌ 缺少: [列出需要實作的部分]

## 階段 3: 建立 Scaffold 計畫

按以下結構輸出 markdown 計畫檔案：

```markdown
# Feature Scaffold 計畫: [功能名稱]

## 概述

[簡短描述]

## 參考文件

- ISA Feature: [路徑]
- API Spec: [路徑]
- DB Schema: [路徑]

## 現況分析

### 已存在

- ✅ [項目]

### 缺少

- ❌ [項目]

## API 需求

### [端點名稱]

- **方法與路徑**: GET/POST/PUT/DELETE /path
- **請求範例**: [JSON]
- **回應範例**: [JSON]
- **業務規則**: [關鍵邏輯]

## 實作步驟

### 步驟 1: [標題]

**檔案**: [路徑] (建立/修改)
**內容**:

- [具體指示]
- [程式碼結構範例]

**依賴**: [先決條件]

[重複其他步驟...]

## 關鍵檔案

**建立**: [列表]
**修改**: [列表]

## 驗證計畫

1. 建置: `make sqlc && make fmt && make lint`
2. 測試: `make test`
3. 場景: [列出每個 ISA 場景]

## 成功標準

- [ ] 所有 ISA 場景通過
- [ ] 符合 swagger 規格
- [ ] 遵循 CLAUDE.md 模式
```

### 步驟排序原則

1. Domain Models (無依賴)
2. DTOs (依賴 models)
3. SQLc Queries (無依賴)
4. Repository (依賴 SQLc 生成的程式碼)
5. Service (依賴 repository)
6. Handler (依賴 service)
7. 錯誤碼 (可提前)
8. 路由註冊 (依賴 handler)
9. 依賴注入 (依賴所有元件)

## 階段 4: 審查

確認計畫：

- ✅ 涵蓋所有 ISA 場景
- ✅ 符合 swagger 規格
- ✅ 符合 db-schema 定義
- ✅ 遵循 CLAUDE.md 慣例
- ✅ 包含所有需要的檔案
- ✅ 包含錯誤處理與授權檢查
- ✅ 包含驗證步驟
