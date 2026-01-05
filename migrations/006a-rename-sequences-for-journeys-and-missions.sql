--liquibase formatted sql

--changeset liquibase:006a-rename-sequences-for-journeys-and-missions
--comment: Rename sequences to match the renamed tables (courses -> journeys, units -> missions)

-- Rename courses_id_seq to journeys_id_seq
ALTER SEQUENCE courses_id_seq RENAME TO journeys_id_seq;

-- Rename units_id_seq to missions_id_seq
ALTER SEQUENCE units_id_seq RENAME TO missions_id_seq;

--rollback ALTER SEQUENCE missions_id_seq RENAME TO units_id_seq;
--rollback ALTER SEQUENCE journeys_id_seq RENAME TO courses_id_seq;
