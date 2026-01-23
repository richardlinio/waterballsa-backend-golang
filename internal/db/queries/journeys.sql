-- name: ListJourneys :many
SELECT id, title, slug, description, cover_image_url, teacher_name, price, created_at, updated_at
FROM journeys
WHERE deleted_at IS NULL
ORDER BY created_at ASC;
