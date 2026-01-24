-- +goose Up
-- Fix incorrect duration_seconds for video resources
-- Correct durations based on actual video lengths:
-- ID 1: 39:50 = 2390 seconds
-- ID 2: 9:53 = 593 seconds
-- ID 3: 2:06:34 = 7594 seconds
-- ID 4: 17:13 = 1033 seconds
-- ID 5: 15:10 = 910 seconds
UPDATE mission_resources
SET
    duration_seconds = 2390,
    updated_at = NOW()
WHERE
    id = 1;

UPDATE mission_resources
SET
    duration_seconds = 593,
    updated_at = NOW()
WHERE
    id = 2;

UPDATE mission_resources
SET
    duration_seconds = 7594,
    updated_at = NOW()
WHERE
    id = 3;

UPDATE mission_resources
SET
    duration_seconds = 1033,
    updated_at = NOW()
WHERE
    id = 4;

UPDATE mission_resources
SET
    duration_seconds = 910,
    updated_at = NOW()
WHERE
    id = 5;

-- +goose Down
-- Revert to original incorrect values (for rollback purposes)
UPDATE mission_resources
SET
    duration_seconds = 256,
    updated_at = NOW()
WHERE
    id = 1;

UPDATE mission_resources
SET
    duration_seconds = 180,
    updated_at = NOW()
WHERE
    id = 2;

UPDATE mission_resources
SET
    duration_seconds = 420,
    updated_at = NOW()
WHERE
    id = 3;

UPDATE mission_resources
SET
    duration_seconds = 360,
    updated_at = NOW()
WHERE
    id = 4;

UPDATE mission_resources
SET
    duration_seconds = 540,
    updated_at = NOW()
WHERE
    id = 5;
