--liquibase formatted sql
--changeset liquibase:015-update-video-duration-to-26-seconds
--comment: Update all video content duration to 26 seconds to match actual test video length

UPDATE mission_contents
SET duration_seconds = 26,
    updated_at = NOW()
WHERE content_type = 'VIDEO'
  AND deleted_at IS NULL;

--rollback UPDATE mission_contents SET duration_seconds = CASE id WHEN 1 THEN 256 WHEN 2 THEN 180 WHEN 3 THEN 420 WHEN 4 THEN 360 WHEN 5 THEN 540 END, updated_at = NOW() WHERE content_type = 'VIDEO' AND deleted_at IS NULL;
