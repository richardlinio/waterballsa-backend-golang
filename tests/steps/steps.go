package steps

import (
	"sync"

	"github.com/linporu/waterballsa-backend-golang/tests/testutil"
)

// contextKey is the type for context keys used in step definitions
type contextKey string

// Context keys for storing test state
const (
	contextKeyTestServer   contextKey = "testServer"
	contextKeyRequestBody  contextKey = "requestBody"
	contextKeyResponse     contextKey = "response"
	contextKeyResponseBody contextKey = "responseBody"
)

// TestServerWrapper wraps the test server for sharing across scenarios
type TestServerWrapper struct {
	Server *testutil.TestServer
	mu     sync.Mutex
}
