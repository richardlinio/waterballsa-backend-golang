# Language: zh-TW
# DSL Layer (L1): Business Domain Language
# Source: Release-3-Spec.md - 3.4 訂單列表與查詢

Feature: 訂單列表與查詢 (Release 3.4)
  作為一個使用者
  我想要查看和管理我的訂單
  以便追蹤我的購買記錄和付款狀態

  Scenario: 使用者查看自己的訂單列表
    Given 系統中存在一位用戶 "Alice" 密碼為 "Test1234!"
    And 系統中存在一個旅程 "設計模式精通之旅" 價格為 1999 元
    And 系統中存在一個旅程 "Java 基礎課程" 價格為 2999 元
    And 系統中存在一個旅程 "Python 入門" 價格為 1599 元
    And "Alice" 已有旅程 "設計模式精通之旅" 的未付款訂單
    And "Alice" 已有旅程 "Java 基礎課程" 的已付款訂單
    And "Alice" 已有旅程 "Python 入門" 的已過期訂單
    When "Alice" 查看訂單列表
    Then 她應該看到 3 個訂單
    And 每個訂單應顯示訂單編號
    And 每個訂單應顯示課程名稱
    And 每個訂單應顯示課程價格
    And 每個訂單應顯示訂單狀態

  Scenario: 使用者沒有訂單時查看空列表
    Given 系統中存在一位用戶 "Bob" 密碼為 "Secure123!"
    And "Bob" 沒有任何訂單
    When "Bob" 查看訂單列表
    Then 應該顯示空列表

  Scenario: 點擊未付款訂單進入付款頁面
    Given 系統中存在一位用戶 "Charlie" 密碼為 "Pass1234!"
    And 系統中存在一個旅程 "JavaScript 完整指南" 價格為 2999 元
    And "Charlie" 已有該旅程的未付款訂單
    And "Charlie" 查看訂單列表
    When "Charlie" 點擊該未付款訂單
    Then 系統應該導向該訂單的付款頁面

  Scenario: 點擊已付款訂單查看訂單詳情
    Given 系統中存在一位用戶 "Diana" 密碼為 "Diana123!"
    And 系統中存在一個旅程 "React 完整課程" 價格為 3999 元
    And "Diana" 已有該旅程的已付款訂單
    And "Diana" 查看訂單列表
    When "Diana" 點擊該已付款訂單
    Then 系統應該顯示該訂單的詳細資訊
    And 訂單詳情應包含訂單編號
    And 訂單詳情應包含使用者名稱
    And 訂單詳情應包含課程名稱
    And 訂單詳情應包含課程價格
    And 訂單詳情應包含訂單狀態
    And 訂單詳情應包含建立時間
    And 訂單詳情應包含付款時間

  Scenario: 點擊已過期訂單查看訂單詳情且不顯示支付按鈕
    Given 系統中存在一位用戶 "Eva" 密碼為 "Eva12345!"
    And 系統中存在一個旅程 "Vue.js 實戰課程" 價格為 2499 元
    And "Eva" 已有該旅程的已過期訂單
    And "Eva" 查看訂單列表
    When "Eva" 點擊該已過期訂單
    Then 系統應該顯示該訂單的詳細資訊
    And 訂單狀態應顯示為「已過期」
    And 不應該顯示「立即支付」按鈕
