--liquibase formatted sql
--changeset liquibase:016-fix-experience-values-and-access-levels
--comment: Fix test data - set experience to 0 for PUBLIC/ARTICLE/QUESTIONNAIRE missions, and change all AUTHENTICATED missions to PURCHASED

-- Update Mission 2 from PUBLIC to PURCHASED
UPDATE missions
SET access_level = 'PURCHASED',
    updated_at = NOW()
WHERE id = 2
  AND deleted_at IS NULL;

-- Update all AUTHENTICATED missions to PURCHASED
UPDATE missions
SET access_level = 'PURCHASED',
    updated_at = NOW()
WHERE access_level = 'AUTHENTICATED'
  AND deleted_at IS NULL;

-- Set experience to 0 for all PUBLIC missions
UPDATE rewards
SET reward_value = 0,
    updated_at = NOW()
WHERE reward_type = 'EXPERIENCE'
  AND mission_id IN (
    SELECT id FROM missions
    WHERE access_level = 'PUBLIC'
    AND deleted_at IS NULL
  )
  AND deleted_at IS NULL;

-- Set experience to 0 for all ARTICLE type missions
UPDATE rewards
SET reward_value = 0,
    updated_at = NOW()
WHERE reward_type = 'EXPERIENCE'
  AND mission_id IN (
    SELECT id FROM missions
    WHERE type = 'ARTICLE'
    AND deleted_at IS NULL
  )
  AND deleted_at IS NULL;

-- Set experience to 0 for all QUESTIONNAIRE type missions
UPDATE rewards
SET reward_value = 0,
    updated_at = NOW()
WHERE reward_type = 'EXPERIENCE'
  AND mission_id IN (
    SELECT id FROM missions
    WHERE type = 'QUESTIONNAIRE'
    AND deleted_at IS NULL
  )
  AND deleted_at IS NULL;

--rollback UPDATE missions SET access_level = 'PUBLIC', updated_at = NOW() WHERE id = 2 AND deleted_at IS NULL;
--rollback UPDATE missions SET access_level = 'AUTHENTICATED', updated_at = NOW() WHERE id IN (3,4,8) AND deleted_at IS NULL;
--rollback UPDATE rewards SET reward_value = 100, updated_at = NOW() WHERE id IN (1,2,4,7) AND reward_type = 'EXPERIENCE' AND deleted_at IS NULL;
--rollback UPDATE rewards SET reward_value = 50, updated_at = NOW() WHERE id = 8 AND reward_type = 'EXPERIENCE' AND deleted_at IS NULL;
