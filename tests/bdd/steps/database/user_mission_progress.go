package database

import (
	"context"
	"fmt"

	"github.com/cucumber/godog"
	"github.com/linporu/waterballsa-backend-golang/tests/bdd/testcontext"
	"github.com/linporu/waterballsa-backend-golang/tests/testutil"
)

// theDatabaseHasAUserMissionProgress creates a user mission progress record in the database from a Gherkin data table
// Expected table format:
//
//	| user_id                | 1                 |
//	| mission_id             | {{lastMissionId}} |
//	| status                 | UNCOMPLETED       |
//	| watch_position_seconds | 30                |
func theDatabaseHasAUserMissionProgress(ctx context.Context, table *godog.Table) (context.Context, error) {
	// Get test server from suite context
	testServer, ok := ctx.Value(testcontext.ContextKeyTestServer).(*testcontext.TestServerWrapper)
	if !ok {
		return ctx, fmt.Errorf("test server not found in context")
	}

	// Parse table into a map
	progressData, err := parseTableToMap(table)
	if err != nil {
		return ctx, fmt.Errorf("failed to parse user mission progress data table: %w", err)
	}

	// Extract user_id
	userIDStr, ok := progressData["user_id"]
	if !ok {
		return ctx, fmt.Errorf("user_id not found in table")
	}

	var userID int64
	if _, err := fmt.Sscanf(userIDStr, "%d", &userID); err != nil {
		return ctx, fmt.Errorf("failed to parse user_id '%s': %w", userIDStr, err)
	}

	// Extract mission_id (with variable substitution support)
	missionIDStr, ok := progressData["mission_id"]
	if !ok {
		return ctx, fmt.Errorf("mission_id not found in table")
	}
	missionIDStr, err = replaceVariables(ctx, missionIDStr)
	if err != nil {
		return ctx, fmt.Errorf("failed to replace variables in mission_id: %w", err)
	}

	var missionID int64
	if _, err := fmt.Sscanf(missionIDStr, "%d", &missionID); err != nil {
		return ctx, fmt.Errorf("failed to parse mission_id '%s': %w", missionIDStr, err)
	}

	// Extract status (UNCOMPLETED, COMPLETED, DELIVERED)
	status, ok := progressData["status"]
	if !ok {
		return ctx, fmt.Errorf("status not found in table")
	}

	// Extract watch_position_seconds
	watchPositionStr, ok := progressData["watch_position_seconds"]
	if !ok {
		return ctx, fmt.Errorf("watch_position_seconds not found in table")
	}

	var watchPositionSeconds int
	if _, err := fmt.Sscanf(watchPositionStr, "%d", &watchPositionSeconds); err != nil {
		return ctx, fmt.Errorf("failed to parse watch_position_seconds '%s': %w", watchPositionStr, err)
	}

	// Create test user mission progress in database
	progressID, err := testutil.CreateTestUserMissionProgress(
		ctx,
		testServer.Server.Pool,
		userID,
		missionID,
		status,
		watchPositionSeconds,
	)
	if err != nil {
		return ctx, fmt.Errorf("failed to create test user mission progress: %w", err)
	}

	// Optionally store progress ID in context if needed by other steps
	_ = progressID

	return ctx, nil
}

// RegisterUserMissionProgressSteps registers user mission progress-related step definitions
func RegisterUserMissionProgressSteps(sc *godog.ScenarioContext) {
	sc.Step(`^the database has a user mission progress:$`, theDatabaseHasAUserMissionProgress)
}
