package dto

type HealthCheckResponse struct {
	Status   string `json:"status"`
	Database string `json:"database"`
}
