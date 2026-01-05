--liquibase formatted sql

--changeset liquibase:002-create-access-tokens-table
--comment: Create access tokens (JWT blacklist) table (Authentication Module)

-- Create access_tokens table
CREATE TABLE access_tokens (
    id BIGSERIAL PRIMARY KEY,
    token_jti VARCHAR(64) NOT NULL UNIQUE,
    user_id BIGINT NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    invalidated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_access_tokens_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Create indexes
CREATE UNIQUE INDEX idx_invalid_tokens_jti ON access_tokens(token_jti);
CREATE INDEX idx_invalid_tokens_user_id ON access_tokens(user_id);
CREATE INDEX idx_invalid_tokens_expires ON access_tokens(expires_at);

--rollback DROP TABLE IF EXISTS access_tokens;
