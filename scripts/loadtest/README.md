# Load Testing Tool for API Endpoints

This tool performs load testing on API endpoints for both Golang and Java Spring Boot backends.

## Supported Endpoints

1. **GET /healthz** - Health check endpoint (database connectivity test)
2. **POST /auth/register** - User registration endpoint (writes to database)

## Features

- Fixed request count testing (configurable)
- Concurrent worker pool for parallel requests
- Keep-Alive connection pooling for realistic performance testing
- Detailed statistics including percentiles (P50, P90, P95, P99)
- Automatic comparison between Golang and Java backends
- Real-time progress indicator
- Support for both GET and POST requests with JSON body

## Configuration

Edit the constants in `main.go` to customize your test:

```go
const (
    // Endpoint selection: "healthz" or "register"
    TestEndpoint = "register"

    // Request timeout duration
    RequestTimeout = 5 * time.Second

    // Configuration for /healthz endpoint
    HealthzRequests       = 10000
    HealthzWorkers        = 100
    GolangHealthzURL      = "https://waterballsa-backend-golang.zeabur.app/healthz"
    JavaHealthzURL        = "https://waterballsa-backend.zeabur.app/healthz"

    // Configuration for /auth/register endpoint
    RegisterRequests      = 100
    RegisterWorkers       = 10
    RegisterUsername      = "aaaa"
    RegisterPassword      = "aaaaaaaa"
    GolangRegisterURL     = "https://waterballsa-backend-golang.zeabur.app/auth/register"
    JavaRegisterURL       = "https://waterballsa-backend.zeabur.app/auth/register"
)
```

### Switching Between Endpoints

Simply change the `TestEndpoint` constant:
- Set to `"healthz"` for testing the health check endpoint
- Set to `"register"` for testing the user registration endpoint

## Usage

### Build and Run

```bash
# Build the binary
go build -o loadtest scripts/loadtest/*.go

# Run the load test
./loadtest
```

### Run Directly

```bash
go run scripts/loadtest/*.go
```

## Test Output

The tool provides:

1. **Individual Backend Reports**: Detailed statistics for each backend
   - Response time metrics (Min, Max, Mean, Median, P90, P95, P99)
   - Throughput (RPS - Requests Per Second)
   - Success/failure counts and rates
   - Error breakdown by type

2. **Comparison Table**: Side-by-side comparison highlighting the winner for each metric

3. **Overall Summary**: Declares overall winner based on key metrics

## Example Output

```
╔════════════════════════════════════════════════════════════════════╗
║          Load Testing Tool for /healthz Endpoint                  ║
╚════════════════════════════════════════════════════════════════════╝

Test Configuration:
  Total Requests:      10000
  Concurrent Workers:  100
  Request Timeout:     5s

Targets:
  1. Golang Backend:   https://waterballsa-backend-golang.zeabur.app/healthz
  2. Java Backend:     https://waterballsa-backend.zeabur.app/healthz

[Progress indicators and results...]
```

## Metrics Explained

- **Min/Max/Mean**: Basic response time statistics in milliseconds
- **Median (P50)**: 50% of requests completed faster than this time
- **P90**: 90% of requests completed faster than this time (good indicator of typical performance)
- **P95**: 95% of requests completed faster than this time
- **P99**: 99% of requests completed faster than this time (helps identify outliers)
- **RPS**: Requests per second (throughput)
- **Success Rate**: Percentage of successful requests (200-299 status codes)

## Recommendations

- Start with smaller values (e.g., 100 requests, 10 workers) for initial testing
- Gradually increase load to find performance limits
- Run multiple tests to get consistent results
- Consider network latency when interpreting results
