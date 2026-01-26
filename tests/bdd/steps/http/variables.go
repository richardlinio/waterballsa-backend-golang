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
	varLastOrderID   = "lastOrderId"
)

// contextKeyMapping defines the mapping from variable names to context keys
var contextKeyMapping = map[string]testcontext.ContextKey{
	varLastJourneyID: testcontext.ContextKeyLastJourneyID,
	varLastChapterID: testcontext.ContextKeyLastChapterID,
	varLastMissionID: testcontext.ContextKeyLastMissionID,
	varLastUserID:    testcontext.ContextKeyLastUserID,
	varLastOrderID:   testcontext.ContextKeyLastOrderID,
}

// resolveVariable resolves a variable value from context
// It first checks context keys (using contextKeyMapping), then checks stored variables
// Returns (value, found) where found is true if the variable exists
func resolveVariable(ctx context.Context, varName string) (any, bool) {
	// Try context keys first
	if contextKey, exists := contextKeyMapping[varName]; exists {
		if val, ok := ctx.Value(contextKey).(int64); ok {
			return val, true
		}
	}

	// Try stored variables
	if storedVars, ok := ctx.Value(testcontext.ContextKeyStoredVariables).(map[string]any); ok {
		if val, exists := storedVars[varName]; exists {
			return val, true
		}
	}

	return nil, false
}

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
	// Use the unified resolveVariable function
	value, found := resolveVariable(ctx, sourceVarName)
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
