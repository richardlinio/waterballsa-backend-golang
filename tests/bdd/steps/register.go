package steps

import (
	"github.com/cucumber/godog"
	"github.com/linporu/waterballsa-backend-golang/tests/bdd/steps/database"
	"github.com/linporu/waterballsa-backend-golang/tests/bdd/steps/http"
)

// RegisterSteps registers all step definitions (database and HTTP)
func RegisterSteps(sc *godog.ScenarioContext) {
	database.RegisterSteps(sc)
	http.RegisterSteps(sc)
}
