# ISA Layer (L2): Backend API E2E Test
# Source: swagger.yaml - /users/{userId}/missions/{missionId}/progress
# Test Scope: HTTP Request → API Handler → Database → HTTP Response
# Maps DSL scenarios to concrete API calls

@isa
Feature: Watch Video and Track Progress API Implementation

  Scenario: 首次觀看影片任務應從頭開始
    # Setup: Create journey, chapter, mission, and video resource
    Given the database has a journey:
      | title   | Java 基礎課程 |
      | slug    | java-basics  |
      | teacher | 水球老師      |
      | price   | 1999         |

    And the database has a chapter:
      | journey_id  | {{lastJourneyId}} |
      | title       | 第一章：變數與型別 |
      | order_index | 1                 |

    And the database has a mission:
      | chapter_id   | {{lastChapterId}} |
      | title        | 認識變數           |
      | type         | VIDEO             |
      | access_level | PURCHASED         |
      | order_index  | 1                 |

    And the database has a mission resource:
      | mission_id       | {{lastMissionId}}         |
      | type             | VIDEO                     |
      | resource_url     | https://example.com/v.mp4 |
      | content_order    | 0                         |
      | duration_seconds | 100                       |

    And the database has a user:
      | username   | Alice    |
      | password   | Test1234! |
      | experience | 0        |

    # Authentication: Login to get access token
    And I set request body to:
      """
      {
        "username": "Alice",
        "password": "Test1234!"
      }
      """
    When I send "POST" request to "/auth/login"
    And I store the response field "accessToken" as "accessToken"
    Given I set Authorization header to "{{accessToken}}"

    # Action: Get mission progress (no existing progress)
    When I send "GET" request to "/users/1/missions/{{lastMissionId}}/progress"

    # Verification: HTTP layer
    Then the response status code should be 200

    # Verification: Response structure
    And the response body should contain field "missionId"
    And the response body should contain field "status"
    And the response body should contain field "watchPositionSeconds"

    # Verification: Response values - should start from 0
    And the response body field "status" should equal string "UNCOMPLETED"
    And the response body field "watchPositionSeconds" should equal number 0

  Scenario: 繼續觀看影片應從上次位置開始
    # Setup: Create journey, chapter, mission, and video resource
    Given the database has a journey:
      | title   | Java 基礎課程 |
      | slug    | java-basics  |
      | teacher | 水球老師      |
      | price   | 1999         |

    And the database has a chapter:
      | journey_id  | {{lastJourneyId}} |
      | title       | 第一章：變數與型別 |
      | order_index | 1                 |

    And the database has a mission:
      | chapter_id   | {{lastChapterId}} |
      | title        | 認識變數           |
      | type         | VIDEO             |
      | access_level | PURCHASED         |
      | order_index  | 1                 |

    And the database has a mission resource:
      | mission_id       | {{lastMissionId}}         |
      | type             | VIDEO                     |
      | resource_url     | https://example.com/v.mp4 |
      | content_order    | 0                         |
      | duration_seconds | 100                       |

    And the database has a user:
      | username   | Bob      |
      | password   | Test1234! |
      | experience | 0        |

    # Setup: User has watched video to 30 seconds
    And the database has a user mission progress:
      | user_id                | 1  |
      | mission_id             | {{lastMissionId}} |
      | status                 | UNCOMPLETED |
      | watch_position_seconds | 30 |

    # Authentication: Login to get access token
    And I set request body to:
      """
      {
        "username": "Bob",
        "password": "Test1234!"
      }
      """
    When I send "POST" request to "/auth/login"
    And I store the response field "accessToken" as "accessToken"
    Given I set Authorization header to "{{accessToken}}"

    # Action: Get mission progress (should return previous position)
    When I send "GET" request to "/users/1/missions/{{lastMissionId}}/progress"

    # Verification: HTTP layer
    Then the response status code should be 200

    # Verification: Response values - should continue from 30 seconds
    And the response body field "status" should equal string "UNCOMPLETED"
    And the response body field "watchPositionSeconds" should equal number 30

  Scenario: 觀看完成後應自動標記為已完成
    # Setup: Create journey, chapter, mission, and video resource
    Given the database has a journey:
      | title   | Java 基礎課程 |
      | slug    | java-basics  |
      | teacher | 水球老師      |
      | price   | 1999         |

    And the database has a chapter:
      | journey_id  | {{lastJourneyId}} |
      | title       | 第一章：變數與型別 |
      | order_index | 1                 |

    And the database has a mission:
      | chapter_id   | {{lastChapterId}} |
      | title        | 認識變數           |
      | type         | VIDEO             |
      | access_level | PURCHASED         |
      | order_index  | 1                 |

    And the database has a mission resource:
      | mission_id       | {{lastMissionId}}         |
      | type             | VIDEO                     |
      | resource_url     | https://example.com/v.mp4 |
      | content_order    | 0                         |
      | duration_seconds | 100                       |

    And the database has a user:
      | username   | Charlie  |
      | password   | Test1234! |
      | experience | 0        |

    # Authentication: Login to get access token
    And I set request body to:
      """
      {
        "username": "Charlie",
        "password": "Test1234!"
      }
      """
    When I send "POST" request to "/auth/login"
    And I store the response field "accessToken" as "accessToken"
    Given I set Authorization header to "{{accessToken}}"

    # Setup: Prepare request body to update progress to 100 seconds
    And I set request body to:
      """
      {
        "watchPositionSeconds": 100
      }
      """

    # Action: Update progress to 100 seconds (video completion)
    When I send "PUT" request to "/users/1/missions/{{lastMissionId}}/progress"

    # Verification: HTTP layer
    Then the response status code should be 200

    # Verification: Response structure
    And the response body should contain field "missionId"
    And the response body should contain field "status"
    And the response body should contain field "watchPositionSeconds"

    # Verification: Response values - status should be COMPLETED
    And the response body field "status" should equal string "COMPLETED"
    And the response body field "watchPositionSeconds" should equal number 100

  Scenario: 播放中每 10 秒應記錄觀看進度
    # 前端行為：前端在播放期間每 10 秒呼叫 PUT API
    # 後端測試：驗證 PUT API 能正確更新進度到 10 秒

    # Setup: Create journey, chapter, mission, and video resource
    Given the database has a journey:
      | title   | Java 基礎課程 |
      | slug    | java-basics  |
      | teacher | 水球老師      |
      | price   | 1999         |

    And the database has a chapter:
      | journey_id  | {{lastJourneyId}} |
      | title       | 第一章：變數與型別 |
      | order_index | 1                 |

    And the database has a mission:
      | chapter_id   | {{lastChapterId}} |
      | title        | 認識變數           |
      | type         | VIDEO             |
      | access_level | PURCHASED         |
      | order_index  | 1                 |

    And the database has a mission resource:
      | mission_id       | {{lastMissionId}}         |
      | type             | VIDEO                     |
      | resource_url     | https://example.com/v.mp4 |
      | content_order    | 0                         |
      | duration_seconds | 100                       |

    And the database has a user:
      | username   | Diana    |
      | password   | Test1234! |
      | experience | 0        |

    # Authentication: Login to get access token
    And I set request body to:
      """
      {
        "username": "Diana",
        "password": "Test1234!"
      }
      """
    When I send "POST" request to "/auth/login"
    And I store the response field "accessToken" as "accessToken"
    Given I set Authorization header to "{{accessToken}}"

    # Setup: Prepare request body to update progress to 10 seconds
    And I set request body to:
      """
      {
        "watchPositionSeconds": 10
      }
      """

    # Action: Update progress to 10 seconds
    When I send "PUT" request to "/users/1/missions/{{lastMissionId}}/progress"

    # Verification: HTTP layer
    Then the response status code should be 200

    # Verification: Response values - progress saved at 10 seconds
    And the response body field "status" should equal string "UNCOMPLETED"
    And the response body field "watchPositionSeconds" should equal number 10

  Scenario: 暫停影片時應記錄觀看進度
    # 前端行為：前端在暫停時呼叫 PUT API
    # 後端測試：驗證 PUT API 能正確更新進度到 25 秒

    # Setup: Create journey, chapter, mission, and video resource
    Given the database has a journey:
      | title   | Java 基礎課程 |
      | slug    | java-basics  |
      | teacher | 水球老師      |
      | price   | 1999         |

    And the database has a chapter:
      | journey_id  | {{lastJourneyId}} |
      | title       | 第一章：變數與型別 |
      | order_index | 1                 |

    And the database has a mission:
      | chapter_id   | {{lastChapterId}} |
      | title        | 認識變數           |
      | type         | VIDEO             |
      | access_level | PURCHASED         |
      | order_index  | 1                 |

    And the database has a mission resource:
      | mission_id       | {{lastMissionId}}         |
      | type             | VIDEO                     |
      | resource_url     | https://example.com/v.mp4 |
      | content_order    | 0                         |
      | duration_seconds | 100                       |

    And the database has a user:
      | username   | Eve      |
      | password   | Test1234! |
      | experience | 0        |

    # Authentication: Login to get access token
    And I set request body to:
      """
      {
        "username": "Eve",
        "password": "Test1234!"
      }
      """
    When I send "POST" request to "/auth/login"
    And I store the response field "accessToken" as "accessToken"
    Given I set Authorization header to "{{accessToken}}"

    # Setup: Prepare request body to update progress to 25 seconds
    And I set request body to:
      """
      {
        "watchPositionSeconds": 25
      }
      """

    # Action: Update progress to 25 seconds (when user pauses)
    When I send "PUT" request to "/users/1/missions/{{lastMissionId}}/progress"

    # Verification: HTTP layer
    Then the response status code should be 200

    # Verification: Response values - progress saved at 25 seconds
    And the response body field "status" should equal string "UNCOMPLETED"
    And the response body field "watchPositionSeconds" should equal number 25

  Scenario: 關閉頁面時應記錄觀看進度
    # 前端行為：前端在關閉頁面前呼叫 PUT API
    # 後端測試：驗證 PUT API 能正確更新進度到 45 秒

    # Setup: Create journey, chapter, mission, and video resource
    Given the database has a journey:
      | title   | Java 基礎課程 |
      | slug    | java-basics  |
      | teacher | 水球老師      |
      | price   | 1999         |

    And the database has a chapter:
      | journey_id  | {{lastJourneyId}} |
      | title       | 第一章：變數與型別 |
      | order_index | 1                 |

    And the database has a mission:
      | chapter_id   | {{lastChapterId}} |
      | title        | 認識變數           |
      | type         | VIDEO             |
      | access_level | PURCHASED         |
      | order_index  | 1                 |

    And the database has a mission resource:
      | mission_id       | {{lastMissionId}}         |
      | type             | VIDEO                     |
      | resource_url     | https://example.com/v.mp4 |
      | content_order    | 0                         |
      | duration_seconds | 100                       |

    And the database has a user:
      | username   | Frank    |
      | password   | Test1234! |
      | experience | 0        |

    # Authentication: Login to get access token
    And I set request body to:
      """
      {
        "username": "Frank",
        "password": "Test1234!"
      }
      """
    When I send "POST" request to "/auth/login"
    And I store the response field "accessToken" as "accessToken"
    Given I set Authorization header to "{{accessToken}}"

    # Setup: Prepare request body to update progress to 45 seconds
    And I set request body to:
      """
      {
        "watchPositionSeconds": 45
      }
      """

    # Action: Update progress to 45 seconds (when user closes page)
    When I send "PUT" request to "/users/1/missions/{{lastMissionId}}/progress"

    # Verification: HTTP layer
    Then the response status code should be 200

    # Verification: Response values - progress saved at 45 seconds
    And the response body field "status" should equal string "UNCOMPLETED"
    And the response body field "watchPositionSeconds" should equal number 45

  Scenario: 已完成的影片任務可以重複觀看
    # Setup: Create journey, chapter, mission, and video resource
    Given the database has a journey:
      | title   | Java 基礎課程 |
      | slug    | java-basics  |
      | teacher | 水球老師      |
      | price   | 1999         |

    And the database has a chapter:
      | journey_id  | {{lastJourneyId}} |
      | title       | 第一章：變數與型別 |
      | order_index | 1                 |

    And the database has a mission:
      | chapter_id   | {{lastChapterId}} |
      | title        | 認識變數           |
      | type         | VIDEO             |
      | access_level | PURCHASED         |
      | order_index  | 1                 |

    And the database has a mission resource:
      | mission_id       | {{lastMissionId}}         |
      | type             | VIDEO                     |
      | resource_url     | https://example.com/v.mp4 |
      | content_order    | 0                         |
      | duration_seconds | 100                       |

    And the database has a user:
      | username   | Grace    |
      | password   | Test1234! |
      | experience | 0        |

    # Setup: User has completed the mission
    And the database has a user mission progress:
      | user_id                | 1         |
      | mission_id             | {{lastMissionId}} |
      | status                 | COMPLETED |
      | watch_position_seconds | 100       |

    # Authentication: Login to get access token
    And I set request body to:
      """
      {
        "username": "Grace",
        "password": "Test1234!"
      }
      """
    When I send "POST" request to "/auth/login"
    And I store the response field "accessToken" as "accessToken"
    Given I set Authorization header to "{{accessToken}}"

    # Setup: Prepare request body to re-watch from 50 seconds
    And I set request body to:
      """
      {
        "watchPositionSeconds": 50
      }
      """

    # Action: Update progress to 50 seconds (re-watching)
    When I send "PUT" request to "/users/1/missions/{{lastMissionId}}/progress"

    # Verification: HTTP layer
    Then the response status code should be 200

    # Verification: Response values - allows re-watching, status remains COMPLETED
    And the response body field "watchPositionSeconds" should equal number 50
    And the response body field "status" should equal string "COMPLETED"

  Scenario: 觀看進度精確達到 100% 時應標記為已完成
    # Setup: Create journey, chapter, mission, and video resource
    Given the database has a journey:
      | title   | Java 基礎課程 |
      | slug    | java-basics  |
      | teacher | 水球老師      |
      | price   | 1999         |

    And the database has a chapter:
      | journey_id  | {{lastJourneyId}} |
      | title       | 第一章：變數與型別 |
      | order_index | 1                 |

    And the database has a mission:
      | chapter_id   | {{lastChapterId}} |
      | title        | 認識變數           |
      | type         | VIDEO             |
      | access_level | PURCHASED         |
      | order_index  | 1                 |

    And the database has a mission resource:
      | mission_id       | {{lastMissionId}}         |
      | type             | VIDEO                     |
      | resource_url     | https://example.com/v.mp4 |
      | content_order    | 0                         |
      | duration_seconds | 120                       |

    And the database has a user:
      | username   | Henry    |
      | password   | Test1234! |
      | experience | 0        |

    # Authentication: Login to get access token
    And I set request body to:
      """
      {
        "username": "Henry",
        "password": "Test1234!"
      }
      """
    When I send "POST" request to "/auth/login"
    And I store the response field "accessToken" as "accessToken"
    Given I set Authorization header to "{{accessToken}}"

    # Setup: Prepare request body to update progress to 120 seconds (100%)
    And I set request body to:
      """
      {
        "watchPositionSeconds": 120
      }
      """

    # Action: Update progress to 120 seconds (exactly 100%)
    When I send "PUT" request to "/users/1/missions/{{lastMissionId}}/progress"

    # Verification: HTTP layer
    Then the response status code should be 200

    # Verification: Response values - status should be COMPLETED at 100%
    And the response body field "status" should equal string "COMPLETED"
    And the response body field "watchPositionSeconds" should equal number 120

  Scenario: 觀看進度未達 100% 時應保持未完成狀態
    # Setup: Create journey, chapter, mission, and video resource
    Given the database has a journey:
      | title   | Java 基礎課程 |
      | slug    | java-basics  |
      | teacher | 水球老師      |
      | price   | 1999         |

    And the database has a chapter:
      | journey_id  | {{lastJourneyId}} |
      | title       | 第一章：變數與型別 |
      | order_index | 1                 |

    And the database has a mission:
      | chapter_id   | {{lastChapterId}} |
      | title        | 認識變數           |
      | type         | VIDEO             |
      | access_level | PURCHASED         |
      | order_index  | 1                 |

    And the database has a mission resource:
      | mission_id       | {{lastMissionId}}         |
      | type             | VIDEO                     |
      | resource_url     | https://example.com/v.mp4 |
      | content_order    | 0                         |
      | duration_seconds | 120                       |

    And the database has a user:
      | username   | Iris     |
      | password   | Test1234! |
      | experience | 0        |

    # Authentication: Login to get access token
    And I set request body to:
      """
      {
        "username": "Iris",
        "password": "Test1234!"
      }
      """
    When I send "POST" request to "/auth/login"
    And I store the response field "accessToken" as "accessToken"
    Given I set Authorization header to "{{accessToken}}"

    # Setup: Prepare request body to update progress to 119 seconds (not 100%)
    And I set request body to:
      """
      {
        "watchPositionSeconds": 119
      }
      """

    # Action: Update progress to 119 seconds (not 100% yet)
    When I send "PUT" request to "/users/1/missions/{{lastMissionId}}/progress"

    # Verification: HTTP layer
    Then the response status code should be 200

    # Verification: Response values - status should remain UNCOMPLETED
    And the response body field "status" should equal string "UNCOMPLETED"
    And the response body field "watchPositionSeconds" should equal number 119
