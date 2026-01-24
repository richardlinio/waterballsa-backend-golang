package apperror

var errorMessages = map[string]string{
	// Validation errors (400 - Bad Request) - both use same message for security
	CodeValidationFailed:     "使用者名稱或密碼格式無效",
	CodePasswordTooLong:      "使用者名稱或密碼格式無效",
	CodeInvalidWatchPosition: "觀看位置無效：必須介於 0 和影片長度之間",
	CodeMissionNotCompleted:  "任務尚未完成",

	// Auth errors (401 - Unauthorized) - use generic message to avoid leaking information
	CodeAuthFailed:   "使用者名稱或密碼無效", // Same message for all auth failures (user not found, wrong password, rate limited)
	CodeUnauthorized: "未授權或權杖無效",

	// Forbidden errors (403 - Forbidden)
	CodeProgressForbidden: "無法存取其他使用者的進度",

	// Not found errors (404 - Not Found)
	CodeJourneyNotFound:  "查無此旅程",
	CodeMissionNotFound:  "查無此任務",
	CodeProgressNotFound: "找不到進度記錄",

	// Conflict errors (409 - Conflict)
	CodeUsernameExists:          "使用者名稱已存在",
	CodeMissionAlreadyDelivered: "任務已經交付過了",

	// Rate limiting errors (429)
	CodeRateLimitExceeded: "請求次數過多,請稍後再試",

	// Server errors (500)
	CodeInternalServerError: "伺服器內部錯誤",
	CodeDatabaseError:       "資料庫操作失敗",
	CodeAuthStateError:      "認證狀態異常，請重新登入",

	// Service unavailable errors (503)
	CodeServiceUnavailable: "服務暫時無法使用",
}

func GetMessage(code string) string {
	if msg, ok := errorMessages[code]; ok {
		return msg
	}
	return "未知錯誤"
}
