# Language: zh-TW
# DSL Layer (L1): Business Domain Language
# Source: Release-2-Spec.md - 2.1 課程列表與瀏覽

Feature: 課程列表與瀏覽 (Release 2.1)
  作為一個訪客或使用者
  我想要瀏覽所有可用的課程
  以便了解平台提供的學習內容並選擇感興趣的課程

  Scenario: 訪客瀏覽包含多門課程的列表
    Given 系統中存在一個旅程 "軟體設計模式精通之旅" 由 "水球潘" 教授，價格為 1999 元
    And 該旅程的簡介為 "用 C.A. 模式大大提昇系統思維能力"
    And 該旅程的封面圖為 "https://example.com/cover1.jpg"
    And 系統中存在一個旅程 "Web 開發入門" 由 "林老師" 教授，價格為 999 元
    And 該旅程的簡介為 "從零開始學習網頁開發"
    And 該旅程的封面圖為 "https://example.com/cover2.jpg"
    When 訪客查看課程列表
    Then 她應該看到 2 門課程
    And 課程 "軟體設計模式精通之旅" 應該顯示老師名稱 "水球潘"
    And 課程 "軟體設計模式精通之旅" 應該顯示價格 1999 元
    And 課程 "軟體設計模式精通之旅" 應該顯示簡介 "用 C.A. 模式大大提昇系統思維能力"
    And 課程 "軟體設計模式精通之旅" 應該顯示封面圖 "https://example.com/cover1.jpg"
    And 課程 "Web 開發入門" 應該顯示老師名稱 "林老師"
    And 課程 "Web 開發入門" 應該顯示價格 999 元

  Scenario: 系統中沒有課程時顯示空列表
    Given 系統中不存在任何旅程
    When 訪客查看課程列表
    Then 她應該看到 0 門課程

  Scenario: 訪客瀏覽單一課程
    Given 系統中存在一個旅程 "Java 基礎課程" 由 "陳老師" 教授，價格為 1500 元
    And 該旅程的簡介為 "深入淺出學習 Java 程式設計"
    And 該旅程的封面圖為 "https://example.com/java-cover.jpg"
    When 訪客查看課程列表
    Then 她應該看到 1 門課程
    And 課程 "Java 基礎課程" 應該顯示老師名稱 "陳老師"
    And 課程 "Java 基礎課程" 應該顯示價格 1500 元
    And 課程 "Java 基礎課程" 應該顯示簡介 "深入淺出學習 Java 程式設計"
    And 課程 "Java 基礎課程" 應該顯示封面圖 "https://example.com/java-cover.jpg"
