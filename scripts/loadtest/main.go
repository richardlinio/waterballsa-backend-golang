package main

import (
	"fmt"
	"time"
)

// Test configuration constants - adjust these to change test parameters
const (
	// Endpoint type constants
	endpointHealthz  = "healthz"
	endpointRegister = "register"

	// Endpoint selection: "healthz" or "register"
	TestEndpoint = endpointRegister

	// Request timeout duration
	RequestTimeout = 5 * time.Second

	// Configuration for /healthz endpoint
	HealthzRequests  = 10000
	HealthzWorkers   = 100
	GolangHealthzURL = "https://waterballsa-backend-golang.zeabur.app/healthz"
	JavaHealthzURL   = "https://waterballsa-backend.zeabur.app/healthz"

	// Configuration for /auth/register endpoint
	RegisterRequests  = 100000
	RegisterWorkers   = 1000
	RegisterUsername  = "aaaa"
	RegisterPassword  = "aaaaaaaa"
	GolangRegisterURL = "https://waterballsa-backend-golang.zeabur.app/auth/register"
	JavaRegisterURL   = "https://waterballsa-backend.zeabur.app/auth/register"
)

func main() {
	// Determine configuration based on selected endpoint
	var (
		endpointName   string
		totalRequests  int
		workers        int
		golangURL      string
		javaURL        string
		additionalInfo string
	)

	switch TestEndpoint {
	case endpointHealthz:
		endpointName = "/healthz Endpoint (GET)"
		totalRequests = HealthzRequests
		workers = HealthzWorkers
		golangURL = GolangHealthzURL
		javaURL = JavaHealthzURL
	case endpointRegister:
		endpointName = "/auth/register Endpoint (POST)"
		totalRequests = RegisterRequests
		workers = RegisterWorkers
		golangURL = GolangRegisterURL
		javaURL = JavaRegisterURL
		additionalInfo = fmt.Sprintf("  Test Credentials:    username=%s, password=%s\n", RegisterUsername, RegisterPassword)
	default:
		fmt.Printf("Error: Unknown endpoint '%s'. Use 'healthz' or 'register'\n", TestEndpoint)
		return
	}

	// Print header
	fmt.Println("╔════════════════════════════════════════════════════════════════════╗")
	fmt.Printf("║          Load Testing Tool - %-36s ║\n", endpointName)
	fmt.Println("╚════════════════════════════════════════════════════════════════════╝")
	fmt.Println()
	fmt.Println("Test Configuration:")
	fmt.Printf("  Endpoint:            %s\n", endpointName)
	fmt.Printf("  Total Requests:      %d\n", totalRequests)
	fmt.Printf("  Concurrent Workers:  %d\n", workers)
	fmt.Printf("  Request Timeout:     %s\n", RequestTimeout)
	if additionalInfo != "" {
		fmt.Print(additionalInfo)
	}
	fmt.Println()
	fmt.Println("Targets:")
	fmt.Printf("  1. Golang Backend:   %s\n", golangURL)
	fmt.Printf("  2. Java Backend:     %s\n", javaURL)
	fmt.Println()
	fmt.Println("Starting tests in 2 seconds...")
	time.Sleep(2 * time.Second)

	// Create load tester
	tester := NewLoadTester(RequestTimeout, workers)

	var golangResult, javaResult *TestResult

	// Test Golang backend
	fmt.Println()
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("Testing Golang Backend...")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	if TestEndpoint == endpointRegister {
		golangResult = tester.RunRegisterTest(golangURL, RegisterUsername, RegisterPassword, totalRequests, workers)
	} else {
		golangResult = tester.RunTest(golangURL, totalRequests, workers)
	}
	PrintTestResult(golangResult, "Golang Backend", endpointName)

	// Wait a bit before next test
	fmt.Println("Waiting 3 seconds before next test...")
	time.Sleep(3 * time.Second)

	// Test Java backend
	fmt.Println()
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("Testing Java Backend...")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	if TestEndpoint == endpointRegister {
		javaResult = tester.RunRegisterTest(javaURL, RegisterUsername, RegisterPassword, totalRequests, workers)
	} else {
		javaResult = tester.RunTest(javaURL, totalRequests, workers)
	}
	PrintTestResult(javaResult, "Java Backend", endpointName)

	// Print comparison
	PrintComparison(golangResult, javaResult)

	fmt.Println("✓ Load testing completed!")
}
