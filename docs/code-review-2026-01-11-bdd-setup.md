# Code Review - BDD 測試框架設定

**日期**: 2026-01-11
**審查範圍**: feat/setup-BDD 分支
**審查者**: Claude Code Review

---

## 📋 概述

此次審查涵蓋 BDD (Behavior-Driven Development) 測試框架的設定,主要使用 Godog 和 Testcontainers 建立端對端測試基礎設施。

### 變更檔案清單

- `Makefile` - 新增 BDD 測試指令
- `go.mod` - 新增測試相關依賴
- `internal/app/app.go` - 應用程式啟動邏輯
- `internal/infrastructure/database/postgres.go` - 資料庫連線池設定
- `tests/steps/isa_database_steps.go` - 資料庫相關測試步驟
- `tests/steps/isa_http_steps.go` - HTTP 請求相關測試步驟
- `tests/steps/steps.go` - 共享資料結構
- `tests/steps/suite_test.go` - 測試套件主程式
- `tests/testutil/database.go` - 資料庫測試工具函式
- `tests/testutil/server.go` - 測試伺服器
- `tests/testutil/testcontainer.go` - Testcontainer 設定

---

## ✅ 整體優點

1. **架構清晰**: 測試程式碼結構良好,職責分離明確
2. **錯誤處理完整**: 包含詳細的錯誤訊息,有助於除錯
3. **註解詳細**: 程式碼可讀性高,維護性好
4. **安全性考量**: 資料庫密碼處理、bcrypt 雜湊等符合安全最佳實踐
5. **資源管理**: 使用 `sync.Once`、`defer`、`errgroup` 等正確管理資源生命週期

---

## 🔧 改進建議

### 高優先級 (建議立即修正)

#### 1. HTTP Client 重複使用

**檔案**: [tests/steps/isa_http_steps.go:59](tests/steps/isa_http_steps.go#L59)

**問題**:

```go
// 每次請求都建立新的 HTTP Client
client := &http.Client{}
```

**影響**: 效能浪費,每次請求都建立新的 TCP 連線

**建議方案**:

```go
// 在 steps.go 中定義共享的 HTTP Client
var defaultHTTPClient = &http.Client{
    Timeout: 30 * time.Second,
}

// 在 iSendRequestTo 中使用
func iSendRequestTo(ctx context.Context, method, path string) (context.Context, error) {
    // ... existing code ...

    // 使用共享的 client
    resp, err := defaultHTTPClient.Do(req)
    // ...
}
```

---

#### 2. 伺服器啟動等待機制

**檔案**: [tests/testutil/server.go:145](tests/testutil/server.go#L145)

**問題**:

```go
// 硬編碼等待時間,不可靠
time.Sleep(100 * time.Millisecond)
```

**影響**:

- 在慢速環境可能導致測試失敗
- 在快速環境浪費時間

**建議方案**:

```go
// Start 方法中改用主動健康檢查
func (ts *TestServer) Start() error {
    addr := fmt.Sprintf("%s:%d", ts.Config.Server.Host, ts.Config.Server.Port)

    ts.Server = &http.Server{
        Addr:           addr,
        Handler:        ts.Engine,
        ReadTimeout:    ts.Config.Server.ReadTimeout,
        WriteTimeout:   ts.Config.Server.WriteTimeout,
        MaxHeaderBytes: 1 << 20,
    }

    // Start server in background
    go func() {
        if err := ts.Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            ts.Logger.Error("Test server failed to start", "error", err)
        }
    }()

    // 主動檢查伺服器是否就緒
    if err := ts.waitForReady(5 * time.Second); err != nil {
        return fmt.Errorf("server failed to become ready: %w", err)
    }

    ts.Logger.Info("Test server started", "address", addr)
    return nil
}

// waitForReady 輪詢 health endpoint 直到成功或超時
func (ts *TestServer) waitForReady(timeout time.Duration) error {
    ctx, cancel := context.WithTimeout(context.Background(), timeout)
    defer cancel()

    healthURL := ts.BaseURL() + "/health"
    client := &http.Client{Timeout: 1 * time.Second}

    ticker := time.NewTicker(50 * time.Millisecond)
    defer ticker.Stop()

    for {
        select {
        case <-ctx.Done():
            return fmt.Errorf("timeout waiting for server to be ready")
        case <-ticker.C:
            resp, err := client.Get(healthURL)
            if err == nil && resp.StatusCode == http.StatusOK {
                resp.Body.Close()
                return nil
            }
            if resp != nil {
                resp.Body.Close()
            }
        }
    }
}
```

---

#### 3. Response Body Close 錯誤處理

**檔案**: [tests/steps/isa_http_steps.go:67-68](tests/steps/isa_http_steps.go#L67-L68)

**問題**:

```go
if closeErr := resp.Body.Close(); closeErr != nil {
    return ctx, fmt.Errorf("failed to close response body: %w", closeErr)
}
```

**影響**: 在讀取 body 失敗時,close error 會覆蓋讀取錯誤

**建議方案**:

```go
// 讀取 response body
responseBody, err := io.ReadAll(resp.Body)
defer resp.Body.Close() // defer 確保一定會關閉

if err != nil {
    return ctx, fmt.Errorf("failed to read response body: %w", err)
}

// 繼續處理...
```

---

### 中優先級 (建議後續版本修正)

#### 4. 驗證器註冊冪等性

**檔案**: [tests/testutil/server.go:78-79](tests/testutil/server.go#L78-L79)

**問題**:

```go
// 註解聲稱是冪等的,但可能不是
// This is safe to call multiple times as it's idempotent
if err := validator.RegisterAuthValidators(); err != nil {
    return nil, fmt.Errorf("failed to register validators: %w", err)
}
```

**建議方案 1**: 確保 `RegisterAuthValidators` 真正是冪等的

```go
// internal/validator/auth.go
var registerOnce sync.Once
var registerErr error

func RegisterAuthValidators() error {
    registerOnce.Do(func() {
        // 註冊驗證器的邏輯...
    })
    return registerErr
}
```

**建議方案 2**: 在 suite 層級只註冊一次

```go
// tests/steps/suite_test.go
func TestMain(m *testing.M) {
    // 在所有測試之前註冊驗證器
    if err := validator.RegisterAuthValidators(); err != nil {
        panic(fmt.Sprintf("Failed to register validators: %v", err))
    }

    // ... 其他設定
}
```

---

#### 5. 未使用的變數處理

**檔案**: [tests/steps/isa_database_steps.go:47](tests/steps/isa_database_steps.go#L47)

**問題**:

```go
// Optionally store user ID in context if needed by other steps
_ = userID
```

**建議方案 1**: 如果確實不需要,直接忽略返回值

```go
_, err := testutil.CreateTestUser(ctx, testServer.Server.Pool, username, password)
if err != nil {
    return ctx, fmt.Errorf("failed to create test user: %w", err)
}
```

**建議方案 2**: 將 userID 存入 context 供後續步驟使用

```go
userID, err := testutil.CreateTestUser(ctx, testServer.Server.Pool, username, password)
if err != nil {
    return ctx, fmt.Errorf("failed to create test user: %w", err)
}

// 存入 context 供其他步驟使用 (例如驗證使用者資料)
ctx = context.WithValue(ctx, contextKeyCreatedUserID, userID)

return ctx, nil
```

並在 `steps.go` 中新增常數:

```go
const (
    contextKeyTestServer     contextKey = "testServer"
    contextKeyRequestBody    contextKey = "requestBody"
    contextKeyResponse       contextKey = "response"
    contextKeyResponseBody   contextKey = "responseBody"
    contextKeyCreatedUserID  contextKey = "createdUserID"  // 新增
)
```

---

#### 6. 測試併發數可配置化

**檔案**: [tests/steps/suite_test.go:32](tests/steps/suite_test.go#L32)

**問題**:

```go
// 硬編碼併發數
// Concurrency: 4,
```

**建議方案**:

```go
// 在 suite_test.go 開頭加入輔助函式
func getConcurrency() int {
    if concStr := os.Getenv("TEST_CONCURRENCY"); concStr != "" {
        if conc, err := strconv.Atoi(concStr); err == nil && conc > 0 {
            return conc
        }
    }
    return 1 // 預設不併發執行
}

// 在 TestFeatures 中使用
suite := godog.TestSuite{
    Name:                "ISA Auth Registration",
    ScenarioInitializer: InitializeScenario,
    Options: &godog.Options{
        Format:      "pretty",
        Paths:       []string{"../features/isa/auth"},
        Concurrency: getConcurrency(),
        TestingT:    t,
    },
}
```

並在文件中說明:

```bash
# 單執行緒執行 (預設)
make test-bdd-isa

# 或使用 4 個並發執行
TEST_CONCURRENCY=4 make test-bdd-isa
```

---

#### 7. 日誌輸出統一

**檔案**: [tests/steps/suite_test.go:100-101](tests/steps/suite_test.go#L100-L101)

**問題**:

```go
// 混用 fmt.Printf 和 Logger
fmt.Printf("PostgreSQL Testcontainer started at %s:%s\n",
    postgresContainer.Host, postgresContainer.Port)
```

**建議方案**: 統一使用 testing.T 的日誌方法

```go
// 修改函式簽名接收 *testing.M
func TestMain(m *testing.M) {
    ctx := context.Background()

    // 建立臨時 logger 用於 setup
    setupLogger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
        Level: slog.LevelInfo,
    }))

    // Start PostgreSQL Testcontainer
    var err error
    postgresContainer, err = testutil.NewPostgresContainer(ctx)
    if err != nil {
        setupLogger.Error("Failed to start postgres container", "error", err)
        os.Exit(1)
    }

    setupLogger.Info("PostgreSQL Testcontainer started",
        "host", postgresContainer.Host,
        "port", postgresContainer.Port)

    // ... 其他程式碼
}
```

---

### 低優先級 (可選改進)

#### 8. 檔案路徑安全驗證

**檔案**: [tests/testutil/testcontainer.go:89](tests/testutil/testcontainer.go#L89)

**問題**:

```go
content, err := os.ReadFile(filepath.Join(migrationPath, fileName)) // #nosec G304
```

**建議方案**: 加入路徑遍歷防護

```go
func processMigrationFiles(migrationPath, tempDir string, fileNames []string) ([]string, error) {
    var processedFiles []string

    // 確保 migrationPath 是絕對路徑
    absMigrationPath, err := filepath.Abs(migrationPath)
    if err != nil {
        return nil, fmt.Errorf("failed to get absolute migration path: %w", err)
    }

    for _, fileName := range fileNames {
        fullPath := filepath.Join(absMigrationPath, fileName)

        // 驗證檔案路徑在 migration 目錄內
        absFullPath, err := filepath.Abs(fullPath)
        if err != nil {
            return nil, fmt.Errorf("failed to get absolute path for %s: %w", fileName, err)
        }

        if !strings.HasPrefix(absFullPath, absMigrationPath) {
            return nil, fmt.Errorf("invalid migration file path (path traversal detected): %s", fileName)
        }

        content, err := os.ReadFile(absFullPath)
        if err != nil {
            return nil, fmt.Errorf("failed to read migration file %s: %w", fileName, err)
        }

        // ... 繼續處理
    }
    return processedFiles, nil
}
```

---

#### 9. 臨時目錄清理時機

**檔案**: [tests/testutil/testcontainer.go:126-130](tests/testutil/testcontainer.go#L126-L130)

**問題**:

```go
defer func() {
    if err := os.RemoveAll(tempDir); err != nil {
        fmt.Printf("Warning: failed to remove temp directory %s: %v\n", tempDir, err)
    }
}()
```

**說明**:
當前實作在函式返回時立即清理臨時目錄。這在正常情況下沒問題,因為 Testcontainer 的 `WithInitScripts` 會在容器啟動時讀取並執行這些檔案。但如果需要保留檔案以供除錯,可以考慮以下方案。

**建議方案** (可選):

```go
// 新增環境變數控制是否保留臨時檔案
keepTempFiles := os.Getenv("KEEP_TEST_TEMP_FILES") == "true"

if !keepTempFiles {
    defer func() {
        if err := os.RemoveAll(tempDir); err != nil {
            fmt.Printf("Warning: failed to remove temp directory %s: %v\n", tempDir, err)
        }
    }()
} else {
    fmt.Printf("Keeping temp directory for debugging: %s\n", tempDir)
}
```

---

#### 10. 改進測試資料表格解析註解

**檔案**: [tests/steps/isa_database_steps.go:11-15](tests/steps/isa_database_steps.go#L11-L15)

**建議**: 讓註解更明確說明設計考量

```go
// theDatabaseHasAUser creates a test user in the database from a Gherkin data table
//
// Table format uses key-value pairs for flexible field mapping:
// This allows adding new user fields in the future without breaking existing tests
//
// Expected format:
//	| username | Bob         |
//	| password | Secure123!  |
func theDatabaseHasAUser(ctx context.Context, table *godog.Table) (context.Context, error) {
    // ...
}
```

---

## 📊 改進優先級總結

### 🔴 高優先級 (建議立即處理)

1. **HTTP Client 重複使用** - 效能優化
2. **伺服器啟動等待機制** - 可靠性提升
3. **Response Body Close 處理** - 錯誤處理正確性

### 🟡 中優先級 (建議下個迭代處理)

4. **驗證器註冊冪等性** - 避免潛在的重複註冊問題
5. **未使用變數處理** - 程式碼清晰度
6. **測試併發數可配置** - 靈活性提升
7. **日誌輸出統一** - 程式碼一致性

### 🟢 低優先級 (可選改進)

8. **檔案路徑安全驗證** - 安全性增強
9. **臨時目錄清理時機** - 除錯便利性
10. **註解改進** - 文件完整性

---

## 🎯 建議的實作順序

### 第一階段: 快速修正 (預估 30-45 分鐘)

1. HTTP Client 重複使用 (10 分鐘)
2. Response Body Close 處理 (5 分鐘)
3. 未使用變數處理 (5 分鐘)
4. 伺服器啟動等待機制 (15-20 分鐘)

### 第二階段: 基礎架構改進 (預估 30-45 分鐘)

5. 驗證器註冊冪等性 (15-20 分鐘)
6. 日誌輸出統一 (10 分鐘)
7. 測試併發數可配置 (10 分鐘)

### 第三階段: 安全與文件 (預估 30 分鐘)

8. 檔案路徑安全驗證 (15 分鐘)
9. 註解改進 (5 分鐘)
10. 臨時目錄清理機制 (10 分鐘,可選)

---

## 📝 其他建議

### 測試覆蓋率

建議為 `tests/testutil/` 套件新增單元測試:

- `parseTableToMap` 函式測試
- `extractUpMigration` 函式測試
- `toFloat64` 型別轉換測試

### CI/CD 整合

建議在 CI pipeline 中加入:

```yaml
# .github/workflows/test.yml (範例)
- name: Run BDD Tests
  run: make test-bdd-isa
  env:
    TEST_CONCURRENCY: 2 # CI 環境可使用併發
```

### 文件更新

建議更新 [CLAUDE.md](CLAUDE.md) 加入 BDD 測試相關說明:

````markdown
## BDD Testing

### Running BDD Tests

```bash
# Run BDD tests
make test-bdd-isa

# Run with concurrency
TEST_CONCURRENCY=4 make test-bdd-isa

# Keep temp files for debugging
KEEP_TEST_TEMP_FILES=true make test-bdd-isa
```
````

### Writing BDD Tests

BDD tests use Godog and are located in:

- Feature files: `tests/features/`
- Step definitions: `tests/steps/`
- Test utilities: `tests/testutil/`

See existing tests for examples.

```

---

## ✅ 結論

整體而言,這是一個品質良好的 BDD 測試框架實作。建議的改進主要集中在:
1. **效能優化** (HTTP Client 重用)
2. **可靠性提升** (主動健康檢查)
3. **程式碼品質** (錯誤處理、日誌統一)

遵循建議的實作順序,可以在約 2 小時內完成所有高優先級和中優先級的改進,大幅提升測試框架的品質與可維護性。
```
