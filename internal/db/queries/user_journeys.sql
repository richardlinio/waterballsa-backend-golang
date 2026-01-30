-- name: GetUserJourneysByUserID :many
SELECT
    uj.user_id,
    uj.journey_id,
    uj.order_id,
    uj.purchased_at,
    j.title AS journey_title,
    j.slug AS journey_slug,
    j.cover_image_url,
    j.teacher_name,
    o.order_number
FROM user_journeys uj
INNER JOIN journeys j ON uj.journey_id = j.id
INNER JOIN orders o ON uj.order_id = o.id
WHERE uj.user_id = $1
  AND uj.deleted_at IS NULL
  AND j.deleted_at IS NULL
  AND o.deleted_at IS NULL
  AND o.status = 'PAID'
ORDER BY uj.purchased_at DESC;
