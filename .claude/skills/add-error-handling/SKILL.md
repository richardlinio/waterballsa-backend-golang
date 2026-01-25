---
name: add-error-handling
description: 在應用程式中新增一個新的自訂業務錯誤。當現有的 apperror 不足時，使用此 skill 來定義新的錯誤碼、訊息、HTTP 狀態並在服務中回傳。
allowed-tools: Read, Edit
---

# 新增錯誤處理

本 skill 提供在專案中新增自訂錯誤類型的完整步驟。

## 新增錯誤類型的 5 步驟

### 步驟 1: 新增錯誤碼常數

**檔案**: `internal/apperror/codes.go`

在 `const` 區塊中新增錯誤碼。
**重要**: 錯誤碼必須 **依 HTTP 狀態碼分組並排序**。

```go
const (
    // 404 Not Found
    CodeUserNotFound    = "ERR_USER_NOT_FOUND"
    CodeJourneyNotFound = "ERR_JOURNEY_NOT_FOUND" // New
    CodeMissionNotFound = "ERR_MISSION_NOT_FOUND"
)
```

### 步驟 2: 新增中文錯誤訊息

**檔案**: `internal/apperror/messages.go`

在 `errorMessages` map 中新增對應的中文訊息。

```go
var errorMessages = map[string]string{
    // 404 Not Found
    CodeUserNotFound:    "使用者不存在",
    CodeJourneyNotFound: "旅程不存在", // New
    CodeMissionNotFound: "任務不存在",
}
```

### 步驟 3: 對應 HTTP 狀態碼

**檔案**: `internal/apperror/status.go`

在 `httpStatusMap` map 中設定錯誤碼對應的 HTTP 狀態。

```go
var httpStatusMap = map[string]int{
    // 404 Not Found
    CodeUserNotFound:    http.StatusNotFound,
    CodeJourneyNotFound: http.StatusNotFound, // New
    CodeMissionNotFound: http.StatusNotFound,
}
```

### 步驟 4: 建立 Constructor Function

**檔案**: `internal/apperror/error.go`

建立一個 PascalCase 的 constructor function 來方便地生成錯誤。

```go
// 404 Not Found
func UserNotFound() *AppError {
    return New(CodeUserNotFound)
}

func JourneyNotFound() *AppError { // New
    return New(CodeJourneyNotFound)
}
```

### 步驟 5: 在 Service Layer 使用

**檔案**: `internal/service/*.go`

在適當的業務邏輯中，回傳新建的 `AppError`。

```go
if errors.Is(err, repository.ErrResourceNotFound) {
    return nil, apperror.JourneyNotFound() // Use the new error
}
```

---

## 檢查清單

新增錯誤類型時，確認已完成以下所有步驟：

- [ ] 錯誤碼已加入 `codes.go` (依 HTTP status 排序)
- [ ] 錯誤訊息已加入 `messages.go` (依 HTTP status 排序)
- [ ] HTTP 狀態碼已對應於 `status.go` (依 HTTP status 排序)
- [ ] Constructor function 已建立於 `error.go` (依 HTTP status 分組)
- [ ] Service layer 正確使用新錯誤。

## 延伸閱讀

- **架構與格式**: [錯誤處理架構](references/architecture.md)
- **使用方法**: [在程式碼中使用錯誤與記錄](references/usage-patterns.md)
- **完整範例**: [新增 "Journey Not Found" 錯誤](references/example.md)
- **常用範本**: [常見錯誤類型範本](references/templates.md)
