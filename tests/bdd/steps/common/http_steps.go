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

	// Add cookies to request if they were set
	if cookies, ok := ctx.Value(testcontext.ContextKeyCookies).(map[string]string); ok {
		for name, value := range cookies {
			req.AddCookie(&http.Cookie{
				Name:  name,
				Value: value,
			})
		}
	}

	// Set Authorization header if present in context
	if authHeader, ok := ctx.Value(testcontext.ContextKeyAuthHeader).(string); ok {
		req.Header.Set("Authorization", authHeader)
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
// Supports nested field access using dot notation (e.g., "user.id")
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
	var jsonBody map[string]interface{}
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
	var jsonBody map[string]interface{}
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

// Helper function to get nested field from JSON object using dot notation
// Supports accessing nested fields like "user.id" or "data.user.profile.name"
func getNestedField(data map[string]interface{}, fieldPath string) (interface{}, bool) {
	// Split the field path by dots
	parts := strings.Split(fieldPath, ".")

	var current interface{} = data

	// Navigate through each part of the path
	for _, part := range parts {
		// Check if current is a map
		currentMap, ok := current.(map[string]interface{})
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

	return current, true
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

// iSetCookie sets a cookie to be sent with the next HTTP request
// Supports "stored cookieName" syntax to use a previously extracted cookie
func iSetCookie(ctx context.Context, cookieName, cookieValue string) (context.Context, error) {
	// Check if cookieValue is a reference to a stored cookie (format: "stored cookie_name")
	if strings.HasPrefix(cookieValue, "stored ") {
		storedCookieName := strings.TrimPrefix(cookieValue, "stored ")
		storedCookies, _ := ctx.Value(testcontext.ContextKeyStoredCookies).(map[string]string)
		if storedCookies == nil {
			return ctx, fmt.Errorf("no stored cookies found")
		}

		storedValue, ok := storedCookies[storedCookieName]
		if !ok {
			return ctx, fmt.Errorf("stored cookie '%s' not found", storedCookieName)
		}
		cookieValue = storedValue
	}

	// Get existing cookies or create new map
	cookies, _ := ctx.Value(testcontext.ContextKeyCookies).(map[string]string)
	if cookies == nil {
		cookies = make(map[string]string)
	}

	// Set the cookie
	cookies[cookieName] = cookieValue

	return context.WithValue(ctx, testcontext.ContextKeyCookies, cookies), nil
}

// iExtractCookieFromResponse extracts a cookie from the response and stores it for later use
func iExtractCookieFromResponse(ctx context.Context, cookieName string) (context.Context, error) {
	resp, ok := ctx.Value(testcontext.ContextKeyResponse).(*http.Response)
	if !ok {
		return ctx, fmt.Errorf("response not found in context")
	}

	// Find the cookie in response
	var foundCookie *http.Cookie
	for _, cookie := range resp.Cookies() {
		if cookie.Name == cookieName {
			foundCookie = cookie
			break
		}
	}

	if foundCookie == nil {
		return ctx, fmt.Errorf("cookie '%s' not found in response", cookieName)
	}

	// Get existing stored cookies or create new map
	storedCookies, _ := ctx.Value(testcontext.ContextKeyStoredCookies).(map[string]string)
	if storedCookies == nil {
		storedCookies = make(map[string]string)
	}

	// Store the cookie value
	storedCookies[cookieName] = foundCookie.Value

	return context.WithValue(ctx, testcontext.ContextKeyStoredCookies, storedCookies), nil
}

// iUseStoredCookie sets a cookie using a previously extracted value
func iUseStoredCookie(ctx context.Context, cookieName string) (context.Context, error) {
	storedCookies, _ := ctx.Value(testcontext.ContextKeyStoredCookies).(map[string]string)
	if storedCookies == nil {
		return ctx, fmt.Errorf("no stored cookies found")
	}

	cookieValue, ok := storedCookies[cookieName]
	if !ok {
		return ctx, fmt.Errorf("stored cookie '%s' not found", cookieName)
	}

	// Get existing cookies or create new map
	cookies, _ := ctx.Value(testcontext.ContextKeyCookies).(map[string]string)
	if cookies == nil {
		cookies = make(map[string]string)
	}

	cookies[cookieName] = cookieValue

	return context.WithValue(ctx, testcontext.ContextKeyCookies, cookies), nil
}

// theResponseShouldSetCookie verifies that a cookie is set in the response
func theResponseShouldSetCookie(ctx context.Context, cookieName string) error {
	resp, ok := ctx.Value(testcontext.ContextKeyResponse).(*http.Response)
	if !ok {
		return fmt.Errorf("response not found in context")
	}

	// Check if cookie exists in response
	for _, cookie := range resp.Cookies() {
		if cookie.Name == cookieName {
			return nil
		}
	}

	return fmt.Errorf("cookie '%s' not found in response", cookieName)
}

// cookieShouldHaveAttribute verifies that a cookie has a specific attribute
func cookieShouldHaveAttribute(ctx context.Context, cookieName, attribute string) error {
	resp, ok := ctx.Value(testcontext.ContextKeyResponse).(*http.Response)
	if !ok {
		return fmt.Errorf("response not found in context")
	}

	// Find the cookie
	var foundCookie *http.Cookie
	for _, cookie := range resp.Cookies() {
		if cookie.Name == cookieName {
			foundCookie = cookie
			break
		}
	}

	if foundCookie == nil {
		return fmt.Errorf("cookie '%s' not found in response", cookieName)
	}

	// Check the attribute
	switch attribute {
	case "HttpOnly":
		if !foundCookie.HttpOnly {
			return fmt.Errorf("cookie '%s' does not have HttpOnly attribute", cookieName)
		}
	case "SameSite=Strict":
		if foundCookie.SameSite != http.SameSiteStrictMode {
			return fmt.Errorf("cookie '%s' does not have SameSite=Strict attribute (got %v)", cookieName, foundCookie.SameSite)
		}
	case "SameSite=Lax":
		if foundCookie.SameSite != http.SameSiteLaxMode {
			return fmt.Errorf("cookie '%s' does not have SameSite=Lax attribute (got %v)", cookieName, foundCookie.SameSite)
		}
	case "Secure":
		if !foundCookie.Secure {
			return fmt.Errorf("cookie '%s' does not have Secure attribute", cookieName)
		}
	default:
		return fmt.Errorf("unsupported cookie attribute: %s", attribute)
	}

	return nil
}

// iStoreTheResponseFieldAs stores a field value from the response body for later use
func iStoreTheResponseFieldAs(ctx context.Context, fieldName, variableName string) (context.Context, error) {
	body, ok := ctx.Value(testcontext.ContextKeyResponseBody).([]byte)
	if !ok {
		return ctx, fmt.Errorf("response body not found in context")
	}

	// Parse JSON response
	var jsonBody map[string]interface{}
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
	var storedVars map[string]interface{}
	if existingVars, ok := ctx.Value(testcontext.ContextKeyStoredVariables).(map[string]interface{}); ok {
		storedVars = existingVars
	} else {
		storedVars = make(map[string]interface{})
	}

	// Store the value
	storedVars[variableName] = value

	// Update context with stored variables
	return context.WithValue(ctx, testcontext.ContextKeyStoredVariables, storedVars), nil
}

// iSetAuthorizationHeaderTo sets the Authorization header with a token value
// Supports variable substitution using {{variableName}} syntax
func iSetAuthorizationHeaderTo(ctx context.Context, tokenPlaceholder string) (context.Context, error) {
	// Check if tokenPlaceholder contains a variable reference (e.g., "{{token}}")
	if strings.HasPrefix(tokenPlaceholder, "{{") && strings.HasSuffix(tokenPlaceholder, "}}") {
		// Extract variable name (remove "{{" and "}}")
		variableName := strings.TrimSuffix(strings.TrimPrefix(tokenPlaceholder, "{{"), "}}")

		// Get stored variables from context
		storedVars, ok := ctx.Value(testcontext.ContextKeyStoredVariables).(map[string]interface{})
		if !ok || storedVars == nil {
			return ctx, fmt.Errorf("no stored variables found in context. Did you forget to store the variable '%s'?", variableName)
		}

		// Get the actual token value
		tokenValue, exists := storedVars[variableName]
		if !exists {
			return ctx, fmt.Errorf("variable '%s' not found in stored variables. Available variables: %v",
				variableName, getStoredVariableNames(storedVars))
		}

		// Convert token value to string
		tokenStr := fmt.Sprintf("%v", tokenValue)

		// Store as Authorization header with "Bearer " prefix
		authHeader := "Bearer " + tokenStr
		return context.WithValue(ctx, testcontext.ContextKeyAuthHeader, authHeader), nil
	}

	// If not a variable reference, use the value directly
	authHeader := "Bearer " + tokenPlaceholder
	return context.WithValue(ctx, testcontext.ContextKeyAuthHeader, authHeader), nil
}

// Helper function to get all stored variable names for error messages
func getStoredVariableNames(vars map[string]interface{}) []string {
	names := make([]string, 0, len(vars))
	for name := range vars {
		names = append(names, name)
	}
	return names
}

// RegisterHTTPSteps registers all HTTP-related step definitions
func RegisterHTTPSteps(sc *godog.ScenarioContext) {
	sc.Step(`^I set request body to:$`, iSetRequestBodyTo)
	sc.Step(`^I send "([^"]*)" request to "([^"]*)"$`, iSendRequestTo)
	sc.Step(`^the response status code should be (\d+)$`, theResponseStatusCodeShouldBe)
	sc.Step(`^the response body should contain field "([^"]*)"$`, theResponseBodyShouldContainField)
	sc.Step(`^the response body field "([^"]*)" should equal string "([^"]*)"$`, theResponseBodyFieldShouldEqualString)
	sc.Step(`^the response body field "([^"]*)" should equal number (.+)$`, theResponseBodyFieldShouldEqualNumber)
	sc.Step(`^I set cookie "([^"]*)" to "([^"]*)"$`, iSetCookie)
	sc.Step(`^I use stored cookie "([^"]*)"$`, iUseStoredCookie)
	sc.Step(`^I extract cookie "([^"]*)" from response$`, iExtractCookieFromResponse)
	sc.Step(`^the response should set cookie "([^"]*)"$`, theResponseShouldSetCookie)
	sc.Step(`^cookie "([^"]*)" should have attribute "([^"]*)"$`, cookieShouldHaveAttribute)
	sc.Step(`^I store the response field "([^"]*)" as "([^"]*)"$`, iStoreTheResponseFieldAs)
	sc.Step(`^I set Authorization header to "([^"]*)"$`, iSetAuthorizationHeaderTo)
}
