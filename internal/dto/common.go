package dto

type ErrorResponse struct {
	Code    string `json:"code,omitempty"`
	Error   string `json:"error"`
	Details any    `json:"details,omitempty"`
}
