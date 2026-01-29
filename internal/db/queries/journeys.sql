-- name: ListJourneys :many
SELECT
    id,
    title,
    slug,
    description,
    cover_image_url,
    teacher_name,
    price,
    created_at,
    updated_at
FROM
    journeys
WHERE
    deleted_at IS NULL
ORDER BY
    created_at ASC;

-- name: GetJourneyByID :one
SELECT
    id,
    title,
    slug,
    description,
    cover_image_url,
    teacher_name,
    price,
    created_at,
    updated_at
FROM
    journeys
WHERE
    id = $1
    AND deleted_at IS NULL;

-- name: ListChaptersByJourneyID :many
SELECT
    id,
    journey_id,
    title,
    order_index,
    created_at,
    updated_at
FROM
    chapters
WHERE
    journey_id = $1
    AND deleted_at IS NULL
ORDER BY
    order_index ASC;

-- name: ListMissionsByChapterIDs :many
SELECT
    id,
    chapter_id,
    title,
TYPE,
access_level,
order_index,
created_at,
updated_at
FROM
    missions
WHERE
    chapter_id = ANY ($1::BIGINT[])
    AND deleted_at IS NULL
ORDER BY
    chapter_id ASC,
    order_index ASC;

-- name: GetJourneyTitleByID :one
SELECT title FROM journeys WHERE id = $1 AND deleted_at IS NULL;

-- name: GetJourneyTitlesByIDs :many
SELECT id, title FROM journeys WHERE id = ANY($1::BIGINT[]) AND deleted_at IS NULL;

-- name: UpdateJourneyPrice :exec
UPDATE journeys SET price = $2, updated_at = NOW() WHERE id = $1 AND deleted_at IS NULL;