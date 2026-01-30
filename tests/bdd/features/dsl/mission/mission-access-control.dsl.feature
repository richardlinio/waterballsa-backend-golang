# Language: zh-TW
# DSL Layer (L1): Business Domain Language
# Source: Release-3-Spec.md - 3.5 課程權限控制

Feature: 任務權限控制 (Release 3.5)
  作為系統
  我需要根據使用者的購買狀態控制任務的存取權限
  以便保護付費內容並提供正確的課程體驗

  # 場景組 I: 購買後單元解鎖

  Scenario: 使用者購買旅程後所有任務解鎖
    Given 系統中存在一位用戶 "Alice" 密碼為 "Test1234!"
    And 系統中存在一個旅程 "設計模式精通之旅"
    And 該旅程包含一個章節 "進階課程"
    And 該章節包含一個影片任務 "SOLID 原則" 存取權限為 "需購買"
    And 該章節包含一個影片任務 "設計模式導論" 存取權限為 "需購買"
    And 該章節包含一個影片任務 "工廠模式" 存取權限為 "需購買"
    And "Alice" 已購買該旅程
    When "Alice" 查看該旅程的詳細資訊
    Then 任務 "SOLID 原則" 應顯示為已解鎖
    And 任務 "設計模式導論" 應顯示為已解鎖
    And 任務 "工廠模式" 應顯示為已解鎖
    And 所有任務不應該顯示鎖定圖示

  Scenario: 使用者購買旅程後可以觀看所有需購買任務
    Given 系統中存在一位用戶 "Bob" 密碼為 "Secure123!"
    And 系統中存在一個旅程 "Java 基礎課程"
    And 該旅程包含一個章節 "物件導向程式設計"
    And 該章節包含一個影片任務 "類別與物件" 存取權限為 "需購買" 長度為 300 秒
    And "Bob" 已購買該旅程
    When "Bob" 點擊任務 "類別與物件"
    Then 應該顯示影片播放器
    And 影片應該可以正常播放

  # 場景組 J: 未購買課程的鎖定顯示

  Scenario: 未購買旅程的需購買任務顯示鎖定圖示
    Given 系統中存在一位用戶 "Charlie" 密碼為 "Pass1234!"
    And 系統中存在一個旅程 "Python 入門"
    And 該旅程包含一個章節 "進階主題"
    And 該章節包含一個影片任務 "裝飾器" 存取權限為 "需購買"
    And 該章節包含一個影片任務 "生成器" 存取權限為 "需購買"
    And 該章節包含一個影片任務 "上下文管理器" 存取權限為 "需購買"
    And "Charlie" 未購買該旅程
    When "Charlie" 查看該旅程的詳細資訊
    Then 任務 "裝飾器" 應該顯示鎖定圖示
    And 任務 "生成器" 應該顯示鎖定圖示
    And 任務 "上下文管理器" 應該顯示鎖定圖示

  Scenario: 未購買旅程點擊需購買任務顯示付費提示
    Given 系統中存在一位用戶 "Diana" 密碼為 "Diana123!"
    And 系統中存在一個旅程 "JavaScript 完整指南"
    And 該旅程包含一個章節 "非同步程式設計"
    And 該章節包含一個影片任務 "Promise 詳解" 存取權限為 "需購買" 長度為 400 秒
    And "Diana" 未購買該旅程
    When "Diana" 點擊任務 "Promise 詳解"
    Then 應該顯示「此為付費內容，請先購買課程」提示訊息
    And 不應該顯示影片播放器

  Scenario: 未登入訪客查看需購買任務顯示鎖定圖示
    Given 系統中存在一個旅程 "React 完整課程"
    And 該旅程包含一個章節 "React Hooks"
    And 該章節包含一個影片任務 "useState 詳解" 存取權限為 "需購買"
    And 該章節包含一個影片任務 "useEffect 詳解" 存取權限為 "需購買"
    When 訪客查看該旅程的詳細資訊
    Then 任務 "useState 詳解" 應該顯示鎖定圖示
    And 任務 "useEffect 詳解" 應該顯示鎖定圖示

  Scenario: 未登入訪客點擊需購買任務顯示付費提示
    Given 系統中存在一個旅程 "Vue.js 實戰課程"
    And 該旅程包含一個章節 "Vuex 狀態管理"
    And 該章節包含一個影片任務 "Vuex 核心概念" 存取權限為 "需購買" 長度為 350 秒
    When 訪客點擊任務 "Vuex 核心概念"
    Then 應該顯示「此為付費內容，請先購買課程」提示訊息
    And 不應該顯示影片播放器
