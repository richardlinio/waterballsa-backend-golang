# ISA Layer (L2): Backend API E2E Test
# Source: swagger.yaml - /orders, /orders/{orderId}, /orders/{orderId}/action/pay, /users/{userId}/orders
# Test Scope: HTTP Request → API Handler → Database → HTTP Response
# Maps DSL scenarios to concrete API calls

@isa
Feature: Purchase Authorization API Implementation

  Background:
    # Database is cleaned before each scenario (handled by @Before hook)

  # ============================================================
  # Authorization Tests for Order Management APIs
  # ============================================================

  Scenario: Unauthenticated user cannot create order
    # Setup: Create test journey
    Given the database has a journey:
      | title       | 設計模式精通之旅        |
      | slug        | design-patterns-master |
      | description | 學習設計模式            |
      | teacher     | 水球老師                |
      | price       | 1999.00                |

    # Action: Attempt to create order without authentication
    And I set request body to:
      """
      {
        "items": [
          {
            "journeyId": {{lastJourneyId}},
            "quantity": 1
          }
        ]
      }
      """
    When I send "POST" request to "/orders"

    # Verification: HTTP layer - should return 401
    Then the response status code should be 401

    # Verification: Error message
    And the response body should contain field "error"
    And the response body field "error" should equal string "未授權或權杖無效"

  Scenario: User can only view their own orders
    # Setup: Create two users
    Given the database has a user:
      | username   | Alice     |
      | password   | Test1234! |
      | experience | 0         |

    # Store Alice's user ID
    And I copy variable "lastUserId" to "aliceUserId"

    Given the database has a user:
      | username   | Bob        |
      | password   | Secure123! |
      | experience | 0          |

    # Setup: Create test journey
    And the database has a journey:
      | title       | Java 基礎課程     |
      | slug        | java-basics      |
      | description | 學習 Java 基礎    |
      | teacher     | 水球老師          |
      | price       | 2999.00          |

    # Setup: Create Alice's unpaid order
    And the database has an order:
      | user_id    | {{aliceUserId}}   |
      | journey_id | {{lastJourneyId}} |
      | status     | UNPAID            |

    # Store Alice's order ID (lastOrderId is automatically set by database step)
    And I copy variable "lastOrderId" to "aliceOrderId"

    # Setup: Login as Bob
    And I login as "Bob" with password "Secure123!"

    # Action: Bob attempts to view Alice's order
    And I set Authorization header to "{{accessToken}}"
    When I send "GET" request to "/orders/{{aliceOrderId}}"

    # Verification: HTTP layer - should return 404 (not found, for security reasons)
    Then the response status code should be 404

    # Verification: Error message
    And the response body should contain field "error"
    And the response body field "error" should equal string "找不到訂單"

  Scenario: Unauthenticated user cannot pay for order
    # Setup: Create test user
    Given the database has a user:
      | username   | Charlie   |
      | password   | Pass1234! |
      | experience | 0         |

    # Setup: Create test journey
    And the database has a journey:
      | title       | Python 入門         |
      | slug        | python-intro        |
      | description | 學習 Python 程式設計 |
      | teacher     | 水球老師             |
      | price       | 1599.00             |

    # Setup: Create Charlie's unpaid order
    And the database has an order:
      | user_id    | {{lastUserId}}    |
      | journey_id | {{lastJourneyId}} |
      | status     | UNPAID            |

    # Action: Attempt to pay without authentication (no login step)
    # Note: lastOrderId is automatically set by the database step above
    When I send "POST" request to "/orders/{{lastOrderId}}/action/pay"

    # Verification: HTTP layer - should return 401
    Then the response status code should be 401

    # Verification: Error message
    And the response body should contain field "error"
    And the response body field "error" should equal string "未授權或權杖無效"

  Scenario: User can only pay for their own orders
    # Setup: Create two users
    Given the database has a user:
      | username   | Diana     |
      | password   | Diana123! |
      | experience | 0         |

    # Store Diana's user ID
    And I copy variable "lastUserId" to "dianaUserId"

    Given the database has a user:
      | username   | Eva       |
      | password   | Eva12345! |
      | experience | 0         |

    # Setup: Create test journey
    And the database has a journey:
      | title       | JavaScript 完整指南 |
      | slug        | javascript-guide   |
      | description | 學習 JavaScript    |
      | teacher     | 水球老師            |
      | price       | 2999.00            |

    # Setup: Create Diana's unpaid order
    And the database has an order:
      | user_id    | {{dianaUserId}}   |
      | journey_id | {{lastJourneyId}} |
      | status     | UNPAID            |

    # Store Diana's order ID (lastOrderId is automatically set by database step)
    And I copy variable "lastOrderId" to "dianaOrderId"

    # Setup: Login as Eva
    And I login as "Eva" with password "Eva12345!"

    # Action: Eva attempts to pay Diana's order
    And I set Authorization header to "{{accessToken}}"
    When I send "POST" request to "/orders/{{dianaOrderId}}/action/pay"

    # Verification: HTTP layer - should return 404 (not found, for security reasons)
    Then the response status code should be 404

    # Verification: Error message
    And the response body should contain field "error"
    And the response body field "error" should equal string "找不到訂單"

  Scenario: Unauthenticated user cannot view order list
    # Action: Attempt to view order list without authentication
    # Note: Using userId=1 as a test value; the actual userId doesn't matter for 401 test
    When I send "GET" request to "/users/1/orders"

    # Verification: HTTP layer - should return 401
    Then the response status code should be 401

    # Verification: Error message
    And the response body should contain field "error"
    And the response body field "error" should equal string "未授權或權杖無效"
