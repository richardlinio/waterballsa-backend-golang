package bdd

import (
	"context"
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/cucumber/godog"
	"github.com/linporu/waterballsa-backend-golang/tests/bdd/steps/common"
	"github.com/linporu/waterballsa-backend-golang/tests/bdd/testcontext"
	"github.com/linporu/waterballsa-backend-golang/tests/testutil"
)

var (
	// Global test server instance (initialized once per test suite)
	testServerInstance *testcontext.TestServerWrapper
	// Global Testcontainer PostgreSQL instance
	postgresContainer *testutil.PostgresContainer
	testServerOnce    sync.Once
)

// TestFeatures is the entry point for running godog BDD tests
func TestFeatures(t *testing.T) {
	suite := godog.TestSuite{
		Name:                "ISA Features",
		ScenarioInitializer: InitializeScenario,
		Options: &godog.Options{
			Format: "pretty",
			Paths:  []string{"features/isa"},
			// Uncomment below to run tests concurrently
			// Concurrency: 4,
			TestingT: t,
		},
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}
}

// InitializeScenario registers all step definitions and hooks for each scenario
func InitializeScenario(sc *godog.ScenarioContext) {
	// Initialize test server once for the entire suite
	testServerOnce.Do(func() {
		server, err := testutil.NewTestServer(
			context.Background(),
			postgresContainer.Host,
			postgresContainer.Port,
		)
		if err != nil {
			panic(fmt.Sprintf("Failed to create test server: %v", err))
		}

		// Start the HTTP server
		if err := server.Start(); err != nil {
			panic(fmt.Sprintf("Failed to start test server: %v", err))
		}

		testServerInstance = &testcontext.TestServerWrapper{
			Server: server,
		}
	})

	// Before each scenario: clean database and add server to context
	sc.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		// Clean database to ensure isolation between scenarios
		cleanCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		testServerInstance.Mu.Lock()
		defer testServerInstance.Mu.Unlock()

		if err := testutil.CleanDatabase(cleanCtx, testServerInstance.Server.Pool); err != nil {
			return ctx, fmt.Errorf("failed to clean database before scenario '%s': %w", sc.Name, err)
		}

		// Add test server to scenario context
		ctx = context.WithValue(ctx, testcontext.ContextKeyTestServer, testServerInstance)

		return ctx, nil
	})

	// Register all step definitions
	common.RegisterHTTPSteps(sc)
	common.RegisterDatabaseSteps(sc)
}

// TestMain handles suite-level setup and teardown
func TestMain(m *testing.M) {
	ctx := context.Background()

	// 1. Start PostgreSQL Testcontainer
	var err error
	postgresContainer, err = testutil.NewPostgresContainer(ctx)
	if err != nil {
		panic(fmt.Sprintf("Failed to start postgres container: %v", err))
	}

	fmt.Fprintf(os.Stderr, "[BDD Setup] PostgreSQL Testcontainer started at %s:%s\n",
		postgresContainer.Host, postgresContainer.Port)

	// 2. Run tests
	status := m.Run()

	// 3. Cleanup: shutdown test server if it was initialized
	if testServerInstance != nil && testServerInstance.Server != nil {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

		if err := testServerInstance.Server.Shutdown(shutdownCtx); err != nil {
			fmt.Fprintf(os.Stderr, "[BDD Cleanup] Error shutting down test server: %v\n", err)
		}
		cancel()
	}

	// 4. Cleanup: terminate PostgreSQL Testcontainer
	if postgresContainer != nil {
		terminateCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

		if err := postgresContainer.Terminate(terminateCtx); err != nil {
			fmt.Fprintf(os.Stderr, "[BDD Cleanup] Error terminating postgres container: %v\n", err)
		}
		cancel()
	}

	// 5. Exit with test status
	os.Exit(status)
}
