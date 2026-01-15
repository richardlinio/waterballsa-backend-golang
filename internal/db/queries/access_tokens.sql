-- name: InvalidateToken :exec
INSERT INTO access_tokens (token_jti, user_id, expires_at)
VALUES ($1, $2, $3);

-- name: IsTokenInvalidated :one
SELECT EXISTS(SELECT 1 FROM access_tokens WHERE token_jti = $1);

-- name: DeleteExpiredTokens :exec
DELETE FROM access_tokens WHERE expires_at < NOW();
