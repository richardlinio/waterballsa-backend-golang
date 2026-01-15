-- +goose Up
-- Create refresh_tokens table for storing refresh tokens (dual token pattern)
CREATE TABLE refresh_tokens (
    id BIGSERIAL PRIMARY KEY,
    token_jti VARCHAR(64) NOT NULL UNIQUE,
    user_id BIGINT NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    revoked_at TIMESTAMP,  -- NULL = valid, non-NULL = revoked
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_refresh_tokens_user FOREIGN KEY (user_id)
        REFERENCES users (id) ON DELETE CASCADE
);

-- Create indexes
CREATE UNIQUE INDEX idx_refresh_tokens_jti ON refresh_tokens (token_jti);
CREATE INDEX idx_refresh_tokens_user_id ON refresh_tokens (user_id);
CREATE INDEX idx_refresh_tokens_expires ON refresh_tokens (expires_at);

-- +goose Down
DROP TABLE IF EXISTS refresh_tokens;
