# ISA Layer (L2): Backend API E2E Test
# Source: swagger.yaml - /orders/{orderId}/action/pay
# Test Scope: HTTP Request → API Handler → Database → HTTP Response
# Maps DSL scenarios to concrete API calls

@isa
Feature: Order Payment API Implementation

  Background:
    # Database is cleaned before each scenario (handled by @Before hook)

  # ============================================================
  # 場景組 E: 訂單付款成功流程
  # ============================================================

  Scenario: User successfully pays for an unpaid order
    # Setup: Create test user in database
    Given the database has a user:
      | username   | Alice     |
      | password   | Test1234! |
      | experience | 0         |

    # Setup: Create test journey in database
    And the database has a journey:
      | title       | 設計模式精通之旅          |
      | slug        | design-patterns-mastery |
      | description | 學習設計模式              |
      | teacher     | 水球老師                  |
      | price       | 1999.00                  |

    # Setup: Create unpaid order in database
    And the database has an order:
      | user_id    | {{lastUserId}}    |
      | journey_id | {{lastJourneyId}} |
      | status     | UNPAID            |

    # Setup: Login to get access token
    And I login as "Alice" with password "Test1234!"

    # Setup: Set authorization header
    And I set Authorization header to "{{accessToken}}"

    # Action: Pay for the order
    When I send "POST" request to "/orders/{{lastOrderId}}/action/pay"

    # Verification: HTTP layer
    Then the response status code should be 200

    # Verification: Response structure
    And the response body should contain field "id"
    And the response body should contain field "orderNumber"
    And the response body should contain field "status"
    And the response body should contain field "price"
    And the response body should contain field "paidAt"
    And the response body should contain field "message"

    # Verification: Response values
    And the response body field "status" should equal string "PAID"
    And the response body field "price" should equal decimal "1999.00"
    And the response body field "message" should equal string "付款完成"
    And the response body should have non-null field "paidAt"

  @frontend
  Scenario: User clicks "立即上課" button after payment to navigate to first task
    # 前端場景：此場景涉及頁面導航和 UI 互動
    # 後端 API 不負責頁面導航邏輯
    # 前端應該在付款完成後：
    # 1. 透過 GET /journeys/{journeyId} 取得旅程的第一個章節
    # 2. 從第一個章節取得第一個任務
    # 3. 導航至該任務頁面
    # 可考慮使用 Playwright/Cypress 等前端測試工具

  # ============================================================
  # 場景組 F: 訂單付款錯誤處理
  # ============================================================

  Scenario: Already paid order cannot be paid again
    # Setup: Create test user in database
    Given the database has a user:
      | username   | Charlie   |
      | password   | Pass1234! |
      | experience | 0         |

    # Setup: Create test journey in database
    And the database has a journey:
      | title       | Python 入門              |
      | slug        | python-basics           |
      | description | 學習 Python 程式設計     |
      | teacher     | 水球老師                 |
      | price       | 1599.00                 |

    # Setup: Create PAID order in database
    And the database has an order:
      | user_id    | {{lastUserId}}    |
      | journey_id | {{lastJourneyId}} |
      | status     | PAID              |

    # Setup: Login to get access token
    And I login as "Charlie" with password "Pass1234!"

    # Setup: Set authorization header
    And I set Authorization header to "{{accessToken}}"

    # Action: Attempt to pay again
    When I send "POST" request to "/orders/{{lastOrderId}}/action/pay"

    # Verification: HTTP layer - should return conflict
    Then the response status code should be 409

    # Verification: Error message
    And the response body should contain field "error"
    And the response body field "error" should equal string "訂單已經付款"

  Scenario: Expired order cannot be paid
    # Setup: Create test user in database
    Given the database has a user:
      | username   | Diana     |
      | password   | Diana123! |
      | experience | 0         |

    # Setup: Create test journey in database
    And the database has a journey:
      | title       | JavaScript 完整指南      |
      | slug        | javascript-complete     |
      | description | 學習 JavaScript         |
      | teacher     | 水球老師                 |
      | price       | 2999.00                 |

    # Setup: Create EXPIRED order in database
    And the database has an order:
      | user_id    | {{lastUserId}}    |
      | journey_id | {{lastJourneyId}} |
      | status     | EXPIRED           |

    # Setup: Login to get access token
    And I login as "Diana" with password "Diana123!"

    # Setup: Set authorization header
    And I set Authorization header to "{{accessToken}}"

    # Action: Attempt to pay expired order
    When I send "POST" request to "/orders/{{lastOrderId}}/action/pay"

    # Verification: HTTP layer - should return conflict
    Then the response status code should be 409

    # Verification: Error message
    And the response body should contain field "error"
    And the response body field "error" should equal string "訂單已過期"

  # ============================================================
  # 場景組 G: 訂單狀態轉換
  # ============================================================

  Scenario: Order status is UNPAID after creation
    # Setup: Create test user in database
    Given the database has a user:
      | username   | Eva       |
      | password   | Eva12345! |
      | experience | 0         |

    # Setup: Create test journey in database
    And the database has a journey:
      | title       | React 完整課程          |
      | slug        | react-complete         |
      | description | 學習 React             |
      | teacher     | 水球老師                |
      | price       | 3999.00                |

    # Setup: Login to get access token
    And I login as "Eva" with password "Eva12345!"

    # Setup: Set authorization header
    And I set Authorization header to "{{accessToken}}"

    # Setup: Prepare create order request
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

    # Action: Create order
    When I send "POST" request to "/orders"

    # Verification: HTTP layer
    Then the response status code should be 201

    # Verification: Order status is UNPAID
    And the response body field "status" should equal string "UNPAID"

    # Verification: paidAt should be null for unpaid order
    And the response body field "paidAt" should be null

  Scenario: Order status changes to PAID after payment
    # Setup: Create test user in database
    Given the database has a user:
      | username   | Frank     |
      | password   | Frank123! |
      | experience | 0         |

    # Setup: Create test journey in database
    And the database has a journey:
      | title       | Vue.js 實戰課程         |
      | slug        | vuejs-practice         |
      | description | 學習 Vue.js            |
      | teacher     | 水球老師                |
      | price       | 2499.00                |

    # Setup: Create UNPAID order in database
    And the database has an order:
      | user_id    | {{lastUserId}}    |
      | journey_id | {{lastJourneyId}} |
      | status     | UNPAID            |

    # Setup: Login to get access token
    And I login as "Frank" with password "Frank123!"

    # Setup: Set authorization header
    And I set Authorization header to "{{accessToken}}"

    # Action: Pay for the order
    When I send "POST" request to "/orders/{{lastOrderId}}/action/pay"

    # Verification: HTTP layer
    Then the response status code should be 200

    # Verification: Order status changed to PAID
    And the response body field "status" should equal string "PAID"

    # Verification: paidAt should be set
    And the response body should have non-null field "paidAt"

  @skip
  Scenario: Order status changes to EXPIRED after 72 hours
    # 此場景涉及時間推移，需要：
    # 1. 排程任務 (Scheduled Job) 來自動更新過期訂單狀態
    # 2. 或在測試中手動修改訂單的 created_at 時間戳
    #
    # 建議實作方式：
    # - 後端應有排程任務每小時檢查過期訂單
    # - 測試可透過直接修改資料庫 created_at 欄位來模擬時間推移
    # - 然後呼叫 GET /orders/{orderId} 驗證狀態已變更為 EXPIRED
    #
    # 此場景標記為 @skip，等待排程任務實作後再啟用

  Scenario: User can still view expired order details
    # Setup: Create test user in database
    Given the database has a user:
      | username   | Henry     |
      | password   | Henry123! |
      | experience | 0         |

    # Setup: Create test journey in database
    And the database has a journey:
      | title       | Go 語言入門             |
      | slug        | golang-basics          |
      | description | 學習 Go 語言            |
      | teacher     | 水球老師                |
      | price       | 1999.00                |

    # Setup: Create EXPIRED order in database
    And the database has an order:
      | user_id    | {{lastUserId}}    |
      | journey_id | {{lastJourneyId}} |
      | status     | EXPIRED           |

    # Setup: Login to get access token
    And I login as "Henry" with password "Henry123!"

    # Setup: Set authorization header
    And I set Authorization header to "{{accessToken}}"

    # Action: View the expired order
    When I send "GET" request to "/orders/{{lastOrderId}}"

    # Verification: HTTP layer - should be able to view
    Then the response status code should be 200

    # Verification: Order details are returned
    And the response body should contain field "id"
    And the response body should contain field "orderNumber"
    And the response body should contain field "status"

    # Verification: Order status is EXPIRED
    And the response body field "status" should equal string "EXPIRED"

    # Verification: paidAt should be null for expired order
    And the response body field "paidAt" should be null
