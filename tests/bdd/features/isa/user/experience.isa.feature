# ISA Layer (L2): Backend API E2E Test
# Source: swagger.yaml - /users/me endpoint
# Test Scope: HTTP Request → API Handler → Database → HTTP Response
# Maps DSL scenarios to concrete API calls

@isa
Feature: 經驗值系統 API 測試

  Background:
    # Database is cleaned before each scenario (handled by @Before hook)

  Scenario: 新註冊使用者的初始經驗值應為 0
    # Setup: Create user with 0 experience
    Given the database has a user:
      | username   | Alice     |
      | password   | Test1234! |
      | experience | 0         |

    # Setup: Login to get access token
    And I set request body to:
      """
      {
        "username": "Alice",
        "password": "Test1234!"
      }
      """
    When I send "POST" request to "/auth/login"
    And I store the response field "accessToken" as "aliceToken"

    # Action: Get user profile
    Given I set Authorization header to "{{aliceToken}}"
    When I send "GET" request to "/users/me"

    # Verification: HTTP layer
    Then the response status code should be 200

    # Verification: Response structure
    And the response body should contain field "experiencePoints"

    # Verification: Response values
    And the response body field "experiencePoints" should equal number 0

  Scenario: 使用者交付任務後經驗值應增加
    # Setup: Create complete learning context
    Given the database has a user:
      | username   | Bob        |
      | password   | Secure123! |
      | experience | 0          |
    And the database has a journey:
      | title   | Java 基礎課程 |
      | slug    | java-basics  |
      | teacher | 水球老師      |
      | price   | 1999.00      |
    And the database has a chapter:
      | journey_id  | {{lastJourneyId}}  |
      | title       | 第一章:變數與型別 |
      | order_index | 1                  |
    And the database has a mission:
      | chapter_id   | {{lastChapterId}} |
      | title        | 認識變數           |
      | type         | VIDEO             |
      | access_level | PURCHASED         |
      | order_index  | 1                 |
    And the database has a mission resource:
      | mission_id       | {{lastMissionId}}           |
      | type             | VIDEO                       |
      | resource_url     | https://example.com/vid.mp4 |
      | content_order    | 0                           |
      | duration_seconds | 100                         |
    And the database has a reward:
      | mission_id   | {{lastMissionId}} |
      | reward_type  | EXPERIENCE        |
      | reward_value | 100               |
    And the database has an order:
      | user_id    | 1                 |
      | journey_id | {{lastJourneyId}} |
      | status     | PAID              |
    # Setup: User has completed the mission (100% watched)
    And the database has a user mission progress:
      | user_id                | 1                 |
      | mission_id             | {{lastMissionId}} |
      | status                 | COMPLETED         |
      | watch_position_seconds | 100               |

    # Setup: Login to get access token
    And I set request body to:
      """
      {
        "username": "Bob",
        "password": "Secure123!"
      }
      """
    When I send "POST" request to "/auth/login"
    And I store the response field "accessToken" as "bobToken"

    # Action: Deliver the mission to gain experience
    Given I set Authorization header to "{{bobToken}}"
    When I send "POST" request to "/users/1/missions/{{lastMissionId}}/progress/deliver"

    # Verification: Delivery successful
    Then the response status code should be 200
    And the response body should contain field "experienceGained"
    And the response body field "experienceGained" should equal number 100

    # Action: Get user profile to verify total experience
    Given I set Authorization header to "{{bobToken}}"
    When I send "GET" request to "/users/me"

    # Verification: HTTP layer
    Then the response status code should be 200

    # Verification: Response structure
    And the response body should contain field "experiencePoints"

    # Verification: Response values
    And the response body field "experiencePoints" should equal number 100

  Scenario: 未登入使用者無法查詢經驗值
    # Action: Attempt to get user profile without token
    When I send "GET" request to "/users/me"

    # Verification: HTTP layer
    Then the response status code should be 401

    # Verification: Error response
    And the response body should contain field "error"

  Scenario: 使用者完成任務但未交付時經驗值保持為 0
    # Setup: Create complete learning context (mission completed but not delivered)
    Given the database has a user:
      | username   | Diana     |
      | password   | Diana123! |
      | experience | 0         |
    And the database has a journey:
      | title   | JavaScript 入門  |
      | slug    | javascript-intro |
      | teacher | 水球老師          |
      | price   | 1999.00          |
    And the database has a chapter:
      | journey_id  | {{lastJourneyId}} |
      | title       | 第一章:基礎語法   |
      | order_index | 1                 |
    And the database has a mission:
      | chapter_id   | {{lastChapterId}} |
      | title        | 認識 JavaScript   |
      | type         | VIDEO             |
      | access_level | PURCHASED         |
      | order_index  | 1                 |
    And the database has a mission resource:
      | mission_id       | {{lastMissionId}}           |
      | type             | VIDEO                       |
      | resource_url     | https://example.com/vid.mp4 |
      | content_order    | 0                           |
      | duration_seconds | 120                         |
    And the database has a reward:
      | mission_id   | {{lastMissionId}} |
      | reward_type  | EXPERIENCE        |
      | reward_value | 100               |
    And the database has an order:
      | user_id    | 1                 |
      | journey_id | {{lastJourneyId}} |
      | status     | PAID              |
    And the database has a user mission progress:
      | user_id                | 1                 |
      | mission_id             | {{lastMissionId}} |
      | status                 | COMPLETED         |
      | watch_position_seconds | 120               |

    # Setup: Login to get access token
    And I set request body to:
      """
      {
        "username": "Diana",
        "password": "Diana123!"
      }
      """
    When I send "POST" request to "/auth/login"
    And I store the response field "accessToken" as "dianaToken"

    # Action: Get user profile to check experience
    Given I set Authorization header to "{{dianaToken}}"
    When I send "GET" request to "/users/me"

    # Verification: HTTP layer
    Then the response status code should be 200

    # Verification: Response structure
    And the response body should contain field "experiencePoints"

    # Verification: Response values - Should still be 0 because mission was not delivered
    And the response body field "experiencePoints" should equal number 0
