package database

import (
	"context"
	"fmt"

	"github.com/cucumber/godog"
	"github.com/linporu/waterballsa-backend-golang/tests/bdd/testcontext"
	"github.com/linporu/waterballsa-backend-golang/tests/testutil"
)

// theDatabaseHasAMission creates a test mission in the database from a Gherkin data table
// Expected table format:
//
//	| chapter_id   | {{lastChapterId}}              |
//	| title        | 這門課手把手帶你成為架構設計的高手  |
//	| type         | VIDEO                          |
//	| access_level | PUBLIC                         |
//	| order_index  | 1                              |
func theDatabaseHasAMission(ctx context.Context, table *godog.Table) (context.Context, error) {
	// Get test server from suite context
	testServer, ok := ctx.Value(testcontext.ContextKeyTestServer).(*testcontext.TestServerWrapper)
	if !ok {
		return ctx, fmt.Errorf("test server not found in context")
	}

	// Parse table into a map
	missionData, err := parseTableToMap(table)
	if err != nil {
		return ctx, fmt.Errorf("failed to parse mission data table: %w", err)
	}

	// Extract chapter_id (with variable substitution support)
	chapterIDStr, ok := missionData["chapter_id"]
	if !ok {
		return ctx, fmt.Errorf("chapter_id not found in table")
	}
	chapterIDStr, err = replaceVariables(ctx, chapterIDStr)
	if err != nil {
		return ctx, fmt.Errorf("failed to replace variables in chapter_id: %w", err)
	}

	var chapterID int64
	if _, err := fmt.Sscanf(chapterIDStr, "%d", &chapterID); err != nil {
		return ctx, fmt.Errorf("failed to parse chapter_id '%s': %w", chapterIDStr, err)
	}

	// Extract title
	title, ok := missionData["title"]
	if !ok {
		return ctx, fmt.Errorf("title not found in table")
	}

	// Extract description (optional field)
	description, ok := missionData["description"]
	if !ok {
		description = "" // Optional field
	}

	// Extract type (VIDEO, ARTICLE, QUESTIONNAIRE)
	missionType, ok := missionData["type"]
	if !ok {
		return ctx, fmt.Errorf("type not found in table")
	}

	// Extract access_level (PUBLIC, AUTHENTICATED, PURCHASED)
	accessLevel, ok := missionData["access_level"]
	if !ok {
		return ctx, fmt.Errorf("access_level not found in table")
	}

	// Extract order_index
	orderIndexStr, ok := missionData["order_index"]
	if !ok {
		return ctx, fmt.Errorf("order_index not found in table")
	}

	var orderIndex int
	if _, err := fmt.Sscanf(orderIndexStr, "%d", &orderIndex); err != nil {
		return ctx, fmt.Errorf("failed to parse order_index '%s': %w", orderIndexStr, err)
	}

	// Create test mission in database
	missionID, err := testutil.CreateTestMission(
		ctx,
		testServer.Server.Pool,
		chapterID,
		title,
		description,
		missionType,
		accessLevel,
		orderIndex,
	)
	if err != nil {
		return ctx, fmt.Errorf("failed to create test mission: %w", err)
	}

	// Store mission ID in context for use by other steps
	ctx = context.WithValue(ctx, testcontext.ContextKeyLastMissionID, missionID)

	return ctx, nil
}

// RegisterMissionSteps registers mission-related step definitions
func RegisterMissionSteps(sc *godog.ScenarioContext) {
	sc.Step(`^the database has a mission:$`, theDatabaseHasAMission)
}
