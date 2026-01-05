-- +goose Up
-- Rename courses table to journeys and units table to missions

-- Rename units table to missions
ALTER TABLE units RENAME TO missions;

-- Rename units indexes
ALTER INDEX idx_units_chapter_order RENAME TO idx_missions_chapter_order;

-- Rename courses table to journeys
ALTER TABLE courses RENAME TO journeys;

-- Rename courses indexes
ALTER INDEX idx_courses_title RENAME TO idx_journeys_title;

-- Rename foreign key column in chapters table
ALTER TABLE chapters RENAME COLUMN course_id TO journey_id;

-- Rename chapters constraint and index
ALTER TABLE chapters DROP CONSTRAINT fk_chapters_course;
ALTER TABLE chapters ADD CONSTRAINT fk_chapters_journey FOREIGN KEY (journey_id) REFERENCES journeys(id) ON DELETE CASCADE;
ALTER INDEX idx_chapters_course_order RENAME TO idx_chapters_journey_order;

-- Rename user_unit_progress table to user_mission_progress
ALTER TABLE user_unit_progress RENAME TO user_mission_progress;

-- Rename foreign key column in user_mission_progress table
ALTER TABLE user_mission_progress RENAME COLUMN unit_id TO mission_id;

-- Rename user_mission_progress constraints and indexes
ALTER TABLE user_mission_progress DROP CONSTRAINT fk_progress_unit;
ALTER TABLE user_mission_progress ADD CONSTRAINT fk_progress_mission FOREIGN KEY (mission_id) REFERENCES missions(id) ON DELETE CASCADE;
ALTER INDEX idx_progress_user_unit RENAME TO idx_progress_user_mission;
ALTER INDEX idx_progress_unit_id RENAME TO idx_progress_mission_id;

-- Rename user_unit_submissions table to user_mission_submissions
ALTER TABLE user_unit_submissions RENAME TO user_mission_submissions;

-- Rename foreign key column in user_mission_submissions table
ALTER TABLE user_mission_submissions RENAME COLUMN unit_id TO mission_id;

-- Rename user_mission_submissions constraints and indexes
ALTER TABLE user_mission_submissions DROP CONSTRAINT fk_submissions_unit;
ALTER TABLE user_mission_submissions ADD CONSTRAINT fk_submissions_mission FOREIGN KEY (mission_id) REFERENCES missions(id) ON DELETE CASCADE;
ALTER INDEX idx_submissions_user_unit RENAME TO idx_submissions_user_mission;

-- +goose Down
ALTER INDEX idx_submissions_user_mission RENAME TO idx_submissions_user_unit;
ALTER TABLE user_mission_submissions DROP CONSTRAINT fk_submissions_mission;
ALTER TABLE user_mission_submissions ADD CONSTRAINT fk_submissions_unit FOREIGN KEY (unit_id) REFERENCES units(id) ON DELETE CASCADE;
ALTER TABLE user_mission_submissions RENAME COLUMN mission_id TO unit_id;
ALTER TABLE user_mission_submissions RENAME TO user_unit_submissions;
ALTER INDEX idx_progress_mission_id RENAME TO idx_progress_unit_id;
ALTER INDEX idx_progress_user_mission RENAME TO idx_progress_user_unit;
ALTER TABLE user_mission_progress DROP CONSTRAINT fk_progress_mission;
ALTER TABLE user_mission_progress ADD CONSTRAINT fk_progress_unit FOREIGN KEY (unit_id) REFERENCES units(id) ON DELETE CASCADE;
ALTER TABLE user_mission_progress RENAME COLUMN mission_id TO unit_id;
ALTER TABLE user_mission_progress RENAME TO user_unit_progress;
ALTER INDEX idx_chapters_journey_order RENAME TO idx_chapters_course_order;
ALTER TABLE chapters DROP CONSTRAINT fk_chapters_journey;
ALTER TABLE chapters ADD CONSTRAINT fk_chapters_course FOREIGN KEY (course_id) REFERENCES courses(id) ON DELETE CASCADE;
ALTER TABLE chapters RENAME COLUMN journey_id TO course_id;
ALTER INDEX idx_journeys_title RENAME TO idx_courses_title;
ALTER TABLE journeys RENAME TO courses;
ALTER INDEX idx_missions_chapter_order RENAME TO idx_units_chapter_order;
ALTER TABLE missions RENAME TO units;
