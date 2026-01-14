# ISA Layer (L2): Implementation/API Level
# Source: swagger.yaml - /auth/register endpoint
# Maps DSL scenarios to concrete HTTP requests

@isa
Feature: User Registration API Implementation

  Background:
    # Database is cleaned before each scenario (handled by @Before hook)

  Scenario: Successful registration with valid credentials
    # No setup needed - user doesn't exist yet

    # Setup: Prepare registration request body
    Given I set request body to:
      """
      {
        "username": "Alice",
        "password": "Test1234!"
      }
      """

    # Action: Register new user with valid credentials
    When I send "POST" request to "/auth/register"

    # Verification: HTTP layer
    Then the response status code should be 201

    # Verification: Response structure
    And the response body should contain field "message"
    And the response body should contain field "userId"

    # Verification: Response values
    And the response body field "message" should equal string "註冊成功"

  Scenario: Attempt to register with existing username
    # Setup: Create test user in database
    Given the database has a user:
      | username | Bob         |
      | password | Secure123!  |

    # Setup: Prepare registration request body with existing username
    And I set request body to:
      """
      {
        "username": "Bob",
        "password": "NewPassword456!"
      }
      """

    # Action: Attempt to register with same username
    When I send "POST" request to "/auth/register"

    # Verification: HTTP layer
    Then the response status code should be 409

    # Verification: Error response
    And the response body should contain field "error"
    And the response body field "error" should equal string "使用者名稱已存在"

  Scenario: Register with minimum username length (3 characters)
    # No setup needed - user doesn't exist yet

    # Setup: Prepare registration request body with 3-character username
    Given I set request body to:
      """
      {
        "username": "Tom",
        "password": "Valid123!"
      }
      """

    # Action: Register with 3-character username
    When I send "POST" request to "/auth/register"

    # Verification: HTTP layer
    Then the response status code should be 201

    # Verification: Response structure
    And the response body should contain field "message"
    And the response body should contain field "userId"

    # Verification: Response values
    And the response body field "message" should equal string "註冊成功"

  Scenario: Register with maximum username length (50 characters)
    # No setup needed - user doesn't exist yet

    # Setup: Prepare registration request body with 50-character username
    Given I set request body to:
      """
      {
        "username": "alice_chen_waterball_student_learning_java_202512",
        "password": "Strong123!"
      }
      """

    # Action: Register user
    When I send "POST" request to "/auth/register"

    # Verification: HTTP layer
    Then the response status code should be 201

    # Verification: Response structure
    And the response body should contain field "message"
    And the response body should contain field "userId"

    # Verification: Response values
    And the response body field "message" should equal string "註冊成功"

  Scenario: Register with minimum password length (8 characters)
    # No setup needed - user doesn't exist yet

    # Setup: Prepare registration request body with 8-character password
    Given I set request body to:
      """
      {
        "username": "Charlie",
        "password": "Pass123!"
      }
      """

    # Action: Register user
    When I send "POST" request to "/auth/register"

    # Verification: HTTP layer
    Then the response status code should be 201

    # Verification: Response structure
    And the response body should contain field "message"
    And the response body should contain field "userId"

    # Verification: Response values
    And the response body field "message" should equal string "註冊成功"

  Scenario: Register with maximum password length (72 characters)
    # No setup needed - user doesn't exist yet

    # Setup: Prepare registration request body with 72-character password
    Given I set request body to:
      """
      {
        "username": "Diana",
        "password": "A1!bcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@$%*"
      }
      """

    # Action: Register user
    When I send "POST" request to "/auth/register"

    # Verification: HTTP layer
    Then the response status code should be 201

    # Verification: Response structure
    And the response body should contain field "message"
    And the response body should contain field "userId"

    # Verification: Response values
    And the response body field "message" should equal string "註冊成功"

  Scenario: Failed registration with username too short (2 characters)
    # No setup needed - testing validation

    # Setup: Prepare registration request body with 2-character username
    Given I set request body to:
      """
      {
        "username": "Ed",
        "password": "Valid123!"
      }
      """

    # Action: Attempt to register
    When I send "POST" request to "/auth/register"

    # Verification: HTTP layer
    Then the response status code should be 400

    # Verification: Error response
    And the response body should contain field "error"
    And the response body field "error" should equal string "使用者名稱或密碼格式無效"

  Scenario: Failed registration with username too long (51 characters)
    # No setup needed - testing validation

    # Setup: Prepare registration request body with 51-character username
    Given I set request body to:
      """
      {
        "username": "frank_johnson_waterball_student_learning_python_2025",
        "password": "Valid123!"
      }
      """

    # Action: Attempt to register
    When I send "POST" request to "/auth/register"

    # Verification: HTTP layer
    Then the response status code should be 400

    # Verification: Error response
    And the response body should contain field "error"
    And the response body field "error" should equal string "使用者名稱或密碼格式無效"

  Scenario: Failed registration with invalid username characters
    # No setup needed - testing validation

    # Setup: Prepare registration request body with invalid username characters
    Given I set request body to:
      """
      {
        "username": "alice@chen",
        "password": "Valid123!"
      }
      """

    # Action: Attempt to register
    When I send "POST" request to "/auth/register"

    # Verification: HTTP layer
    Then the response status code should be 400

    # Verification: Error response
    And the response body should contain field "error"
    And the response body field "error" should equal string "使用者名稱或密碼格式無效"

  Scenario: Failed registration with password too short (7 characters)
    # No setup needed - testing validation

    # Setup: Prepare registration request body with 7-character password
    Given I set request body to:
      """
      {
        "username": "Grace",
        "password": "Pass12!"
      }
      """

    # Action: Attempt to register
    When I send "POST" request to "/auth/register"

    # Verification: HTTP layer
    Then the response status code should be 400

    # Verification: Error response
    And the response body should contain field "error"
    And the response body field "error" should equal string "使用者名稱或密碼格式無效"

  Scenario: Failed registration with password too long (73 characters)
    # No setup needed - testing validation

    # Setup: Prepare registration request body with 73-character password
    Given I set request body to:
      """
      {
        "username": "Henry",
        "password": "A1!bcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@$%*?&#A"
      }
      """

    # Action: Attempt to register
    When I send "POST" request to "/auth/register"

    # Verification: HTTP layer
    Then the response status code should be 400

    # Verification: Error response
    And the response body should contain field "error"
    And the response body field "error" should equal string "使用者名稱或密碼格式無效"
