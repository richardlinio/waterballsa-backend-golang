# Language: zh-TW
# DSL Layer (L1): Business Domain Language
# Source: Release-1-Spec.md - Token Refresh Functionality

Feature: 刷新存取令牌 (Refresh Token)
  作為一個已登入的使用者
  我想要使用刷新令牌來獲取新的存取令牌
  以便在存取令牌過期後繼續使用系統

  Scenario: 使用有效的刷新令牌獲取新的存取令牌
    Given 系統中存在一位用戶 "Alice" 密碼為 "Test1234!"
    And "Alice" 已經成功登入並獲得刷新令牌
    When "Alice" 使用刷新令牌請求新的存取令牌
    Then 刷新應該成功
    And 她應該收到一組新的存取 Token
    And 系統應該顯示她的帳號為 "Alice"

  Scenario: 使用無效的刷新令牌
    When 使用者使用無效的刷新令牌請求新的存取令牌
    Then 刷新應該失敗
    And 系統應該提示 "未授權或權杖無效"

  Scenario: 未提供刷新令牌
    When 使用者未提供刷新令牌就請求新的存取令牌
    Then 刷新應該失敗
    And 系統應該提示 "未授權或權杖無效"

  Scenario: 使用已被撤銷的刷新令牌
    Given 系統中存在一位用戶 "Bob" 密碼為 "Secure123!"
    And "Bob" 已經成功登入並獲得刷新令牌
    And "Bob" 已使用該刷新令牌獲取新令牌
    When "Bob" 再次使用舊的刷新令牌請求存取令牌
    Then 刷新應該失敗
    And 系統應該提示 "未授權或權杖無效"
