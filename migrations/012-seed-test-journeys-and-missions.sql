--liquibase formatted sql
--changeset liquibase:012-seed-test-journeys-and-missions
--comment: Seed test data for journeys, chapters, missions, contents, and rewards
-- Insert test journey
INSERT INTO
    journeys (id, title, description, cover_image_url, teacher_name, price, created_at, updated_at, deleted_at)
VALUES
    (
        1,
        '軟體設計模式精通之旅',
        '用 C.A. 模式大大提昇系統思維能力',
        'https://images.unsplash.com/photo-1516116216624-53e697fedbea?w=800',
        '水球潘',
        3000.00,
        NOW(),
        NOW(),
        NULL
    );


-- Insert test chapters
INSERT INTO
    chapters (id, journey_id, title, order_index, created_at, updated_at, deleted_at)
VALUES
    (1, 1, '課程介紹', 1, NOW(), NOW(), NULL),
    (2, 1, '架構思維的 C.A. 模式', 2, NOW(), NOW(), NULL),
    (3, 1, '設計模式實戰', 3, NOW(), NOW(), NULL);


-- Insert test missions
INSERT INTO
    missions (
        id,
        chapter_id,
        TYPE,
        title,
        description,
        access_level,
        order_index,
        created_at,
        updated_at,
        deleted_at
    )
VALUES
    -- Chapter 1: 課程介紹
    (1, 1, 'VIDEO', '課程介紹：這門課手把手帶你成為架構設計的高手', '思維升級：Christopher Alexander 模式六大要素及模式語言', 'PUBLIC', 1, NOW(), NOW(), NULL),
    (2, 1, 'VIDEO', '你該知道：在 AI 的時代下，只會下 prompt 絕對寫不出好 Code', '深入探討 AI 輔助開發的本質與限制', 'PUBLIC', 2, NOW(), NOW(), NULL),
    -- Chapter 2: 架構思維的 C.A. 模式
    (3, 2, 'VIDEO', '架構師該如何思考', '從 Christopher Alexander 的建築思維學習軟體架構', 'AUTHENTICATED', 1, NOW(), NOW(), NULL),
    (4, 2, 'ARTICLE', 'UML 不是英文縮寫字', 'UML 統一塑模語言完整介紹', 'AUTHENTICATED', 2, NOW(), NOW(), NULL),
    (5, 2, 'VIDEO', '模式語言：從建築到軟體', '探索模式語言如何應用在軟體設計中', 'PURCHASED', 3, NOW(), NOW(), NULL),
    -- Chapter 3: 設計模式實戰
    (6, 3, 'VIDEO', 'Strategy Pattern 策略模式實作', '透過實例學習策略模式的應用場景', 'PURCHASED', 1, NOW(), NOW(), NULL),
    (7, 3, 'ARTICLE', 'Observer Pattern 觀察者模式深度解析', '理解 Observer Pattern 的核心概念與實作技巧', 'PURCHASED', 2, NOW(), NOW(), NULL),
    (8, 3, 'QUESTIONNAIRE', '課程回饋問卷', '請填寫課程回饋問卷，幫助我們改進課程內容', 'AUTHENTICATED', 3, NOW(), NOW(), NULL);


-- Insert mission contents
INSERT INTO
    mission_contents (id, mission_id, content_type, content_url, content_order, duration_seconds, created_at, updated_at, deleted_at)
VALUES
    -- Video contents (using public YouTube videos)
    (1, 1, 'VIDEO', 'https://youtu.be/iJHzesWhDj0', 0, 256, NOW(), NOW(), NULL),
    (2, 2, 'VIDEO', 'https://youtu.be/upAg_4p-6p8', 0, 180, NOW(), NOW(), NULL),
    (3, 3, 'VIDEO', 'https://youtu.be/1wODH0FT3pM', 0, 420, NOW(), NOW(), NULL),
    (4, 5, 'VIDEO', 'https://youtu.be/A99IECKCJF8', 0, 360, NOW(), NOW(), NULL),
    (5, 6, 'VIDEO', 'https://youtu.be/NBcP5XPgcUo', 0, 540, NOW(), NOW(), NULL),
    -- Article contents
    (6, 4, 'ARTICLE', 'https://cdn.waterballsa.tw/articles/uml-intro.html', 0, NULL, NOW(), NOW(), NULL),
    (7, 7, 'ARTICLE', 'https://cdn.waterballsa.tw/articles/observer-pattern.html', 0, NULL, NOW(), NOW(), NULL),
    -- Form/Questionnaire content
    (8, 8, 'FORM', 'https://forms.waterballsa.tw/feedback-form-1', 0, NULL, NOW(), NOW(), NULL);


-- Insert rewards for all missions
INSERT INTO
    rewards (id, mission_id, reward_type, reward_value, created_at, updated_at, deleted_at)
VALUES
    (1, 1, 'EXPERIENCE', 100, NOW(), NOW(), NULL),
    (2, 2, 'EXPERIENCE', 100, NOW(), NOW(), NULL),
    (3, 3, 'EXPERIENCE', 100, NOW(), NOW(), NULL),
    (4, 4, 'EXPERIENCE', 100, NOW(), NOW(), NULL),
    (5, 5, 'EXPERIENCE', 100, NOW(), NOW(), NULL),
    (6, 6, 'EXPERIENCE', 100, NOW(), NOW(), NULL),
    (7, 7, 'EXPERIENCE', 100, NOW(), NOW(), NULL),
    (8, 8, 'EXPERIENCE', 50, NOW(), NOW(), NULL);


-- Reset sequences to prevent ID conflicts
SELECT
    SETVAL(
        'journeys_id_seq',
        (
            SELECT
                MAX(id)
            FROM
                journeys
        )
    );


SELECT
    SETVAL(
        'chapters_id_seq',
        (
            SELECT
                MAX(id)
            FROM
                chapters
        )
    );


SELECT
    SETVAL(
        'missions_id_seq',
        (
            SELECT
                MAX(id)
            FROM
                missions
        )
    );


SELECT
    SETVAL(
        'mission_contents_id_seq',
        (
            SELECT
                MAX(id)
            FROM
                mission_contents
        )
    );


SELECT
    SETVAL(
        'rewards_id_seq',
        (
            SELECT
                MAX(id)
            FROM
                rewards
        )
    );


--rollback DELETE FROM rewards WHERE mission_id IN (1,2,3,4,5,6,7,8);
--rollback DELETE FROM mission_contents WHERE mission_id IN (1,2,3,4,5,6,7,8);
--rollback DELETE FROM missions WHERE chapter_id IN (1,2,3);
--rollback DELETE FROM chapters WHERE journey_id = 1;
--rollback DELETE FROM journeys WHERE id = 1;