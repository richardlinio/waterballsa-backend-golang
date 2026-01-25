# ISA Layer (L2): Backend API E2E Test
# Source: swagger.yaml - /orders endpoint
# Test Scope: HTTP Request → API Handler → Database → HTTP Response
# Maps DSL scenarios to concrete API calls

@isa
Feature: Create Order API Implementation

  Background:
    # Database is cleaned before each scenario (handled by @Before hook)

  # ============================================================
  # Scenario Group A & B: Purchase Button Display and Navigation
  # These are FRONTEND scenarios - cannot be tested via backend API
  # ============================================================

  @frontend
  Scenario: Already purchased journey shows "Continue Learning" button
    # 前端場景：此場景測試 UI 按鈕顯示邏輯
    # 後端 API /journeys/{journeyId} 不返回購買狀態或按鈕顯示資訊
    # 建議使用前端測試工具（如 Playwright/Cypress）進行測試

  @frontend
  Scenario: Not purchased journey shows "Join Course Now" button
    # 前端場景：此場景測試 UI 按鈕顯示邏輯
    # 後端 API /journeys/{journeyId} 不返回購買狀態或按鈕顯示資訊
    # 建議使用前端測試工具（如 Playwright/Cypress）進行測試

  @frontend
  Scenario: Guest viewing journey shows "Join Course Now" button
    # 前端場景：此場景測試 UI 按鈕顯示邏輯
    # 訪客查看旅程屬於前端渲染邏輯
    # 建議使用前端測試工具（如 Playwright/Cypress）進行測試

  @frontend
  Scenario: Already purchased journey clicking "Continue Learning" navigates to first task
    # 前端場景：此場景測試導航行為
    # 導航邏輯屬於前端路由處理
    # 建議使用前端測試工具（如 Playwright/Cypress）進行測試

  @frontend
  Scenario: Guest clicking purchase button navigates to login page
    # 前端場景：此場景測試導航行為
    # 導航邏輯屬於前端路由處理
    # 建議使用前端測試工具（如 Playwright/Cypress）進行測試

  # ============================================================
  # Scenario Group C: Successful Order Creation Flow
  # ============================================================

  Scenario: User successfully creates an order
    # Setup: Create test user and journey
    Given the database has a user:
      | username   | Diana       |
      | password   | Diana123!   |
      | experience | 0           |
    And the database has a journey:
      | title            | Vue.js 實戰課程                     |
      | slug             | vue-js-practical                   |
      | description      | 學習 Vue.js 框架實戰開發            |
      | teacher          | 王老師                              |
      | price            | 2499.00                            |
      | cover_image_url  | https://example.com/vue-cover.jpg  |

    # Setup: Login as Diana and get access token
    And I login as "Diana" with password "Diana123!"

    # Action: Create order for the journey
    And I set Authorization header to "{{accessToken}}"
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

    # Verification: HTTP layer
    Then the response status code should be 201

    # Verification: Response structure - Fields that must have values
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
    And the response body should have non-null field "expiredAt"

    # Verification: Response structure - Fields that can be null
    And the response body should contain field "paidAt"

    # Verification: Response values - Order basic info
    And the response body field "username" should equal string "Diana"
    And the response body field "status" should equal string "UNPAID"
    And the response body field "originalPrice" should equal decimal "2499.00"
    And the response body field "discount" should equal decimal "0.00"
    And the response body field "price" should equal decimal "2499.00"
    And the response body field "paidAt" should be null

    # Verification: Response values - Order items
    And the response body field "items" should have size 1
    And the response body field "items[0].journeyTitle" should equal string "Vue.js 實戰課程"
    And the response body field "items[0].quantity" should equal number 1
    And the response body field "items[0].originalPrice" should equal decimal "2499.00"
    And the response body field "items[0].discount" should equal decimal "0.00"
    And the response body field "items[0].price" should equal decimal "2499.00"

  Scenario: Order number format is correct and contains user ID
    # Setup: Create test user with specific ID
    Given the database has a user:
      | username   | Eva         |
      | password   | Eva12345!   |
      | experience | 0           |
    And the database has a journey:
      | title            | Node.js 後端開發                    |
      | slug             | nodejs-backend                     |
      | description      | 學習 Node.js 後端開發               |
      | teacher          | 張老師                              |
      | price            | 3499.00                            |
      | cover_image_url  | https://example.com/node-cover.jpg |

    # Setup: Login and create order
    And I login as "Eva" with password "Eva12345!"
    And I set Authorization header to "{{accessToken}}"
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

    # Verification: HTTP layer
    Then the response status code should be 201

    # Verification: Order number contains user ID
    # Order number format: {timestamp(10位)}{userId}{randomCode(5位)}
    And the response body should contain field "orderNumber"
    And the response body field "orderNumber" should contain "{{lastUserId}}"

  Scenario: Order price is locked at creation time
    # Setup: Create test user and journey
    Given the database has a user:
      | username   | Frank       |
      | password   | Frank123!   |
      | experience | 0           |
    And the database has a journey:
      | title            | Go 語言入門                         |
      | slug             | go-language-intro                  |
      | description      | 學習 Go 語言基礎                    |
      | teacher          | 陳老師                              |
      | price            | 1999.00                            |
      | cover_image_url  | https://example.com/go-cover.jpg   |

    # Setup: Login and create first order
    And I login as "Frank" with password "Frank123!"
    And I set Authorization header to "{{accessToken}}"
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

    # Verification: First order created successfully with original price
    Then the response status code should be 201
    And the response body field "price" should equal decimal "1999.00"
    And the response body field "items[0].price" should equal decimal "1999.00"

    # Store the first order ID for later verification
    And I store the response field "id" as "firstOrderId"

    # Action: Update journey price in database
    Given I update the journey with id "{{lastJourneyId}}" to price "2999.00"

    # Verification: Retrieve first order and verify price is still locked at 1999.00
    And I set Authorization header to "{{accessToken}}"
    When I send "GET" request to "/orders/{{firstOrderId}}"
    Then the response status code should be 200
    And the response body field "price" should equal decimal "1999.00"
    And the response body field "items[0].originalPrice" should equal decimal "1999.00"

  # ============================================================
  # Scenario Group D: Order Creation Error Handling
  # ============================================================

  Scenario: User cannot create order for already purchased journey
    # Setup: Create test user, journey, and purchase record
    Given the database has a user:
      | username   | Grace       |
      | password   | Grace123!   |
      | experience | 0           |
    And the database has a journey:
      | title            | Docker 容器化技術                   |
      | slug             | docker-containerization            |
      | description      | 學習 Docker 容器化技術              |
      | teacher          | 李老師                              |
      | price            | 2799.00                            |
      | cover_image_url  | https://example.com/docker.jpg     |
    And the database has an order:
      | user_id    | {{lastUserId}}    |
      | journey_id | {{lastJourneyId}} |
      | status     | PAID              |

    # Setup: Login as Grace
    And I login as "Grace" with password "Grace123!"

    # Action: Attempt to create order for already purchased journey
    And I set Authorization header to "{{accessToken}}"
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

    # Verification: HTTP layer - should return 409 Conflict
    Then the response status code should be 409

    # Verification: Error message
    And the response body should contain field "error"
    And the response body field "error" should equal string "你已經購買此課程"

  Scenario: User with existing unpaid order receives existing order instead of creating new one
    # Setup: Create test user, journey, and unpaid order
    Given the database has a user:
      | username   | Henry       |
      | password   | Henry123!   |
      | experience | 0           |
    And the database has a journey:
      | title            | Kubernetes 實戰                     |
      | slug             | kubernetes-practical               |
      | description      | 學習 Kubernetes 容器編排            |
      | teacher          | 周老師                              |
      | price            | 3999.00                            |
      | cover_image_url  | https://example.com/k8s.jpg        |
    And the database has an order:
      | user_id    | {{lastUserId}}    |
      | journey_id | {{lastJourneyId}} |
      | status     | UNPAID            |

    # Setup: Login as Henry
    And I login as "Henry" with password "Henry123!"

    # Action: Attempt to create another order for the same journey
    And I set Authorization header to "{{accessToken}}"
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

    # Verification: HTTP layer - should return 200 OK (existing order)
    Then the response status code should be 200

    # Verification: Returns existing order
    And the response body should contain field "id"
    And the response body should contain field "orderNumber"
    And the response body field "status" should equal string "UNPAID"
    And the response body field "price" should equal decimal "3999.00"
    And the response body field "items[0].journeyTitle" should equal string "Kubernetes 實戰"

    # Verification: Database should still have only 1 unpaid order
    And the database should have 1 unpaid order for user "{{lastUserId}}" and journey "{{lastJourneyId}}"

  Scenario: User can have multiple unpaid orders for different journeys
    # Setup: Create test user and two different journeys
    Given the database has a user:
      | username   | Ivy         |
      | password   | Ivy12345!   |
      | experience | 0           |
    And the database has a journey:
      | title            | AWS 雲端服務                        |
      | slug             | aws-cloud-services                 |
      | description      | 學習 AWS 雲端服務                   |
      | teacher          | 吳老師                              |
      | price            | 4999.00                            |
      | cover_image_url  | https://example.com/aws.jpg        |

    # Store first journey ID before creating second journey
    And I copy variable "lastJourneyId" to "awsJourneyId"

    # Create second journey
    And the database has a journey:
      | title            | Azure 雲端架構                      |
      | slug             | azure-cloud-architecture           |
      | description      | 學習 Azure 雲端架構                 |
      | teacher          | 鄭老師                              |
      | price            | 4599.00                            |
      | cover_image_url  | https://example.com/azure.jpg      |

    # Note: lastJourneyId is now Azure journey ID, awsJourneyId was stored earlier

    # Setup: Create unpaid order for first journey (AWS)
    And the database has an order:
      | user_id    | {{lastUserId}}  |
      | journey_id | {{awsJourneyId}} |
      | status     | UNPAID          |

    # Setup: Login as Ivy
    And I login as "Ivy" with password "Ivy12345!"

    # Action: Create order for second journey (Azure, using lastJourneyId)
    And I set Authorization header to "{{accessToken}}"
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

    # Verification: HTTP layer - should return 201 Created
    Then the response status code should be 201

    # Verification: Order created successfully
    And the response body should contain field "id"
    And the response body field "status" should equal string "UNPAID"
    And the response body field "price" should equal decimal "4599.00"
    And the response body field "items[0].journeyTitle" should equal string "Azure 雲端架構"

    # Verification: Database should have 2 unpaid orders for this user
    And the database should have 2 unpaid orders for user "{{lastUserId}}"
