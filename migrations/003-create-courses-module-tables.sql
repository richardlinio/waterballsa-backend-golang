--liquibase formatted sql

--changeset liquibase:003-create-courses-module-tables
--comment: Create course module tables (courses, chapters, units)

-- Create courses table
CREATE TABLE courses (
    id BIGSERIAL PRIMARY KEY,
    title VARCHAR(200) NOT NULL,
    description TEXT,
    cover_image_url VARCHAR(500),
    price DECIMAL(10,2) NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP NULL
);

-- Create courses indexes
CREATE INDEX idx_courses_title ON courses(title);

-- Create chapters table
CREATE TABLE chapters (
    id BIGSERIAL PRIMARY KEY,
    course_id BIGINT NOT NULL,
    title VARCHAR(255) NOT NULL,
    order_index INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP NULL,
    CONSTRAINT fk_chapters_course FOREIGN KEY (course_id) REFERENCES courses(id) ON DELETE CASCADE
);

-- Create chapters indexes
CREATE UNIQUE INDEX idx_chapters_course_order ON chapters(course_id, order_index);

-- Create units table
CREATE TABLE units (
    id BIGSERIAL PRIMARY KEY,
    chapter_id BIGINT NOT NULL,
    title VARCHAR(255) NOT NULL,
    video_url VARCHAR(500) NOT NULL,
    duration_seconds INTEGER NOT NULL,
    experience_reward INTEGER NOT NULL DEFAULT 100,
    is_free_preview BOOLEAN NOT NULL DEFAULT false,
    order_index INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP NULL,
    CONSTRAINT fk_units_chapter FOREIGN KEY (chapter_id) REFERENCES chapters(id) ON DELETE CASCADE
);

-- Create units indexes
CREATE UNIQUE INDEX idx_units_chapter_order ON units(chapter_id, order_index);

--rollback DROP TABLE IF EXISTS units;
--rollback DROP TABLE IF EXISTS chapters;
--rollback DROP TABLE IF EXISTS courses;
