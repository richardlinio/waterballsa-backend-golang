# Language: zh-TW
# DSL Layer (L1): Business Domain Language
# Source: Release-3-Spec.md - 3.1 課程購買入口, 3.2 訂單建立

Feature: 購買旅程流程 (Release 3.1, 3.2)
  作為一個使用者
  我想要透過課程詳情頁購買旅程
  以便開始學習課程內容

  # 場景組 A: 購買按鈕顯示狀態

  Scenario: 已購買旅程顯示「繼續學習」按鈕
    Given 系統中存在一位用戶 "Alice" 密碼為 "Test1234!"
    And 系統中存在一個旅程 "設計模式精通之旅" 價格為 1999 元
    And "Alice" 已購買該旅程
    When "Alice" 查看該旅程的詳細資訊
    Then 她應該看到「繼續學習」按鈕
    And 她不應該看到「立即加入課程」按鈕

  Scenario: 未購買旅程顯示「立即加入課程」按鈕
    Given 系統中存在一位用戶 "Bob" 密碼為 "Secure123!"
    And 系統中存在一個旅程 "Java 基礎課程" 價格為 2999 元
    And "Bob" 未購買該旅程
    When "Bob" 查看該旅程的詳細資訊
    Then 他應該看到「立即加入課程」按鈕
    And 他不應該看到「繼續學習」按鈕

  Scenario: 訪客查看旅程應顯示「立即加入課程」按鈕
    Given 系統中存在一個旅程 "Python 入門" 價格為 1599 元
    When 訪客查看該旅程的詳細資訊
    Then 應該顯示「立即加入課程」按鈕
    And 不應該顯示「繼續學習」按鈕

  # 場景組 B: 購買按鈕點擊行為

  Scenario: 已購買旅程點擊「繼續學習」導向第一個任務
    Given 系統中存在一位用戶 "Charlie" 密碼為 "Pass1234!"
    And 系統中存在一個旅程 "JavaScript 完整指南"
    And 該旅程包含一個章節 "第一章：JavaScript 基礎" 順序為 1
    And 該章節包含一個影片任務 "認識 JavaScript" 順序為 1
    And "Charlie" 已購買該旅程
    And "Charlie" 查看該旅程的詳細資訊
    When "Charlie" 點擊「繼續學習」按鈕
    Then 系統應該導向任務 "認識 JavaScript" 的頁面

  Scenario: 未登入訪客點擊購買按鈕導向登入頁面
    Given 系統中存在一個旅程 "React 完整課程" 價格為 3999 元
    And 訪客查看該旅程的詳細資訊
    When 訪客點擊「立即加入課程」按鈕
    Then 系統應該導向登入頁面

  # 場景組 C: 訂單建立成功流程

  Scenario: 使用者成功建立訂單
    Given 系統中存在一位用戶 "Diana" 密碼為 "Diana123!"
    And 系統中存在一個旅程 "Vue.js 實戰課程" 價格為 2499 元
    And "Diana" 已登入
    And "Diana" 未購買該旅程
    And "Diana" 沒有該旅程的未付款訂單
    When "Diana" 點擊「立即加入課程」並建立訂單
    Then 訂單建立應該成功
    And 訂單狀態應為「未付款」
    And 訂單應顯示使用者名稱為 "Diana"
    And 訂單應顯示課程名稱為 "Vue.js 實戰課程"
    And 訂單應顯示課程價格為 2499 元
    And 訂單應顯示建立時間
    And 系統應該進入「完成支付」階段

  Scenario: 訂單編號格式正確
    Given 系統中存在一位用戶 "Eva" 密碼為 "Eva12345!" ID 為 42
    And 系統中存在一個旅程 "Node.js 後端開發" 價格為 3499 元
    And "Eva" 已登入
    When "Eva" 建立該旅程的訂單
    Then 訂單編號格式應符合「時間戳(10位) + 使用者ID + 隨機碼(5位)」
    And 訂單編號應包含使用者 ID "42"

  Scenario: 訂單價格鎖定建立當下的價格
    Given 系統中存在一位用戶 "Frank" 密碼為 "Frank123!"
    And 系統中存在一個旅程 "Go 語言入門" 價格為 1999 元
    And "Frank" 已登入
    And "Frank" 建立該旅程的訂單
    When 系統將該旅程的價格調整為 2999 元
    Then "Frank" 的訂單價格應保持為 1999 元
    And 新建立的訂單價格應為 2999 元

  # 場景組 D: 訂單建立錯誤處理

  Scenario: 使用者已購買該旅程無法建立訂單
    Given 系統中存在一位用戶 "Grace" 密碼為 "Grace123!"
    And 系統中存在一個旅程 "Docker 容器化技術" 價格為 2799 元
    And "Grace" 已購買該旅程
    When "Grace" 嘗試建立該旅程的訂單
    Then 訂單建立應該失敗
    And 系統應該提示「你已經購買此課程」

  Scenario: 使用者已有該旅程的未付款訂單時導向付款頁面
    Given 系統中存在一位用戶 "Henry" 密碼為 "Henry123!"
    And 系統中存在一個旅程 "Kubernetes 實戰" 價格為 3999 元
    And "Henry" 已有該旅程的未付款訂單
    When "Henry" 嘗試再次建立該旅程的訂單
    Then 系統應該導向該未付款訂單的付款頁面
    And 不應該建立新訂單

  Scenario: 使用者可以同時擁有多個不同旅程的未付款訂單
    Given 系統中存在一位用戶 "Ivy" 密碼為 "Ivy12345!"
    And 系統中存在一個旅程 "AWS 雲端服務" 價格為 4999 元
    And 系統中存在一個旅程 "Azure 雲端架構" 價格為 4599 元
    And "Ivy" 已有旅程 "AWS 雲端服務" 的未付款訂單
    When "Ivy" 建立旅程 "Azure 雲端架構" 的訂單
    Then 訂單建立應該成功
    And "Ivy" 應該同時擁有 2 個未付款訂單
