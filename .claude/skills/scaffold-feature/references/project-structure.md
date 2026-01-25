# 專案結構查詢

**策略:** 使用分層查詢避免 token 浪費

### 快速導覽 (優先使用)

```bash
# 僅顯示主要目錄結構 (1-2 層)
tree -L 2 -d -I 'vendor|node_modules|.git|tmp' --dirsfirst

# 查看特定模組的檔案
tree internal/handler -I '*_test.go'
tree internal/service -I '*_test.go'
tree internal/repository -I '*_test.go'
```

### 常見路徑速查

| 需求       | 路徑                        | 說明            |
| ---------- | --------------------------- | --------------- |
| HTTP 處理  | `internal/handler/`         | Gin 路由處理器  |
| 業務邏輯   | `internal/service/`         | 服務層實作      |
| 資料存取   | `internal/repository/`      | Repository 層   |
| 資料結構   | `internal/dto/`             | API 請求/回應   |
| 領域模型   | `internal/model/`           | 資料庫實體      |
| SQL 查詢   | `internal/db/queries/`      | SQLc 查詢定義   |
| 錯誤定義   | `internal/apperror/`        | 統一錯誤處理    |
| 路由註冊   | `internal/router/router.go` | 路由設定        |
| DI 配置    | `internal/app/app.go`       | 依賴注入        |
| 資料庫遷移 | `migrations/`               | Goose migration |
| BDD 測試   | `tests/features/`           | Gherkin 功能檔  |

### 完整結構 (僅必要時使用)

```bash
# ⚠️ 會產生大量輸出，僅在需要全貌時使用
tree -L 3 -I 'vendor|node_modules|.git|tmp' --dirsfirst
```
