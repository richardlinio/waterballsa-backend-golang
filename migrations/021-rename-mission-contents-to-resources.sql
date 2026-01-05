--liquibase formatted sql

--changeset liquibase:021-rename-mission-contents-to-resources
--comment: Rename mission_contents table and related fields to mission_resources

-- Step 1: Rename the enum type
ALTER TYPE content_type RENAME TO resource_type;

-- Step 2: Rename the table
ALTER TABLE mission_contents RENAME TO mission_resources;

-- Step 3: Rename the column from content_type to resource_type
ALTER TABLE mission_resources RENAME COLUMN content_type TO resource_type;

-- Step 4: Rename the column from content_url to resource_url
ALTER TABLE mission_resources RENAME COLUMN content_url TO resource_url;

-- Step 5: Make resource_url nullable (optional)
ALTER TABLE mission_resources ALTER COLUMN resource_url DROP NOT NULL;

-- Step 6: Add the new resource_content column (optional)
ALTER TABLE mission_resources ADD COLUMN resource_content TEXT NULL;

-- Step 7: Rename the foreign key constraint
ALTER TABLE mission_resources RENAME CONSTRAINT fk_contents_mission TO fk_resources_mission;

-- Step 8: Rename the indexes
ALTER INDEX idx_contents_mission_id RENAME TO idx_resources_mission_id;
ALTER INDEX idx_contents_mission_order RENAME TO idx_resources_mission_order;

--rollback ALTER INDEX idx_resources_mission_order RENAME TO idx_contents_mission_order;
--rollback ALTER INDEX idx_resources_mission_id RENAME TO idx_contents_mission_id;
--rollback ALTER TABLE mission_resources RENAME CONSTRAINT fk_resources_mission TO fk_contents_mission;
--rollback ALTER TABLE mission_resources DROP COLUMN resource_content;
--rollback ALTER TABLE mission_resources ALTER COLUMN resource_url SET NOT NULL;
--rollback ALTER TABLE mission_resources RENAME COLUMN resource_url TO content_url;
--rollback ALTER TABLE mission_resources RENAME COLUMN resource_type TO content_type;
--rollback ALTER TABLE mission_resources RENAME TO mission_contents;
--rollback ALTER TYPE resource_type RENAME TO content_type;
