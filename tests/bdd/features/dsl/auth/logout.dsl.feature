# Language: zh-TW
# DSL Layer (L1): Business Domain Language
# Source: Release-1-Spec.md - 1.4 使用者登出

Feature: 使用者登出 (Release 1.4)
  作為一個已登入的使用者
  我想要主動登出系統
  以便確保我的帳號安全，防止 Token 被竊取後遭第三人使用

  Scenario: 已登入使用者成功登出
    Given 系統中存在一位用戶 "Alice" 密碼為 "Test1234!"
    And "Alice" 已使用正確密碼登入
    When "Alice" 執行登出操作
    Then 登出應該成功

  Scenario: 未登入使用者嘗試登出成功
    When 未登入的使用者嘗試執行登出操作
    Then 登出應該成功
    And 系統應該提示 "登出成功"

