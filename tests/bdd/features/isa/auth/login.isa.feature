# ISA Layer (L2): Implementation/API Level
# Source: swagger.yaml - /auth/login endpoint
# Maps DSL scenarios to concrete HTTP requests

@isa
Feature: User Login API Implementation

  Background:
    # Database is cleaned before each scenario (handled by @Before hook)

  Scenario: Successful login with valid credentials
    # Setup: Create test user directly in database
    Given the database has a user:
      | username   | Alice       |
      | password   | Test1234!   |
      | experience | 0           |

    # Setup: Prepare login request body
    And I set request body to:
      """
      {
        "username": "Alice",
        "password": "Test1234!"
      }
      """

    # Action: Login with valid credentials
    When I send "POST" request to "/auth/login"

    # Verification: HTTP layer
    Then the response status code should be 200

    # Verification: Response structure
    And the response body should contain field "accessToken"
    And the response body should contain field "user.id"
    And the response body should contain field "user.username"
    And the response body should contain field "user.experience"

    # Verification: Response values
    And the response body field "user.username" should equal string "Alice"
    And the response body field "user.experience" should equal number 0

  Scenario: Failed login with wrong password
    # Setup: Create test user
    Given the database has a user:
      | username | Bob         |
      | password | Correct123! |

    # Setup: Prepare login request body with wrong password
    And I set request body to:
      """
      {
        "username": "Bob",
        "password": "WrongPassword!"
      }
      """

    # Action: Attempt login with wrong password
    When I send "POST" request to "/auth/login"

    # Verification: HTTP layer
    Then the response status code should be 401

    # Verification: Error response
    And the response body should contain field "error"
    And the response body field "error" should equal string "帳號或密碼錯誤"

  Scenario: Failed login with non-existent username
    # No setup needed - user doesn't exist

    # Setup: Prepare login request body with non-existent user
    Given I set request body to:
      """
      {
        "username": "NonExistentUser",
        "password": "AnyPassword123!"
      }
      """

    # Action: Attempt login with non-existent user
    When I send "POST" request to "/auth/login"

    # Verification: HTTP layer
    Then the response status code should be 401

    # Verification: Error response
    And the response body should contain field "error"
    And the response body field "error" should equal string "帳號或密碼錯誤"
