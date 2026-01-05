-- liquibase formatted sql

-- changeset waterballsa:020-add-order-expired-status
-- comment: Add EXPIRED status to order_status enum and add expired_at column to orders table

-- Add EXPIRED value to order_status enum
ALTER TYPE order_status ADD VALUE 'EXPIRED';

-- Add expired_at column to orders table
ALTER TABLE orders ADD COLUMN expired_at TIMESTAMP NULL;

-- Add index for expired_at to optimize scheduled task queries
CREATE INDEX idx_orders_expired_at ON orders(expired_at);

-- Add comment for the new column
COMMENT ON COLUMN orders.expired_at IS '訂單過期時間（建立時設為 created_at + 3天）';

-- rollback ALTER TABLE orders DROP COLUMN expired_at;
-- rollback DROP INDEX idx_orders_expired_at;
-- rollback -- Note: PostgreSQL does not support removing enum values directly. Manual intervention required.
