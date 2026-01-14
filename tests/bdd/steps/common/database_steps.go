package common

import (
	"context"
	"fmt"

	"github.com/cucumber/godog"
	"github.com/linporu/waterballsa-backend-golang/tests/bdd/testcontext"
	"github.com/linporu/waterballsa-backend-golang/tests/testutil"
)

// theDatabaseHasAUser creates a test user in the database from a Gherkin data table
// Expected table format:
//
//	| username | Bob         |
//	| password | Secure123!  |
func theDatabaseHasAUser(ctx context.Context, table *godog.Table) (context.Context, error) {
	// Get test server from suite context
	testServer, ok := ctx.Value(testcontext.ContextKeyTestServer).(*testcontext.TestServerWrapper)
	if !ok {
		return ctx, fmt.Errorf("test server not found in context")
	}

	// Parse table into a map
	userData, err := parseTableToMap(table)
	if err != nil {
		return ctx, fmt.Errorf("failed to parse user data table: %w", err)
	}

	// Extract username and password
	username, ok := userData["username"]
	if !ok {
		return ctx, fmt.Errorf("username not found in table")
	}

	password, ok := userData["password"]
	if !ok {
		return ctx, fmt.Errorf("password not found in table")
	}

	// Create test user in database
	userID, err := testutil.CreateTestUser(ctx, testServer.Server.Pool, username, password)
	if err != nil {
		return ctx, fmt.Errorf("failed to create test user: %w", err)
	}

	// Optionally store user ID in context if needed by other steps
	_ = userID

	return ctx, nil
}

// parseTableToMap converts a Gherkin table to a map[string]string
// Expected format:
//
//	| key1 | value1 |
//	| key2 | value2 |
func parseTableToMap(table *godog.Table) (map[string]string, error) {
	if table == nil {
		return nil, fmt.Errorf("table is nil")
	}

	if len(table.Rows) == 0 {
		return nil, fmt.Errorf("table has no rows")
	}

	result := make(map[string]string)

	for _, row := range table.Rows {
		if len(row.Cells) != 2 {
			return nil, fmt.Errorf("expected 2 columns per row, got %d", len(row.Cells))
		}

		key := row.Cells[0].Value
		value := row.Cells[1].Value
		result[key] = value
	}

	return result, nil
}

// RegisterDatabaseSteps registers all database-related step definitions
func RegisterDatabaseSteps(sc *godog.ScenarioContext) {
	sc.Step(`^the database has a user:$`, theDatabaseHasAUser)
}
