package database

import (
	"context"
	"fmt"
	"strconv"

	"github.com/cucumber/godog"
	"github.com/richardlinio/waterballsa-backend-golang/tests/bdd/testcontext"
	"github.com/richardlinio/waterballsa-backend-golang/tests/testutil"
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

	// Resolve template variable if present (supports {{lastUserId}})
	userIDStr, err = replaceVariables(ctx, userIDStr)
	if err != nil {
		return ctx, fmt.Errorf("failed to resolve user_id variable: %w", err)
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

	// Store order ID in context for variable substitution ({{lastOrderId}})
	ctx = context.WithValue(ctx, testcontext.ContextKeyLastOrderID, orderID)

	return ctx, nil
}

// theDatabaseShouldHaveUnpaidOrderForUserAndJourney verifies the count of unpaid orders for a specific user and journey
// Supports variable substitution for user ID and journey ID
func theDatabaseShouldHaveUnpaidOrderForUserAndJourney(ctx context.Context, expectedCount int, userIDStr, journeyIDStr string) error {
	// Get test server from context
	testServer, ok := ctx.Value(testcontext.ContextKeyTestServer).(*testcontext.TestServerWrapper)
	if !ok {
		return fmt.Errorf("test server not found in context")
	}

	// Resolve template variables if present
	userIDStr, err := replaceVariables(ctx, userIDStr)
	if err != nil {
		return fmt.Errorf("failed to resolve user ID variable: %w", err)
	}

	journeyIDStr, err = replaceVariables(ctx, journeyIDStr)
	if err != nil {
		return fmt.Errorf("failed to resolve journey ID variable: %w", err)
	}

	// Parse user ID
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid user ID: %w", err)
	}

	// Parse journey ID
	journeyID, err := strconv.ParseInt(journeyIDStr, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid journey ID: %w", err)
	}

	// Query database for unpaid orders
	query := `
		SELECT COUNT(*)
		FROM orders o
		INNER JOIN order_items oi ON o.id = oi.order_id
		WHERE o.user_id = $1 AND oi.journey_id = $2 AND o.status = 'UNPAID'
	`

	var actualCount int
	err = testServer.Server.Pool.QueryRow(ctx, query, userID, journeyID).Scan(&actualCount)
	if err != nil {
		return fmt.Errorf("failed to query unpaid orders: %w", err)
	}

	// Verify count
	if actualCount != expectedCount {
		return fmt.Errorf("expected %d unpaid order(s) for user %d and journey %d, but found %d",
			expectedCount, userID, journeyID, actualCount)
	}

	return nil
}

// theDatabaseShouldHaveUnpaidOrdersForUser verifies the count of unpaid orders for a specific user
// Supports variable substitution for user ID
func theDatabaseShouldHaveUnpaidOrdersForUser(ctx context.Context, expectedCount int, userIDStr string) error {
	// Get test server from context
	testServer, ok := ctx.Value(testcontext.ContextKeyTestServer).(*testcontext.TestServerWrapper)
	if !ok {
		return fmt.Errorf("test server not found in context")
	}

	// Resolve template variable if present
	userIDStr, err := replaceVariables(ctx, userIDStr)
	if err != nil {
		return fmt.Errorf("failed to resolve user ID variable: %w", err)
	}

	// Parse user ID
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid user ID: %w", err)
	}

	// Query database for unpaid orders
	query := `
		SELECT COUNT(*)
		FROM orders
		WHERE user_id = $1 AND status = 'UNPAID'
	`

	var actualCount int
	err = testServer.Server.Pool.QueryRow(ctx, query, userID).Scan(&actualCount)
	if err != nil {
		return fmt.Errorf("failed to query unpaid orders: %w", err)
	}

	// Verify count
	if actualCount != expectedCount {
		return fmt.Errorf("expected %d unpaid order(s) for user %d, but found %d",
			expectedCount, userID, actualCount)
	}

	return nil
}

// RegisterOrderSteps registers order-related step definitions
func RegisterOrderSteps(sc *godog.ScenarioContext) {
	sc.Step(`^the database has an order:$`, theDatabaseHasAnOrder)
	sc.Step(`^the database should have (\d+) unpaid order for user "([^"]*)" and journey "([^"]*)"$`, theDatabaseShouldHaveUnpaidOrderForUserAndJourney)
	sc.Step(`^the database should have (\d+) unpaid orders for user "([^"]*)"$`, theDatabaseShouldHaveUnpaidOrdersForUser)
}
