package database

import (
	"context"
	"fmt"

	"github.com/cucumber/godog"
	"github.com/linporu/waterballsa-backend-golang/tests/bdd/testcontext"
	"github.com/linporu/waterballsa-backend-golang/tests/testutil"
)

// theDatabaseHasAMissionResource creates a test mission resource in the database from a Gherkin data table
// Expected table format:
//
//	| mission_id       | {{lastMissionId}}              |
//	| type             | VIDEO                          |
//	| resource_url     | https://example.com/video.m3u8 |
//	| content_order    | 0                              |
//	| duration_seconds | 256                            |
func theDatabaseHasAMissionResource(ctx context.Context, table *godog.Table) (context.Context, error) {
	// Get test server from suite context
	testServer, ok := ctx.Value(testcontext.ContextKeyTestServer).(*testcontext.TestServerWrapper)
	if !ok {
		return ctx, fmt.Errorf("test server not found in context")
	}

	// Parse table into a map
	resourceData, err := parseTableToMap(table)
	if err != nil {
		return ctx, fmt.Errorf("failed to parse mission resource data table: %w", err)
	}

	// Extract mission_id (with variable substitution support)
	missionIDStr, ok := resourceData["mission_id"]
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

	// Extract type (maps to resource_type in database)
	resourceType, ok := resourceData["type"]
	if !ok {
		return ctx, fmt.Errorf("type not found in table")
	}

	// Extract resource_url
	resourceURL, ok := resourceData["resource_url"]
	if !ok {
		return ctx, fmt.Errorf("resource_url not found in table")
	}

	// Extract content_order
	contentOrderStr, ok := resourceData["content_order"]
	if !ok {
		return ctx, fmt.Errorf("content_order not found in table")
	}

	var contentOrder int
	if _, err := fmt.Sscanf(contentOrderStr, "%d", &contentOrder); err != nil {
		return ctx, fmt.Errorf("failed to parse content_order '%s': %w", contentOrderStr, err)
	}

	// Extract duration_seconds (optional, nullable field)
	var durationSeconds *int
	if durationSecondsStr, ok := resourceData["duration_seconds"]; ok {
		var duration int
		if _, err := fmt.Sscanf(durationSecondsStr, "%d", &duration); err != nil {
			return ctx, fmt.Errorf("failed to parse duration_seconds '%s': %w", durationSecondsStr, err)
		}
		durationSeconds = &duration
	}

	// Create test mission resource in database
	_, err = testutil.CreateTestMissionResource(
		ctx,
		testServer.Server.Pool,
		missionID,
		resourceType,
		resourceURL,
		contentOrder,
		durationSeconds,
	)
	if err != nil {
		return ctx, fmt.Errorf("failed to create test mission resource: %w", err)
	}

	// No need to store resource ID in context (not referenced in tests)

	return ctx, nil
}

// RegisterMissionResourceSteps registers mission resource-related step definitions
func RegisterMissionResourceSteps(sc *godog.ScenarioContext) {
	sc.Step(`^the database has a mission resource:$`, theDatabaseHasAMissionResource)
}
