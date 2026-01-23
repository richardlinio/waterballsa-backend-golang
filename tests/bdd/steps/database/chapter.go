package database

import (
	"context"
	"fmt"

	"github.com/cucumber/godog"
	"github.com/linporu/waterballsa-backend-golang/tests/bdd/testcontext"
	"github.com/linporu/waterballsa-backend-golang/tests/testutil"
)

// theDatabaseHasAChapter creates a test chapter in the database from a Gherkin data table
// Expected table format:
//
//	| journey_id   | {{lastJourneyId}} |
//	| title        | 課程介紹           |
//	| order_index  | 1                 |
func theDatabaseHasAChapter(ctx context.Context, table *godog.Table) (context.Context, error) {
	// Get test server from suite context
	testServer, ok := ctx.Value(testcontext.ContextKeyTestServer).(*testcontext.TestServerWrapper)
	if !ok {
		return ctx, fmt.Errorf("test server not found in context")
	}

	// Parse table into a map
	chapterData, err := parseTableToMap(table)
	if err != nil {
		return ctx, fmt.Errorf("failed to parse chapter data table: %w", err)
	}

	// Extract journey_id (with variable substitution support)
	journeyIDStr, ok := chapterData["journey_id"]
	if !ok {
		return ctx, fmt.Errorf("journey_id not found in table")
	}
	journeyIDStr, err = replaceVariables(ctx, journeyIDStr)
	if err != nil {
		return ctx, fmt.Errorf("failed to replace variables in journey_id: %w", err)
	}

	var journeyID int64
	if _, err := fmt.Sscanf(journeyIDStr, "%d", &journeyID); err != nil {
		return ctx, fmt.Errorf("failed to parse journey_id '%s': %w", journeyIDStr, err)
	}

	// Extract title
	title, ok := chapterData["title"]
	if !ok {
		return ctx, fmt.Errorf("title not found in table")
	}

	// Extract order_index
	orderIndexStr, ok := chapterData["order_index"]
	if !ok {
		return ctx, fmt.Errorf("order_index not found in table")
	}

	var orderIndex int
	if _, err := fmt.Sscanf(orderIndexStr, "%d", &orderIndex); err != nil {
		return ctx, fmt.Errorf("failed to parse order_index '%s': %w", orderIndexStr, err)
	}

	// Create test chapter in database
	chapterID, err := testutil.CreateTestChapter(
		ctx,
		testServer.Server.Pool,
		journeyID,
		title,
		orderIndex,
	)
	if err != nil {
		return ctx, fmt.Errorf("failed to create test chapter: %w", err)
	}

	// Store chapter ID in context for use by other steps (e.g., creating missions)
	ctx = context.WithValue(ctx, testcontext.ContextKeyLastChapterID, chapterID)

	return ctx, nil
}

// RegisterChapterSteps registers chapter-related step definitions
func RegisterChapterSteps(sc *godog.ScenarioContext) {
	sc.Step(`^the database has a chapter:$`, theDatabaseHasAChapter)
}
