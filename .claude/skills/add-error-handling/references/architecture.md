# 錯誤處理架構

## 錯誤流程

專案使用 **統一錯誤處理系統**，透過集中式 middleware 確保一致的錯誤回應與記錄。

```
Handler → Service → Repository
         ↓ returns AppError
      Middleware 攔截並回應 JSON
```

## 套件結構

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
