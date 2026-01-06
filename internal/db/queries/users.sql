-- name: CreateUser :one
INSERT INTO users (username, password_hash, role, experience_points, level)
VALUES ($1, $2, 'STUDENT', 0, 1)
RETURNING id;

-- name: GetUserByUsername :one
SELECT id, username, password_hash, role, experience_points, level, created_at, updated_at
FROM users
WHERE username = $1 AND deleted_at IS NULL;
