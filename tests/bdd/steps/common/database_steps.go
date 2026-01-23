package common

import (
	"context"
	"fmt"
	"regexp"
	"strings"

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

// replaceVariables replaces {{variableName}} placeholders with values from context
// Supports: lastJourneyId, lastChapterId, lastMissionId
func replaceVariables(ctx context.Context, value string) (string, error) {
	// Regular expression to find {{variableName}} patterns
	re := regexp.MustCompile(`\{\{([^}]+)\}\}`)

	// Track if any variable was not found
	var missingVars []string

	result := re.ReplaceAllStringFunc(value, func(match string) string {
		// Extract variable name (remove {{ and }})
		varName := strings.Trim(match, "{}")

		// Get variable value from context based on name
		switch varName {
		case "lastJourneyId":
			if val, ok := ctx.Value(testcontext.ContextKeyLastJourneyID).(int64); ok {
				return fmt.Sprintf("%d", val)
			}
		case "lastChapterId":
			if val, ok := ctx.Value(testcontext.ContextKeyLastChapterID).(int64); ok {
				return fmt.Sprintf("%d", val)
			}
		case "lastMissionId":
			if val, ok := ctx.Value(testcontext.ContextKeyLastMissionID).(int64); ok {
				return fmt.Sprintf("%d", val)
			}
		}

		// Variable not found, track it
		missingVars = append(missingVars, varName)
		return match // Keep original if not found
	})

	// If any variables were not found, return error
	if len(missingVars) > 0 {
		return "", fmt.Errorf("variables not found in context: %v", missingVars)
	}

	return result, nil
}

// RegisterDatabaseSteps registers all database-related step definitions
func RegisterDatabaseSteps(sc *godog.ScenarioContext) {
	sc.Step(`^the database has a user:$`, theDatabaseHasAUser)
	sc.Step(`^the database has a journey:$`, theDatabaseHasAJourney)
	sc.Step(`^the database has a chapter:$`, theDatabaseHasAChapter)
	sc.Step(`^the database has a mission:$`, theDatabaseHasAMission)
}
