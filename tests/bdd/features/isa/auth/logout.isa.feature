# ISA Layer (L2): Implementation/API Level
# Source: swagger.yaml - /auth/logout endpoint
# Maps DSL scenarios to concrete HTTP requests

@isa
Feature: User Logout API Implementation

  Background:
    # Database is cleaned before each scenario (handled by @Before hook)

  Scenario: Logged-in user successfully logs out
    # Setup: Create test user directly in database
    Given the database has a user:
      | username   | Alice     |
      | password   | Test1234! |
      | experience | 0         |

    # Setup: Prepare login request body
    And I set request body to:
      """
      {
        "username": "Alice",
        "password": "Test1234!"
      }
      """

    # Action: Login to get access token
    When I send "POST" request to "/auth/login"

    # Verification: Login successful
    Then the response status code should be 200
    And the response body should contain field "accessToken"

    # Store token for logout
    And I store the response field "accessToken" as "token"

    # Setup: Set Authorization header for logout
    Given I set Authorization header to "{{token}}"

    # Action: Logout with valid token
    When I send "POST" request to "/auth/logout"

    # Verification: HTTP layer
    Then the response status code should be 200

    # Verification: Response structure
    And the response body should contain field "message"

    # Verification: Response values
    And the response body field "message" should equal string "登出成功"

  Scenario: Non-logged-in user attempts to logout and fails
    # No setup needed - no user logged in

    # Action: Attempt logout without Authorization header
    When I send "POST" request to "/auth/logout"

    # Verification: HTTP layer
    Then the response status code should be 401

    # Verification: Error response
    And the response body should contain field "error"
    And the response body field "error" should equal string "未授權或權杖無效"
