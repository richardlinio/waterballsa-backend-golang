# ISA Layer (L2): Backend API E2E Test
# Source: swagger.yaml - /journeys/{journeyId}
# Test Scope: HTTP Request → API Handler → Database → HTTP Response
# Maps DSL scenarios to concrete API calls

@isa
Feature: Journey Detail API Implementation

  Scenario: Get journey detail with chapters and missions
    # Setup: Create journey with chapters and missions
    Given the database has a journey:
      | title            | 設計模式精通之旅                    |
      | slug             | design-patterns-mastery            |
      | description      | 用 C.A. 模式大大提昇系統思維能力      |
      | teacher          | 水球潘                             |
      | price            | 1999                               |
      | cover_image_url  | https://example.com/cover1.jpg     |

    And the database has a chapter:
      | journey_id   | {{lastJourneyId}} |
      | title        | 課程介紹           |
      | order_index  | 1                 |

    And the database has a mission:
      | chapter_id   | {{lastChapterId}}                      |
      | title        | 這門課手把手帶你成為架構設計的高手        |
      | type         | VIDEO                                  |
      | access_level | PUBLIC                                 |
      | order_index  | 1                                      |

    And the database has a mission:
      | chapter_id   | {{lastChapterId}}                                                  |
      | title        | 你該知道：在 AI 的時代下，只會下 prompt 絕對寫不出好 Code              |
      | type         | VIDEO                                                              |
      | access_level | PUBLIC                                                             |
      | order_index  | 2                                                                  |

    And the database has a chapter:
      | journey_id   | {{lastJourneyId}}      |
      | title        | 架構思維的 C.A. 模式    |
      | order_index  | 2                      |

    And the database has a mission:
      | chapter_id   | {{lastChapterId}} |
      | title        | 架構師該如何思考   |
      | type         | VIDEO             |
      | access_level | PURCHASED         |
      | order_index  | 1                 |

    # Action: Get journey detail
    When I send "GET" request to "/journeys/{{lastJourneyId}}"

    # Verification: HTTP layer
    Then the response status code should be 200

    # Verification: Journey basic info
    And the response body should contain field "id"
    And the response body should contain field "title"
    And the response body field "title" should equal string "設計模式精通之旅"
    And the response body field "description" should equal string "用 C.A. 模式大大提昇系統思維能力"
    And the response body field "teacherName" should equal string "水球潘"
    And the response body field "price" should equal decimal "1999.0"
    And the response body field "coverImageUrl" should equal string "https://example.com/cover1.jpg"

    # Verification: Chapters
    And the response body should contain field "chapters"
    And the response body field "chapters" should have size 2

    # Verification: First chapter missions
    And the response body field "chapters[0].title" should equal string "課程介紹"
    And the response body field "chapters[0].missions" should have size 2

    # Verification: Second chapter missions
    And the response body field "chapters[1].title" should equal string "架構思維的 C.A. 模式"
    And the response body field "chapters[1].missions" should have size 1

    # Verification: Mission access levels
    And the response body field "chapters[0].missions[0].accessLevel" should equal string "PUBLIC"
    And the response body field "chapters[1].missions[0].accessLevel" should equal string "PURCHASED"

  Scenario: Get non-existent journey returns 404
    # No setup needed - journey doesn't exist

    # Action: Get non-existent journey
    When I send "GET" request to "/journeys/999"

    # Verification: HTTP layer
    Then the response status code should be 404

    # Verification: Error response
    And the response body should contain field "error"
    And the response body field "error" should equal string "查無此旅程"
