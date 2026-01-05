--liquibase formatted sql

--changeset liquibase:018-create-order-items-table
--comment: Create order_items table to store order line items

-- Create order_items table
CREATE TABLE order_items (
    id BIGSERIAL PRIMARY KEY,
    order_id BIGINT NOT NULL,
    journey_id BIGINT NOT NULL,
    quantity INTEGER NOT NULL DEFAULT 1,
    original_price DECIMAL(10,2) NOT NULL,
    discount DECIMAL(10,2) NOT NULL DEFAULT 0,
    price DECIMAL(10,2) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP NULL,
    CONSTRAINT fk_order_items_order FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE CASCADE,
    CONSTRAINT fk_order_items_journey FOREIGN KEY (journey_id) REFERENCES journeys(id) ON DELETE CASCADE
);

-- Create order_items indexes
CREATE INDEX idx_order_items_order_id ON order_items(order_id);
CREATE INDEX idx_order_items_journey_id ON order_items(journey_id);
CREATE UNIQUE INDEX idx_order_items_order_journey ON order_items(order_id, journey_id);

--rollback DROP TABLE IF EXISTS order_items;
