---
name: add-error-handling
description: 新增錯誤處理的完整實作步驟，包含錯誤碼定義、訊息對應、HTTP 狀態碼設定等 5 步驟流程。
allowed-tools: Read, Edit
---

# 新增錯誤處理

本 skill 提供新增自訂錯誤類型的完整步驟。

## 錯誤處理架構概覽

專案使用 **統一錯誤處理系統**，透過集中式 middleware 確保一致的錯誤回應與記錄。

### 錯誤流程

```
Handler → Service → Repository
         ↓ returns AppError
      Middleware 攔截並回應 JSON
```

### 套件結構

```
internal/apperror/
├── error.go      → AppError struct 與 constructor functions
├── codes.go      → 錯誤碼常數 (ERR_*)
├── messages.go   → 中文錯誤訊息
└── status.go     → HTTP 狀態碼對應

internal/middleware/
└── error_handler.go → 集中式錯誤處理與記錄
```

## 錯誤回應格式

所有錯誤回應遵循以下 JSON 結構：

```json
{
	"code": "ERR_USERNAME_EXISTS",
	"error": "使用者名稱已存在",
	"details": {
		"Username": "長度至少需要 3 個字元"
	}
}
```

- `code`: 錯誤碼常數 (例: `ERR_USERNAME_EXISTS`)
- `error`: 使用者友善的中文錯誤訊息
- `details`: 選擇性欄位，用於驗證錯誤詳情 (欄位層級錯誤)

## 新增錯誤類型的 5 步驟

### 步驟 1: 新增錯誤碼常數

**檔案**: [internal/apperror/codes.go](internal/apperror/codes.go)

**重要**: 錯誤碼必須 **依 HTTP 狀態碼分組並排序**

```go
const (
    // 404 Not Found
    CodeUserNotFound    = "ERR_USER_NOT_FOUND"
    CodeCourseNotFound  = "ERR_COURSE_NOT_FOUND"  // 新增
    CodeMissionNotFound = "ERR_MISSION_NOT_FOUND"
)
```

**HTTP 狀態碼順序**:
- 400 (Bad Request)
- 401 (Unauthorized)
- 403 (Forbidden)
- 404 (Not Found)
- 409 (Conflict)
- 429 (Too Many Requests)
- 500 (Internal Server Error)
- 503 (Service Unavailable)

### 步驟 2: 新增中文錯誤訊息

**檔案**: [internal/apperror/messages.go](internal/apperror/messages.go)

```go
var errorMessages = map[string]string{
    // 404 Not Found
    CodeUserNotFound:    "使用者不存在",
    CodeCourseNotFound:  "課程不存在",  // 新增
    CodeMissionNotFound: "任務不存在",
}
```

### 步驟 3: 對應 HTTP 狀態碼

**檔案**: [internal/apperror/status.go](internal/apperror/status.go)

```go
var httpStatusMap = map[string]int{
    // 404 Not Found
    CodeUserNotFound:    http.StatusNotFound,
    CodeCourseNotFound:  http.StatusNotFound,  // 新增
    CodeMissionNotFound: http.StatusNotFound,
}
```

### 步驟 4: 建立 Constructor Function

**檔案**: [internal/apperror/error.go](internal/apperror/error.go)

```go
// 404 Not Found
func UserNotFound() *AppError {
    return New(CodeUserNotFound)
}

func CourseNotFound() *AppError {  // 新增
    return New(CodeCourseNotFound)
}

func MissionNotFound() *AppError {
    return New(CodeMissionNotFound)
}
```

**命名規則**: 函數名稱應移除 `Code` 前綴並使用 PascalCase

### 步驟 5: 在 Service Layer 使用

**檔案**: `internal/service/*.go`

```go
func (s *Service) GetCourse(ctx context.Context, courseID int64) (*model.Course, error) {
    course, err := s.repository.GetCourseByID(ctx, courseID)
    if err != nil {
        // Repository 回傳 domain error
        if errors.Is(err, repository.ErrResourceNotFound) {
            return nil, apperror.CourseNotFound()  // 使用新的錯誤
        }
        return nil, apperror.DatabaseError(err)
    }
    return course, nil
}
```

## 在程式碼中使用

### Service Layer

Service 應針對所有業務邏輯錯誤回傳 `*apperror.AppError`：

```go
// 回傳預定義錯誤
if exists {
    return 0, apperror.UsernameExists()
}

// 包裝底層錯誤 (資料庫、內部錯誤)
if err != nil {
    return 0, apperror.DatabaseError(err)
}
```

### Handler Layer

Handler 應 **避免** 手動建立 JSON 回應。使用 `c.Error()` 將錯誤傳遞給 middleware：

```go
// 驗證錯誤
if err := c.ShouldBindJSON(&req); err != nil {
    _ = c.Error(apperror.NewWithError(apperror.CodeValidationFailed, err))
    return
}

// Service 錯誤 (service 已回傳 AppError)
if err := h.service.SomeMethod(ctx, req); err != nil {
    _ = c.Error(err)
    return
}
```

### 關鍵模式

- ✅ Service 回傳 `*apperror.AppError`
- ✅ Handler 呼叫 `c.Error(err)` 並提早返回
- ✅ Middleware 自動建立 JSON 回應
- ❌ **避免** 在 handler 中手動呼叫 `c.JSON()` 處理錯誤

## 錯誤記錄

錯誤記錄由 error handler middleware **自動處理**：

- **AppError 帶底層錯誤** (`Err` 欄位有值): ERROR level 記錄，包含底層錯誤詳情
- **AppError 無底層錯誤**: WARN level 記錄，包含 code、message、status、path、method
- **非預期錯誤** (非 AppError): ERROR level 記錄，包含完整錯誤詳情

**Handler 中無需手動記錄錯誤** - middleware 會自動處理所有錯誤記錄。

## 完整範例：新增 "Journey Not Found" 錯誤

### 1. codes.go

```go
const (
    // 404 Not Found
    CodeJourneyNotFound = "ERR_JOURNEY_NOT_FOUND"
)
```

### 2. messages.go

```go
var errorMessages = map[string]string{
    CodeJourneyNotFound: "旅程不存在",
}
```

### 3. status.go

```go
var httpStatusMap = map[string]int{
    CodeJourneyNotFound: http.StatusNotFound,
}
```

### 4. error.go

```go
func JourneyNotFound() *AppError {
    return New(CodeJourneyNotFound)
}
```

### 5. service/journey.go

```go
func (s *Service) GetJourney(ctx context.Context, journeyID int64) (*model.Journey, error) {
    journey, err := s.repository.GetJourneyByID(ctx, journeyID)
    if err != nil {
        if errors.Is(err, repository.ErrResourceNotFound) {
            return nil, apperror.JourneyNotFound()
        }
        return nil, apperror.DatabaseError(err)
    }
    return journey, nil
}
```

## 檢查清單

新增錯誤類型時，確認：

- [ ] 錯誤碼已加入 `codes.go` (依 HTTP status 排序)
- [ ] 錯誤訊息已加入 `messages.go` (依 HTTP status 排序)
- [ ] HTTP 狀態碼已對應於 `status.go` (依 HTTP status 排序)
- [ ] Constructor function 已建立於 `error.go` (依 HTTP status 分組)
- [ ] Service layer 正確使用新錯誤
- [ ] Service 回傳 `*apperror.AppError`
- [ ] Handler 使用 `c.Error()` 而非 `c.JSON()`
- [ ] Repository 將資料庫錯誤轉換為 domain errors (如 `repository.ErrResourceNotFound`)

## 常見錯誤類型範本

### Resource Not Found (404)

```go
// codes.go
CodeXxxNotFound = "ERR_XXX_NOT_FOUND"

// messages.go
CodeXxxNotFound: "XXX不存在"

// status.go
CodeXxxNotFound: http.StatusNotFound

// error.go
func XxxNotFound() *AppError {
    return New(CodeXxxNotFound)
}
```

### Validation Failed (400)

```go
// codes.go
CodeInvalidXxx = "ERR_INVALID_XXX"

// messages.go
CodeInvalidXxx: "XXX格式不正確"

// status.go
CodeInvalidXxx: http.StatusBadRequest

// error.go
func InvalidXxx() *AppError {
    return New(CodeInvalidXxx)
}
```

### Conflict (409)

```go
// codes.go
CodeXxxExists = "ERR_XXX_EXISTS"

// messages.go
CodeXxxExists: "XXX已存在"

// status.go
CodeXxxExists: http.StatusConflict

// error.go
func XxxExists() *AppError {
    return New(CodeXxxExists)
}
```

### Unauthorized (401)

```go
// codes.go
CodeInvalidCredentials = "ERR_INVALID_CREDENTIALS"

// messages.go
CodeInvalidCredentials: "帳號或密碼錯誤"

// status.go
CodeInvalidCredentials: http.StatusUnauthorized

// error.go
func InvalidCredentials() *AppError {
    return New(CodeInvalidCredentials)
}
```
