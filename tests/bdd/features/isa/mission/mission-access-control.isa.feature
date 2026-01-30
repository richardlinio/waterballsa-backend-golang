# ISA Layer (L2): Backend API E2E Test
# Source: swagger.yaml - /journeys/{journeyId}, /journeys/{journeyId}/missions/{missionId}
# Test Scope: HTTP Request → API Handler → Database → HTTP Response
# Maps DSL scenarios to concrete API calls

@isa
Feature: Mission Access Control API Implementation

  Scenario: 已購買使用者查看旅程詳情應看到所有任務解鎖
    # Setup: Create journey with chapter and 3 PURCHASED missions
    Given the database has a journey:
      | title           | 設計模式精通之旅                 |
      | slug            | design-patterns-mastery         |
      | description     | 用 C.A. 模式大大提昇系統思維能力 |
      | teacher         | 水球潘                          |
      | price           | 3990                            |
      | cover_image_url | https://example.com/cover.jpg   |

    And the database has a chapter:
      | journey_id  | {{lastJourneyId}} |
      | title       | 進階課程           |
      | order_index | 1                 |

    And the database has a mission:
      | chapter_id   | {{lastChapterId}} |
      | title        | SOLID 原則        |
      | type         | VIDEO             |
      | access_level | PURCHASED         |
      | order_index  | 1                 |

    And the database has a mission:
      | chapter_id   | {{lastChapterId}} |
      | title        | 設計模式導論       |
      | type         | VIDEO             |
      | access_level | PURCHASED         |
      | order_index  | 2                 |

    And the database has a mission:
      | chapter_id   | {{lastChapterId}} |
      | title        | 工廠模式          |
      | type         | VIDEO             |
      | access_level | PURCHASED         |
      | order_index  | 3                 |

    # Setup: Create user Alice and purchase the journey
    And the database has a user:
      | username   | Alice     |
      | password   | Test1234! |
      | experience | 0         |

    And the database has an order:
      | user_id    | {{lastUserId}}    |
      | journey_id | {{lastJourneyId}} |
      | status     | PAID              |

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

    # Action: Get journey detail with authentication
    When I send "GET" request to "/journeys/{{lastJourneyId}}"

    # Verification: HTTP layer
    Then the response status code should be 200

    # Verification: Journey basic info
    And the response body should contain field "id"
    And the response body should contain field "title"
    And the response body field "title" should equal string "設計模式精通之旅"

    # Verification: Chapters and missions structure
    And the response body should contain field "chapters"
    And the response body field "chapters" should have size 1
    And the response body field "chapters[0].title" should equal string "進階課程"
    And the response body field "chapters[0].missions" should have size 3

    # Verification: All missions have PURCHASED access level
    And the response body field "chapters[0].missions[0].title" should equal string "SOLID 原則"
    And the response body field "chapters[0].missions[0].accessLevel" should equal string "PURCHASED"
    And the response body field "chapters[0].missions[1].title" should equal string "設計模式導論"
    And the response body field "chapters[0].missions[1].accessLevel" should equal string "PURCHASED"
    And the response body field "chapters[0].missions[2].title" should equal string "工廠模式"
    And the response body field "chapters[0].missions[2].accessLevel" should equal string "PURCHASED"

  Scenario: 已購買使用者可以成功查看需購買任務詳情
    # Setup: Create journey with chapter and PURCHASED mission with resource
    Given the database has a journey:
      | title           | Java 基礎課程                   |
      | slug            | java-basics                    |
      | description     | 學習 Java 程式設計基礎          |
      | teacher         | 水球老師                        |
      | price           | 1999                           |
      | cover_image_url | https://example.com/java.jpg   |

    And the database has a chapter:
      | journey_id  | {{lastJourneyId}}   |
      | title       | 物件導向程式設計     |
      | order_index | 1                   |

    And the database has a mission:
      | chapter_id   | {{lastChapterId}}     |
      | title        | 類別與物件            |
      | type         | VIDEO                 |
      | description  | 學習類別與物件的概念   |
      | access_level | PURCHASED             |
      | order_index  | 1                     |

    And the database has a mission resource:
      | mission_id       | {{lastMissionId}}                   |
      | type             | VIDEO                               |
      | resource_url     | https://example.com/class-obj.m3u8  |
      | content_order    | 0                                   |
      | duration_seconds | 300                                 |

    And the database has a reward:
      | mission_id   | {{lastMissionId}} |
      | reward_type  | EXPERIENCE        |
      | reward_value | 100               |

    # Setup: Create user Bob and purchase the journey
    And the database has a user:
      | username   | Bob       |
      | password   | Secure123! |
      | experience | 0         |

    And the database has an order:
      | user_id    | {{lastUserId}}    |
      | journey_id | {{lastJourneyId}} |
      | status     | PAID              |

    # Authentication: Login to get access token
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

    # Action: Get mission detail with authentication
    When I send "GET" request to "/journeys/{{lastJourneyId}}/missions/{{lastMissionId}}"

    # Verification: HTTP layer
    Then the response status code should be 200

    # Verification: Mission basic info
    And the response body should contain field "id"
    And the response body should contain field "title"
    And the response body field "title" should equal string "類別與物件"
    And the response body field "description" should equal string "學習類別與物件的概念"
    And the response body field "type" should equal string "VIDEO"
    And the response body field "accessLevel" should equal string "PURCHASED"

    # Verification: Resources are accessible
    And the response body should contain field "resource"
    And the response body field "resource" should have size 1
    And the response body field "resource[0].type" should equal string "video"
    And the response body field "resource[0].resourceUrl" should equal string "https://example.com/class-obj.m3u8"
    And the response body field "resource[0].durationSeconds" should equal number 300

    # Verification: Reward info
    And the response body should contain field "reward"
    And the response body should contain field "reward.exp"
    And the response body field "reward.exp" should equal number 100

  Scenario: 未購買使用者查看旅程詳情應看到任務標記為需購買
    # Setup: Create journey with chapter and 3 PURCHASED missions
    Given the database has a journey:
      | title           | Python 入門                      |
      | slug            | python-intro                    |
      | description     | 從零開始學習 Python              |
      | teacher         | 水球老師                         |
      | price           | 2490                            |
      | cover_image_url | https://example.com/python.jpg  |

    And the database has a chapter:
      | journey_id  | {{lastJourneyId}} |
      | title       | 進階主題           |
      | order_index | 1                 |

    And the database has a mission:
      | chapter_id   | {{lastChapterId}} |
      | title        | 裝飾器            |
      | type         | VIDEO             |
      | access_level | PURCHASED         |
      | order_index  | 1                 |

    And the database has a mission:
      | chapter_id   | {{lastChapterId}} |
      | title        | 生成器            |
      | type         | VIDEO             |
      | access_level | PURCHASED         |
      | order_index  | 2                 |

    And the database has a mission:
      | chapter_id   | {{lastChapterId}} |
      | title        | 上下文管理器       |
      | type         | VIDEO             |
      | access_level | PURCHASED         |
      | order_index  | 3                 |

    # Setup: Create user Charlie but DO NOT purchase the journey
    And the database has a user:
      | username   | Charlie   |
      | password   | Pass1234! |
      | experience | 0         |

    # Authentication: Login to get access token
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

    # Action: Get journey detail with authentication (but not purchased)
    When I send "GET" request to "/journeys/{{lastJourneyId}}"

    # Verification: HTTP layer
    Then the response status code should be 200

    # Verification: Journey basic info
    And the response body should contain field "id"
    And the response body should contain field "title"
    And the response body field "title" should equal string "Python 入門"

    # Verification: Chapters and missions structure
    And the response body should contain field "chapters"
    And the response body field "chapters" should have size 1
    And the response body field "chapters[0].title" should equal string "進階主題"
    And the response body field "chapters[0].missions" should have size 3

    # Verification: All missions show PURCHASED access level
    # Frontend will render lock icons based on this access level
    And the response body field "chapters[0].missions[0].title" should equal string "裝飾器"
    And the response body field "chapters[0].missions[0].accessLevel" should equal string "PURCHASED"
    And the response body field "chapters[0].missions[1].title" should equal string "生成器"
    And the response body field "chapters[0].missions[1].accessLevel" should equal string "PURCHASED"
    And the response body field "chapters[0].missions[2].title" should equal string "上下文管理器"
    And the response body field "chapters[0].missions[2].accessLevel" should equal string "PURCHASED"

  Scenario: 未購買使用者嘗試查看需購買任務詳情返回 403
    # Setup: Create journey with chapter and PURCHASED mission
    Given the database has a journey:
      | title           | JavaScript 完整指南              |
      | slug            | javascript-guide                |
      | description     | 深入學習 JavaScript             |
      | teacher         | 水球老師                         |
      | price           | 2990                            |
      | cover_image_url | https://example.com/js.jpg      |

    And the database has a chapter:
      | journey_id  | {{lastJourneyId}} |
      | title       | 非同步程式設計     |
      | order_index | 1                 |

    And the database has a mission:
      | chapter_id   | {{lastChapterId}}           |
      | title        | Promise 詳解                |
      | type         | VIDEO                       |
      | description  | 深入理解 JavaScript Promise |
      | access_level | PURCHASED                   |
      | order_index  | 1                           |

    # Setup: Create user Diana but DO NOT purchase the journey
    And the database has a user:
      | username   | Diana     |
      | password   | Diana123! |
      | experience | 0         |

    # Authentication: Login to get access token
    And I set request body to:
      """
      {
        "username": "Diana",
        "password": "Diana123!"
      }
      """
    When I send "POST" request to "/auth/login"
    And I store the response field "accessToken" as "accessToken"
    Given I set Authorization header to "{{accessToken}}"

    # Action: Attempt to get PURCHASED mission detail without purchasing
    When I send "GET" request to "/journeys/{{lastJourneyId}}/missions/{{lastMissionId}}"

    # Verification: HTTP layer
    Then the response status code should be 403

    # Verification: Error response
    And the response body should contain field "error"
    And the response body field "error" should equal string "禁止訪問"

  Scenario: 訪客查看旅程詳情應看到任務標記為需購買
    # Setup: Create journey with chapter and 2 PURCHASED missions
    Given the database has a journey:
      | title           | React 完整課程                   |
      | slug            | react-complete                  |
      | description     | 從基礎到進階的 React 課程        |
      | teacher         | 水球老師                         |
      | price           | 3490                            |
      | cover_image_url | https://example.com/react.jpg   |

    And the database has a chapter:
      | journey_id  | {{lastJourneyId}} |
      | title       | React Hooks       |
      | order_index | 1                 |

    And the database has a mission:
      | chapter_id   | {{lastChapterId}} |
      | title        | useState 詳解     |
      | type         | VIDEO             |
      | access_level | PURCHASED         |
      | order_index  | 1                 |

    And the database has a mission:
      | chapter_id   | {{lastChapterId}} |
      | title        | useEffect 詳解    |
      | type         | VIDEO             |
      | access_level | PURCHASED         |
      | order_index  | 2                 |

    # No user created - testing guest access

    # Action: Get journey detail without authentication (guest)
    When I send "GET" request to "/journeys/{{lastJourneyId}}"

    # Verification: HTTP layer
    Then the response status code should be 200

    # Verification: Journey basic info
    And the response body should contain field "id"
    And the response body should contain field "title"
    And the response body field "title" should equal string "React 完整課程"

    # Verification: Chapters and missions structure
    And the response body should contain field "chapters"
    And the response body field "chapters" should have size 1
    And the response body field "chapters[0].title" should equal string "React Hooks"
    And the response body field "chapters[0].missions" should have size 2

    # Verification: All missions show PURCHASED access level
    # Frontend will render lock icons based on this access level
    And the response body field "chapters[0].missions[0].title" should equal string "useState 詳解"
    And the response body field "chapters[0].missions[0].accessLevel" should equal string "PURCHASED"
    And the response body field "chapters[0].missions[1].title" should equal string "useEffect 詳解"
    And the response body field "chapters[0].missions[1].accessLevel" should equal string "PURCHASED"

    # Verification: Guest users see null status (not authenticated)
    And the response body field "chapters[0].missions[0].status" should be null
    And the response body field "chapters[0].missions[1].status" should be null

  Scenario: 訪客嘗試查看需購買任務詳情返回 401
    # Setup: Create journey with chapter and PURCHASED mission
    Given the database has a journey:
      | title           | Vue.js 實戰課程                  |
      | slug            | vuejs-practical                 |
      | description     | Vue.js 從入門到實戰              |
      | teacher         | 水球老師                         |
      | price           | 2790                            |
      | cover_image_url | https://example.com/vue.jpg     |

    And the database has a chapter:
      | journey_id  | {{lastJourneyId}} |
      | title       | Vuex 狀態管理     |
      | order_index | 1                 |

    And the database has a mission:
      | chapter_id   | {{lastChapterId}}         |
      | title        | Vuex 核心概念             |
      | type         | VIDEO                     |
      | description  | 理解 Vuex 的核心概念      |
      | access_level | PURCHASED                 |
      | order_index  | 1                         |

    # No user created - testing guest access

    # Action: Attempt to get PURCHASED mission detail without authentication
    When I send "GET" request to "/journeys/{{lastJourneyId}}/missions/{{lastMissionId}}"

    # Verification: HTTP layer
    Then the response status code should be 401

    # Verification: Error response
    And the response body should contain field "error"
    And the response body field "error" should equal string "未授權或權杖無效"
