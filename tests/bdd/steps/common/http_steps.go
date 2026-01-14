package common

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strings"
	"time"

	"github.com/cucumber/godog"
	"github.com/linporu/waterballsa-backend-golang/tests/bdd/testcontext"
)

// defaultHTTPClient is a shared HTTP client for all test requests
// It reuses TCP connections via connection pooling for better performance
var defaultHTTPClient = &http.Client{
	Timeout: 30 * time.Second,
}

// iSetRequestBodyTo stores the request body from a doc string for later use
func iSetRequestBodyTo(ctx context.Context, docString *godog.DocString) (context.Context, error) {
	if docString == nil {
		return ctx, fmt.Errorf("doc string is nil")
	}

	// Validate that the content is valid JSON
	var js json.RawMessage
	if err := json.Unmarshal([]byte(docString.Content), &js); err != nil {
		return ctx, fmt.Errorf("invalid JSON in request body: %w", err)
	}

	return context.WithValue(ctx, testcontext.ContextKeyRequestBody, docString.Content), nil
}

// iSendRequestTo makes an HTTP request with the specified method and path
func iSendRequestTo(ctx context.Context, method, path string) (context.Context, error) {
	// Get test server from suite context
	testServer, ok := ctx.Value(testcontext.ContextKeyTestServer).(*testcontext.TestServerWrapper)
	if !ok {
		return ctx, fmt.Errorf("test server not found in context")
	}

	// Get request body if it was set
	var requestBody io.Reader
	if body, ok := ctx.Value(testcontext.ContextKeyRequestBody).(string); ok {
		requestBody = strings.NewReader(body)
	}

	// Build full URL
	url := testServer.Server.BaseURL() + path

	// Create HTTP request
	req, err := http.NewRequest(method, url, requestBody)
	if err != nil {
		return ctx, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	// Set Content-Type header for JSON
	if requestBody != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	// Send request using shared HTTP client
	resp, err := defaultHTTPClient.Do(req)
	if err != nil {
		return ctx, fmt.Errorf("failed to send HTTP request: %w", err)
	}

	// Read response body
	responseBody, err := io.ReadAll(resp.Body)
	if closeErr := resp.Body.Close(); closeErr != nil && err == nil {
		return ctx, fmt.Errorf("failed to close response body: %w", closeErr)
	}
	if err != nil {
		return ctx, fmt.Errorf("failed to read response body: %w", err)
	}

	// Store response and response body in context
	ctx = context.WithValue(ctx, testcontext.ContextKeyResponse, resp)
	ctx = context.WithValue(ctx, testcontext.ContextKeyResponseBody, responseBody)

	return ctx, nil
}

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
func theResponseBodyShouldContainField(ctx context.Context, fieldName string) error {
	body, ok := ctx.Value(testcontext.ContextKeyResponseBody).([]byte)
	if !ok {
		return fmt.Errorf("response body not found in context")
	}

	// Parse JSON response
	var jsonBody map[string]interface{}
	if err := json.Unmarshal(body, &jsonBody); err != nil {
		return fmt.Errorf("failed to parse JSON response: %w. Body: %s", err, string(body))
	}

	// Check if field exists
	if _, exists := jsonBody[fieldName]; !exists {
		return fmt.Errorf("field '%s' not found in response body. Available fields: %v",
			fieldName, getMapKeys(jsonBody))
	}

	return nil
}

// theResponseBodyFieldShouldEqualString asserts that a response field equals a specific string value
func theResponseBodyFieldShouldEqualString(ctx context.Context, fieldName, expectedValue string) error {
	body, ok := ctx.Value(testcontext.ContextKeyResponseBody).([]byte)
	if !ok {
		return fmt.Errorf("response body not found in context")
	}

	// Parse JSON response
	var jsonBody map[string]interface{}
	if err := json.Unmarshal(body, &jsonBody); err != nil {
		return fmt.Errorf("failed to parse JSON response: %w. Body: %s", err, string(body))
	}

	// Get field value
	actualValue, exists := jsonBody[fieldName]
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
func theResponseBodyFieldShouldEqualNumber(ctx context.Context, fieldName string, expectedValue float64) error {
	body, ok := ctx.Value(testcontext.ContextKeyResponseBody).([]byte)
	if !ok {
		return fmt.Errorf("response body not found in context")
	}

	// Parse JSON response
	var jsonBody map[string]interface{}
	if err := json.Unmarshal(body, &jsonBody); err != nil {
		return fmt.Errorf("failed to parse JSON response: %w. Body: %s", err, string(body))
	}

	// Get field value
	actualValue, exists := jsonBody[fieldName]
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

// Helper function to get all keys from a map for error messages
func getMapKeys(m map[string]interface{}) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// Helper function to convert various numeric types to float64
func toFloat64(val interface{}) (float64, bool) {
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

// RegisterHTTPSteps registers all HTTP-related step definitions
func RegisterHTTPSteps(sc *godog.ScenarioContext) {
	sc.Step(`^I set request body to:$`, iSetRequestBodyTo)
	sc.Step(`^I send "([^"]*)" request to "([^"]*)"$`, iSendRequestTo)
	sc.Step(`^the response status code should be (\d+)$`, theResponseStatusCodeShouldBe)
	sc.Step(`^the response body should contain field "([^"]*)"$`, theResponseBodyShouldContainField)
	sc.Step(`^the response body field "([^"]*)" should equal string "([^"]*)"$`, theResponseBodyFieldShouldEqualString)
	sc.Step(`^the response body field "([^"]*)" should equal number (.+)$`, theResponseBodyFieldShouldEqualNumber)
}
