# Reference Guide: Feature Scaffolding

本參考指南補充 SKILL.md 的執行步驟，提供實戰範例與速查表。

## 內容索引

- **[專案結構與查詢](project-structure.md)**
  - `tree` 指令用法
  - 常用路徑速查表

- **[實作模式](implementation-patterns.md)**
  - Slice 初始化、回傳值設計、授權檢查
  - Not Found 處理、UPSERT 模式
  - SQLc 命名慣例、路由註冊

- **[依賴注入](dependency-injection.md)**
  - 生產環境 (`app.go`)
  - 測試環境 (`server.go`)

- **[測試與驗證](testing.md)**
  - ISA/BDD 測試範例
  - Scaffold 開發檢查清單

- **[常見陷阱與最佳實踐](common-pitfalls.md)**
  - 應避免的模式 (Anti-patterns)
  - 建議遵循的實踐 (Best Practices)

## 相關文件與 Skills

- **專案架構**: `/CLAUDE.md` - 分層架構、開發流程
- **Go 慣用寫法**: `/go-idiomatic` - 回傳值設計、Aggregate Structures、命名慣例
- **錯誤處理**: `/add-error-handling` - 5 步驟新增錯誤類型