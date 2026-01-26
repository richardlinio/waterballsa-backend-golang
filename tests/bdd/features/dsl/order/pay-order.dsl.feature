# Language: zh-TW
# DSL Layer (L1): Business Domain Language
# Source: Release-3-Spec.md - 3.2 訂單付款, 3.3 訂單狀態管理

Feature: 訂單付款與狀態管理 (Release 3.2, 3.3)
  作為一個使用者
  我想要完成訂單付款
  以便取得課程的存取權限

  # 場景組 E: 訂單付款成功流程

  Scenario: 使用者成功付款訂單
    Given 系統中存在一位用戶 "Alice" 密碼為 "Test1234!"
    And 系統中存在一個旅程 "設計模式精通之旅" 價格為 1999 元
    And "Alice" 已有該旅程的未付款訂單
    When "Alice" 點擊「立即支付」
    Then 付款應該成功
    And 訂單狀態應變更為「已付款」
    And 訂單應記錄付款時間
    And 訂單應顯示使用者名稱為 "Alice"
    And 訂單應顯示課程名稱為 "設計模式精通之旅"
    And 訂單應顯示課程價格為 1999 元
    And 訂單應顯示付款金額為 1999 元
    And 訂單應顯示建立時間
    And 訂單應顯示付款時間
    And 訂單應顯示「立即上課」按鈕
    And "Alice" 應該擁有旅程 "設計模式精通之旅" 的存取權

  Scenario: 付款後點擊「立即上課」導向第一個任務
    Given 系統中存在一位用戶 "Bob" 密碼為 "Secure123!"
    And 系統中存在一個旅程 "Java 基礎課程"
    And 該旅程包含一個章節 "第一章：變數與型別" 順序為 1
    And 該章節包含一個影片任務 "認識變數" 順序為 1
    And "Bob" 剛完成該旅程的付款
    When "Bob" 點擊「立即上課」按鈕
    Then 系統應該導向任務 "認識變數" 的頁面

  # 場景組 F: 訂單付款錯誤處理

  Scenario: 已付款訂單無法重複付款
    Given 系統中存在一位用戶 "Charlie" 密碼為 "Pass1234!"
    And 系統中存在一個旅程 "Python 入門" 價格為 1599 元
    And "Charlie" 已有該旅程的已付款訂單
    When "Charlie" 嘗試再次付款該訂單
    Then 付款應該失敗
    And 系統應該提示「訂單已付款」

  Scenario: 過期訂單無法付款
    Given 系統中存在一位用戶 "Diana" 密碼為 "Diana123!"
    And 系統中存在一個旅程 "JavaScript 完整指南" 價格為 2999 元
    And "Diana" 已有該旅程的訂單，建立時間為 73 小時前
    And 該訂單狀態為「已過期」
    When "Diana" 嘗試付款該訂單
    Then 付款應該失敗
    And 系統應該提示「訂單已過期」

  # 場景組 G: 訂單狀態轉換

  Scenario: 訂單建立後狀態為「未付款」
    Given 系統中存在一位用戶 "Eva" 密碼為 "Eva12345!"
    And 系統中存在一個旅程 "React 完整課程" 價格為 3999 元
    When "Eva" 建立該旅程的訂單
    Then 訂單狀態應為「未付款」

  Scenario: 訂單付款後狀態為「已付款」
    Given 系統中存在一位用戶 "Frank" 密碼為 "Frank123!"
    And 系統中存在一個旅程 "Vue.js 實戰課程" 價格為 2499 元
    And "Frank" 已有該旅程的未付款訂單
    When "Frank" 完成付款
    Then 訂單狀態應變更為「已付款」

  Scenario: 訂單建立 72 小時後狀態為「已過期」
    Given 系統中存在一位用戶 "Grace" 密碼為 "Grace123!"
    And 系統中存在一個旅程 "Node.js 後端開發" 價格為 3499 元
    And "Grace" 已建立該旅程的訂單
    When 經過 72 小時後
    Then 訂單狀態應變更為「已過期」

  Scenario: 使用者仍可查看已過期訂單
    Given 系統中存在一位用戶 "Henry" 密碼為 "Henry123!"
    And 系統中存在一個旅程 "Go 語言入門" 價格為 1999 元
    And "Henry" 已有該旅程的已過期訂單
    When "Henry" 查看該訂單
    Then 應該能看到該訂單
    And 訂單狀態應顯示為「已過期」
