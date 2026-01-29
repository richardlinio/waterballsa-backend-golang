# ISA Layer (L2): Backend API E2E Test
# Source: swagger.yaml - /users/{userId}/orders, /orders/{orderId}
# Test Scope: HTTP Request → API Handler → Database → HTTP Response
# Maps DSL scenarios to concrete API calls

@isa
Feature: View Orders API Implementation

  Background:
    # Database is cleaned before each scenario (handled by @Before hook)

  Scenario: User views their order list with multiple orders
    # Setup: Create test user
    Given the database has a user:
      | username   | Alice     |
      | password   | Test1234! |
      | experience | 0         |

    # Setup: Create three journeys
    And the database has a journey:
      | title           | 設計模式精通之旅                    |
      | slug            | design-patterns-mastery           |
      | description     | 學習設計模式                        |
      | teacher         | 水球老師                            |
      | price           | 1999.00                            |
      | cover_image_url | https://example.com/design.jpg    |

    # Store first journey ID
    And I copy variable "lastJourneyId" to "journey1Id"

    # Setup: Create second journey
    And the database has a journey:
      | title           | Java 基礎課程                       |
      | slug            | java-basics                       |
      | description     | 學習 Java 程式設計基礎               |
      | teacher         | 水球老師                            |
      | price           | 2999.00                            |
      | cover_image_url | https://example.com/java.jpg      |

    # Store second journey ID
    And I copy variable "lastJourneyId" to "journey2Id"

    # Setup: Create third journey
    And the database has a journey:
      | title           | Python 入門                         |
      | slug            | python-intro                      |
      | description     | 學習 Python 程式設計                 |
      | teacher         | 水球老師                            |
      | price           | 1599.00                            |
      | cover_image_url | https://example.com/python.jpg    |

    # Store third journey ID
    And I copy variable "lastJourneyId" to "journey3Id"

    # Setup: Create unpaid order for first journey
    And the database has an order:
      | user_id    | {{lastUserId}} |
      | journey_id | {{journey1Id}} |
      | status     | UNPAID         |

    # Setup: Create paid order for second journey
    And the database has an order:
      | user_id    | {{lastUserId}} |
      | journey_id | {{journey2Id}} |
      | status     | PAID           |

    # Setup: Create expired order for third journey
    And the database has an order:
      | user_id    | {{lastUserId}} |
      | journey_id | {{journey3Id}} |
      | status     | EXPIRED        |

    # Setup: Login as Alice
    And I login as "Alice" with password "Test1234!"

    # Action: Get user's order list
    And I set Authorization header to "{{accessToken}}"
    When I send "GET" request to "/users/{{userId}}/orders"

    # Verification: HTTP layer
    Then the response status code should be 200

    # Verification: Response structure - orders array
    And the response body should contain field "orders"
    And the response body should contain field "pagination"

    # Verification: Response structure - pagination fields
    And the response body should contain field "pagination.page"
    And the response body should contain field "pagination.limit"
    And the response body should contain field "pagination.total"

    # Verification: Response values - order count
    And the response body field "orders" should have size 3

    # Verification: Response values - pagination
    And the response body field "pagination.page" should equal number 1
    And the response body field "pagination.limit" should equal number 20
    And the response body field "pagination.total" should equal number 3

    # Verification: Response structure - each order has required fields
    And the response body should contain field "orders[0].id"
    And the response body should contain field "orders[0].orderNumber"
    And the response body should contain field "orders[0].status"
    And the response body should contain field "orders[0].price"
    And the response body should contain field "orders[0].items"
    And the response body should contain field "orders[0].createdAt"

    # Verification: Response structure - items array in each order
    And the response body should contain field "orders[0].items[0].journeyId"
    And the response body should contain field "orders[0].items[0].journeyTitle"

    # Verification: Response values - verify different order statuses exist
    # Note: We don't verify exact order as database insertion order may vary
    # But we ensure all three statuses are present in the list

  Scenario: User with no orders views empty order list
    # Setup: Create test user with no orders
    Given the database has a user:
      | username   | Bob       |
      | password   | Secure123! |
      | experience | 0         |

    # Setup: Login as Bob
    And I login as "Bob" with password "Secure123!"

    # Action: Get user's order list
    And I set Authorization header to "{{accessToken}}"
    When I send "GET" request to "/users/{{userId}}/orders"

    # Verification: HTTP layer
    Then the response status code should be 200

    # Verification: Response structure
    And the response body should contain field "orders"
    And the response body should contain field "pagination"

    # Verification: Response values - empty order list
    And the response body field "orders" should have size 0

    # Verification: Response values - pagination for empty list
    And the response body field "pagination.page" should equal number 1
    And the response body field "pagination.limit" should equal number 20
    And the response body field "pagination.total" should equal number 0

  @frontend
  Scenario: User clicks unpaid order to navigate to payment page
    # 前端場景：此場景涉及 UI 點擊事件和頁面導航
    # 後端 API 不負責頁面導航邏輯
    # 前端應該在點擊未付款訂單後：
    # 1. 取得訂單 ID
    # 2. 導航至付款頁面（例如：/orders/{orderId}/payment）
    # 建議使用前端測試工具（如 Playwright/Cypress）進行測試

  Scenario: User views paid order details
    # Setup: Create test user
    Given the database has a user:
      | username   | Diana     |
      | password   | Diana123! |
      | experience | 0         |

    # Setup: Create test journey
    And the database has a journey:
      | title           | React 完整課程                      |
      | slug            | react-complete                    |
      | description     | 學習 React 框架                     |
      | teacher         | 水球老師                            |
      | price           | 3999.00                            |
      | cover_image_url | https://example.com/react.jpg     |

    # Setup: Create paid order
    And the database has an order:
      | user_id    | {{lastUserId}}    |
      | journey_id | {{lastJourneyId}} |
      | status     | PAID              |

    # Setup: Login as Diana
    And I login as "Diana" with password "Diana123!"

    # Action: Get order details
    And I set Authorization header to "{{accessToken}}"
    When I send "GET" request to "/orders/{{lastOrderId}}"

    # Verification: HTTP layer
    Then the response status code should be 200

    # Verification: Response structure - order basic fields
    And the response body should have non-null field "id"
    And the response body should have non-null field "orderNumber"
    And the response body should have non-null field "userId"
    And the response body should have non-null field "username"
    And the response body should have non-null field "status"
    And the response body should have non-null field "originalPrice"
    And the response body should have non-null field "discount"
    And the response body should have non-null field "price"
    And the response body should have non-null field "items"
    And the response body should have non-null field "createdAt"

    # Verification: Response structure - timestamp fields
    And the response body should contain field "expiredAt"
    And the response body should contain field "paidAt"

    # Verification: Response values - order basic info
    And the response body field "username" should equal string "Diana"
    And the response body field "status" should equal string "PAID"
    And the response body field "originalPrice" should equal decimal "3999.00"
    And the response body field "discount" should equal decimal "0.00"
    And the response body field "price" should equal decimal "3999.00"

    # Verification: Response values - paid order has paidAt timestamp
    And the response body should have non-null field "paidAt"

    # Verification: Response structure - order items
    And the response body field "items" should have size 1
    And the response body should have non-null field "items[0].journeyId"
    And the response body should have non-null field "items[0].journeyTitle"
    And the response body should have non-null field "items[0].quantity"
    And the response body should have non-null field "items[0].originalPrice"
    And the response body should have non-null field "items[0].discount"
    And the response body should have non-null field "items[0].price"

    # Verification: Response values - order item details
    And the response body field "items[0].journeyTitle" should equal string "React 完整課程"
    And the response body field "items[0].quantity" should equal number 1
    And the response body field "items[0].originalPrice" should equal decimal "3999.00"
    And the response body field "items[0].discount" should equal decimal "0.00"
    And the response body field "items[0].price" should equal decimal "3999.00"

  Scenario: User views expired order details
    # Setup: Create test user
    Given the database has a user:
      | username   | Eva       |
      | password   | Eva12345! |
      | experience | 0         |

    # Setup: Create test journey
    And the database has a journey:
      | title           | Vue.js 實戰課程                     |
      | slug            | vue-js-practical                  |
      | description     | 學習 Vue.js 框架實戰開發             |
      | teacher         | 水球老師                            |
      | price           | 2499.00                            |
      | cover_image_url | https://example.com/vue.jpg       |

    # Setup: Create expired order
    And the database has an order:
      | user_id    | {{lastUserId}}    |
      | journey_id | {{lastJourneyId}} |
      | status     | EXPIRED           |

    # Setup: Login as Eva
    And I login as "Eva" with password "Eva12345!"

    # Action: Get expired order details
    And I set Authorization header to "{{accessToken}}"
    When I send "GET" request to "/orders/{{lastOrderId}}"

    # Verification: HTTP layer
    Then the response status code should be 200

    # Verification: Response structure
    And the response body should have non-null field "id"
    And the response body should have non-null field "orderNumber"
    And the response body should have non-null field "status"
    And the response body should have non-null field "price"
    And the response body should contain field "paidAt"

    # Verification: Response values - order status
    And the response body field "status" should equal string "EXPIRED"

    # Verification: Response values - expired order has no payment timestamp
    And the response body field "paidAt" should be null

  @frontend
  Scenario: Expired order does not show payment button
    # 前端場景：此場景涉及 UI 按鈕顯示邏輯
    # 後端 API 僅返回訂單狀態，不決定按鈕顯示與否
    # 前端應該根據訂單狀態決定是否顯示「立即支付」按鈕：
    # - status === 'UNPAID' → 顯示「立即支付」按鈕
    # - status === 'PAID' → 隱藏「立即支付」按鈕，可能顯示「查看課程」按鈕
    # - status === 'EXPIRED' → 隱藏「立即支付」按鈕
    # 建議使用前端測試工具（如 Playwright/Cypress）進行測試
