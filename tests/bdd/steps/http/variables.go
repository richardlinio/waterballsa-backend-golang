package http

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/cucumber/godog"
	"github.com/richardlinio/waterballsa-backend-golang/tests/bdd/testcontext"
)

// iStoreTheResponseFieldAs stores a field value from the response body for later use
func iStoreTheResponseFieldAs(ctx context.Context, fieldName, variableName string) (context.Context, error) {
	body, ok := ctx.Value(testcontext.ContextKeyResponseBody).([]byte)
	if !ok {
		return ctx, fmt.Errorf("response body not found in context")
	}

	// Parse JSON response
	var jsonBody map[string]any
	if err := json.Unmarshal(body, &jsonBody); err != nil {
		return ctx, fmt.Errorf("failed to parse JSON response: %w. Body: %s", err, string(body))
	}

	// Get field value (supports nested fields with dot notation)
	value, exists := getNestedField(jsonBody, fieldName)
	if !exists {
		return ctx, fmt.Errorf("field '%s' not found in response body. Available fields: %v",
			fieldName, getMapKeys(jsonBody))
	}

	// Get or create stored variables map
	var storedVars map[string]any
	if existingVars, ok := ctx.Value(testcontext.ContextKeyStoredVariables).(map[string]any); ok {
		storedVars = existingVars
	} else {
		storedVars = make(map[string]any)
	}

	// Store the value
	storedVars[variableName] = value

	// Update context with stored variables
	return context.WithValue(ctx, testcontext.ContextKeyStoredVariables, storedVars), nil
}

// Helper function to get all stored variable names for error messages
func getStoredVariableNames(vars map[string]any) []string {
	names := make([]string, 0, len(vars))
	for name := range vars {
		names = append(names, name)
	}
	return names
}

// RegisterVariableSteps registers variable-related step definitions
func RegisterVariableSteps(sc *godog.ScenarioContext) {
	sc.Step(`^I store the response field "([^"]*)" as "([^"]*)"$`, iStoreTheResponseFieldAs)
}
