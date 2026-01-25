package http

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/cucumber/godog"
	"github.com/richardlinio/waterballsa-backend-golang/tests/bdd/testcontext"
)

// theResponseStatusCodeShouldBe asserts that the response status code matches expected value
func theResponseStatusCodeShouldBe(ctx context.Context, expectedStatusCode int) error {
	resp, ok := ctx.Value(testcontext.ContextKeyResponse).(*http.Response)
	if !ok {
		return fmt.Errorf("response not found in context")
	}

	if resp.StatusCode != expectedStatusCode {
		// Include response body in error for debugging
		body, ok := ctx.Value(testcontext.ContextKeyResponseBody).([]byte)
		if ok {
			return fmt.Errorf("expected status code %d, got %d. Response body: %s",
				expectedStatusCode, resp.StatusCode, string(body))
		}
		return fmt.Errorf("expected status code %d, got %d",
			expectedStatusCode, resp.StatusCode)
	}

	return nil
}

// theResponseBodyShouldContainField asserts that the response body contains the specified field
// Supports nested field access using dot notation (e.g., "user.id")
func theResponseBodyShouldContainField(ctx context.Context, fieldName string) error {
	body, ok := ctx.Value(testcontext.ContextKeyResponseBody).([]byte)
	if !ok {
		return fmt.Errorf("response body not found in context")
	}

	// Parse JSON response
	var jsonBody map[string]any
	if err := json.Unmarshal(body, &jsonBody); err != nil {
		return fmt.Errorf("failed to parse JSON response: %w. Body: %s", err, string(body))
	}

	// Check if field exists (supports nested fields with dot notation)
	_, exists := getNestedField(jsonBody, fieldName)
	if !exists {
		return fmt.Errorf("field '%s' not found in response body. Available fields: %v",
			fieldName, getMapKeys(jsonBody))
	}

	return nil
}

// theResponseBodyFieldShouldEqualString asserts that a response field equals a specific string value
// Supports nested field access using dot notation (e.g., "user.username")
func theResponseBodyFieldShouldEqualString(ctx context.Context, fieldName, expectedValue string) error {
	body, ok := ctx.Value(testcontext.ContextKeyResponseBody).([]byte)
	if !ok {
		return fmt.Errorf("response body not found in context")
	}

	// Parse JSON response
	var jsonBody map[string]any
	if err := json.Unmarshal(body, &jsonBody); err != nil {
		return fmt.Errorf("failed to parse JSON response: %w. Body: %s", err, string(body))
	}

	// Get field value (supports nested fields with dot notation)
	actualValue, exists := getNestedField(jsonBody, fieldName)
	if !exists {
		return fmt.Errorf("field '%s' not found in response body. Available fields: %v",
			fieldName, getMapKeys(jsonBody))
	}

	// Convert actual value to string for comparison
	actualStr := fmt.Sprintf("%v", actualValue)

	// Compare values
	if actualStr != expectedValue {
		return fmt.Errorf("field '%s' expected to be '%s', but got '%s'",
			fieldName, expectedValue, actualStr)
	}

	return nil
}

// theResponseBodyFieldShouldEqualNumber asserts that a response field equals a specific number value
// Supports nested field access using dot notation (e.g., "user.experience")
func theResponseBodyFieldShouldEqualNumber(ctx context.Context, fieldName string, expectedValue float64) error {
	body, ok := ctx.Value(testcontext.ContextKeyResponseBody).([]byte)
	if !ok {
		return fmt.Errorf("response body not found in context")
	}

	// Parse JSON response
	var jsonBody map[string]any
	if err := json.Unmarshal(body, &jsonBody); err != nil {
		return fmt.Errorf("failed to parse JSON response: %w. Body: %s", err, string(body))
	}

	// Get field value (supports nested fields with dot notation)
	actualValue, exists := getNestedField(jsonBody, fieldName)
	if !exists {
		return fmt.Errorf("field '%s' not found in response body. Available fields: %v",
			fieldName, getMapKeys(jsonBody))
	}

	// Convert to float64 for comparison
	actualFloat, ok := toFloat64(actualValue)
	if !ok {
		return fmt.Errorf("field '%s' is not a number, got type %v with value '%v'",
			fieldName, reflect.TypeOf(actualValue), actualValue)
	}

	// Compare values
	if actualFloat != expectedValue {
		return fmt.Errorf("field '%s' expected to be %v, but got %v",
			fieldName, expectedValue, actualFloat)
	}

	return nil
}

// theResponseBodyFieldShouldHaveSize asserts that a response field (array) has a specific size
// Supports nested field access using dot notation (e.g., "data.users")
func theResponseBodyFieldShouldHaveSize(ctx context.Context, fieldName string, expectedSize int) error {
	body, ok := ctx.Value(testcontext.ContextKeyResponseBody).([]byte)
	if !ok {
		return fmt.Errorf("response body not found in context")
	}

	// Parse JSON response
	var jsonBody map[string]any
	if err := json.Unmarshal(body, &jsonBody); err != nil {
		return fmt.Errorf("failed to parse JSON response: %w. Body: %s", err, string(body))
	}

	// Get field value (supports nested fields with dot notation)
	actualValue, exists := getNestedField(jsonBody, fieldName)
	if !exists {
		return fmt.Errorf("field '%s' not found in response body. Available fields: %v",
			fieldName, getMapKeys(jsonBody))
	}

	// Check if the field is an array/slice
	actualSlice, ok := actualValue.([]any)
	if !ok {
		return fmt.Errorf("field '%s' is not an array, got type %v with value '%v'",
			fieldName, reflect.TypeOf(actualValue), actualValue)
	}

	// Compare sizes
	actualSize := len(actualSlice)
	if actualSize != expectedSize {
		return fmt.Errorf("field '%s' expected to have size %d, but got %d",
			fieldName, expectedSize, actualSize)
	}

	return nil
}

// theResponseBodyShouldHaveNonNullField asserts that the response body contains the specified field AND it is not null
// Supports nested field access using dot notation (e.g., "user.id")
func theResponseBodyShouldHaveNonNullField(ctx context.Context, fieldName string) error {
	body, ok := ctx.Value(testcontext.ContextKeyResponseBody).([]byte)
	if !ok {
		return fmt.Errorf("response body not found in context")
	}

	// Parse JSON response
	var jsonBody map[string]any
	if err := json.Unmarshal(body, &jsonBody); err != nil {
		return fmt.Errorf("failed to parse JSON response: %w. Body: %s", err, string(body))
	}

	// Check if field exists (supports nested fields with dot notation)
	value, exists := getNestedField(jsonBody, fieldName)
	if !exists {
		return fmt.Errorf("field '%s' not found in response body. Available fields: %v",
			fieldName, getMapKeys(jsonBody))
	}

	// Check if value is null
	if value == nil {
		return fmt.Errorf("field '%s' exists but is null", fieldName)
	}

	return nil
}

// theResponseBodyFieldShouldBeNull asserts that a response field exists AND its value is null
// Supports nested field access using dot notation (e.g., "user.deletedAt")
func theResponseBodyFieldShouldBeNull(ctx context.Context, fieldName string) error {
	body, ok := ctx.Value(testcontext.ContextKeyResponseBody).([]byte)
	if !ok {
		return fmt.Errorf("response body not found in context")
	}

	// Parse JSON response
	var jsonBody map[string]any
	if err := json.Unmarshal(body, &jsonBody); err != nil {
		return fmt.Errorf("failed to parse JSON response: %w. Body: %s", err, string(body))
	}

	// Check if field exists (supports nested fields with dot notation)
	value, exists := getNestedField(jsonBody, fieldName)
	if !exists {
		return fmt.Errorf("field '%s' not found in response body. Available fields: %v",
			fieldName, getMapKeys(jsonBody))
	}

	// Check if value is null
	if value != nil {
		return fmt.Errorf("field '%s' expected to be null, but got '%v' (type: %T)",
			fieldName, value, value)
	}

	return nil
}

// theResponseBodyFieldShouldContain asserts that a response field value contains a specific substring
// Supports nested field access using dot notation and variable substitution
func theResponseBodyFieldShouldContain(ctx context.Context, fieldName, expectedSubstring string) error {
	body, ok := ctx.Value(testcontext.ContextKeyResponseBody).([]byte)
	if !ok {
		return fmt.Errorf("response body not found in context")
	}

	// Replace variables in expected substring
	expectedSubstring, err := replaceVariablesInString(ctx, expectedSubstring)
	if err != nil {
		return fmt.Errorf("failed to replace variables in expected substring: %w", err)
	}

	// Parse JSON response
	var jsonBody map[string]any
	if err := json.Unmarshal(body, &jsonBody); err != nil {
		return fmt.Errorf("failed to parse JSON response: %w. Body: %s", err, string(body))
	}

	// Get field value (supports nested fields with dot notation)
	actualValue, exists := getNestedField(jsonBody, fieldName)
	if !exists {
		return fmt.Errorf("field '%s' not found in response body. Available fields: %v",
			fieldName, getMapKeys(jsonBody))
	}

	// Convert actual value to string for comparison
	actualStr := fmt.Sprintf("%v", actualValue)

	// Check if substring is contained
	if !strings.Contains(actualStr, expectedSubstring) {
		return fmt.Errorf("field '%s' expected to contain '%s', but got '%s'",
			fieldName, expectedSubstring, actualStr)
	}

	return nil
}

// RegisterResponseSteps registers response-related step definitions
func RegisterResponseSteps(sc *godog.ScenarioContext) {
	sc.Step(`^the response status code should be (\d+)$`, theResponseStatusCodeShouldBe)
	sc.Step(`^the response body should contain field "([^"]*)"$`, theResponseBodyShouldContainField)
	sc.Step(`^the response body should have non-null field "([^"]*)"$`, theResponseBodyShouldHaveNonNullField)
	sc.Step(`^the response body field "([^"]*)" should equal string "([^"]*)"$`, theResponseBodyFieldShouldEqualString)
	sc.Step(`^the response body field "([^"]*)" should equal number (.+)$`, theResponseBodyFieldShouldEqualNumber)
	sc.Step(`^the response body field "([^"]*)" should equal decimal "([^"]*)"$`, theResponseBodyFieldShouldEqualNumber)
	sc.Step(`^the response body field "([^"]*)" should have size (\d+)$`, theResponseBodyFieldShouldHaveSize)
	sc.Step(`^the response body field "([^"]*)" should be null$`, theResponseBodyFieldShouldBeNull)
	sc.Step(`^the response body field "([^"]*)" should contain "([^"]*)"$`, theResponseBodyFieldShouldContain)
}
