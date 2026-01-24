-- name: GetUserMissionProgress :one
SELECT
    COALESCE(p.id, 0) as id,
    $1::bigint as user_id,
    m.id as mission_id,
    COALESCE(p.status, 'UNCOMPLETED'::progress_status) as status,
    COALESCE(p.watch_position_seconds, 0) as watch_position_seconds,
    COALESCE(p.created_at, NOW()) as created_at,
    COALESCE(p.updated_at, NOW()) as updated_at
FROM
    missions m
LEFT JOIN
    user_mission_progress p
    ON m.id = p.mission_id AND p.user_id = $1 AND p.deleted_at IS NULL
WHERE
    m.id = $2
    AND m.deleted_at IS NULL;

-- name: UpsertUserMissionProgress :one
INSERT INTO
    user_mission_progress (user_id, mission_id, status, watch_position_seconds, created_at, updated_at)
VALUES
    ($1, $2, $3, $4, NOW(), NOW()) ON CONFLICT (user_id, mission_id)
DO
UPDATE
SET
    status = EXCLUDED.status,
    watch_position_seconds = EXCLUDED.watch_position_seconds,
    updated_at = NOW()
RETURNING
    id,
    user_id,
    mission_id,
    status,
    watch_position_seconds,
    created_at,
    updated_at;