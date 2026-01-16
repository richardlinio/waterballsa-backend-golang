package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/linporu/waterballsa-backend-golang/internal/db"
	"github.com/linporu/waterballsa-backend-golang/internal/model"
)

// ErrUserNotFound is returned when a user is not found in the database
var ErrUserNotFound = errors.New("user not found")

// UserRepository implements user data access operations using sqlc generated queries.
type UserRepository struct {
	queries db.Querier
}

// NewUserRepository creates a new instance of UserRepository
func NewUserRepository(queries db.Querier) *UserRepository {
	return &UserRepository{
		queries: queries,
	}
}

// GetByID retrieves a user by ID
func (r *UserRepository) GetByID(ctx context.Context, id int64) (*model.User, error) {
	row, err := r.queries.GetUserByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	user := &model.User{
		ID:               row.ID,
		Username:         row.Username,
		PasswordHash:     row.PasswordHash,
		Role:             string(row.Role),
		ExperiencePoints: row.ExperiencePoints,
		Level:            row.Level,
		CreatedAt:        row.CreatedAt.Time,
		UpdatedAt:        row.UpdatedAt.Time,
	}

	return user, nil
}

// GetByUsername retrieves a user by username
func (r *UserRepository) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	row, err := r.queries.GetUserByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	// Convert db row to domain model
	user := &model.User{
		ID:               row.ID,
		Username:         row.Username,
		PasswordHash:     row.PasswordHash,
		Role:             string(row.Role),
		ExperiencePoints: row.ExperiencePoints,
		Level:            row.Level,
		CreatedAt:        row.CreatedAt.Time,
		UpdatedAt:        row.UpdatedAt.Time,
	}

	return user, nil
}

// Create creates a new user with the given username and password hash
func (r *UserRepository) Create(ctx context.Context, username, passwordHash string) (int64, error) {
	userID, err := r.queries.CreateUser(ctx, db.CreateUserParams{
		Username:     username,
		PasswordHash: passwordHash,
	})
	if err != nil {
		return 0, err
	}

	return userID, nil
}

// ExistsByUsername checks if a user with the given username exists (excluding soft-deleted users)
func (r *UserRepository) ExistsByUsername(ctx context.Context, username string) (bool, error) {
	return r.queries.ExistsUserByUsername(ctx, username)
}
