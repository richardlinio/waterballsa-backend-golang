--liquibase formatted sql

--changeset liquibase:013-add-slug-to-journeys
--comment: Add slug column to journeys table for URL-friendly identifiers

-- Add slug column to journeys table
ALTER TABLE journeys ADD COLUMN slug VARCHAR(200);

-- Set slug for existing journey
UPDATE journeys
SET slug = 'software-design-pattern'
WHERE id = 1;

-- Make the column NOT NULL after setting values
ALTER TABLE journeys ALTER COLUMN slug SET NOT NULL;

-- Add unique constraint to slug column
ALTER TABLE journeys ADD CONSTRAINT uk_journeys_slug UNIQUE (slug);

-- Create index on slug for faster lookups
CREATE INDEX idx_journeys_slug ON journeys(slug);

--rollback DROP INDEX IF EXISTS idx_journeys_slug;
--rollback ALTER TABLE journeys DROP CONSTRAINT IF EXISTS uk_journeys_slug;
--rollback ALTER TABLE journeys DROP COLUMN slug;
