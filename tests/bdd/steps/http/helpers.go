package http

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"
)

// defaultHTTPClient is a shared HTTP client for all test requests
// It reuses TCP connections via connection pooling for better performance
var defaultHTTPClient = &http.Client{
	Timeout: 30 * time.Second,
}

// Helper function to get nested field from JSON object using dot notation
// Supports accessing nested fields like "user.id" or "data.user.profile.name"
// Also supports array indexing like "journeys[0].title" or "data.users[2].name"
func getNestedField(data map[string]any, fieldPath string) (any, bool) {
	// Split the field path by dots
	parts := strings.Split(fieldPath, ".")

	var current any = data

	// Navigate through each part of the path
	for _, part := range parts {
		// Check if part contains array index (e.g., "journeys[0]")
		if strings.Contains(part, "[") && strings.HasSuffix(part, "]") {
			// Split into field name and index
			openBracket := strings.Index(part, "[")
			if openBracket == -1 {
				return nil, false
			}
			fieldName := part[:openBracket]
			indexStr := part[openBracket+1 : len(part)-1]

			// Parse index
			var index int
			if _, err := fmt.Sscanf(indexStr, "%d", &index); err != nil {
				return nil, false
			}

			// Get the array field first
			currentMap, ok := current.(map[string]any)
			if !ok {
				return nil, false
			}

			arrayField, exists := currentMap[fieldName]
			if !exists {
				return nil, false
			}

			// Check if it's an array
			arrayValue, ok := arrayField.([]any)
			if !ok {
				return nil, false
			}

			// Check index bounds
			if index < 0 || index >= len(arrayValue) {
				return nil, false
			}

			// Get the element at index
			current = arrayValue[index]
		} else {
			// Regular field access (no array indexing)
			currentMap, ok := current.(map[string]any)
			if !ok {
				return nil, false
			}

			// Get the next level
			next, exists := currentMap[part]
			if !exists {
				return nil, false
			}

			current = next
		}
	}

	return current, true
}

// Helper function to get all keys from a map for error messages
func getMapKeys(m map[string]any) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// Helper function to convert various numeric types to float64
func toFloat64(val any) (float64, bool) {
	switch v := val.(type) {
	case float64:
		return v, true
	case float32:
		return float64(v), true
	case int:
		return float64(v), true
	case int64:
		return float64(v), true
	case int32:
		return float64(v), true
	case int16:
		return float64(v), true
	case int8:
		return float64(v), true
	case uint:
		return float64(v), true
	case uint64:
		return float64(v), true
	case uint32:
		return float64(v), true
	case uint16:
		return float64(v), true
	case uint8:
		return float64(v), true
	default:
		return 0, false
	}
}

// replaceVariables is the unified variable replacement function
// It replaces {{variableName}} placeholders with values from context
// Supports all context keys and stored variables
func replaceVariables(ctx context.Context, input string) (string, error) {
	// Regular expression to find {{variableName}} patterns
	re := regexp.MustCompile(`\{\{([^}]+)\}\}`)

	result := re.ReplaceAllStringFunc(input, func(match string) string {
		// Extract variable name (remove {{ and }})
		varName := strings.Trim(match, "{}")

		// Use the resolveVariable function from variables.go
		if val, found := resolveVariable(ctx, varName); found {
			return fmt.Sprintf("%v", val)
		}

		// Variable not found, keep original placeholder
		return match
	})

	return result, nil
}

// replaceVariablesInPath replaces {{variableName}} placeholders in URL paths with values from context
// Supports: lastJourneyId, lastChapterId, lastMissionId, lastUserId, lastOrderId, and stored variables
func replaceVariablesInPath(ctx context.Context, path string) (string, error) {
	return replaceVariables(ctx, path)
}

// replaceVariablesInRequestBody replaces {{variableName}} placeholders in request bodies with values from context
// Supports: lastJourneyId, lastChapterId, lastMissionId, lastUserId, lastOrderId, and stored variables
func replaceVariablesInRequestBody(ctx context.Context, body string) (string, error) {
	return replaceVariables(ctx, body)
}

// replaceVariablesInString replaces {{variableName}} placeholders in strings with values from context
// Supports: lastJourneyId, lastChapterId, lastMissionId, lastUserId, lastOrderId, and stored variables
func replaceVariablesInString(ctx context.Context, str string) (string, error) {
	return replaceVariables(ctx, str)
}
