# 常見錯誤類型範本

## Resource Not Found (404)

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

## Validation Failed (400)

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

## Conflict (409)

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

## Unauthorized (401)

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
