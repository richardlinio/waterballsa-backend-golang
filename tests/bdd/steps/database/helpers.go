package database

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/cucumber/godog"
	"github.com/linporu/waterballsa-backend-golang/tests/bdd/testcontext"
)

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
