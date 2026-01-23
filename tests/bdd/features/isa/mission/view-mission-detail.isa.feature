# ISA Layer (L2): Backend API E2E Test
# Source: swagger.yaml - /journeys/{journeyId}/missions/{missionId}
# Test Scope: HTTP Request → API Handler → Database → HTTP Response
# Maps DSL scenarios to concrete API calls

@isa
Feature: Mission Detail API Implementation

  Scenario: Guest views public mission detail
    # Setup: Create journey, chapter, and public mission with resource
    Given the database has a journey:
      | title   | 設計模式精通之旅    |
      | slug    | design-patterns    |
      | teacher | 水球潘             |
      | price   | 1999               |

    And the database has a chapter:
      | journey_id   | {{lastJourneyId}} |
      | title        | 課程介紹           |
      | order_index  | 1                 |

    And the database has a mission:
      | chapter_id   | {{lastChapterId}}                                  |
      | title        | 這門課手把手帶你成為架構設計的高手                    |
      | type         | VIDEO                                              |
      | description  | 思維升級：Christopher Alexander 模式六大要素及模式語言 |
      | access_level | PUBLIC                                             |
      | order_index  | 1                                                  |

    And the database has a reward:
      | mission_id   | {{lastMissionId}} |
      | reward_type  | EXPERIENCE        |
      | reward_value | 100               |

    And the database has a mission resource:
      | mission_id       | {{lastMissionId}}                     |
      | type             | VIDEO                                 |
      | resource_url     | https://example.com/video.m3u8        |
      | content_order    | 0                                     |
      | duration_seconds | 256                                   |

    # Action: Get mission detail (no authentication)
    When I send "GET" request to "/journeys/{{lastJourneyId}}/missions/{{lastMissionId}}"

    # Verification: HTTP layer
    Then the response status code should be 200

    # Verification: Mission basic info
    And the response body should contain field "id"
    And the response body should contain field "title"
    And the response body field "title" should equal string "這門課手把手帶你成為架構設計的高手"
    And the response body field "description" should equal string "思維升級：Christopher Alexander 模式六大要素及模式語言"
    And the response body field "type" should equal string "VIDEO"
    And the response body field "accessLevel" should equal string "PUBLIC"

    # Verification: Reward
    And the response body should contain field "reward"
    And the response body should contain field "reward.exp"
    And the response body field "reward.exp" should equal number 100

    # Verification: Resources
    And the response body should contain field "resource"
    And the response body field "resource" should have size 1
    And the response body field "resource[0].type" should equal string "video"

  Scenario: Guest attempts to view authenticated mission returns 401
    # Setup: Create authenticated mission
    Given the database has a journey:
      | title   | 設計模式精通之旅 |
      | slug    | design-patterns |
      | teacher | 水球潘          |
      | price   | 1999            |

    And the database has a chapter:
      | journey_id   | {{lastJourneyId}} |
      | title        | 進階課程           |
      | order_index  | 1                 |

    And the database has a mission:
      | chapter_id   | {{lastChapterId}} |
      | title        | 進階架構設計       |
      | type         | VIDEO             |
      | access_level | AUTHENTICATED     |
      | order_index  | 1                 |

    # Action: Attempt to get authenticated mission without login
    When I send "GET" request to "/journeys/{{lastJourneyId}}/missions/{{lastMissionId}}"

    # Verification: HTTP layer
    Then the response status code should be 401

    # Verification: Error response
    And the response body should contain field "error"
    And the response body field "error" should equal string "登入資料已過期"

  Scenario: Guest attempts to view purchased mission returns 401
    # Setup: Create purchased mission
    Given the database has a journey:
      | title   | 設計模式精通之旅 |
      | slug    | design-patterns |
      | teacher | 水球潘          |
      | price   | 1999            |

    And the database has a chapter:
      | journey_id   | {{lastJourneyId}} |
      | title        | 付費課程           |
      | order_index  | 1                 |

    And the database has a mission:
      | chapter_id   | {{lastChapterId}} |
      | title        | 專業架構設計       |
      | type         | VIDEO             |
      | access_level | PURCHASED         |
      | order_index  | 1                 |

    # Action: Attempt to get purchased mission without login
    When I send "GET" request to "/journeys/{{lastJourneyId}}/missions/{{lastMissionId}}"

    # Verification: HTTP layer
    Then the response status code should be 401

    # Verification: Error response
    And the response body should contain field "error"
    And the response body field "error" should equal string "登入資料已過期"

  Scenario: Get non-existent mission returns 404
    # Setup: Create journey but no mission with ID 999
    Given the database has a journey:
      | title   | 設計模式精通之旅 |
      | slug    | design-patterns |
      | teacher | 水球潘          |
      | price   | 1999            |

    # Action: Attempt to get non-existent mission
    When I send "GET" request to "/journeys/{{lastJourneyId}}/missions/999"

    # Verification: HTTP layer
    Then the response status code should be 404

    # Verification: Error response
    And the response body should contain field "error"
    And the response body field "error" should equal string "查無此任務"
