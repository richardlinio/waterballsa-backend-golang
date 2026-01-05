--liquibase formatted sql

--changeset liquibase:007-add-teacher-name-to-journeys
--comment: Add teacher_name column to journeys table

-- Add teacher_name column to journeys table
ALTER TABLE journeys ADD COLUMN teacher_name VARCHAR(100);

-- Update existing records with a default value (you may need to adjust this based on your data)
UPDATE journeys SET teacher_name = 'Unknown Teacher' WHERE teacher_name IS NULL;

-- Make the column NOT NULL after setting default values
ALTER TABLE journeys ALTER COLUMN teacher_name SET NOT NULL;

--rollback ALTER TABLE journeys DROP COLUMN teacher_name;
