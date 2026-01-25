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

## Scaffold 檢查清單

實作順序 (依賴關係):

1.  **Domain Models** (無依賴) - `internal/model/`
2.  **DTOs** (依賴 models) - `internal/dto/`
3.  **SQLc Queries** (無依賴) - `internal/db/queries/*.sql`
4.  **錯誤碼** (可提前) - `internal/apperror/`
5.  **Repository** (依賴 SQLc) - `internal/repository/`
6.  **Service** (依賴 repository) - `internal/service/`
7.  **Handler** (依賴 service) - `internal/handler/`
8.  **路由註冊** (依賴 handler) - `internal/router/router.go`
9.  **依賴注入** (依賴所有元件)
   - `internal/app/app.go` - 生產環境
   - `tests/testutil/server.go` - 測試環境 ⚠️ **必須同步更新**

必檢項目:

- [ ] 所有層級已建立
- [ ] 錯誤處理已加入 (參考 `/add-error-handling`)
- [ ] 授權檢查已實作 (如需要)
- [ ] 路由已註冊
- [ ] 依賴注入已配置 (`internal/app/app.go`)
- [ ] **測試伺服器已更新** (`tests/testutil/server.go`) ⚠️
- [ ] Migration 已建立 (如需要)
- [ ] SQLc 已執行 (`make sqlc`)
- [ ] 測試步驟已定義
- [ ] 遵循 Go idiomatic patterns (參考 `/go-idiomatic`)
