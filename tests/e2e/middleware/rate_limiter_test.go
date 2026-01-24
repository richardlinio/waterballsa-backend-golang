package e2e

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/richardlinio/waterballsa-backend-golang/internal/config"
	"github.com/richardlinio/waterballsa-backend-golang/internal/middleware"
)

const (
	testIPDefault = "192.168.1.1"
)

// errorResponse represents the expected error response structure
type errorResponse struct {
	Code  string `json:"code"`
	Error string `json:"error"`
}

// setupTestServer creates a test HTTP server with rate limiter middleware
func setupTestServer(rateLimitConfig config.RateLimitConfig) *httptest.Server {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Create a no-op logger for tests
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	// Create a default JWT config for tests
	jwtConfig := config.JWTConfig{
		Realm: "test",
	}

	// Add middlewares in the same order as the real application
	router.Use(gin.Recovery())
	router.Use(middleware.ErrorHandler(logger, jwtConfig))
	router.Use(middleware.RateLimit(rateLimitConfig, logger))

	// Setup test endpoints
	router.GET("/api/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	router.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})

	return httptest.NewServer(router)
}

// makeRequest makes an HTTP GET request to the server with optional custom IP
func makeRequest(server *httptest.Server, path string, customIP string) (*http.Response, error) {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	req, err := http.NewRequest(http.MethodGet, server.URL+path, nil)
	if err != nil {
		return nil, err
	}

	// Set custom IP via X-Forwarded-For header if provided
	if customIP != "" {
		req.Header.Set("X-Forwarded-For", customIP)
	}

	return client.Do(req)
}

// parseErrorResponse parses the JSON error response
func parseErrorResponse(resp *http.Response) (*errorResponse, error) {
	body, err := io.ReadAll(resp.Body)
	if closeErr := resp.Body.Close(); closeErr != nil && err == nil {
		err = closeErr
	}
	if err != nil {
		return nil, err
	}

	var errResp errorResponse
	if err := json.Unmarshal(body, &errResp); err != nil {
		return nil, err
	}

	return &errResp, nil
}

// testRequestsWithinLimit verifies that requests within the rate limit are allowed
// and requests exceeding the limit are blocked with proper error response
func testRequestsWithinLimit(t *testing.T, server *httptest.Server) {
	// Send 3 requests - all should succeed
	for i := 0; i < 3; i++ {
		resp, err := makeRequest(server, "/api/test", testIPDefault)
		if err != nil {
			t.Fatalf("Request %d failed: %v", i+1, err)
		}
		defer func() {
			if err := resp.Body.Close(); err != nil {
				t.Logf("Failed to close response body: %v", err)
			}
		}()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Request %d: expected status 200, got %d", i+1, resp.StatusCode)
		}
	}

	// 4th request should be blocked
	resp, err := makeRequest(server, "/api/test", testIPDefault)
	if err != nil {
		t.Fatalf("4th request failed: %v", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			t.Logf("Failed to close response body: %v", err)
		}
	}()

	if resp.StatusCode != http.StatusTooManyRequests {
		t.Errorf("4th request: expected status 429, got %d", resp.StatusCode)
	}

	// Verify error response
	errResp, err := parseErrorResponse(resp)
	if err != nil {
		t.Fatalf("Failed to parse error response: %v", err)
	}

	if errResp.Code != "ERR_RATE_LIMIT_EXCEEDED" {
		t.Errorf("Expected error code ERR_RATE_LIMIT_EXCEEDED, got %s", errResp.Code)
	}

	if errResp.Error != "請求次數過多,請稍後再試" {
		t.Errorf("Expected error message '請求次數過多,請稍後再試', got '%s'", errResp.Error)
	}
}

// testWindowReset verifies that the rate limit window resets after the configured duration
func testWindowReset(t *testing.T, server *httptest.Server) {
	// Send 2 requests - should succeed
	for i := 0; i < 2; i++ {
		resp, err := makeRequest(server, "/api/test", testIPDefault)
		if err != nil {
			t.Fatalf("Request %d failed: %v", i+1, err)
		}
		defer func() {
			if err := resp.Body.Close(); err != nil {
				t.Logf("Failed to close response body: %v", err)
			}
		}()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Request %d: expected status 200, got %d", i+1, resp.StatusCode)
		}
	}

	// 3rd request should be blocked
	resp, err := makeRequest(server, "/api/test", testIPDefault)
	if err != nil {
		t.Fatalf("3rd request failed: %v", err)
	}
	if err := resp.Body.Close(); err != nil {
		t.Errorf("Failed to close response body: %v", err)
	}

	if resp.StatusCode != http.StatusTooManyRequests {
		t.Errorf("3rd request: expected status 429, got %d", resp.StatusCode)
	}

	// Wait for window to reset
	time.Sleep(150 * time.Millisecond)

	// New request should succeed after window reset
	resp, err = makeRequest(server, "/api/test", testIPDefault)
	if err != nil {
		t.Fatalf("Request after window reset failed: %v", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			t.Logf("Failed to close response body: %v", err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Request after window reset: expected status 200, got %d", resp.StatusCode)
	}
}

// testMultipleIPsIndependent verifies that different IP addresses have independent rate limits
func testMultipleIPsIndependent(t *testing.T, server *httptest.Server) {
	ip1 := testIPDefault
	ip2 := "192.168.1.2"

	// IP1: Send 2 requests - should succeed
	for i := 0; i < 2; i++ {
		resp, err := makeRequest(server, "/api/test", ip1)
		if err != nil {
			t.Fatalf("IP1 request %d failed: %v", i+1, err)
		}
		if err := resp.Body.Close(); err != nil {
			t.Errorf("Failed to close response body: %v", err)
		}

		if resp.StatusCode != http.StatusOK {
			t.Errorf("IP1 request %d: expected status 200, got %d", i+1, resp.StatusCode)
		}
	}

	// IP1: 3rd request should be blocked
	resp, err := makeRequest(server, "/api/test", ip1)
	if err != nil {
		t.Fatalf("IP1 3rd request failed: %v", err)
	}
	if err := resp.Body.Close(); err != nil {
		t.Errorf("Failed to close response body: %v", err)
	}

	if resp.StatusCode != http.StatusTooManyRequests {
		t.Errorf("IP1 3rd request: expected status 429, got %d", resp.StatusCode)
	}

	// IP2: Should have independent limit - 2 requests should succeed
	for i := 0; i < 2; i++ {
		resp, err := makeRequest(server, "/api/test", ip2)
		if err != nil {
			t.Fatalf("IP2 request %d failed: %v", i+1, err)
		}
		if err := resp.Body.Close(); err != nil {
			t.Errorf("Failed to close response body: %v", err)
		}

		if resp.StatusCode != http.StatusOK {
			t.Errorf("IP2 request %d: expected status 200, got %d", i+1, resp.StatusCode)
		}
	}

	// IP2: 3rd request should be blocked
	resp, err = makeRequest(server, "/api/test", ip2)
	if err != nil {
		t.Fatalf("IP2 3rd request failed: %v", err)
	}
	if err := resp.Body.Close(); err != nil {
		t.Errorf("Failed to close response body: %v", err)
	}

	if resp.StatusCode != http.StatusTooManyRequests {
		t.Errorf("IP2 3rd request: expected status 429, got %d", resp.StatusCode)
	}
}

// testConcurrentRequests verifies that concurrent requests from the same IP respect rate limits
func testConcurrentRequests(t *testing.T, server *httptest.Server) {
	numRequests := 20
	var wg sync.WaitGroup
	results := make(chan int, numRequests)

	// Launch concurrent requests
	for i := 0; i < numRequests; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			resp, err := makeRequest(server, "/api/test", testIPDefault)
			if err != nil {
				t.Errorf("Concurrent request failed: %v", err)
				return
			}
			defer func() {
				if err := resp.Body.Close(); err != nil {
					t.Logf("Failed to close response body: %v", err)
				}
			}()
			results <- resp.StatusCode
		}()
	}

	wg.Wait()
	close(results)

	// Count successful and rate-limited requests
	var successCount, rateLimitedCount int
	for statusCode := range results {
		switch statusCode {
		case http.StatusOK:
			successCount++
		case http.StatusTooManyRequests:
			rateLimitedCount++
		default:
			t.Errorf("Unexpected status code: %d", statusCode)
		}
	}

	// Verify exactly 10 requests succeeded
	if successCount != 10 {
		t.Errorf("Expected 10 successful requests, got %d", successCount)
	}

	// Verify exactly 10 requests were rate-limited
	if rateLimitedCount != 10 {
		t.Errorf("Expected 10 rate-limited requests, got %d", rateLimitedCount)
	}
}

// testRateLimitingDisabled verifies that when rate limiting is disabled, all requests succeed
func testRateLimitingDisabled(t *testing.T, server *httptest.Server) {
	// Send 20 requests - all should succeed when rate limiting is disabled
	for i := 0; i < 20; i++ {
		resp, err := makeRequest(server, "/api/test", testIPDefault)
		if err != nil {
			t.Fatalf("Request %d failed: %v", i+1, err)
		}
		defer func() {
			if err := resp.Body.Close(); err != nil {
				t.Logf("Failed to close response body: %v", err)
			}
		}()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Request %d: expected status 200, got %d (rate limiting should be disabled)", i+1, resp.StatusCode)
		}
	}
}

// testHealthCheckBypass verifies that health check endpoints bypass rate limiting
func testHealthCheckBypass(t *testing.T, server *httptest.Server) {
	// Send 2 requests to regular endpoint - should succeed
	for i := 0; i < 2; i++ {
		resp, err := makeRequest(server, "/api/test", testIPDefault)
		if err != nil {
			t.Fatalf("Regular endpoint request %d failed: %v", i+1, err)
		}
		if err := resp.Body.Close(); err != nil {
			t.Errorf("Failed to close response body: %v", err)
		}

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Regular endpoint request %d: expected status 200, got %d", i+1, resp.StatusCode)
		}
	}

	// 3rd request to regular endpoint should be blocked
	resp, err := makeRequest(server, "/api/test", testIPDefault)
	if err != nil {
		t.Fatalf("Regular endpoint 3rd request failed: %v", err)
	}
	if err := resp.Body.Close(); err != nil {
		t.Errorf("Failed to close response body: %v", err)
	}

	if resp.StatusCode != http.StatusTooManyRequests {
		t.Errorf("Regular endpoint 3rd request: expected status 429, got %d", resp.StatusCode)
	}

	// Health check endpoint should always work regardless of rate limit
	for i := 0; i < 10; i++ {
		resp, err := makeRequest(server, "/healthz", testIPDefault)
		if err != nil {
			t.Fatalf("Health check request %d failed: %v", i+1, err)
		}
		defer func() {
			if err := resp.Body.Close(); err != nil {
				t.Logf("Failed to close response body: %v", err)
			}
		}()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Health check request %d: expected status 200, got %d (should bypass rate limiting)", i+1, resp.StatusCode)
		}
	}

	// Regular endpoint should still be rate-limited
	resp, err = makeRequest(server, "/api/test", testIPDefault)
	if err != nil {
		t.Fatalf("Final regular endpoint request failed: %v", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			t.Logf("Failed to close response body: %v", err)
		}
	}()

	if resp.StatusCode != http.StatusTooManyRequests {
		t.Errorf("Final regular endpoint request: expected status 429, got %d", resp.StatusCode)
	}
}

func TestRateLimiterE2E(t *testing.T) {
	tests := []struct {
		name     string
		config   config.RateLimitConfig
		testFunc func(*testing.T, *httptest.Server)
	}{
		{
			name: "requests within limit are allowed",
			config: config.RateLimitConfig{
				Enabled:           true,
				RequestsPerMinute: 3,
				WindowDuration:    100 * time.Millisecond,
				CleanupInterval:   1 * time.Minute,
				RecordExpiry:      2 * time.Minute,
			},
			testFunc: testRequestsWithinLimit,
		},
		{
			name: "window reset allows new requests",
			config: config.RateLimitConfig{
				Enabled:           true,
				RequestsPerMinute: 2,
				WindowDuration:    100 * time.Millisecond,
				CleanupInterval:   1 * time.Minute,
				RecordExpiry:      2 * time.Minute,
			},
			testFunc: testWindowReset,
		},
		{
			name: "multiple IPs have independent rate limits",
			config: config.RateLimitConfig{
				Enabled:           true,
				RequestsPerMinute: 2,
				WindowDuration:    100 * time.Millisecond,
				CleanupInterval:   1 * time.Minute,
				RecordExpiry:      2 * time.Minute,
			},
			testFunc: testMultipleIPsIndependent,
		},
		{
			name: "concurrent requests from same IP respect rate limit",
			config: config.RateLimitConfig{
				Enabled:           true,
				RequestsPerMinute: 10,
				WindowDuration:    500 * time.Millisecond,
				CleanupInterval:   1 * time.Minute,
				RecordExpiry:      2 * time.Minute,
			},
			testFunc: testConcurrentRequests,
		},
		{
			name: "rate limiting disabled allows unlimited requests",
			config: config.RateLimitConfig{
				Enabled:           false,
				RequestsPerMinute: 5,
				WindowDuration:    100 * time.Millisecond,
				CleanupInterval:   1 * time.Minute,
				RecordExpiry:      2 * time.Minute,
			},
			testFunc: testRateLimitingDisabled,
		},
		{
			name: "health check endpoint bypasses rate limiting",
			config: config.RateLimitConfig{
				Enabled:           true,
				RequestsPerMinute: 2,
				WindowDuration:    100 * time.Millisecond,
				CleanupInterval:   1 * time.Minute,
				RecordExpiry:      2 * time.Minute,
			},
			testFunc: testHealthCheckBypass,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := setupTestServer(tt.config)
			defer server.Close()

			tt.testFunc(t, server)
		})
	}
}

func TestMain(m *testing.M) {
	// Set Gin to test mode for all tests
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}
