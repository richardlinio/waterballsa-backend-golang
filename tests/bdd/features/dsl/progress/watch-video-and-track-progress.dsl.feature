# Language: zh-TW
# DSL Layer (L1): Business Domain Language
# Source: Release-X-Spec.md - 2.3 影片觀看與進度追蹤

Feature: 影片觀看與進度追蹤 (Release 2.3)
  作為一個學習者
  我想要觀看課程影片並自動記錄進度
  以便下次可以從上次觀看的位置繼續學習

  Scenario: 首次觀看影片任務應從頭開始
    Given 系統中存在一個旅程 "Java 基礎課程"
    And 該旅程包含一個章節 "第一章：變數與型別"
    And 該章節包含一個影片任務 "認識變數" 長度為 100 秒
    And 使用者 "Alice" 尚未觀看過該影片任務
    When "Alice" 開始觀看該影片任務
    Then 觀看位置應該從 0 秒開始

  Scenario: 繼續觀看影片應從上次位置開始
    Given 系統中存在一個旅程 "Java 基礎課程"
    And 該旅程包含一個章節 "第一章：變數與型別"
    And 該章節包含一個影片任務 "認識變數" 長度為 100 秒
    And 使用者 "Bob" 之前觀看該影片到 30 秒
    When "Bob" 再次觀看該影片任務
    Then 觀看位置應該從 30 秒開始

  Scenario: 觀看完成後應自動標記為已完成
    Given 系統中存在一個旅程 "Java 基礎課程"
    And 該旅程包含一個章節 "第一章：變數與型別"
    And 該章節包含一個影片任務 "認識變數" 長度為 100 秒
    And 使用者 "Charlie" 正在觀看該影片任務
    When "Charlie" 的觀看進度達到 100 秒
    Then 該任務狀態應該自動變更為 "已完成"

  Scenario: 播放中每 10 秒應記錄觀看進度
    Given 系統中存在一個旅程 "Java 基礎課程"
    And 該旅程包含一個章節 "第一章：變數與型別"
    And 該章節包含一個影片任務 "認識變數" 長度為 100 秒
    And 使用者 "Diana" 正在觀看該影片任務
    When 影片播放 10 秒後
    Then 系統應該記錄 "Diana" 的觀看進度為 10 秒

  Scenario: 暫停影片時應記錄觀看進度
    Given 系統中存在一個旅程 "Java 基礎課程"
    And 該旅程包含一個章節 "第一章：變數與型別"
    And 該章節包含一個影片任務 "認識變數" 長度為 100 秒
    And 使用者 "Eve" 正在觀看該影片任務
    And "Eve" 已觀看到 25 秒
    When "Eve" 暫停影片
    Then 系統應該記錄 "Eve" 的觀看進度為 25 秒

  Scenario: 關閉頁面時應記錄觀看進度
    Given 系統中存在一個旅程 "Java 基礎課程"
    And 該旅程包含一個章節 "第一章：變數與型別"
    And 該章節包含一個影片任務 "認識變數" 長度為 100 秒
    And 使用者 "Frank" 正在觀看該影片任務
    And "Frank" 已觀看到 45 秒
    When "Frank" 關閉頁面
    Then 系統應該記錄 "Frank" 的觀看進度為 45 秒

  Scenario: 已完成的影片任務可以重複觀看
    Given 系統中存在一個旅程 "Java 基礎課程"
    And 該旅程包含一個章節 "第一章：變數與型別"
    And 該章節包含一個影片任務 "認識變數" 長度為 100 秒
    And 使用者 "Grace" 已完成該影片任務
    When "Grace" 再次觀看該影片任務
    Then 應該允許 "Grace" 重複觀看
    And 該任務狀態應該仍為 "已完成"

  Scenario: 觀看進度精確達到 100% 時應標記為已完成
    Given 系統中存在一個旅程 "Java 基礎課程"
    And 該旅程包含一個章節 "第一章：變數與型別"
    And 該章節包含一個影片任務 "認識變數" 長度為 120 秒
    And 使用者 "Henry" 正在觀看該影片任務
    When "Henry" 的觀看進度達到 120 秒
    Then 該任務狀態應該變更為 "已完成"

  Scenario: 觀看進度未達 100% 時應保持未完成狀態
    Given 系統中存在一個旅程 "Java 基礎課程"
    And 該旅程包含一個章節 "第一章：變數與型別"
    And 該章節包含一個影片任務 "認識變數" 長度為 120 秒
    And 使用者 "Iris" 正在觀看該影片任務
    When "Iris" 的觀看進度為 119 秒
    Then 該任務狀態應該保持為 "未完成"
