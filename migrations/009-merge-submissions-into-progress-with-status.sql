-- +goose Up
-- Add status column to user_mission_progress and drop user_mission_submissions table

-- Add status column to user_mission_progress table
ALTER TABLE user_mission_progress ADD COLUMN status mission_status NOT NULL DEFAULT 'UNCOMPLETED';

-- Create index for status
CREATE INDEX idx_progress_status ON user_mission_progress(status);

-- Migrate data from user_mission_submissions to user_mission_progress
-- Mark all submitted missions as DELIVERED
UPDATE user_mission_progress
SET status = 'DELIVERED', updated_at = NOW()
WHERE (user_id, mission_id) IN (
    SELECT user_id, mission_id FROM user_mission_submissions
);

-- Drop user_mission_submissions table as its functionality is now merged into user_mission_progress
DROP TABLE IF EXISTS user_mission_submissions;

-- +goose Down
CREATE TABLE user_mission_submissions (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    mission_id BIGINT NOT NULL,
    submitted_at TIMESTAMP NOT NULL DEFAULT NOW(),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP NULL,
    CONSTRAINT fk_submissions_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_submissions_mission FOREIGN KEY (mission_id) REFERENCES missions(id) ON DELETE CASCADE
);
CREATE UNIQUE INDEX idx_submissions_user_mission ON user_mission_submissions(user_id, mission_id);
CREATE INDEX idx_submissions_user_id ON user_mission_submissions(user_id);
CREATE INDEX idx_submissions_time ON user_mission_submissions(submitted_at);
DROP INDEX IF EXISTS idx_progress_status;
ALTER TABLE user_mission_progress DROP COLUMN status;
