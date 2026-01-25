# 測試與驗證 (Testing & Validation)

本文件提供 BDD/ISA 測試範例和功能開發的檢查清單。

## ISA Feature 測試範例

```gherkin
@isa
Feature: 功能名稱

  Scenario: 場景描述
    # 資料準備
    Given the database has a user:
      | username | Alice |
      | password | Test1234! |

    # 認證
    When I send "POST" request to "/auth/login"
    And I store the response field "accessToken" as "accessToken"
    Given I set Authorization header to "{{accessToken}}"

    # 執行動作
    When I send "GET" request to "/users/1/progress"

    # 驗證回應
    Then the response status code should be 200
    And the response body field "status" should equal string "UNCOMPLETED"
```

