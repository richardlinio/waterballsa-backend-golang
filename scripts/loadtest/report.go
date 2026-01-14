package main

import (
	"fmt"
	"strings"
)

// PrintTestResult prints the detailed test result for a single backend
func PrintTestResult(result *TestResult, backendName, endpointName string) {
	stats := result.Calculate()

	fmt.Println()
	fmt.Println(strings.Repeat("=", 70))
	fmt.Printf("Load Test Results - %s\n", backendName)
	fmt.Println(strings.Repeat("=", 70))
	fmt.Printf("Endpoint: %s\n", endpointName)
	fmt.Printf("Target: %s\n", result.TargetURL)
	fmt.Println()

	fmt.Println("Response Time (ms):")
	fmt.Printf("  Min:       %8.2f\n", stats.Min)
	fmt.Printf("  Max:       %8.2f\n", stats.Max)
	fmt.Printf("  Mean:      %8.2f\n", stats.Mean)
	fmt.Printf("  Median:    %8.2f\n", stats.Median)
	fmt.Printf("  P90:       %8.2f\n", stats.P90)
	fmt.Printf("  P95:       %8.2f\n", stats.P95)
	fmt.Printf("  P99:       %8.2f\n", stats.P99)
	fmt.Println()

	fmt.Println("Throughput:")
	fmt.Printf("  RPS:           %10.2f\n", stats.RPS)
	fmt.Printf("  Success:       %10d (%.2f%%)\n", result.SuccessfulRequests, stats.SuccessRate)
	fmt.Printf("  Failed:        %10d (%.2f%%)\n", result.FailedRequests, 100-stats.SuccessRate)
	fmt.Printf("  Duration:      %10s\n", result.Duration.Round(100*1000).String()) // Round to 0.1s
	fmt.Println()

	if len(result.Errors) > 0 {
		fmt.Println("Errors by Type:")
		for errType, count := range result.Errors {
			fmt.Printf("  %-15s: %d\n", errType, count)
		}
		fmt.Println()
	}
}

// PrintComparison prints a side-by-side comparison of two test results
func PrintComparison(golangResult, javaResult *TestResult) {
	golangStats := golangResult.Calculate()
	javaStats := javaResult.Calculate()

	fmt.Println()
	fmt.Println(strings.Repeat("=", 70))
	fmt.Println("Performance Comparison")
	fmt.Println(strings.Repeat("=", 70))
	fmt.Printf("%-25s %12s %12s %12s\n", "Metric", "Golang", "Java", "Winner")
	fmt.Println(strings.Repeat("-", 70))

	// Helper function to determine winner (lower is better for response times)
	printComparisonRow := func(metric string, golangVal, javaVal float64, lowerIsBetter bool) {
		var winner string
		if lowerIsBetter {
			switch {
			case golangVal < javaVal:
				winner = "Golang ⭐"
			case javaVal < golangVal:
				winner = "Java ⭐"
			default:
				winner = "Tie"
			}
		} else {
			switch {
			case golangVal > javaVal:
				winner = "Golang ⭐"
			case javaVal > golangVal:
				winner = "Java ⭐"
			default:
				winner = "Tie"
			}
		}
		fmt.Printf("%-25s %12.2f %12.2f %12s\n", metric, golangVal, javaVal, winner)
	}

	// Response time comparisons (lower is better)
	printComparisonRow("Min Response (ms)", golangStats.Min, javaStats.Min, true)
	printComparisonRow("Max Response (ms)", golangStats.Max, javaStats.Max, true)
	printComparisonRow("Mean Response (ms)", golangStats.Mean, javaStats.Mean, true)
	printComparisonRow("Median Response (ms)", golangStats.Median, javaStats.Median, true)
	printComparisonRow("P90 Response (ms)", golangStats.P90, javaStats.P90, true)
	printComparisonRow("P95 Response (ms)", golangStats.P95, javaStats.P95, true)
	printComparisonRow("P99 Response (ms)", golangStats.P99, javaStats.P99, true)

	fmt.Println()

	// Throughput comparisons (higher is better)
	printComparisonRow("RPS", golangStats.RPS, javaStats.RPS, false)
	printComparisonRow("Success Rate (%)", golangStats.SuccessRate, javaStats.SuccessRate, false)

	fmt.Println(strings.Repeat("=", 70))
	fmt.Println()

	// Summary
	golangWins := 0
	javaWins := 0

	metrics := []struct {
		name          string
		golangVal     float64
		javaVal       float64
		lowerIsBetter bool
	}{
		{"Mean Response", golangStats.Mean, javaStats.Mean, true},
		{"P99 Response", golangStats.P99, javaStats.P99, true},
		{"RPS", golangStats.RPS, javaStats.RPS, false},
		{"Success Rate", golangStats.SuccessRate, javaStats.SuccessRate, false},
	}

	for _, m := range metrics {
		if m.lowerIsBetter {
			if m.golangVal < m.javaVal {
				golangWins++
			} else if m.javaVal < m.golangVal {
				javaWins++
			}
		} else {
			if m.golangVal > m.javaVal {
				golangWins++
			} else if m.javaVal > m.golangVal {
				javaWins++
			}
		}
	}

	fmt.Printf("Summary: Golang won %d key metrics, Java won %d key metrics\n", golangWins, javaWins)
	switch {
	case golangWins > javaWins:
		fmt.Println("Overall Winner: Golang ⭐")
	case javaWins > golangWins:
		fmt.Println("Overall Winner: Java ⭐")
	default:
		fmt.Println("Overall Result: Tie")
	}
	fmt.Println()
}
