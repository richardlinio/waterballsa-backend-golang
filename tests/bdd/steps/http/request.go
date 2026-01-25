package http

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/cucumber/godog"
	"github.com/richardlinio/waterballsa-backend-golang/tests/bdd/testcontext"
)

// iSetRequestBodyTo stores the request body from a doc string for later use
// Supports variable substitution using {{variableName}} syntax
func iSetRequestBodyTo(ctx context.Context, docString *godog.DocString) (context.Context, error) {
	if docString == nil {
		return ctx, fmt.Errorf("doc string is nil")
	}

	// Replace variables in request body
	requestBody, err := replaceVariablesInRequestBody(ctx, docString.Content)
	if err != nil {
		return ctx, fmt.Errorf("failed to replace variables in request body: %w", err)
	}

	// Validate that the content is valid JSON after variable substitution
	var js json.RawMessage
	if err := json.Unmarshal([]byte(requestBody), &js); err != nil {
		return ctx, fmt.Errorf("invalid JSON in request body: %w", err)
	}

	return context.WithValue(ctx, testcontext.ContextKeyRequestBody, requestBody), nil
}

// iSendRequestTo makes an HTTP request with the specified method and path
func iSendRequestTo(ctx context.Context, method, path string) (context.Context, error) {
	// Get test server from suite context
	testServer, ok := ctx.Value(testcontext.ContextKeyTestServer).(*testcontext.TestServerWrapper)
	if !ok {
		return ctx, fmt.Errorf("test server not found in context")
	}

	// Replace variables in path (e.g., {{lastJourneyId}})
	path, err := replaceVariablesInPath(ctx, path)
	if err != nil {
		return ctx, fmt.Errorf("failed to replace variables in path: %w", err)
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

// iSetAuthorizationHeaderTo sets the Authorization header with a token value
// Supports variable substitution using {{variableName}} syntax
func iSetAuthorizationHeaderTo(ctx context.Context, tokenPlaceholder string) (context.Context, error) {
	// Check if tokenPlaceholder contains a variable reference (e.g., "{{token}}")
	if strings.HasPrefix(tokenPlaceholder, "{{") && strings.HasSuffix(tokenPlaceholder, "}}") {
		// Extract variable name (remove "{{" and "}}")
		variableName := strings.TrimSuffix(strings.TrimPrefix(tokenPlaceholder, "{{"), "}}")

		// Get stored variables from context
		storedVars, ok := ctx.Value(testcontext.ContextKeyStoredVariables).(map[string]any)
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

// iLoginAsWithPassword performs login and stores the access token
// Sends POST request to /auth/login with credentials
// Stores accessToken and lastUserId from response
func iLoginAsWithPassword(ctx context.Context, username, password string) (context.Context, error) {
	// Prepare login request body
	loginBody := fmt.Sprintf(`{"username": "%s", "password": "%s"}`, username, password)
	ctx = context.WithValue(ctx, testcontext.ContextKeyRequestBody, loginBody)

	// Send POST request to /auth/login
	ctx, err := iSendRequestTo(ctx, "POST", "/auth/login")
	if err != nil {
		return ctx, fmt.Errorf("failed to send login request: %w", err)
	}

	// Parse response to extract access token
	body, ok := ctx.Value(testcontext.ContextKeyResponseBody).([]byte)
	if !ok {
		return ctx, fmt.Errorf("response body not found in context")
	}

	// Parse JSON response
	var loginResponse map[string]any
	if err := json.Unmarshal(body, &loginResponse); err != nil {
		return ctx, fmt.Errorf("failed to parse login response: %w. Body: %s", err, string(body))
	}

	// Extract access token
	accessToken, ok := loginResponse["accessToken"].(string)
	if !ok {
		return ctx, fmt.Errorf("accessToken not found in login response. Response: %s", string(body))
	}

	// Extract user object
	userObj, ok := loginResponse["user"].(map[string]any)
	if !ok {
		return ctx, fmt.Errorf("user object not found in login response. Response: %s", string(body))
	}

	// Extract user ID from user object
	userID, ok := userObj["id"]
	if !ok {
		return ctx, fmt.Errorf("user.id not found in login response. Response: %s", string(body))
	}

	// Convert userID to int64 for context storage
	var userIDInt64 int64
	switch v := userID.(type) {
	case float64:
		userIDInt64 = int64(v)
	case int64:
		userIDInt64 = v
	default:
		return ctx, fmt.Errorf("user.id has unexpected type: %T", userID)
	}

	// Store user ID in context (for consistency with database steps)
	ctx = context.WithValue(ctx, testcontext.ContextKeyLastUserID, userIDInt64)

	// Get or create stored variables map
	var storedVars map[string]any
	if existingVars, ok := ctx.Value(testcontext.ContextKeyStoredVariables).(map[string]any); ok {
		storedVars = existingVars
	} else {
		storedVars = make(map[string]any)
	}

	// Store access token and user ID in variables map (for template substitution)
	storedVars["accessToken"] = accessToken
	storedVars["lastUserId"] = userIDInt64

	// Update context
	ctx = context.WithValue(ctx, testcontext.ContextKeyStoredVariables, storedVars)

	return ctx, nil
}

// RegisterRequestSteps registers request-related step definitions
func RegisterRequestSteps(sc *godog.ScenarioContext) {
	sc.Step(`^I set request body to:$`, iSetRequestBodyTo)
	sc.Step(`^I send "([^"]*)" request to "([^"]*)"$`, iSendRequestTo)
	sc.Step(`^I set Authorization header to "([^"]*)"$`, iSetAuthorizationHeaderTo)
	sc.Step(`^I login as "([^"]*)" with password "([^"]*)"$`, iLoginAsWithPassword)
}
