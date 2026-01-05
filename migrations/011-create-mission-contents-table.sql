--liquibase formatted sql

--changeset liquibase:011-create-mission-contents-table
--comment: Create mission contents module tables and enum types

-- Create content_type enum
CREATE TYPE content_type AS ENUM ('VIDEO', 'ARTICLE', 'FORM');

-- Create mission_contents table
CREATE TABLE mission_contents (
    id BIGSERIAL PRIMARY KEY,
    mission_id BIGINT NOT NULL,
    content_type content_type NOT NULL,
    content_url VARCHAR(1000) NOT NULL,
    content_order INTEGER NOT NULL DEFAULT 0,
    duration_seconds INTEGER NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP NULL,
    CONSTRAINT fk_contents_mission FOREIGN KEY (mission_id) REFERENCES missions(id) ON DELETE CASCADE
);

-- Create mission_contents indexes
CREATE INDEX idx_contents_mission_id ON mission_contents(mission_id);
CREATE UNIQUE INDEX idx_contents_mission_order ON mission_contents(mission_id, content_order);

--rollback DROP TABLE IF EXISTS mission_contents;
--rollback DROP TYPE IF EXISTS content_type;
