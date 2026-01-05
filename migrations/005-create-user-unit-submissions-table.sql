--liquibase formatted sql

--changeset liquibase:005-create-user-unit-submissions-table
--comment: Create user unit submissions table for tracking completed units

-- Create user_unit_submissions table
CREATE TABLE user_unit_submissions (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    unit_id BIGINT NOT NULL,
    submitted_at TIMESTAMP NOT NULL DEFAULT NOW(),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP NULL,
    CONSTRAINT fk_submissions_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_submissions_unit FOREIGN KEY (unit_id) REFERENCES units(id) ON DELETE CASCADE
);

-- Create user_unit_submissions indexes
CREATE UNIQUE INDEX idx_submissions_user_unit ON user_unit_submissions(user_id, unit_id);
CREATE INDEX idx_submissions_user_id ON user_unit_submissions(user_id);
CREATE INDEX idx_submissions_time ON user_unit_submissions(submitted_at);

--rollback DROP TABLE IF EXISTS user_unit_submissions;
