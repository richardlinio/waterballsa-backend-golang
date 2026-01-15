package dto

type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50,alphanum_underscore"`
	Password string `json:"password" binding:"required,min=8,max=128,password_charset"`
}

type RegisterResponse struct {
	Message string `json:"message"`
	UserID  int64  `json:"userId"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	AccessToken string   `json:"accessToken"`
	User        UserInfo `json:"user"`
}

type UserInfo struct {
	ID         int64  `json:"id"`
	Username   string `json:"username"`
	Experience int32  `json:"experience"`
}

type LogoutResponse struct {
	Message string `json:"message"`
}

type RefreshResponse struct {
	AccessToken string   `json:"accessToken"`
	User        UserInfo `json:"user"`
}
