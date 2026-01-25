# 在程式碼中使用錯誤

## Service Layer

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

## Handler Layer

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

---

# 錯誤記錄

錯誤記錄由 error handler middleware **自動處理**：

- **AppError 帶底層錯誤** (`Err` 欄位有值): `ERROR` level 記錄，包含底層錯誤詳情。
- **AppError 無底層錯誤**: `WARN` level 記錄，包含 code、message、status、path、method。
- **非預期錯誤** (非 AppError): `ERROR` level 記錄，包含完整錯誤詳情。

**Handler 中無需手動記錄錯誤** - middleware 會自動處理所有錯誤記錄。
