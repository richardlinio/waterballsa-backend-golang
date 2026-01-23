package steps

import (
	"github.com/cucumber/godog"
	"github.com/linporu/waterballsa-backend-golang/tests/bdd/steps/database"
	"github.com/linporu/waterballsa-backend-golang/tests/bdd/steps/http"
)

// RegisterSteps registers all step definitions (database and HTTP)
func RegisterSteps(sc *godog.ScenarioContext) {
	// Database steps
	database.RegisterUserSteps(sc)
	database.RegisterJourneySteps(sc)
	database.RegisterChapterSteps(sc)
	database.RegisterMissionSteps(sc)
	database.RegisterRewardSteps(sc)
	database.RegisterMissionResourceSteps(sc)

	// HTTP steps
	http.RegisterRequestSteps(sc)
	http.RegisterResponseSteps(sc)
	http.RegisterCookieSteps(sc)
	http.RegisterVariableSteps(sc)
}
