package database

import (
	"context"
	"fmt"
	"strconv"

	"github.com/cucumber/godog"
	"github.com/linporu/waterballsa-backend-golang/tests/bdd/testcontext"
	"github.com/linporu/waterballsa-backend-golang/tests/testutil"
)

// theDatabaseHasAnOrder creates a test order in the database from a Gherkin data table
// Expected table format:
//
//	| user_id    | 1                 |
//	| journey_id | {{lastJourneyId}} |
//	| status     | PAID              |
func theDatabaseHasAnOrder(ctx context.Context, table *godog.Table) (context.Context, error) {
	// Get test server from context
	testServer, ok := ctx.Value(testcontext.ContextKeyTestServer).(*testcontext.TestServerWrapper)
	if !ok {
		return ctx, fmt.Errorf("test server not found in context")
	}

	// Parse table into a map
	orderData, err := parseTableToMap(table)
	if err != nil {
		return ctx, fmt.Errorf("failed to parse order data table: %w", err)
	}

	// Extract user_id
	userIDStr, ok := orderData["user_id"]
	if !ok {
		return ctx, fmt.Errorf("user_id not found in table")
	}
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		return ctx, fmt.Errorf("invalid user_id: %w", err)
	}

	// Extract journey_id (may contain template variable)
	journeyIDStr, ok := orderData["journey_id"]
	if !ok {
		return ctx, fmt.Errorf("journey_id not found in table")
	}

	// Resolve template variable if present
	journeyIDStr, err = replaceVariables(ctx, journeyIDStr)
	if err != nil {
		return ctx, fmt.Errorf("failed to resolve journey_id variable: %w", err)
	}

	journeyID, err := strconv.ParseInt(journeyIDStr, 10, 64)
	if err != nil {
		return ctx, fmt.Errorf("invalid journey_id: %w", err)
	}

	// Extract status
	status, ok := orderData["status"]
	if !ok {
		return ctx, fmt.Errorf("status not found in table")
	}

	// Create test order in database
	orderID, err := testutil.CreateTestOrder(
		ctx,
		testServer.Server.Pool,
		userID,
		journeyID,
		status,
	)
	if err != nil {
		return ctx, fmt.Errorf("failed to create test order: %w", err)
	}

	// Store order ID in context if needed
	_ = orderID

	return ctx, nil
}

// RegisterOrderSteps registers order-related step definitions
func RegisterOrderSteps(sc *godog.ScenarioContext) {
	sc.Step(`^the database has an order:$`, theDatabaseHasAnOrder)
}
