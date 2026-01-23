package database

import (
	"context"
	"fmt"

	"github.com/cucumber/godog"
	"github.com/linporu/waterballsa-backend-golang/tests/bdd/testcontext"
	"github.com/linporu/waterballsa-backend-golang/tests/testutil"
)

// theDatabaseHasAReward creates a test reward in the database from a Gherkin data table
// Expected table format:
//
//	| mission_id   | {{lastMissionId}} |
//	| reward_type  | EXPERIENCE        |
//	| reward_value | 100               |
func theDatabaseHasAReward(ctx context.Context, table *godog.Table) (context.Context, error) {
	// Get test server from suite context
	testServer, ok := ctx.Value(testcontext.ContextKeyTestServer).(*testcontext.TestServerWrapper)
	if !ok {
		return ctx, fmt.Errorf("test server not found in context")
	}

	// Parse table into a map
	rewardData, err := parseTableToMap(table)
	if err != nil {
		return ctx, fmt.Errorf("failed to parse reward data table: %w", err)
	}

	// Extract mission_id (with variable substitution support)
	missionIDStr, ok := rewardData["mission_id"]
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

	// Extract reward_type
	rewardType, ok := rewardData["reward_type"]
	if !ok {
		return ctx, fmt.Errorf("reward_type not found in table")
	}

	// Extract reward_value
	rewardValueStr, ok := rewardData["reward_value"]
	if !ok {
		return ctx, fmt.Errorf("reward_value not found in table")
	}

	var rewardValue int
	if _, err := fmt.Sscanf(rewardValueStr, "%d", &rewardValue); err != nil {
		return ctx, fmt.Errorf("failed to parse reward_value '%s': %w", rewardValueStr, err)
	}

	// Create test reward in database
	_, err = testutil.CreateTestReward(
		ctx,
		testServer.Server.Pool,
		missionID,
		rewardType,
		rewardValue,
	)
	if err != nil {
		return ctx, fmt.Errorf("failed to create test reward: %w", err)
	}

	// No need to store reward ID in context (not referenced in tests)

	return ctx, nil
}

// RegisterRewardSteps registers reward-related step definitions
func RegisterRewardSteps(sc *godog.ScenarioContext) {
	sc.Step(`^the database has a reward:$`, theDatabaseHasAReward)
}
