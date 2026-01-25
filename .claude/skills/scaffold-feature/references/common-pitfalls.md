# 常見陷阱 (Common Pitfalls)

本文件列出開發過程中應避免的常見錯誤和必須遵循的最佳實踐。

❌ **避免:**

- Service 方法回傳太多值 (參考 `/go-idiomatic`)
- **Service 層檢查 `pgx.ErrNoRows`** (參考 `/go-idiomatic`)
- Service 依賴 DTO layer (參考 `/go-idiomatic`)
- 跳過授權檢查
- 忘記註冊路由
- 遺漏依賴注入
- **忘記更新測試伺服器** (`tests/testutil/server.go`)
- 使用縮寫命名 (repo, svc, cfg) (參考 `/go-idiomatic`)
- 錯誤處理未依 HTTP status 排序 (參考 `/add-error-handling`)
- 忘記執行 `make sqlc`
- 忽略檢查現有 migrations
- 忽略檢查現有 step definitions

✅ **務必:**

- Service 回傳多個資料時使用 domain model struct (參考 `/go-idiomatic`)
- **Service 檢查 repository 自訂錯誤，而非 pgx 錯誤** (參考 `/go-idiomatic`)
- 驗證使用者權限
- 透過建構子注入依賴
- **同步更新生產與測試環境的 DI 配置**
- Repository 將 `pgx.ErrNoRows` 轉為 domain error (參考 `/go-idiomatic`)
- 使用完整變數名稱 (repository, service, config) (參考 `/go-idiomatic`)
- 新增錯誤類型遵循 5 步驟流程 (參考 `/add-error-handling`)
- 新增 queries 後執行 `make sqlc`
- 確認資料表結構一致
- 重用現有測試步驟
