# ISA Layer (L2): Implementation/API Level
# Source: /auth/refresh endpoint
# Maps DSL scenarios to concrete HTTP requests

@isa
Feature: Refresh Token API Implementation

  Background:
    # Database is cleaned before each scenario (handled by @Before hook)

  Scenario: Successful refresh with valid token
    # Setup: Create test user and login to get refresh token
    Given the database has a user:
      | username | Alice     |
      | password | Test1234! |
    And I set request body to:
      """
      {
        "username": "Alice",
        "password": "Test1234!"
      }
      """
    And I send "POST" request to "/auth/login"
    And I extract cookie "refresh_token" from response

    # Action: Use refresh token to get new access token
    When I use stored cookie "refresh_token"
    And I send "POST" request to "/auth/refresh"

    # Verification: HTTP layer
    Then the response status code should be 200

    # Verification: Cookies (new tokens issued)
    And the response should set cookie "access_token"
    And the response should set cookie "refresh_token"

    # Verification: Response structure
    And the response body should contain field "accessToken"
    And the response body should contain field "user.id"
    And the response body should contain field "user.username"
    And the response body should contain field "user.experience"

    # Verification: Response values
    And the response body field "user.username" should equal string "Alice"

  Scenario: Refresh without token
    # Action: Attempt refresh without providing refresh token
    When I send "POST" request to "/auth/refresh"

    # Verification: HTTP layer
    Then the response status code should be 401

    # Verification: Error response
    And the response body should contain field "error"
    And the response body field "error" should equal string "未授權或權杖無效"

  Scenario: Refresh with invalid token
    # Setup: Set invalid refresh token
    Given I set cookie "refresh_token" to "invalid_token_12345"

    # Action: Attempt refresh with invalid token
    When I send "POST" request to "/auth/refresh"

    # Verification: HTTP layer
    Then the response status code should be 401

    # Verification: Error response
    And the response body should contain field "error"
    And the response body field "error" should equal string "未授權或權杖無效"

  Scenario: Token rotation - old token becomes invalid
    # Setup: Create test user and login to get initial refresh token
    Given the database has a user:
      | username | Bob         |
      | password | Secure123!  |
    And I set request body to:
      """
      {
        "username": "Bob",
        "password": "Secure123!"
      }
      """
    And I send "POST" request to "/auth/login"
    And I extract cookie "refresh_token" from response

    # Action: Use refresh token once (this rotates the token)
    When I use stored cookie "refresh_token"
    And I send "POST" request to "/auth/refresh"
    Then the response status code should be 200

    # Action: Try to use the old refresh token again
    When I use stored cookie "refresh_token"
    And I send "POST" request to "/auth/refresh"

    # Verification: Old token should be invalid
    Then the response status code should be 401
    And the response body field "error" should equal string "未授權或權杖無效"
