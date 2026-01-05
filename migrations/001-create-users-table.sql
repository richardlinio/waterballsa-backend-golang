--liquibase formatted sql

--changeset liquibase:001-create-users-table
--comment: Create users table (Authentication Module)

-- Create user_role ENUM type
CREATE TYPE user_role AS ENUM ('STUDENT', 'TEACHER', 'ADMIN');

-- Create users table
CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE,
    password_hash VARCHAR(72) NOT NULL,
    role user_role NOT NULL DEFAULT 'STUDENT',
    experience_points INTEGER NOT NULL DEFAULT 0,
    level INTEGER NOT NULL DEFAULT 1,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP NULL
);

-- Create indexes
CREATE UNIQUE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_experience ON users(experience_points);

--rollback DROP TABLE IF EXISTS users;
--rollback DROP TYPE IF EXISTS user_role;
