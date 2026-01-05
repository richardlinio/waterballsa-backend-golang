--liquibase formatted sql

--changeset liquibase:019-create-user-journeys-table
--comment: Create user_journeys table to track user's purchased journeys

-- Create user_journeys table
CREATE TABLE user_journeys (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    journey_id BIGINT NOT NULL,
    order_id BIGINT NOT NULL,
    purchased_at TIMESTAMP NOT NULL DEFAULT NOW(),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP NULL,
    CONSTRAINT fk_user_journeys_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_user_journeys_journey FOREIGN KEY (journey_id) REFERENCES journeys(id) ON DELETE CASCADE,
    CONSTRAINT fk_user_journeys_order FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE CASCADE
);

-- Create user_journeys indexes
CREATE UNIQUE INDEX idx_user_journeys_user_journey ON user_journeys(user_id, journey_id);
CREATE INDEX idx_user_journeys_user_id ON user_journeys(user_id);
CREATE INDEX idx_user_journeys_journey_id ON user_journeys(journey_id);
CREATE INDEX idx_user_journeys_order_id ON user_journeys(order_id);

--rollback DROP TABLE IF EXISTS user_journeys;
