---
name: go-idiomatic
description: Go 語言慣用寫法與設計模式指南，包含回傳值設計、Aggregate Structures、Slice 初始化、命名慣例等。
allowed-tools: Read
---

# Go Idiomatic 寫法指南

本 skill 提供 Go 語言慣用寫法與專案特定的設計模式。在修改或撰寫 Go 程式碼時，請遵循以下核心原則。

## 核心原則與檢查清單

當你需要特定主題的詳細說明與程式碼範例時，請查閱對應的參考文件。

- **[回傳值設計](references/return-values.md)**
  - [ ] 服務層方法的回傳值超過 2 個時，應使用 Struct 進行封裝。

- **[Aggregate 結構設計](references/aggregate-structures.md)**
  - [ ] 服務層 (Service) 的業務結果 Struct `*Result` 不應依賴 DTO 層。
  - [ ] 控制層 (Handler) 負責將 Service 的 `*Result` 轉換為 Response DTO。

- **[Slice 初始化](references/slice-initialization.md)**
  - [ ] 使用 `make([]T, 0, capacity)` 搭配 `append` 模式初始化 Slice。

- **[命名慣例](references/naming-conventions.md)**
  - [ ] 變數名稱使用完整名稱（如 `authService`），避免使用縮寫（如 `authSvc`）。

- **[設定管理](references/config-management.md)**
  - [ ] 設定檔 (Config) 應使用傳值 (pass-by-value) 模式來表達其不可變性。

- **[Repository 錯誤處理](references/repository-errors.md)**
  - [ ] Repository 層負責將資料庫特定錯誤（如 `pgx.ErrNoRows`）轉換為領域錯誤（如 `repository.ErrResourceNotFound`）。
  - [ ] 服務層 (Service) 只檢查 Repository 定義的領域錯誤，而非資料庫的原始錯誤。
