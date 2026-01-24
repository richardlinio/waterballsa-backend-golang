package database

import (
	"context"
	"fmt"

	"github.com/cucumber/godog"
	"github.com/richardlinio/waterballsa-backend-golang/tests/bdd/testcontext"
	"github.com/richardlinio/waterballsa-backend-golang/tests/testutil"
)

// theDatabaseHasAJourney creates a test journey in the database from a Gherkin data table
// Expected table format:
//
//	| title           | 軟體設計模式精通之旅                  |
//	| slug            | design-patterns-mastery            |
//	| description     | 用 C.A. 模式大大提昇系統思維能力      |
//	| teacher         | 水球潘                             |
//	| price           | 1999                               |
//	| cover_image_url | https://example.com/cover1.jpg     |
func theDatabaseHasAJourney(ctx context.Context, table *godog.Table) (context.Context, error) {
	// Get test server from suite context
	testServer, ok := ctx.Value(testcontext.ContextKeyTestServer).(*testcontext.TestServerWrapper)
	if !ok {
		return ctx, fmt.Errorf("test server not found in context")
	}

	// Parse table into a map
	journeyData, err := parseTableToMap(table)
	if err != nil {
		return ctx, fmt.Errorf("failed to parse journey data table: %w", err)
	}

	// Extract required fields
	title, ok := journeyData["title"]
	if !ok {
		return ctx, fmt.Errorf("title not found in table")
	}

	slug, ok := journeyData["slug"]
	if !ok {
		return ctx, fmt.Errorf("slug not found in table")
	}

	description, ok := journeyData["description"]
	if !ok {
		description = "" // Optional field
	}

	teacher, ok := journeyData["teacher"]
	if !ok {
		return ctx, fmt.Errorf("teacher not found in table")
	}

	priceStr, ok := journeyData["price"]
	if !ok {
		return ctx, fmt.Errorf("price not found in table")
	}

	// Parse price as float64
	var price float64
	if _, err := fmt.Sscanf(priceStr, "%f", &price); err != nil {
		return ctx, fmt.Errorf("failed to parse price '%s': %w", priceStr, err)
	}

	coverImageURL, ok := journeyData["cover_image_url"]
	if !ok {
		coverImageURL = "" // Optional field
	}

	// Create test journey in database
	journeyID, err := testutil.CreateTestJourney(
		ctx,
		testServer.Server.Pool,
		title,
		slug,
		description,
		teacher,
		price,
		coverImageURL,
	)
	if err != nil {
		return ctx, fmt.Errorf("failed to create test journey: %w", err)
	}

	// Store journey ID in context for use by other steps (e.g., creating chapters)
	ctx = context.WithValue(ctx, testcontext.ContextKeyLastJourneyID, journeyID)

	return ctx, nil
}

// RegisterJourneySteps registers journey-related step definitions
func RegisterJourneySteps(sc *godog.ScenarioContext) {
	sc.Step(`^the database has a journey:$`, theDatabaseHasAJourney)
}
