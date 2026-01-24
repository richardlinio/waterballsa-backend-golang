-- name: GetUserMissionProgress :one
SELECT
    id,
    user_id,
    mission_id,
    status,
    watch_position_seconds,
    created_at,
    updated_at
FROM
    user_mission_progress
WHERE
    user_id = $1
    AND mission_id = $2
    AND deleted_at IS NULL;

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