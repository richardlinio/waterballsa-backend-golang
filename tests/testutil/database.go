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
		TRUNCATE TABLE users, journeys, chapters, missions, rewards, mission_resources, user_mission_progress RESTART IDENTITY CASCADE;
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

// CreateTestChapter creates a chapter in the database with the given parameters
// Returns the created chapter ID
func CreateTestChapter(
	ctx context.Context,
	pool *pgxpool.Pool,
	journeyID int64,
	title string,
	orderIndex int,
) (int64, error) {
	query := `
		INSERT INTO chapters (journey_id, title, order_index)
		VALUES ($1, $2, $3)
		RETURNING id
	`

	var chapterID int64
	err := pool.QueryRow(ctx, query, journeyID, title, orderIndex).Scan(&chapterID)
	if err != nil {
		return 0, fmt.Errorf("failed to create test chapter: %w", err)
	}

	return chapterID, nil
}

// CreateTestMission creates a mission in the database with the given parameters
// Returns the created mission ID
// missionType should be one of: VIDEO, ARTICLE, QUESTIONNAIRE
// accessLevel should be one of: PUBLIC, AUTHENTICATED, PURCHASED
func CreateTestMission(
	ctx context.Context,
	pool *pgxpool.Pool,
	chapterID int64,
	title string,
	description string,
	missionType string,
	accessLevel string,
	orderIndex int,
) (int64, error) {
	query := `
		INSERT INTO missions (chapter_id, title, description, type, access_level, order_index)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`

	var missionID int64
	err := pool.QueryRow(ctx, query, chapterID, title, description, missionType, accessLevel, orderIndex).Scan(&missionID)
	if err != nil {
		return 0, fmt.Errorf("failed to create test mission: %w", err)
	}

	return missionID, nil
}

// CreateTestReward creates a reward in the database with the given parameters
// Returns the created reward ID
// rewardType should be one of: EXPERIENCE (currently only supported enum value)
func CreateTestReward(
	ctx context.Context,
	pool *pgxpool.Pool,
	missionID int64,
	rewardType string,
	rewardValue int,
) (int64, error) {
	query := `
		INSERT INTO rewards (mission_id, reward_type, reward_value)
		VALUES ($1, $2, $3)
		RETURNING id
	`

	var rewardID int64
	err := pool.QueryRow(ctx, query, missionID, rewardType, rewardValue).Scan(&rewardID)
	if err != nil {
		return 0, fmt.Errorf("failed to create test reward: %w", err)
	}

	return rewardID, nil
}

// CreateTestMissionResource creates a mission resource in the database with the given parameters
// Returns the created resource ID
// resourceType should be one of: VIDEO, ARTICLE, FORM
// durationSeconds is nullable (use nil for non-video resources)
func CreateTestMissionResource(
	ctx context.Context,
	pool *pgxpool.Pool,
	missionID int64,
	resourceType string,
	resourceURL string,
	contentOrder int,
	durationSeconds *int,
) (int64, error) {
	query := `
		INSERT INTO mission_resources (mission_id, resource_type, resource_url, content_order, duration_seconds)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`

	var resourceID int64
	err := pool.QueryRow(ctx, query, missionID, resourceType, resourceURL, contentOrder, durationSeconds).Scan(&resourceID)
	if err != nil {
		return 0, fmt.Errorf("failed to create test mission resource: %w", err)
	}

	return resourceID, nil
}

// CreateTestUserMissionProgress creates a user mission progress record in the database
// Returns the created progress ID
// status should be one of: UNCOMPLETED, COMPLETED, DELIVERED
func CreateTestUserMissionProgress(
	ctx context.Context,
	pool *pgxpool.Pool,
	userID int64,
	missionID int64,
	status string,
	watchPositionSeconds int,
) (int64, error) {
	query := `
		INSERT INTO user_mission_progress (user_id, mission_id, status, watch_position_seconds)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`

	var progressID int64
	err := pool.QueryRow(ctx, query, userID, missionID, status, watchPositionSeconds).Scan(&progressID)
	if err != nil {
		return 0, fmt.Errorf("failed to create test user mission progress: %w", err)
	}

	return progressID, nil
}
