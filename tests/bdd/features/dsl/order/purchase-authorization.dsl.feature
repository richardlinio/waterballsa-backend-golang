# Language: zh-TW
# DSL Layer (L1): Business Domain Language
# Source: Release-3-Spec.md - 購課流程的認證與授權檢查

Feature: 購課授權檢查 (Release 3)
  作為系統
  我需要驗證使用者的身份和權限
  以便保護訂單資料和購課流程的安全性

  Scenario: 未登入使用者無法建立訂單
    Given 系統中存在一個旅程 "設計模式精通之旅" 價格為 1999 元
    When 未登入的使用者嘗試建立該旅程的訂單
    Then 訂單建立應該失敗
    And 系統應該提示「未授權或權杖無效」

  Scenario: 使用者只能查看自己的訂單
    Given 系統中存在一位用戶 "Alice" 密碼為 "Test1234!"
    And 系統中存在一位用戶 "Bob" 密碼為 "Secure123!"
    And 系統中存在一個旅程 "Java 基礎課程" 價格為 2999 元
    And "Alice" 已有該旅程的未付款訂單
    When "Bob" 嘗試查看 "Alice" 的訂單
    Then 查看應該失敗
    And 系統應該提示「找不到訂單」

  Scenario: 未登入使用者無法付款訂單
    Given 系統中存在一位用戶 "Charlie" 密碼為 "Pass1234!"
    And 系統中存在一個旅程 "Python 入門" 價格為 1599 元
    And "Charlie" 已有該旅程的未付款訂單
    And "Charlie" 已登出
    When 未登入的使用者嘗試付款該訂單
    Then 付款應該失敗
    And 系統應該提示「未授權或權杖無效」

  Scenario: 使用者只能付款自己的訂單
    Given 系統中存在一位用戶 "Diana" 密碼為 "Diana123!"
    And 系統中存在一位用戶 "Eva" 密碼為 "Eva12345!"
    And 系統中存在一個旅程 "JavaScript 完整指南" 價格為 2999 元
    And "Diana" 已有該旅程的未付款訂單
    When "Eva" 嘗試付款 "Diana" 的訂單
    Then 付款應該失敗
    And 系統應該提示「找不到訂單」

  Scenario: 未登入使用者無法查看訂單列表
    When 未登入的使用者嘗試查看訂單列表
    Then 查看應該失敗
    And 系統應該提示「未授權或權杖無效」
