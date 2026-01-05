--liquibase formatted sql
--changeset liquibase:022-update-video-urls
--comment: Update video URLs for mission resources

-- Update video URLs for mission resources with IDs 1-5
UPDATE mission_resources SET resource_url = 'https://www.youtube.com/watch?v=Qv92qaIGbDg' WHERE id = 1;
UPDATE mission_resources SET resource_url = 'https://www.youtube.com/watch?v=i7twT3x5yv8' WHERE id = 2;
UPDATE mission_resources SET resource_url = 'https://www.youtube.com/watch?v=iCx3zwK8Ms8' WHERE id = 3;
UPDATE mission_resources SET resource_url = 'https://www.youtube.com/watch?v=BUE-icVYRFU' WHERE id = 4;
UPDATE mission_resources SET resource_url = 'https://www.youtube.com/watch?v=HZGCoVF3YvM' WHERE id = 5;

--rollback UPDATE mission_resources SET resource_url = 'https://youtu.be/iJHzesWhDj0' WHERE id = 1;
--rollback UPDATE mission_resources SET resource_url = 'https://youtu.be/upAg_4p-6p8' WHERE id = 2;
--rollback UPDATE mission_resources SET resource_url = 'https://youtu.be/1wODH0FT3pM' WHERE id = 3;
--rollback UPDATE mission_resources SET resource_url = 'https://youtu.be/A99IECKCJF8' WHERE id = 4;
--rollback UPDATE mission_resources SET resource_url = 'https://youtu.be/NBcP5XPgcUo' WHERE id = 5;
