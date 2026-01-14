# Language: zh-TW
# DSL Layer (L1): Business Domain Language
# Source: Release-1-Spec.md - 1.2 使用者登入

Feature: 使用者登入 (Release 1.2)
  作為一個已註冊的使用者
  我想要使用帳號和密碼登入系統
  以便訪問需要登入的功能

  Scenario: 使用者使用正確的帳密登入成功
    Given 系統中存在一位用戶 "Alice" 密碼為 "Test1234!"
    When "Alice" 嘗試使用 "Test1234!" 進行登入
    Then 登入應該成功
    And 她應該收到一組有效的存取 Token
    And 系統應該顯示她的帳號為 "Alice"
    And 系統應該顯示她的經驗值為 0

  Scenario: 使用者使用錯誤的密碼登入失敗
    Given 系統中存在一位用戶 "Bob" 密碼為 "Correct123!"
    When "Bob" 嘗試使用 "WrongPassword!" 進行登入
    Then 登入應該失敗
    And 系統應該提示 "帳號或密碼錯誤"

  Scenario: 使用者使用不存在的帳號登入失敗
    When "NonExistentUser" 嘗試使用 "AnyPassword123!" 進行登入
    Then 登入應該失敗
    And 系統應該提示 "帳號或密碼錯誤"
