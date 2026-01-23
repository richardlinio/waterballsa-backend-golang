package database

import "github.com/cucumber/godog"

// RegisterSteps registers all database-related step definitions
func RegisterSteps(sc *godog.ScenarioContext) {
	RegisterUserSteps(sc)
	RegisterJourneySteps(sc)
	RegisterChapterSteps(sc)
	RegisterMissionSteps(sc)
	RegisterRewardSteps(sc)
	RegisterMissionResourceSteps(sc)
}
