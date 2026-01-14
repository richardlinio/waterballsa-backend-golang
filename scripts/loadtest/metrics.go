package main

import (
	"sort"
	"time"
)

// TestResult holds the results of a load test
type TestResult struct {
	TargetURL          string
	TotalRequests      int
	SuccessfulRequests int
	FailedRequests     int
	Duration           time.Duration
	ResponseTimes      []float64 // in milliseconds
	Errors             map[string]int
}

// Statistics holds calculated statistical metrics
type Statistics struct {
	Min         float64
	Max         float64
	Mean        float64
	Median      float64
	P90         float64
	P95         float64
	P99         float64
	RPS         float64
	SuccessRate float64
}

// Calculate computes all statistics from the test result
func (r *TestResult) Calculate() Statistics {
	if len(r.ResponseTimes) == 0 {
		return Statistics{}
	}

	// Sort response times for percentile calculation
	sorted := make([]float64, len(r.ResponseTimes))
	copy(sorted, r.ResponseTimes)
	sort.Float64s(sorted)

	stats := Statistics{
		Min:         sorted[0],
		Max:         sorted[len(sorted)-1],
		Mean:        calculateMean(sorted),
		Median:      calculatePercentile(sorted, 50),
		P90:         calculatePercentile(sorted, 90),
		P95:         calculatePercentile(sorted, 95),
		P99:         calculatePercentile(sorted, 99),
		RPS:         calculateRPS(r.TotalRequests, r.Duration),
		SuccessRate: calculateSuccessRate(r.SuccessfulRequests, r.TotalRequests),
	}

	return stats
}

// calculateMean calculates the arithmetic mean of values
func calculateMean(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}

	var sum float64
	for _, v := range values {
		sum += v
	}
	return sum / float64(len(values))
}

// calculatePercentile calculates the nth percentile using linear interpolation
func calculatePercentile(sorted []float64, percentile float64) float64 {
	if len(sorted) == 0 {
		return 0
	}
	if len(sorted) == 1 {
		return sorted[0]
	}

	// Calculate the index
	index := (percentile / 100.0) * float64(len(sorted)-1)
	lower := int(index)
	upper := lower + 1

	if upper >= len(sorted) {
		return sorted[len(sorted)-1]
	}

	// Linear interpolation
	weight := index - float64(lower)
	return sorted[lower]*(1-weight) + sorted[upper]*weight
}

// calculateRPS calculates requests per second
func calculateRPS(totalRequests int, duration time.Duration) float64 {
	if duration == 0 {
		return 0
	}
	return float64(totalRequests) / duration.Seconds()
}

// calculateSuccessRate calculates the success rate as a percentage
func calculateSuccessRate(successful, total int) float64 {
	if total == 0 {
		return 0
	}
	return (float64(successful) / float64(total)) * 100
}
