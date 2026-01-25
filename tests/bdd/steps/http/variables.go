package http

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/cucumber/godog"
	"github.com/richardlinio/waterballsa-backend-golang/tests/bdd/testcontext"
)

const (
	varLastJourneyID = "lastJourneyId"
	varLastChapterID = "lastChapterId"
	varLastMissionID = "lastMissionId"
	varLastUserID    = "lastUserId"
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

// iCopyVariableTo copies a variable from one name to another
// Supports copying from context keys (lastJourneyId, lastUserId, etc.) or stored variables
func iCopyVariableTo(ctx context.Context, sourceVarName, targetVarName string) (context.Context, error) {
	var value any
	var found bool

	// Try context keys first
	switch sourceVarName {
	case varLastJourneyID:
		if val, ok := ctx.Value(testcontext.ContextKeyLastJourneyID).(int64); ok {
			value = val
			found = true
		}
	case varLastChapterID:
		if val, ok := ctx.Value(testcontext.ContextKeyLastChapterID).(int64); ok {
			value = val
			found = true
		}
	case varLastMissionID:
		if val, ok := ctx.Value(testcontext.ContextKeyLastMissionID).(int64); ok {
			value = val
			found = true
		}
	case varLastUserID:
		if val, ok := ctx.Value(testcontext.ContextKeyLastUserID).(int64); ok {
			value = val
			found = true
		}
	}

	// If not found in context keys, try stored variables
	if !found {
		if storedVars, ok := ctx.Value(testcontext.ContextKeyStoredVariables).(map[string]any); ok {
			if val, exists := storedVars[sourceVarName]; exists {
				value = val
				found = true
			}
		}
	}

	if !found {
		return ctx, fmt.Errorf("source variable '%s' not found in context or stored variables", sourceVarName)
	}

	// Get or create stored variables map
	var storedVars map[string]any
	if existingVars, ok := ctx.Value(testcontext.ContextKeyStoredVariables).(map[string]any); ok {
		storedVars = existingVars
	} else {
		storedVars = make(map[string]any)
	}

	// Store the value with the new name
	storedVars[targetVarName] = value

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
	sc.Step(`^I copy variable "([^"]*)" to "([^"]*)"$`, iCopyVariableTo)
}
