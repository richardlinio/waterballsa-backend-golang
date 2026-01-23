-- name: GetMissionByID :one
SELECT
    m.id,
    m.chapter_id,
    m.type,
    m.title,
    m.description,
    m.access_level,
    m.order_index,
    m.created_at,
    m.updated_at,
    c.journey_id
FROM
    missions m
    INNER JOIN chapters c ON m.chapter_id = c.id
WHERE
    m.id = $1
    AND m.deleted_at IS NULL
    AND c.deleted_at IS NULL;

-- name: GetRewardByMissionID :one
SELECT
    id,
    mission_id,
    reward_type,
    reward_value,
    created_at,
    updated_at
FROM
    rewards
WHERE
    mission_id = $1
    AND deleted_at IS NULL;

-- name: ListResourcesByMissionID :many
SELECT
    id,
    mission_id,
    resource_type,
    resource_url,
    resource_content,
    content_order,
    duration_seconds,
    created_at,
    updated_at
FROM
    mission_resources
WHERE
    mission_id = $1
    AND deleted_at IS NULL
ORDER BY
    content_order ASC;