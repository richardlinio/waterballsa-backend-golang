# ISA Layer (L2): Backend API E2E Test
# Source: swagger.yaml - /journeys
# Test Scope: HTTP Request → API Handler → Database → HTTP Response
# Maps DSL scenarios to concrete API calls

@isa
Feature: Journey List API Implementation

  Scenario: Get list with multiple journeys
    # Setup: Create 2 journeys with full details
    Given the database has a journey:
      | title            | 軟體設計模式精通之旅                |
      | slug             | design-patterns-mastery            |
      | description      | 用 C.A. 模式大大提昇系統思維能力      |
      | teacher          | 水球潘                             |
      | price            | 1999                               |
      | cover_image_url  | https://example.com/cover1.jpg     |
    And the database has a journey:
      | title            | Web 開發入門                        |
      | slug             | web-development-intro              |
      | description      | 從零開始學習網頁開發                  |
      | teacher          | 林老師                             |
      | price            | 999                                |
      | cover_image_url  | https://example.com/cover2.jpg     |

    # Action: Get journey list
    When I send "GET" request to "/journeys"

    # Verification: HTTP layer
    Then the response status code should be 200

    # Verification: Response structure
    And the response body should contain field "journeys"
    And the response body field "journeys" should have size 2

    # Verification: First journey values
    And the response body field "journeys[0].title" should equal string "軟體設計模式精通之旅"
    And the response body field "journeys[0].teacherName" should equal string "水球潘"
    And the response body field "journeys[0].price" should equal decimal "1999.0"
    And the response body field "journeys[0].description" should equal string "用 C.A. 模式大大提昇系統思維能力"
    And the response body field "journeys[0].coverImageUrl" should equal string "https://example.com/cover1.jpg"

    # Verification: Second journey values
    And the response body field "journeys[1].title" should equal string "Web 開發入門"
    And the response body field "journeys[1].teacherName" should equal string "林老師"
    And the response body field "journeys[1].price" should equal decimal "999.0"

  Scenario: Get empty list when no journeys exist
    # No setup needed - database is clean

    # Action: Get journey list
    When I send "GET" request to "/journeys"

    # Verification: HTTP layer
    Then the response status code should be 200

    # Verification: Response structure
    And the response body should contain field "journeys"
    And the response body field "journeys" should have size 0

  Scenario: Get list with single journey
    # Setup: Create 1 journey
    Given the database has a journey:
      | title            | Java 基礎課程                       |
      | slug             | java-basics                        |
      | description      | 深入淺出學習 Java 程式設計           |
      | teacher          | 陳老師                             |
      | price            | 1500                               |
      | cover_image_url  | https://example.com/java-cover.jpg |

    # Action: Get journey list
    When I send "GET" request to "/journeys"

    # Verification: HTTP layer
    Then the response status code should be 200

    # Verification: Response structure
    And the response body should contain field "journeys"
    And the response body field "journeys" should have size 1

    # Verification: Journey values
    And the response body field "journeys[0].title" should equal string "Java 基礎課程"
    And the response body field "journeys[0].teacherName" should equal string "陳老師"
    And the response body field "journeys[0].price" should equal decimal "1500.0"
    And the response body field "journeys[0].description" should equal string "深入淺出學習 Java 程式設計"
    And the response body field "journeys[0].coverImageUrl" should equal string "https://example.com/java-cover.jpg"
