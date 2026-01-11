package testutil

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	TestDBName     = "waterballsa_test"
	TestDBUser     = "test_user"
	TestDBPassword = "test_password"
)

// PostgresContainer encapsulates Testcontainer PostgreSQL instance
type PostgresContainer struct {
	Container *postgres.PostgresContainer
	ConnStr   string
	Host      string
	Port      string
}

// NewPostgresContainer creates and starts a PostgreSQL Testcontainer
// with all migrations pre-applied
func NewPostgresContainer(ctx context.Context) (*PostgresContainer, error) {
	// Get absolute path to migrations directory
	migrationPath, err := filepath.Abs("../../migrations")
	if err != nil {
		return nil, fmt.Errorf("failed to get migration path: %w", err)
	}

	// Get migration file names
	migrationFileNames, err := getMigrationFiles(migrationPath)
	if err != nil {
		return nil, err
	}

	// Create temporary directory for processed migration files
	tempDir, err := os.MkdirTemp("", "migrations-*")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer func() {
		if err := os.RemoveAll(tempDir); err != nil {
			fmt.Printf("Warning: failed to remove temp directory %s: %v\n", tempDir, err)
		}
	}()

	// Process migration files
	processedFiles, err := processMigrationFiles(migrationPath, tempDir, migrationFileNames)
	if err != nil {
		return nil, err
	}

	// Start PostgreSQL container with initialization scripts
	postgresContainer, err := postgres.Run(ctx,
		"postgres:16-alpine",
		postgres.WithDatabase(TestDBName),
		postgres.WithUsername(TestDBUser),
		postgres.WithPassword(TestDBPassword),
		postgres.WithInitScripts(processedFiles...),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(60*time.Second)),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to start postgres container: %w", err)
	}

	// Get connection string
	connStr, err := postgresContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		return nil, fmt.Errorf("failed to get connection string: %w", err)
	}

	// Get host and port
	host, err := postgresContainer.Host(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get host: %w", err)
	}

	port, err := postgresContainer.MappedPort(ctx, "5432")
	if err != nil {
		return nil, fmt.Errorf("failed to get port: %w", err)
	}

	return &PostgresContainer{
		Container: postgresContainer,
		ConnStr:   connStr,
		Host:      host,
		Port:      port.Port(),
	}, nil
}

// extractUpMigration extracts only the "Up" portion from a Goose migration file
// Goose migration files contain both "-- +goose Up" and "-- +goose Down" sections
// We only want to execute the Up section for initialization
func extractUpMigration(content string) string {
	lines := strings.Split(content, "\n")
	var upLines []string
	inUpSection := false

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		// Start capturing when we see "-- +goose Up"
		if strings.HasPrefix(trimmed, "-- +goose Up") {
			inUpSection = true
			continue
		}

		// Stop capturing when we see "-- +goose Down"
		if strings.HasPrefix(trimmed, "-- +goose Down") {
			break
		}

		// Collect lines in the Up section
		if inUpSection {
			upLines = append(upLines, line)
		}
	}

	return strings.Join(upLines, "\n")
}

// getMigrationFiles reads and sorts migration file names from the migrations directory
func getMigrationFiles(migrationPath string) ([]string, error) {
	entries, err := os.ReadDir(migrationPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read migrations directory: %w", err)
	}

	var migrationFileNames []string
	for _, entry := range entries {
		if !entry.IsDir() && filepath.Ext(entry.Name()) == ".sql" {
			migrationFileNames = append(migrationFileNames, entry.Name())
		}
	}

	sort.Strings(migrationFileNames)

	if len(migrationFileNames) == 0 {
		return nil, fmt.Errorf("no migration files found in %s", migrationPath)
	}

	return migrationFileNames, nil
}

// processMigrationFiles processes migration files and writes them to a temporary directory
func processMigrationFiles(migrationPath, tempDir string, fileNames []string) ([]string, error) {
	var processedFiles []string
	for _, fileName := range fileNames {
		content, err := os.ReadFile(filepath.Join(migrationPath, fileName)) // #nosec G304 -- reading migration files from project directory in test environment
		if err != nil {
			return nil, fmt.Errorf("failed to read migration file %s: %w", fileName, err)
		}

		upContent := extractUpMigration(string(content))

		tempFile := filepath.Join(tempDir, fileName)
		if err := os.WriteFile(tempFile, []byte(upContent), 0o600); err != nil {
			return nil, fmt.Errorf("failed to write temp migration file %s: %w", fileName, err)
		}

		processedFiles = append(processedFiles, tempFile)
	}
	return processedFiles, nil
}

// Terminate terminates the container
func (pc *PostgresContainer) Terminate(ctx context.Context) error {
	if pc.Container != nil {
		return pc.Container.Terminate(ctx)
	}
	return nil
}
