package database

import (
	"context"
	"fmt"

	"github.com/cucumber/godog"
	"github.com/richardlinio/waterballsa-backend-golang/tests/bdd/testcontext"
	"github.com/richardlinio/waterballsa-backend-golang/tests/testutil"
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

	// Store user ID in context for use by other steps
	ctx = context.WithValue(ctx, testcontext.ContextKeyLastUserID, userID)

	return ctx, nil
}

// RegisterUserSteps registers user-related step definitions
func RegisterUserSteps(sc *godog.ScenarioContext) {
	sc.Step(`^the database has a user:$`, theDatabaseHasAUser)
}
