# ISA Layer (L2): Backend API E2E Test
# Source: swagger.yaml - /users/{userId}/missions/{missionId}/progress/deliver
# Test Scope: HTTP Request → API Handler → Database → HTTP Response
# Maps DSL scenarios to concrete API calls

@isa
Feature: Mission Delivery API Implementation

  Scenario: 影片任務完成後成功交付並獲得經驗值
    # Setup: Create test data
    Given the database has a user:
      | username   | Alice     |
      | password   | Test1234! |
      | experience | 0         |

    And the database has a journey:
      | title   | Java 基礎課程 |
      | slug    | java-basics  |
      | teacher | 水球老師      |
      | price   | 1999         |

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
      | mission_id       | {{lastMissionId}}         |
      | type             | VIDEO                     |
      | resource_url     | https://example.com/v.mp4 |
      | content_order    | 0                         |
      | duration_seconds | 100                       |

    And the database has a reward:
      | mission_id   | {{lastMissionId}} |
      | reward_type  | EXPERIENCE        |
      | reward_value | 100               |

    # Setup: User has purchased the journey
    And the database has an order:
      | user_id    | 1                 |
      | journey_id | {{lastJourneyId}} |
      | status     | PAID              |

    # Setup: User has completed the mission
    And the database has a user mission progress:
      | user_id                | 1                 |
      | mission_id             | {{lastMissionId}} |
      | status                 | COMPLETED         |
      | watch_position_seconds | 100               |

    # Setup: Login to get access token
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

    # Action: Deliver the mission
    When I send "POST" request to "/users/1/missions/{{lastMissionId}}/progress/deliver"

    # Verification: HTTP layer
    Then the response status code should be 200

    # Verification: Response structure
    And the response body should contain field "message"
    And the response body should contain field "experienceGained"
    And the response body should contain field "totalExperience"
    And the response body should contain field "currentLevel"

    # Verification: Response values
    And the response body field "experienceGained" should equal number 100
    And the response body field "totalExperience" should equal number 100
    And the response body field "currentLevel" should equal number 1

  Scenario: 嘗試交付未完成的任務應失敗
    # Setup: Create test data
    Given the database has a user:
      | username   | Bob       |
      | password   | Secure123! |
      | experience | 0         |

    And the database has a journey:
      | title   | Python 入門  |
      | slug    | python-intro |
      | teacher | 水球老師      |
      | price   | 1999         |

    And the database has a chapter:
      | journey_id  | {{lastJourneyId}} |
      | title       | 第一章:基礎語法   |
      | order_index | 1                 |

    And the database has a mission:
      | chapter_id   | {{lastChapterId}} |
      | title        | Hello World       |
      | type         | VIDEO             |
      | access_level | PURCHASED         |
      | order_index  | 1                 |

    And the database has a mission resource:
      | mission_id       | {{lastMissionId}}         |
      | type             | VIDEO                     |
      | resource_url     | https://example.com/v.mp4 |
      | content_order    | 0                         |
      | duration_seconds | 100                       |

    And the database has a reward:
      | mission_id   | {{lastMissionId}} |
      | reward_type  | EXPERIENCE        |
      | reward_value | 100               |

    # Setup: User has purchased the journey
    And the database has an order:
      | user_id    | 1                 |
      | journey_id | {{lastJourneyId}} |
      | status     | PAID              |

    # Setup: User has only watched 50% of the video (UNCOMPLETED)
    And the database has a user mission progress:
      | user_id                | 1                 |
      | mission_id             | {{lastMissionId}} |
      | status                 | UNCOMPLETED       |
      | watch_position_seconds | 50                |

    # Setup: Login to get access token
    And I set request body to:
      """
      {
        "username": "Bob",
        "password": "Secure123!"
      }
      """
    When I send "POST" request to "/auth/login"
    And I store the response field "accessToken" as "accessToken"
    Given I set Authorization header to "{{accessToken}}"

    # Action: Attempt to deliver incomplete mission
    When I send "POST" request to "/users/1/missions/{{lastMissionId}}/progress/deliver"

    # Verification: HTTP layer
    Then the response status code should be 400

    # Verification: Error response
    And the response body should contain field "error"

  Scenario: 嘗試重複交付同一任務應失敗
    # Setup: Create test data
    Given the database has a user:
      | username   | Charlie   |
      | password   | Pass1234! |
      | experience | 100       |

    And the database has a journey:
      | title   | JavaScript 完整指南 |
      | slug    | js-complete        |
      | teacher | 水球老師            |
      | price   | 1999               |

    And the database has a chapter:
      | journey_id  | {{lastJourneyId}}       |
      | title       | 第一章:JavaScript 基礎 |
      | order_index | 1                       |

    And the database has a mission:
      | chapter_id   | {{lastChapterId}} |
      | title        | 變數宣告           |
      | type         | VIDEO             |
      | access_level | PURCHASED         |
      | order_index  | 1                 |

    And the database has a mission resource:
      | mission_id       | {{lastMissionId}}         |
      | type             | VIDEO                     |
      | resource_url     | https://example.com/v.mp4 |
      | content_order    | 0                         |
      | duration_seconds | 120                       |

    And the database has a reward:
      | mission_id   | {{lastMissionId}} |
      | reward_type  | EXPERIENCE        |
      | reward_value | 100               |

    # Setup: User has purchased the journey
    And the database has an order:
      | user_id    | 1                 |
      | journey_id | {{lastJourneyId}} |
      | status     | PAID              |

    # Setup: User has already delivered the mission
    And the database has a user mission progress:
      | user_id                | 1                 |
      | mission_id             | {{lastMissionId}} |
      | status                 | DELIVERED         |
      | watch_position_seconds | 120               |

    # Setup: Login to get access token
    And I set request body to:
      """
      {
        "username": "Charlie",
        "password": "Pass1234!"
      }
      """
    When I send "POST" request to "/auth/login"
    And I store the response field "accessToken" as "accessToken"
    Given I set Authorization header to "{{accessToken}}"

    # Action: Attempt to deliver already delivered mission
    When I send "POST" request to "/users/1/missions/{{lastMissionId}}/progress/deliver"

    # Verification: HTTP layer
    Then the response status code should be 409

    # Verification: Error response
    And the response body should contain field "error"

  Scenario: 未登入使用者無法交付任務
    # Setup: Create test data (no user created, no login)
    Given the database has a journey:
      | title   | React 完整課程 |
      | slug    | react-complete |
      | teacher | 水球老師        |
      | price   | 1999           |

    And the database has a chapter:
      | journey_id  | {{lastJourneyId}} |
      | title       | 第一章:React 基礎 |
      | order_index | 1                 |

    And the database has a mission:
      | chapter_id   | {{lastChapterId}} |
      | title        | 認識 React        |
      | type         | VIDEO             |
      | access_level | PURCHASED         |
      | order_index  | 1                 |

    And the database has a mission resource:
      | mission_id       | {{lastMissionId}}         |
      | type             | VIDEO                     |
      | resource_url     | https://example.com/v.mp4 |
      | content_order    | 0                         |
      | duration_seconds | 150                       |

    And the database has a reward:
      | mission_id   | {{lastMissionId}} |
      | reward_type  | EXPERIENCE        |
      | reward_value | 100               |

    # Action: Attempt to deliver without authentication
    When I send "POST" request to "/users/1/missions/{{lastMissionId}}/progress/deliver"

    # Verification: HTTP layer
    Then the response status code should be 401

    # Verification: Error response
    And the response body should contain field "error"
