package apperror

var errorMessages = map[string]string{
	// Validation errors - both use same message for security
	CodeValidationFailed: "使用者名稱或密碼格式無效",
	CodePasswordTooLong:  "使用者名稱或密碼格式無效",

	// Auth errors - use generic message to avoid leaking information
	CodeAuthFailed:   "使用者名稱或密碼無效", // Same message for all auth failures (user not found, wrong password, rate limited)
	CodeUnauthorized: "未授權或權杖無效",

	// Conflict errors
	CodeUsernameExists: "使用者名稱已存在",

	// Server errors
	CodeInternalError: "伺服器內部錯誤",
	CodeDatabaseError: "資料庫操作失敗",

	// Service unavailable errors
	CodeServiceUnavailable: "服務暫時無法使用",
}

func GetMessage(code string) string {
	if msg, ok := errorMessages[code]; ok {
		return msg
	}
	return "未知錯誤"
}
