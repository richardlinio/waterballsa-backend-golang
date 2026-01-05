--liquibase formatted sql

--changeset liquibase:008-add-mission-enums-and-restructure-missions-table
--comment: Add mission-related enum types and restructure missions table

-- Create mission_type enum
CREATE TYPE mission_type AS ENUM ('VIDEO', 'ARTICLE', 'QUESTIONNAIRE');

-- Create mission_status enum
CREATE TYPE mission_status AS ENUM ('UNCOMPLETED', 'COMPLETED', 'DELIVERED');

-- Create mission_access_level enum
CREATE TYPE mission_access_level AS ENUM ('PUBLIC', 'AUTHENTICATED', 'PURCHASED');

-- Add new columns to missions table
ALTER TABLE missions ADD COLUMN type mission_type NOT NULL DEFAULT 'VIDEO';
ALTER TABLE missions ADD COLUMN description TEXT;
ALTER TABLE missions ADD COLUMN access_level mission_access_level NOT NULL DEFAULT 'PURCHASED';

-- Create index for access_level
CREATE INDEX idx_missions_access_level ON missions(access_level);

-- Remove old columns from missions table
ALTER TABLE missions DROP COLUMN video_url;
ALTER TABLE missions DROP COLUMN duration_seconds;
ALTER TABLE missions DROP COLUMN experience_reward;
ALTER TABLE missions DROP COLUMN is_free_preview;

--rollback ALTER TABLE missions ADD COLUMN is_free_preview BOOLEAN NOT NULL DEFAULT false;
--rollback ALTER TABLE missions ADD COLUMN experience_reward INTEGER NOT NULL DEFAULT 100;
--rollback ALTER TABLE missions ADD COLUMN duration_seconds INTEGER NOT NULL DEFAULT 0;
--rollback ALTER TABLE missions ADD COLUMN video_url VARCHAR(500) NOT NULL DEFAULT '';
--rollback DROP INDEX IF EXISTS idx_missions_access_level;
--rollback ALTER TABLE missions DROP COLUMN access_level;
--rollback ALTER TABLE missions DROP COLUMN description;
--rollback ALTER TABLE missions DROP COLUMN type;
--rollback DROP TYPE IF EXISTS mission_access_level;
--rollback DROP TYPE IF EXISTS mission_status;
--rollback DROP TYPE IF EXISTS mission_type;
