-- +goose Up
-- Rename mission_status enum type to progress_status

-- Rename enum type from mission_status to progress_status
ALTER TYPE mission_status RENAME TO progress_status;

-- +goose Down
ALTER TYPE progress_status RENAME TO mission_status;
