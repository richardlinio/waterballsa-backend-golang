package model

// MissionDeliveryResult represents the outcome of delivering a mission
type MissionDeliveryResult struct {
	Message          string
	ExperienceGained int
	TotalExperience  int32
	CurrentLevel     int32
}
