package testutil

import (
	"context"
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

// CleanDatabase truncates all tables and resets sequences for a clean test state
// This should be called before each test scenario to ensure isolation
func CleanDatabase(ctx context.Context, pool *pgxpool.Pool) error {
	// Truncate all tables and reset identity sequences
	// CASCADE ensures that dependent records in other tables are also deleted
	query := `
		TRUNCATE TABLE users, journeys, chapters, missions, rewards, mission_resources,
		             user_mission_progress, orders, order_items, user_journeys
		RESTART IDENTITY CASCADE;
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

// CreateTestOrder creates an order in the database with the given parameters
// Returns the created order ID
// status should be one of: UNPAID, PAID, EXPIRED
func CreateTestOrder(
	ctx context.Context,
	pool *pgxpool.Pool,
	userID int64,
	journeyID int64,
	status string,
) (int64, error) {
	// Generate order number using pattern: {timestamp}{userId}{randomCode}
	orderNumber := generateOrderNumber(userID)

	// Get journey price for calculating order amounts
	var price float64
	priceQuery := `SELECT price FROM journeys WHERE id = $1`
	err := pool.QueryRow(ctx, priceQuery, journeyID).Scan(&price)
	if err != nil {
		return 0, fmt.Errorf("failed to get journey price: %w", err)
	}

	// Insert order
	// For PAID orders, set paid_at to current timestamp
	var query string
	var orderID int64
	if status == "PAID" {
		query = `
			INSERT INTO orders (order_number, user_id, status, original_price, discount, price, paid_at)
			VALUES ($1, $2, $3, $4, 0, $4, NOW())
			RETURNING id
		`
		err = pool.QueryRow(ctx, query, orderNumber, userID, status, price).Scan(&orderID)
	} else {
		query = `
			INSERT INTO orders (order_number, user_id, status, original_price, discount, price)
			VALUES ($1, $2, $3, $4, 0, $4)
			RETURNING id
		`
		err = pool.QueryRow(ctx, query, orderNumber, userID, status, price).Scan(&orderID)
	}
	if err != nil {
		return 0, fmt.Errorf("failed to create test order: %w", err)
	}

	// Insert order item
	itemQuery := `
		INSERT INTO order_items (order_id, journey_id, quantity, original_price, discount, price)
		VALUES ($1, $2, 1, $3, 0, $3)
	`
	_, err = pool.Exec(ctx, itemQuery, orderID, journeyID, price)
	if err != nil {
		return 0, fmt.Errorf("failed to create order item: %w", err)
	}

	// Insert user_journeys if status is PAID
	if status == "PAID" {
		journeyQuery := `
			INSERT INTO user_journeys (user_id, journey_id, order_id)
			VALUES ($1, $2, $3)
		`
		_, err = pool.Exec(ctx, journeyQuery, userID, journeyID, orderID)
		if err != nil {
			return 0, fmt.Errorf("failed to create user journey: %w", err)
		}
	}

	return orderID, nil
}

// generateOrderNumber generates an order number using pattern: {timestamp}{userId}{randomCode}
func generateOrderNumber(userID int64) string {
	timestamp := time.Now().Unix()

	// Generate cryptographically secure random number between 0-9999
	var randomBytes [2]byte
	_, err := rand.Read(randomBytes[:])
	if err != nil {
		// Fallback to timestamp-based value if crypto/rand fails
		randomCode := timestamp % 10000
		return fmt.Sprintf("%d%d%04d", timestamp, userID, randomCode)
	}
	randomCode := binary.BigEndian.Uint16(randomBytes[:]) % 10000

	return fmt.Sprintf("%d%d%04d", timestamp, userID, randomCode)
}
