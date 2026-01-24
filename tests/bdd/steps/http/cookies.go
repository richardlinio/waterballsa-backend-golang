package http

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/cucumber/godog"
	"github.com/richardlinio/waterballsa-backend-golang/tests/bdd/testcontext"
)

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

// RegisterCookieSteps registers cookie-related step definitions
func RegisterCookieSteps(sc *godog.ScenarioContext) {
	sc.Step(`^I set cookie "([^"]*)" to "([^"]*)"$`, iSetCookie)
	sc.Step(`^I use stored cookie "([^"]*)"$`, iUseStoredCookie)
	sc.Step(`^I extract cookie "([^"]*)" from response$`, iExtractCookieFromResponse)
	sc.Step(`^the response should set cookie "([^"]*)"$`, theResponseShouldSetCookie)
	sc.Step(`^cookie "([^"]*)" should have attribute "([^"]*)"$`, cookieShouldHaveAttribute)
}
