package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

// LoadTester manages load testing configuration and execution
type LoadTester struct {
	client         *http.Client
	requestTimeout time.Duration
}

// requestResult holds the result of a single request
type requestResult struct {
	responseTime float64
	err          error
}

// NewLoadTester creates a new LoadTester with Keep-Alive connection pooling
func NewLoadTester(requestTimeout time.Duration, maxWorkers int) *LoadTester {
	// Configure HTTP transport with connection pooling
	transport := &http.Transport{
		MaxIdleConnsPerHost: maxWorkers,
		IdleConnTimeout:     90 * time.Second,
		DisableKeepAlives:   false,
	}

	client := &http.Client{
		Transport: transport,
		Timeout:   requestTimeout,
	}

	return &LoadTester{
		client:         client,
		requestTimeout: requestTimeout,
	}
}

// RunTest executes a load test against the target URL
func (lt *LoadTester) RunTest(targetURL string, totalRequests, workers int) *TestResult {
	result := &TestResult{
		TargetURL:     targetURL,
		TotalRequests: totalRequests,
		ResponseTimes: make([]float64, 0, totalRequests),
		Errors:        make(map[string]int),
	}

	// Channel for distributing work
	jobs := make(chan int, totalRequests)

	// Channel for collecting results
	results := make(chan requestResult, totalRequests)

	// WaitGroup for worker synchronization
	var wg sync.WaitGroup

	// Start timer
	startTime := time.Now()

	// Start workers
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			lt.worker(targetURL, jobs, results)
		}()
	}

	// Send jobs to workers
	go func() {
		for i := 0; i < totalRequests; i++ {
			jobs <- i
		}
		close(jobs)
	}()

	// Progress tracking in separate goroutine
	done := make(chan bool)
	go func() {
		wg.Wait()
		close(results)
		done <- true
	}()

	// Collect results and show progress
	completed := 0
	for res := range results {
		completed++

		// Show progress every 10% or every 100 requests (whichever is smaller)
		progressInterval := totalRequests / 10
		if progressInterval > 100 {
			progressInterval = 100
		}
		if progressInterval < 1 {
			progressInterval = 1
		}

		if completed%progressInterval == 0 || completed == totalRequests {
			fmt.Printf("\rProgress: %d/%d (%.1f%%)", completed, totalRequests, float64(completed)/float64(totalRequests)*100)
		}

		if res.err != nil {
			result.FailedRequests++
			errType := categorizeError(res.err)
			result.Errors[errType]++
		} else {
			result.SuccessfulRequests++
			result.ResponseTimes = append(result.ResponseTimes, res.responseTime)
		}
	}

	<-done
	fmt.Println() // New line after progress

	// Record duration
	result.Duration = time.Since(startTime)

	return result
}

// worker processes requests from the jobs channel
func (lt *LoadTester) worker(targetURL string, jobs <-chan int, results chan requestResult) {
	for range jobs {
		start := time.Now()
		err := lt.sendRequest(targetURL)
		elapsed := time.Since(start)

		results <- requestResult{
			responseTime: float64(elapsed.Milliseconds()),
			err:          err,
		}
	}
}

// sendRequest sends a single HTTP GET request to the target URL
func (lt *LoadTester) sendRequest(targetURL string) error {
	ctx, cancel := context.WithTimeout(context.Background(), lt.requestTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, targetURL, nil)
	if err != nil {
		return fmt.Errorf("request creation failed: %w", err)
	}

	resp, err := lt.client.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer func() { _ = resp.Body.Close() }() //nolint:errcheck // Close error in deferred call is not critical for load testing

	// Read and discard response body to properly reuse connections
	_, err = io.Copy(io.Discard, resp.Body)
	if err != nil {
		return fmt.Errorf("response read failed: %w", err)
	}

	// Any HTTP response (including 4xx, 5xx) is considered a successful response
	// Only network-level errors (timeout, connection refused) are considered failures
	return nil
}

// categorizeError categorizes errors into types for reporting
func categorizeError(err error) string {
	errMsg := err.Error()

	if contains(errMsg, "timeout") || contains(errMsg, "deadline exceeded") {
		return "Timeout"
	}
	if contains(errMsg, "connection refused") || contains(errMsg, "no such host") {
		return "Connection"
	}
	if contains(errMsg, "status code") {
		return "HTTP Status"
	}
	return "Other"
}

// contains checks if a string contains a substring (case-insensitive helper)
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && findSubstring(s, substr))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// RunRegisterTest executes a load test for the /auth/register endpoint
func (lt *LoadTester) RunRegisterTest(targetURL, username, password string, totalRequests, workers int) *TestResult {
	result := &TestResult{
		TargetURL:     targetURL,
		TotalRequests: totalRequests,
		ResponseTimes: make([]float64, 0, totalRequests),
		Errors:        make(map[string]int),
	}

	// Channel for distributing work
	jobs := make(chan int, totalRequests)

	// Channel for collecting results
	results := make(chan requestResult, totalRequests)

	// WaitGroup for worker synchronization
	var wg sync.WaitGroup

	// Start timer
	startTime := time.Now()

	// Start workers
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			lt.registerWorker(targetURL, username, password, jobs, results)
		}()
	}

	// Send jobs to workers
	go func() {
		for i := 0; i < totalRequests; i++ {
			jobs <- i
		}
		close(jobs)
	}()

	// Progress tracking in separate goroutine
	done := make(chan bool)
	go func() {
		wg.Wait()
		close(results)
		done <- true
	}()

	// Collect results and show progress
	completed := 0
	for res := range results {
		completed++

		// Show progress every 10% or every 10 requests (whichever is smaller for register)
		progressInterval := totalRequests / 10
		if progressInterval > 10 {
			progressInterval = 10
		}
		if progressInterval < 1 {
			progressInterval = 1
		}

		if completed%progressInterval == 0 || completed == totalRequests {
			fmt.Printf("\rProgress: %d/%d (%.1f%%)", completed, totalRequests, float64(completed)/float64(totalRequests)*100)
		}

		if res.err != nil {
			result.FailedRequests++
			errType := categorizeError(res.err)
			result.Errors[errType]++
		} else {
			result.SuccessfulRequests++
			result.ResponseTimes = append(result.ResponseTimes, res.responseTime)
		}
	}

	<-done
	fmt.Println() // New line after progress

	// Record duration
	result.Duration = time.Since(startTime)

	return result
}

// registerWorker processes register requests from the jobs channel
func (lt *LoadTester) registerWorker(targetURL, username, password string, jobs <-chan int, results chan requestResult) {
	for range jobs {
		start := time.Now()
		err := lt.sendRegisterRequest(targetURL, username, password)
		elapsed := time.Since(start)

		results <- requestResult{
			responseTime: float64(elapsed.Milliseconds()),
			err:          err,
		}
	}
}

// sendRegisterRequest sends a POST request to register endpoint with JSON body
func (lt *LoadTester) sendRegisterRequest(targetURL, username, password string) error {
	ctx, cancel := context.WithTimeout(context.Background(), lt.requestTimeout)
	defer cancel()

	// Create request body
	requestBody := map[string]string{
		"username": username,
		"password": password,
	}
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return fmt.Errorf("json marshal failed: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, targetURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		return fmt.Errorf("request creation failed: %w", err)
	}

	// Set Content-Type header
	req.Header.Set("Content-Type", "application/json")

	resp, err := lt.client.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer func() { _ = resp.Body.Close() }() //nolint:errcheck // Close error in deferred call is not critical for load testing

	// Read and discard response body to properly reuse connections
	_, err = io.Copy(io.Discard, resp.Body)
	if err != nil {
		return fmt.Errorf("response read failed: %w", err)
	}

	// Any HTTP response (including 4xx, 5xx) is considered a successful response
	// Only network-level errors (timeout, connection refused) are considered failures
	return nil
}
