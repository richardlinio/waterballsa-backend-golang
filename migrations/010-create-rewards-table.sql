--liquibase formatted sql

--changeset liquibase:010-create-rewards-table
--comment: Create rewards module tables and enum types

-- Create reward_type enum
CREATE TYPE reward_type AS ENUM ('EXPERIENCE');

-- Create rewards table
CREATE TABLE rewards (
    id BIGSERIAL PRIMARY KEY,
    mission_id BIGINT NOT NULL,
    reward_type reward_type NOT NULL,
    reward_value INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP NULL,
    CONSTRAINT fk_rewards_mission FOREIGN KEY (mission_id) REFERENCES missions(id) ON DELETE CASCADE
);

-- Create rewards indexes
CREATE INDEX idx_rewards_mission_id ON rewards(mission_id);
CREATE UNIQUE INDEX idx_rewards_mission_type ON rewards(mission_id, reward_type);

--rollback DROP TABLE IF EXISTS rewards;
--rollback DROP TYPE IF EXISTS reward_type;
