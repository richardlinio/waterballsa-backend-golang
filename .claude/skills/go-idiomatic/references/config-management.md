# 設定管理模式

專案採用 **Pass-by-value (傳值)** 的模式來處理設定，以表達設定在應用程式生命週期中的不可變性。

### Pass-by-value Pattern

此模式的核心是透過傳遞設定物件的副本（值）而非指標，來防止在應用程式執行期間意外修改設定。

### 實作範例

**1. Loader functions 回傳值而非指標**

```go
// config/database.go
func loadDatabaseConfig() (DatabaseConfig, error) {
    // ...
    return DatabaseConfig{...}, nil
}
```

**2. 主 Config struct 儲存值**

```go
// config/config.go
type Config struct {
    Database DatabaseConfig // 直接儲存 struct，而非指標 *DatabaseConfig
}
```

**3. 依賴注入時接收值**

```go
// internal/app/app.go
func New(cfg *config.Config) (*App, error) {
    // ...
    // 當將子設定傳遞給其他元件時，也傳遞值
    db, err := postgres.New(cfg.Database)
    // ...
}

// infrastructure/database/postgres.go
func New(cfg config.DatabaseConfig) (*pgxpool.Pool, error) {
    // ...
}
```

### 原則

- **不可變性 (Immutability)**: 設定一旦在啟動時載入，就不應該在執行期間被任何元件修改。傳值可以從語法層面強化此概念。
- **明確性**: 透過值傳遞，可以清楚地表明該函式或元件僅需讀取設定，而無意修改它。
- **一致性**: 整個專案對設定的處理保持一致的模式，減少混亂。
