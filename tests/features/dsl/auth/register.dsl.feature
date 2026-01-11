# Language: zh-TW
# DSL Layer (L1): Business Domain Language
# Source: Release-1-Spec.md - 1.1 使用者註冊

Feature: 使用者註冊 (Release 1.1)
  作為一個新訪客
  我想要使用帳號和密碼註冊系統
  以便開始學習旅程並累積經驗值

  Scenario: 使用有效的帳號和密碼成功註冊
    When "Alice" 使用密碼 "Test1234!" 進行註冊
    Then 註冊應該成功
    And 系統中應該存在用戶 "Alice"
    And "Alice" 的經驗值應該為 0

  Scenario: 嘗試註冊已存在的帳號
    Given 系統中存在一位用戶 "Bob" 密碼為 "Secure123!"
    When "Bob" 使用密碼 "NewPassword456!" 進行註冊
    Then 註冊應該失敗
    And 系統應該提示 "使用者名稱已存在"

  Scenario: 使用最短帳號長度註冊 (3 字元)
    When "Tom" 使用密碼 "Valid123!" 進行註冊
    Then 註冊應該成功
    And 系統中應該存在用戶 "Tom"
    And "Tom" 的經驗值應該為 0

  Scenario: 使用最長帳號長度註冊 (50 字元)
    When "alice_chen_waterball_student_learning_java_2025_dec" 使用密碼 "Strong123!" 進行註冊
    Then 註冊應該成功
    And 系統中應該存在用戶 "alice_chen_waterball_student_learning_java_2025_dec"

  Scenario: 使用最短密碼長度註冊 (8 字元)
    When "Charlie" 使用密碼 "Pass123!" 進行註冊
    Then 註冊應該成功
    And 系統中應該存在用戶 "Charlie"
    And "Charlie" 的經驗值應該為 0

  Scenario: 使用最長密碼長度註冊 (72 字元)
    When "Diana" 使用密碼 "A1!bcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@$%*" 進行註冊
    Then 註冊應該成功
    And 系統中應該存在用戶 "Diana"

  Scenario: 使用過短的帳號註冊失敗 (2 字元)
    When "Ed" 使用密碼 "Valid123!" 進行註冊
    Then 註冊應該失敗
    And 系統應該提示 "使用者名稱或密碼格式無效"

  Scenario: 使用過長的帳號註冊失敗 (51 字元)
    When "frank_johnson_waterball_student_learning_python_2025" 使用密碼 "Valid123!" 進行註冊
    Then 註冊應該失敗
    And 系統應該提示 "使用者名稱或密碼格式無效"

  Scenario: 使用包含特殊符號的帳號註冊失敗
    When "alice@chen" 使用密碼 "Valid123!" 進行註冊
    Then 註冊應該失敗
    And 系統應該提示 "使用者名稱或密碼格式無效"

  Scenario: 使用過短的密碼註冊失敗 (7 字元)
    When "Grace" 使用密碼 "Pass12!" 進行註冊
    Then 註冊應該失敗
    And 系統應該提示 "使用者名稱或密碼格式無效"

  Scenario: 使用過長的密碼註冊失敗 (73 字元)
    When "Henry" 使用密碼 "A1!bcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@$%*?" 進行註冊
    Then 註冊應該失敗
    And 系統應該提示 "使用者名稱或密碼格式無效"
