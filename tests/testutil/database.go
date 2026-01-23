package testutil

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

// CleanDatabase truncates all tables and resets sequences for a clean test state
// This should be called before each test scenario to ensure isolation
func CleanDatabase(ctx context.Context, pool *pgxpool.Pool) error {
	// Truncate all tables and reset identity sequences
	// CASCADE ensures that dependent records in other tables are also deleted
	query := `
		TRUNCATE TABLE users, journeys RESTART IDENTITY CASCADE;
	`

	_, err := pool.Exec(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to clean database: %w", err)
	}

	return nil
}

// CreateTestUser creates a user in the database with the given username and password
// The password is hashed with bcrypt before storing
// Returns the created user ID
func CreateTestUser(ctx context.Context, pool *pgxpool.Pool, username, password string) (int64, error) {
	// Hash the password using bcrypt with default cost
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return 0, fmt.Errorf("failed to hash password: %w", err)
	}

	// Insert user into database
	// Using default role 'STUDENT' and initial values for experience_points (0) and level (1)
	query := `
		INSERT INTO users (username, password_hash, role, experience_points, level)
		VALUES ($1, $2, 'STUDENT', 0, 1)
		RETURNING id
	`

	var userID int64
	err = pool.QueryRow(ctx, query, username, string(passwordHash)).Scan(&userID)
	if err != nil {
		return 0, fmt.Errorf("failed to create test user: %w", err)
	}

	return userID, nil
}

// CreateTestJourney creates a journey in the database with the given parameters
// Returns the created journey ID
func CreateTestJourney(
	ctx context.Context,
	pool *pgxpool.Pool,
	title string,
	slug string,
	description string,
	teacherName string,
	price float64,
	coverImageURL string,
) (int64, error) {
	query := `
		INSERT INTO journeys (title, slug, description, teacher_name, price, cover_image_url)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`

	var journeyID int64
	err := pool.QueryRow(ctx, query, title, slug, description, teacherName, price, coverImageURL).Scan(&journeyID)
	if err != nil {
		return 0, fmt.Errorf("failed to create test journey: %w", err)
	}

	return journeyID, nil
}
