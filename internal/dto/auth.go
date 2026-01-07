package dto

type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50,alphanum_underscore"`
	Password string `json:"password" binding:"required,min=8,max=128,password_charset"`
}

type RegisterResponse struct {
	Message string `json:"message"`
	UserID  int64  `json:"userId"`
}
