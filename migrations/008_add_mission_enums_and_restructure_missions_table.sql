-- +goose Up
-- Add mission-related enum types and restructure missions table
-- Create mission_type enum
CREATE TYPE mission_type AS ENUM('VIDEO', 'ARTICLE', 'QUESTIONNAIRE');

-- Create mission_status enum
CREATE TYPE mission_status AS ENUM('UNCOMPLETED', 'COMPLETED', 'DELIVERED');

-- Create mission_access_level enum
CREATE TYPE mission_access_level AS ENUM('PUBLIC', 'AUTHENTICATED', 'PURCHASED');

-- Add new columns to missions table
ALTER TABLE missions
ADD COLUMN
TYPE mission_type NOT NULL DEFAULT 'VIDEO';

ALTER TABLE missions
ADD COLUMN description TEXT;

ALTER TABLE missions
ADD COLUMN access_level mission_access_level NOT NULL DEFAULT 'PURCHASED';

-- Create index for access_level
CREATE INDEX idx_missions_access_level ON missions (access_level);

-- Remove old columns from missions table
ALTER TABLE missions
DROP COLUMN video_url;

ALTER TABLE missions
DROP COLUMN duration_seconds;

ALTER TABLE missions
DROP COLUMN experience_reward;

ALTER TABLE missions
DROP COLUMN is_free_preview;

-- +goose Down
ALTER TABLE missions
ADD COLUMN is_free_preview BOOLEAN NOT NULL DEFAULT FALSE;

ALTER TABLE missions
ADD COLUMN experience_reward INTEGER NOT NULL DEFAULT 100;

ALTER TABLE missions
ADD COLUMN duration_seconds INTEGER NOT NULL DEFAULT 0;

ALTER TABLE missions
ADD COLUMN video_url VARCHAR(500) NOT NULL DEFAULT '';

DROP INDEX IF EXISTS idx_missions_access_level;

ALTER TABLE missions
DROP COLUMN access_level;

ALTER TABLE missions
DROP COLUMN description;

ALTER TABLE missions
DROP COLUMN
TYPE;

DROP TYPE IF EXISTS mission_access_level;

DROP TYPE IF EXISTS mission_status;

DROP TYPE IF EXISTS mission_type;