package http

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/cucumber/godog"
	"github.com/linporu/waterballsa-backend-golang/tests/bdd/testcontext"
)

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

// RegisterRequestSteps registers request-related step definitions
func RegisterRequestSteps(sc *godog.ScenarioContext) {
	sc.Step(`^I set request body to:$`, iSetRequestBodyTo)
	sc.Step(`^I send "([^"]*)" request to "([^"]*)"$`, iSendRequestTo)
	sc.Step(`^I set Authorization header to "([^"]*)"$`, iSetAuthorizationHeaderTo)
}
