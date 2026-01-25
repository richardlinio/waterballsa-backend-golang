# 命名慣例

在專案中，應使用完整且具描述性的名稱而非縮寫，以提升程式碼的可讀性與一致性。

### 變數命名

```go
// 推薦 ✓
authService
userRepository
courseRepository
serverConfig
databaseConfig
```

```go
// 避免 ✗
authSvc
userRepo
courseRepo
srvCfg
dbCfg
```

### 原則

- **描述性**: 變數名稱應清楚地描述其用途或所代表的物件。
- **避免縮寫**: 除非是廣為人知的業界標準縮寫（如 `ID`, `URL`, `HTTP`, `DB`），否則應使用全名。
- **一致性**: 在整個專案中保持相同的命名風格，有助於開發者快速理解程式碼。
