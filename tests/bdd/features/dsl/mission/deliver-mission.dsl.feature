# Language: zh-TW
# DSL Layer (L1): Business Domain Language
# Source: Release-2-Spec.md - 2.4 交付功能

Feature: 任務交付 (Release 2.4)
  作為一個學員
  我想要在完成任務後交付任務
  以便獲得經驗值獎勵

  Scenario: 影片任務完成後成功交付並獲得經驗值
    Given 系統中存在一位用戶 "Alice" 密碼為 "Test1234!"
    And 系統中存在一個旅程 "Java 基礎課程"
    And 旅程包含一個章節 "第一章：變數與型別"
    And 章節包含一個影片任務 "認識變數" 長度為 100 秒
    And 該任務的獎勵為 100 經驗值
    And "Alice" 已購買該旅程
    And "Alice" 已將影片任務 "認識變數" 觀看至 100%
    When "Alice" 交付任務 "認識變數"
    Then 交付應該成功
    And 任務 "認識變數" 的狀態應變更為 "已交付"
    And "Alice" 的經驗值應增加 100

  Scenario: 嘗試交付未完成的任務應失敗
    Given 系統中存在一位用戶 "Bob" 密碼為 "Secure123!"
    And 系統中存在一個旅程 "Python 入門"
    And 旅程包含一個章節 "第一章：基礎語法"
    And 章節包含一個影片任務 "Hello World" 長度為 100 秒
    And 該任務的獎勵為 100 經驗值
    And "Bob" 已購買該旅程
    And "Bob" 已將影片任務 "Hello World" 觀看至 50%
    When "Bob" 嘗試交付任務 "Hello World"
    Then 交付應該失敗
    And 系統應該提示 "任務尚未完成"
    And 任務 "Hello World" 的狀態應保持為 "未完成"
    And "Bob" 的經驗值應保持為 0

  Scenario: 嘗試重複交付同一任務應失敗
    Given 系統中存在一位用戶 "Charlie" 密碼為 "Pass1234!"
    And 系統中存在一個旅程 "JavaScript 完整指南"
    And 旅程包含一個章節 "第一章：JavaScript 基礎"
    And 章節包含一個影片任務 "變數宣告" 長度為 120 秒
    And 該任務的獎勵為 100 經驗值
    And "Charlie" 已購買該旅程
    And "Charlie" 已將影片任務 "變數宣告" 觀看至 100%
    And "Charlie" 已交付任務 "變數宣告"
    When "Charlie" 嘗試再次交付任務 "變數宣告"
    Then 交付應該失敗
    And 系統應該提示 "任務已交付"
    And 任務 "變數宣告" 的狀態應保持為 "已交付"
    And "Charlie" 的經驗值應保持為 100

  Scenario: 未登入使用者無法交付任務
    Given 系統中存在一個旅程 "React 完整課程"
    And 旅程包含一個章節 "第一章：React 基礎"
    And 章節包含一個影片任務 "認識 React" 長度為 150 秒
    And 該任務的獎勵為 100 經驗值
    When 未登入的使用者嘗試交付任務 "認識 React"
    Then 交付應該失敗
    And 系統應該提示 "未授權或權杖無效"