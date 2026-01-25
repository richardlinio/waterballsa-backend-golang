# Aggregate Structures 設計原則

專案區分不同類型的 aggregate structures，依據其用途與所在層級。

### Model Layer

包含兩種結構：

1.  **Database Entity**: 單一資料表映射
    - 範例: `Journey`, `Mission`, `Order`, `User`
    - 代表資料庫中的單一 row/document
    - 可包含實體特定的業務方法 (例: `User.CalculateLevel()`)

2.  **Query Result Aggregate**: 多表查詢結果容器
    - 範例: `JourneyDetail`, `MissionDetail`
    - 純資料容器，無業務邏輯
    - 由 Service 層組裝，但定義在 Model 層
    - 目的: 避免 Service 方法回傳太多參數
    - 識別: 組合多個資料表的資料用於讀取操作

### Service Layer

包含 **Business Result Structs**: 業務操作的結果

- 範例: `OrderResult`, `LoginResult`
- 包含業務邏輯產生的資料 (反正規化資料、計算欄位)
- 命名慣例: `*Result` suffix
- 識別: 包含業務邏輯產生的資料，而非純 DB 查詢
- **絕對不可** 依賴 DTO layer (避免 `service → dto` 依賴)

### DTO Layer

包含 **Request/Response DTOs**: API 序列化/反序列化

- 範例: `OrderResponse`, `LoginRequest`, `UserInfo`
- 僅處理 JSON tags、validation tags、API 格式化
- **絕對不可** 被 Service layer 引用
- 可以依賴 Model layer 進行轉換

### Dependency Rules

```
Handler → DTO ← (converts from) ← Service → Model
   ↓                                   ↓
(uses)                              (uses)
   ↓                                   ↓
Service                            Repository
```

- ✅ Handler 可依賴 Service 和 DTO
- ✅ Service 可依賴 Model 和 Repository
- ✅ DTO 可依賴 Model (用於轉換函數)
- ❌ Service **絕對不可** 依賴 DTO (防止循環依賴與層級混淆)

### 正確模式範例

```go
// service/auth.go
type LoginResult struct {
    AccessToken string
    UserID      int64      // ✅ 獨立欄位
    Username    string     // ✅ 無 dto.UserInfo 依賴
    Experience  int32
}

// handler/auth.go
result, _ := authService.Login(ctx, req)
c.JSON(http.StatusOK, dto.LoginResponse{
    AccessToken: result.AccessToken,
    User: dto.UserInfo{  // ✅ Handler 從 Result 建構 DTO
        ID:         result.UserID,
        Username:   result.Username,
        Experience: result.Experience,
    },
})
```

### 反模式 (Anti-pattern)

```go
// ❌ 錯誤: Service 依賴 DTO
type LoginResult struct {
    AccessToken string
    UserInfo    dto.UserInfo  // ❌ 建立 service → dto 依賴
}
```
