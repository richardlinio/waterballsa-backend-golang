# Repository 層錯誤處理

為了維持清晰的分層架構，Service 層不應該知道底層資料庫的實作細節，這包括資料庫驅動程式特定的錯誤。Repository 層的職責之一就是將這些底層錯誤轉換為抽象的領域錯誤。

### ✅ 正確模式

**1. Repository: 轉換為領域錯誤**
Repository 層捕捉資料庫驅動（如 `pgx`）的特定錯誤，並將其轉換為預先定義好的、與領域相關的錯誤。

```go
// repository/user.go
import (
    "errors"
    "github.com/jackc/pgx/v5"
    "your-project/internal/repository" // 引入 repository 層的錯誤定義
)

func (r *userRepository) GetUserByID(ctx context.Context, id int) (*model.User, error) {
    // ...
    if err != nil {
        if errors.Is(err, pgx.ErrNoRows) {
            // 將 pgx 的 ErrNoRows 轉換為自定義的領域錯誤
            return nil, repository.ErrResourceNotFound
        }
        return nil, err // 其他未預期的錯誤直接回傳
    }
    return user, nil
}
```

**2. Service: 檢查領域錯誤**
Service 層只檢查由 Repository 層定義的領域錯誤，完全不知道 `pgx.ErrNoRows` 的存在。

```go
// service/user.go
import (
    "errors"
    "your-project/internal/apperror"
    "your-project/internal/repository"
)

func (s *userService) GetUser(ctx context.Context, id int) (*model.User, error) {
    user, err := s.userRepo.GetUserByID(ctx, id)
    if err != nil {
        if errors.Is(err, repository.ErrResourceNotFound) {
            // 處理資源未找到的情況
            return nil, apperror.ResourceNotFound("user", "id", id)
        }
        return nil, err // 處理其他錯誤
    }
    return user, nil
}
```

### ❌ 反模式 (Anti-pattern)

Service 層直接檢查來自資料庫驅動的錯誤。這破壞了抽象，使得 Service 層與 `pgx` 緊密耦合。如果未來更換資料庫驅動，所有相關的 Service 層程式碼都需要修改。

```go
// service/user.go

// ❌ 錯誤: Service 層不應檢查 pgx.ErrNoRows
if errors.Is(err, pgx.ErrNoRows) {
    return nil, apperror.ResourceNotFound()
}
```

### 原則

- **抽象化**: Repository 層是對資料儲存的抽象。它應該隱藏所有實作細節。
- **單一職責**: 錯誤轉換是 Repository 的職責之一。Service 的職責是執行業務邏輯。
- **鬆耦合**: Service 層不應與特定的資料庫技術（如 `pgx`）耦合。
