--liquibase formatted sql

--changeset liquibase:004-create-user-unit-progress-table
--comment: Create user unit progress table for tracking video watch position

-- Create user_unit_progress table
CREATE TABLE user_unit_progress (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    unit_id BIGINT NOT NULL,
    watch_position_seconds INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP NULL,
    CONSTRAINT fk_progress_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_progress_unit FOREIGN KEY (unit_id) REFERENCES units(id) ON DELETE CASCADE
);

-- Create user_unit_progress indexes
CREATE UNIQUE INDEX idx_progress_user_unit ON user_unit_progress(user_id, unit_id);
CREATE INDEX idx_progress_user_id ON user_unit_progress(user_id);
CREATE INDEX idx_progress_unit_id ON user_unit_progress(unit_id);

--rollback DROP TABLE IF EXISTS user_unit_progress;
