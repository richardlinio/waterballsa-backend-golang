# Language: zh-TW
# DSL Layer (L1): Business Domain Language
# Source: Release-1-Spec.md - 2.5 經驗值系統

Feature: 經驗值系統 (Release 1 - 2.5)
  作為一個學員
  我想要查詢自己的經驗值
  以便追蹤我的學習成長和在未來的排行榜中比較

  Scenario: 新註冊使用者的初始經驗值應為 0
    Given 系統中存在一位用戶 "Alice" 密碼為 "Test1234!"
    When "Alice" 查詢自己的經驗值
    Then 她的經驗值應該為 0

  Scenario: 使用者交付任務後經驗值應增加
    Given 系統中存在一位用戶 "Bob" 密碼為 "Secure123!"
    And 系統中存在一個旅程 "Java 基礎課程"
    And 旅程包含一個章節 "第一章：變數與型別"
    And 章節包含一個影片任務 "認識變數" 長度為 100 秒
    And 該任務的獎勵為 100 經驗值
    And "Bob" 已購買該旅程
    And "Bob" 已將影片任務 "認識變數" 觀看至 100%
    And "Bob" 已交付任務 "認識變數"
    When "Bob" 查詢自己的經驗值
    Then 他的經驗值應該為 100

  Scenario: 未登入使用者無法查詢經驗值
    When 未登入的使用者嘗試查詢經驗值
    Then 系統應該提示 "未授權或權杖無效"

  Scenario: 使用者完成任務但未交付時經驗值保持為 0
    Given 系統中存在一位用戶 "Diana" 密碼為 "Diana123!"
    And 系統中存在一個旅程 "JavaScript 入門"
    And 旅程包含一個章節 "第一章：基礎語法"
    And 章節包含一個影片任務 "認識 JavaScript" 長度為 120 秒
    And 該任務的獎勵為 100 經驗值
    And "Diana" 已購買該旅程
    And "Diana" 已將影片任務 "認識 JavaScript" 觀看至 100%
    When "Diana" 查詢自己的經驗值
    Then 她的經驗值應該為 0
