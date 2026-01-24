package testcontext

import (
	"sync"

	"github.com/richardlinio/waterballsa-backend-golang/tests/testutil"
)

// ContextKey is the type for context keys used in step definitions
type ContextKey string

// Context keys for storing test state
const (
	ContextKeyTestServer      ContextKey = "testServer"
	ContextKeyRequestBody     ContextKey = "requestBody"
	ContextKeyResponse        ContextKey = "response"
	ContextKeyResponseBody    ContextKey = "responseBody"
	ContextKeyCookies         ContextKey = "cookies"         // Cookies to be sent with next request
	ContextKeyStoredCookies   ContextKey = "storedCookies"   // Extracted cookies from response
	ContextKeyStoredVariables ContextKey = "storedVariables" // map[string]any for storing values between steps
	ContextKeyAuthHeader      ContextKey = "authHeader"      // string for Authorization header value
	ContextKeyLastJourneyID   ContextKey = "lastJourneyId"   // int64 for last created journey ID
	ContextKeyLastChapterID   ContextKey = "lastChapterId"   // int64 for last created chapter ID
	ContextKeyLastMissionID   ContextKey = "lastMissionId"   // int64 for last created mission ID
)

// TestServerWrapper wraps the test server for sharing across scenarios
type TestServerWrapper struct {
	Server *testutil.TestServer
	Mu     sync.Mutex
}
