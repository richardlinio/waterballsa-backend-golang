-- name: CreateOrder :one
INSERT INTO orders (
    order_number,
    user_id,
    status,
    original_price,
    discount,
    price,
    expired_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
) RETURNING id, order_number, user_id, status, original_price, discount, price, created_at, expired_at, paid_at, updated_at;

-- name: CreateOrderItem :one
INSERT INTO order_items (
    order_id,
    journey_id,
    quantity,
    original_price,
    discount,
    price
) VALUES (
    $1, $2, $3, $4, $5, $6
) RETURNING id, order_id, journey_id, quantity, original_price, discount, price, created_at;

-- name: GetOrderByID :one
SELECT id, order_number, user_id, status, original_price, discount, price, created_at, expired_at, paid_at, updated_at
FROM orders
WHERE id = $1 AND deleted_at IS NULL;

-- name: GetOrderItemsByOrderID :many
SELECT id, order_id, journey_id, quantity, original_price, discount, price, created_at
FROM order_items
WHERE order_id = $1 AND deleted_at IS NULL;

-- name: GetOrderItemsByOrderIDs :many
SELECT id, order_id, journey_id, quantity, original_price, discount, price, created_at
FROM order_items
WHERE order_id = ANY($1::BIGINT[]) AND deleted_at IS NULL
ORDER BY order_id, id;

-- name: CheckUserHasPurchasedJourney :one
SELECT EXISTS(
    SELECT 1 FROM user_journeys
    WHERE user_id = $1 AND journey_id = $2 AND deleted_at IS NULL
) AS has_purchased;

-- name: GetUnpaidOrderByUserAndJourney :one
SELECT o.id, o.order_number, o.user_id, o.status, o.original_price, o.discount, o.price, o.created_at, o.expired_at, o.paid_at, o.updated_at
FROM orders o
INNER JOIN order_items oi ON o.id = oi.order_id
WHERE o.user_id = $1 AND oi.journey_id = $2 AND o.status = 'UNPAID' AND o.deleted_at IS NULL
LIMIT 1;

-- name: UpdateOrderStatusToPaid :one
UPDATE orders
SET status = 'PAID',
    paid_at = now(),
    updated_at = now()
WHERE id = $1 AND deleted_at IS NULL
RETURNING id, order_number, user_id, status, original_price, discount, price, created_at, expired_at, paid_at, updated_at;

-- name: CreateUserJourney :exec
INSERT INTO user_journeys (user_id, journey_id, order_id, purchased_at)
VALUES ($1, $2, $3, now());

-- name: GetOrdersByUserID :many
SELECT id, order_number, user_id, status, original_price, discount, price,
       created_at, expired_at, paid_at, updated_at
FROM orders
WHERE user_id = $1 AND deleted_at IS NULL
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: CountOrdersByUserID :one
SELECT COUNT(*)
FROM orders
WHERE user_id = $1 AND deleted_at IS NULL;
