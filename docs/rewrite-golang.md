# WaterBall SA Backend - Golang 重寫計畫

> **目標**：將現有 Java Spring Boot 後端重寫為 Golang，作為技術作品集，學習 Golang 後端開發
> **時程**：預計 7-14 天
> **策略**：BDD 驅動開發（復用現有 Cucumber 測試）
> **原則**：保持簡單，基礎優先，功能對齊

---

## 📊 專案現況

- **代碼量**：4,820 行 Java（97 個文件）
- **測試覆蓋**：4,152 行 Gherkin（28 個 .feature 文件）
- **技術棧**：Spring Boot 3.5.7 + PostgreSQL + JPA + JWT
- **功能模塊**：認證、學習旅程、進度追蹤、訂單系統、用戶管理
- **數據庫**：21 個 Liquibase migration files

---

## 🎯 重寫目標

### 功能目標

- ✅ 所有 API endpoints 功能對齊（基於 swagger.yaml）
- ✅ 所有 28 個 BDD 測試通過
- ✅ 資料庫 schema 完全一致
- ✅ JWT 認證授權機制相同

### 技術目標

- ✅ 學習 Golang Web 框架（Gin）
- ✅ 學習 GORM ORM
- ✅ 學習 Golang 測試（godog BDD）
- ✅ 理解 Golang 錯誤處理模式
- ✅ 掌握依賴管理與專案結構

### 非目標（避免過度設計）

- ❌ 不追求極致性能優化
- ❌ 不使用複雜的 goroutine/channel（除非必要）
- ❌ 不重新設計 API 或資料庫 schema
- ❌ 不實現 Java 版本沒有的新功能

---

## 🗂️ 技術選型

| 類別           | Java 版本              | Golang 版本                     | 說明                       |
| -------------- | ---------------------- | ------------------------------- | -------------------------- |
| **Web 框架**   | Spring Boot            | **Gin**                         | 高性能、文檔完善、社群活躍 |
| **ORM**        | JPA/Hibernate          | **GORM**                        | 最接近 JPA 的 Golang ORM   |
| **資料庫**     | PostgreSQL             | PostgreSQL                      | 不變                       |
| **JWT**        | io.jsonwebtoken        | **golang-jwt/jwt v5**           | 官方推薦                   |
| **驗證**       | Spring Validation      | **go-playground/validator v10** | 最流行的驗證庫             |
| **配置管理**   | application.properties | **viper**                       | 支援多格式配置             |
| **資料庫遷移** | Liquibase              | **golang-migrate**              | 直接復用 SQL 文件          |
| **BDD 測試**   | Cucumber JVM           | **godog**                       | Golang 官方 Cucumber 實現  |
| **單元測試**   | JUnit 5                | **testify/assert**              | 標準測試庫                 |
| **HTTP 測試**  | REST Assured           | **httptest**                    | Go 原生 HTTP 測試          |
| **密碼加密**   | BCrypt (Spring)        | **golang.org/x/crypto/bcrypt**  | Go 官方實現                |
| **限流**       | Bucket4j               | **gin middleware**              | 自實現或用 gin-limiter     |
| **日誌**       | Logback                | **logrus** 或 **zap**           | zap 更高性能               |
| **定時任務**   | @Scheduled             | **gocron**                      | 類似 cron 的調度器         |

---

## 📁 專案結構設計

```
waterballsa-backend-go/
├── cmd/
│   └── server/
│       └── main.go              # 應用程式入口點
│
├── internal/                     # 私有應用程式代碼
│   ├── config/
│   │   └── config.go            # 配置管理（viper）
│   │
│   ├── model/                   # Domain entities（對應 Java entity）
│   │   ├── user.go
│   │   ├── journey.go
│   │   ├── mission.go
│   │   ├── order.go
│   │   └── ...
│   │
│   ├── dto/                     # Request/Response DTOs
│   │   ├── auth.go
│   │   ├── journey.go
│   │   └── ...
│   │
│   ├── repository/              # 資料訪問層（對應 Java repository）
│   │   ├── user_repository.go
│   │   ├── journey_repository.go
│   │   └── ...
│   │
│   ├── service/                 # 業務邏輯層（對應 Java service）
│   │   ├── auth_service.go
│   │   ├── journey_service.go
│   │   ├── order_service.go
│   │   └── ...
│   │
│   ├── handler/                 # HTTP handlers（對應 Java controller）
│   │   ├── health_handler.go
│   │   ├── auth_handler.go
│   │   ├── journey_handler.go
│   │   └── ...
│   │
│   ├── middleware/              # HTTP middlewares
│   │   ├── auth.go              # JWT 驗證
│   │   ├── cors.go              # CORS
│   │   ├── rate_limit.go        # 限流
│   │   └── error_handler.go     # 錯誤處理
│   │
│   ├── validator/               # 業務驗證器（對應 Java validator）
│   │   ├── order_validator.go
│   │   └── mission_access_validator.go
│   │
│   ├── util/                    # 工具函數
│   │   ├── jwt.go               # JWT 生成/驗證
│   │   ├── password.go          # 密碼加密
│   │   └── order_number.go      # 訂單號生成
│   │
│   └── scheduler/               # 定時任務
│       └── order_scheduler.go   # 訂單過期處理
│
├── migrations/                   # 資料庫遷移文件
│   ├── 000001_create_users_table.up.sql
│   ├── 000001_create_users_table.down.sql
│   └── ...（從 Java 版本複製過來）
│
├── test/
│   ├── features/                # BDD feature files（從 Java 版本複製）
│   │   ├── isa/
│   │   │   ├── auth/
│   │   │   ├── journey/
│   │   │   ├── mission/
│   │   │   └── order/
│   │   └── dsl/
│   │
│   └── steps/                   # godog step definitions（重新實現）
│       ├── common_steps.go      # 共用 steps
│       ├── database_steps.go    # 資料庫設置 steps
│       └── api_steps.go         # API 請求/驗證 steps
│
├── pkg/                         # 可復用的公開函式庫（如果需要）
│
├── config/
│   ├── config.yaml              # 開發環境配置
│   └── config.prod.yaml         # 生產環境配置
│
├── .env.example                 # 環境變數範例
├── .gitignore
├── go.mod                       # Go modules 依賴管理
├── go.sum
├── Makefile                     # 常用指令快捷鍵
├── Dockerfile                   # Docker 容器化
├── docker-compose.yml           # 本地開發環境（PostgreSQL）
└── README.md
```

### 與 Java 版本的對應關係

| Java                     | Golang                       |
| ------------------------ | ---------------------------- |
| `@RestController`        | `handler/*.go`               |
| `@Service`               | `service/*.go`               |
| `@Repository`            | `repository/*.go`            |
| `@Entity`                | `model/*.go`                 |
| DTO classes              | `dto/*.go`                   |
| `@Component` validators  | `validator/*.go`             |
| `application.properties` | `config/config.yaml` + viper |

---

## 🚀 詳細執行步驟

### **階段 0：專案初始化（1 天）**

#### Step 0.1：建立新的 GitHub Repository

```bash
# 在 GitHub 上建立新 repo：waterballsa-backend-go
# Clone 到本地
git clone https://github.com/YOUR_USERNAME/waterballsa-backend-go.git
cd waterballsa-backend-go
```

#### Step 0.2：初始化 Go Module

```bash
go mod init github.com/YOUR_USERNAME/waterballsa-backend-go
```

#### Step 0.3：安裝核心依賴

```bash
# Web 框架
go get -u github.com/gin-gonic/gin

# ORM
go get -u gorm.io/gorm
go get -u gorm.io/driver/postgres

# JWT
go get -u github.com/golang-jwt/jwt/v5

# 驗證
go get -u github.com/go-playground/validator/v10

# 配置管理
go get -u github.com/spf13/viper

# 密碼加密
go get -u golang.org/x/crypto/bcrypt

# 資料庫遷移
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# 測試相關
go get -u github.com/cucumber/godog/cmd/godog@latest
go get -u github.com/stretchr/testify/assert

# 定時任務
go get -u github.com/go-co-op/gocron

# 日誌
go get -u github.com/sirupsen/logrus

# 熱重載
go install github.com/air-verse/air@latest
```

#### Step 0.4：建立專案目錄結構

```bash
mkdir -p cmd/server
mkdir -p internal/{config,model,dto,repository,service,handler,middleware,validator,util,scheduler}
mkdir -p migrations
mkdir -p test/{features,steps}
mkdir -p config
```

#### Step 0.5：複製資源文件

**從 Java 專案複製：**

1. 複製 `src/test/resources/features/` → `test/features/`（28 個 .feature 文件）
2. 複製 `src/main/resources/db/changelog/migrations/*.sql` → `migrations/`
3. 參考 `docs/api-docs/swagger.yaml` 確認 API 規格

**轉換 Liquibase SQL 為 golang-migrate 格式：**

- 每個 SQL 文件需拆分為 `.up.sql` 和 `.down.sql`
- 命名格式：`000001_create_users_table.up.sql`

#### Step 0.6：建立基礎配置文件

**config/config.yaml：**

```yaml
server:
  port: 8080

database:
  host: localhost
  port: 5432
  user: postgres
  password: postgres
  name: waterballsa
  sslmode: disable

jwt:
  secret: 'your-secret-key-here'
  expiration: 86400 # 1 day in seconds

logging:
  level: debug
```

**.env.example：**

```env
SERVER_PORT=8080
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=waterballsa
JWT_SECRET=your-secret-key
```

**Makefile：**

```makefile
.PHONY: run test migrate-up migrate-down

run:
	go run cmd/server/main.go

test:
	go test ./...

bdd:
	godog run test/features

migrate-up:
	migrate -path migrations -database "postgres://postgres:postgres@localhost:5432/waterballsa?sslmode=disable" up

migrate-down:
	migrate -path migrations -database "postgres://postgres:postgres@localhost:5432/waterballsa?sslmode=disable" down

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down
```

**docker-compose.yml：**

```yaml
version: '3.8'
services:
  postgres:
    image: postgres:15-alpine
    environment:
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: postgres
      POSTGRES_DB: waterballsa
    ports:
      - '5432:5432'
    volumes:
      - postgres_data:/var/lib/postgresql/data

volumes:
  postgres_data:
```

#### Step 0.7：建立 main.go 骨架

```go
// cmd/server/main.go
package main

import (
    "log"
    "github.com/gin-gonic/gin"
)

func main() {
    r := gin.Default()

    r.GET("/healthz", func(c *gin.Context) {
        c.JSON(200, gin.H{"status": "ok"})
    })

    if err := r.Run(":8080"); err != nil {
        log.Fatal("Failed to start server:", err)
    }
}
```

#### Step 0.8：驗證環境

```bash
# 啟動 PostgreSQL
make docker-up

# 執行資料庫遷移（先轉換 SQL 文件）
make migrate-up

# 啟動應用
make run

# 測試 health endpoint
curl http://localhost:8080/healthz
```

**完成標準：**

- ✅ 專案目錄結構完整
- ✅ Go modules 初始化完成
- ✅ PostgreSQL 啟動成功
- ✅ `/healthz` endpoint 可訪問

---

### **階段 1：基礎架構層（1-2 天）**

#### Step 1.1：實作配置管理（internal/config/config.go）

**需要實現：**

- 使用 viper 讀取 `config/config.yaml`
- 支援環境變數覆蓋（如 `DB_HOST`）
- 定義 Config 結構體

**參考：**

- Java: `application.properties`

#### Step 1.2：實作資料庫連接（internal/config/database.go）

**需要實現：**

- GORM 連接 PostgreSQL
- 連接池配置
- 日誌模式（開發環境顯示 SQL）

**參考：**

- Java: `spring.datasource.url`

#### Step 1.3：實作錯誤處理（internal/middleware/error_handler.go）

**需要實現：**

- 統一錯誤響應格式（對應 `ErrorResponse` in swagger.yaml）
- 自定義錯誤類型（如 `NotFoundError`, `UnauthorizedError`）

**參考：**

- Java: `@ControllerAdvice` 異常處理器

#### Step 1.4：實作日誌系統

**需要實現：**

- 整合 logrus
- 設置日誌格式和級別

#### Step 1.5：更新 main.go 整合所有基礎設施

**需要實現：**

- 初始化配置
- 初始化資料庫連接
- 設置中間件（CORS, Logger, ErrorHandler）
- Graceful shutdown

**完成標準：**

- ✅ 配置可從文件和環境變數讀取
- ✅ 資料庫連接成功
- ✅ 日誌輸出正常
- ✅ 錯誤處理中間件運作

---

### **階段 2：認證模塊（1 天）**

#### Step 2.1：實作 User Entity（internal/model/user.go）

**需要實現：**

```go
type User struct {
    ID         int64     `gorm:"primaryKey;autoIncrement"`
    Username   string    `gorm:"unique;not null"`
    Password   string    `gorm:"not null"`  // BCrypt hash
    Experience int       `gorm:"default:0"`
    CreatedAt  time.Time `gorm:"autoCreateTime"`
    UpdatedAt  time.Time `gorm:"autoUpdateTime"`
}
```

**參考：**

- Java: `waterballsa/entity/User.java`

#### Step 2.2：實作 AccessToken Entity（internal/model/access_token.go）

**參考：**

- Java: `waterballsa/entity/AccessToken.java`

#### Step 2.3：實作 Repository 層

**需要實現：**

- `internal/repository/user_repository.go`
  - `FindByUsername(username string) (*User, error)`
  - `Create(user *User) error`
  - `FindByID(id int64) (*User, error)`

**參考：**

- Java: `waterballsa/repository/UserRepository.java`

#### Step 2.4：實作 JWT 工具（internal/util/jwt.go）

**需要實現：**

- `GenerateToken(userId int64, username string) (string, error)`
- `ValidateToken(tokenString string) (*Claims, error)`
- Claims 結構體（包含 userId, username, exp）

**參考：**

- Java: JWT 配置和生成邏輯

#### Step 2.5：實作密碼工具（internal/util/password.go）

**需要實現：**

- `HashPassword(password string) (string, error)`
- `VerifyPassword(hashedPassword, password string) error`

**參考：**

- Java: Spring Security BCrypt

#### Step 2.6：實作 Auth Service（internal/service/auth_service.go）

**需要實現：**

- `Register(username, password string) (*User, error)`
- `Login(username, password string) (string, *User, error)` // 返回 token
- `Logout(token string) error`

**參考：**

- Java: `waterballsa/service/AuthService.java`

#### Step 2.7：實作 Auth DTO（internal/dto/auth.go）

**需要實現：**

```go
type RegisterRequest struct {
    Username string `json:"username" binding:"required,min=3,max=50"`
    Password string `json:"password" binding:"required,min=8"`
}

type LoginRequest struct {
    Username string `json:"username" binding:"required"`
    Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
    AccessToken string   `json:"accessToken"`
    User        UserInfo `json:"user"`
}

type UserInfo struct {
    ID         int64  `json:"id"`
    Username   string `json:"username"`
    Experience int    `json:"experience"`
}
```

**參考：**

- Swagger: `openapi/schemas/auth.yaml`

#### Step 2.8：實作 Auth Handler（internal/handler/auth_handler.go）

**需要實現：**

- `POST /auth/register`
- `POST /auth/login`
- `POST /auth/logout`

**參考：**

- Java: `waterballsa/controller/AuthController.java`
- Swagger: `openapi/paths/auth.yaml`

#### Step 2.9：實作 JWT Middleware（internal/middleware/auth.go）

**需要實現：**

- 從 `Authorization: Bearer <token>` header 提取 token
- 驗證 token
- 將 userId 存入 Gin context

#### Step 2.10：更新 main.go 註冊認證路由

#### Step 2.11：實作認證模塊的 BDD 測試（test/steps/）

**需要實現：**

- 實作 godog step definitions
- 跑通 `test/features/isa/auth/*.feature` 測試

**完成標準：**

- ✅ 所有認證相關 API 實現
- ✅ `test/features/isa/auth/*.feature` 所有測試通過
- ✅ JWT 生成和驗證正確

---

### **階段 3：學習旅程模塊（2 天）**

#### Step 3.1：實作 Entities

**需要實現：**

- `internal/model/journey.go`（對應 Journey entity）
- `internal/model/chapter.go`（對應 Chapter entity）
- `internal/model/mission.go`（對應 Mission entity）
- `internal/model/mission_resource.go`（對應 MissionResource entity）
- `internal/model/enums.go`（MissionType, ResourceType, MissionAccessLevel）

**GORM 關聯關係：**

- Journey has many Chapters
- Chapter has many Missions
- Mission has many MissionResources

**參考：**

- Java entities: `waterballsa/entity/Journey.java`, `Mission.java` 等

#### Step 3.2：實作 Repository 層

**需要實現：**

- `internal/repository/journey_repository.go`
  - `FindAll() ([]Journey, error)`
  - `FindByID(id int64) (*Journey, error)` with preload
- `internal/repository/mission_repository.go`
  - `FindByID(id int64) (*Mission, error)` with preload

**參考：**

- Java: `JourneyRepository.java`, `MissionRepository.java`

#### Step 3.3：實作 Service 層

**需要實現：**

- `internal/service/journey_service.go`
  - `GetAllJourneys(userId int64) ([]JourneyListItem, error)`
  - `GetJourneyDetail(journeyId, userId int64) (*JourneyDetail, error)`
- `internal/service/mission_service.go`
  - `GetMissionDetail(journeyId, missionId, userId int64) (*MissionDetail, error)`

**注意：**

- 需要根據 userId 判斷旅程是否已購買（查詢 user_journeys 表）
- 需要根據 mission access level 判斷可見性

**參考：**

- Java: `JourneyService.java`, `MissionService.java`

#### Step 3.4：實作 DTOs

**需要實現：**

- `internal/dto/journey.go`
  - `JourneyListItem`
  - `JourneyDetail`
  - `ChapterDTO`
  - `MissionSummary`
- `internal/dto/mission.go`
  - `MissionDetail`
  - `MissionReward`
  - `MissionResource`

**參考：**

- Swagger: `openapi/schemas/journeys.yaml`, `missions.yaml`

#### Step 3.5：實作 Handlers

**需要實現：**

- `internal/handler/journey_handler.go`
  - `GET /journeys`
  - `GET /journeys/{journeyId}`
- `internal/handler/mission_handler.go`
  - `GET /journeys/{journeyId}/missions/{missionId}`

**參考：**

- Java: `JourneyController.java`, `MissionController.java`
- Swagger: `openapi/paths/journeys.yaml`, `missions.yaml`

#### Step 3.6：實作 Mission Access Validator

**需要實現：**

- `internal/validator/mission_access_validator.go`
  - 驗證用戶是否有權訪問任務（基於 access level 和購買狀態）

**參考：**

- Java: `waterballsa/validator/MissionAccessValidator.java`

#### Step 3.7：更新路由

#### Step 3.8：BDD 測試

**跑通：**

- `test/features/isa/journey/*.feature`
- `test/features/isa/mission/*.feature`

**完成標準：**

- ✅ 所有旅程和任務 API 實現
- ✅ 相關 BDD 測試通過
- ✅ 權限控制正確

---

### **階段 4：進度追蹤模塊（1 天）**

#### Step 4.1：實作 Entities

**需要實現：**

- `internal/model/user_mission_progress.go`
- `internal/model/enums.go` 添加 `ProgressStatus` enum

**參考：**

- Java: `UserMissionProgress.java`

#### Step 4.2：實作 Repository

**需要實現：**

- `internal/repository/progress_repository.go`
  - `FindByUserAndMission(userId, missionId int64) (*UserMissionProgress, error)`
  - `Upsert(progress *UserMissionProgress) error`

#### Step 4.3：實作 Service

**需要實現：**

- `internal/service/progress_service.go`
  - `GetProgress(userId, missionId int64) (*UserMissionProgressResponse, error)`
  - `UpdateProgress(userId, missionId int64, currentSecond int) error`
  - `DeliverMission(userId, missionId int64) (*DeliverResponse, error)`

**業務邏輯：**

- 更新進度
- 完成任務時授予經驗值
- 更新 user.experience

**參考：**

- Java: `ProgressService.java`

#### Step 4.4：實作 DTOs

**需要實現：**

- `internal/dto/progress.go`

**參考：**

- Swagger: `openapi/schemas/missions.yaml`

#### Step 4.5：實作 Handler

**需要實現：**

- `internal/handler/progress_handler.go`
  - `GET /users/{userId}/missions/{missionId}/progress`
  - `POST /users/{userId}/missions/{missionId}/progress/deliver`

**參考：**

- Java: `ProgressController.java`

#### Step 4.6：BDD 測試

**跑通：**

- `test/features/isa/progress/*.feature`

**完成標準：**

- ✅ 進度追蹤和任務完成功能實現
- ✅ 經驗值授予正確
- ✅ BDD 測試通過

---

### **階段 5：訂單系統（2 天）**

#### Step 5.1：實作 Entities

**需要實現：**

- `internal/model/order.go`
- `internal/model/order_item.go`
- `internal/model/user_journey.go`
- `internal/model/enums.go` 添加 `OrderStatus` enum

**GORM 關聯關係：**

- Order has many OrderItems

**參考：**

- Java: `Order.java`, `OrderItem.java`, `UserJourney.java`

#### Step 5.2：實作 Repository

**需要實現：**

- `internal/repository/order_repository.go`
  - `Create(order *Order) error`
  - `FindByID(orderId int64) (*Order, error)` with preload items
  - `FindByIDAndUserID(orderId, userId int64) (*Order, error)`
  - `FindByUserIDAndStatusAndJourneyID(...) (*Order, error)`
  - `FindByIDAndUserIDForUpdate(orderId, userId int64) (*Order, error)` // 悲觀鎖
  - `FindByStatusAndExpiredAtBefore(status OrderStatus, time time.Time) ([]Order, error)`
- `internal/repository/user_journey_repository.go`
  - `Create(userJourney *UserJourney) error`
  - `FindByUserIDAndJourneyID(userId, journeyId int64) (*UserJourney, error)`

**重點：悲觀鎖實現：**

```go
// GORM pessimistic lock (FOR UPDATE)
db.Clauses(clause.Locking{Strength: "UPDATE"}).
   Where("id = ? AND user_id = ?", orderId, userId).
   First(&order)
```

**參考：**

- Java: `OrderRepository.java`（特別注意 `@Lock(PESSIMISTIC_WRITE)` 的翻譯）

#### Step 5.3：實作 Order Validator

**需要實現：**

- `internal/validator/order_validator.go`
  - `ValidateOrderRequest(req CreateOrderRequest) error`
  - `ValidateJourneyNotPurchased(userId, journeyId int64) error`
  - `ValidateOrderNotPaid(order *Order) error`
  - `ValidateOrderNotExpired(order *Order) error`

**參考：**

- Java: `OrderValidator.java`

#### Step 5.4：實作 Service

**需要實現：**

- `internal/service/order_service.go`
  - `CreateOrder(userId int64, req CreateOrderRequest) (*OrderResponse, bool, error)` // bool 表示是否為新訂單
  - `GetOrderDetail(orderId, userId int64) (*OrderResponse, error)`
  - `PayOrder(orderId, userId int64) (*PayOrderResponse, error)`

**業務邏輯：**

- 創建訂單時鎖定價格
- 冪等性：已有未付款訂單則返回現有訂單
- 支付時使用悲觀鎖
- 支付成功後寫入 user_journeys 表

**參考：**

- Java: `OrderService.java`

#### Step 5.5：實作 Order Number Generator

**需要實現：**

- `internal/util/order_number.go`
  - 生成格式：`ORD{timestamp}{userId}{random}`

**參考：**

- Java: `OrderNumberGenerator.java`

#### Step 5.6：實作 DTOs

**需要實現：**

- `internal/dto/order.go`

**參考：**

- Swagger: `openapi/schemas/orders.yaml`

#### Step 5.7：實作 Handler

**需要實現：**

- `internal/handler/order_handler.go`
  - `POST /orders`
  - `GET /orders/{orderId}`
  - `POST /orders/{orderId}/action/pay`
- `internal/handler/user_handler.go`
  - `GET /users/{userId}/orders`（分頁）
  - `GET /users/{userId}/journeys`
  - `GET /users/me`

**參考：**

- Java: `OrderController.java`, `UserController.java`

#### Step 5.8：實作定時任務（訂單過期）

**需要實現：**

- `internal/scheduler/order_scheduler.go`
  - 使用 gocron
  - 每 10 分鐘執行一次
  - 將超過 3 天未支付的訂單標記為 EXPIRED

**參考：**

- Java: `@Scheduled(cron = "0 */10 * * * *")` 在 `OrderService.java`

#### Step 5.9：整合定時任務到 main.go

#### Step 5.10：BDD 測試

**跑通：**

- `test/features/isa/order/*.feature`
- `test/features/isa/user/*.feature`

**完成標準：**

- ✅ 所有訂單 API 實現
- ✅ 悲觀鎖正確運作
- ✅ 定時任務正確過期訂單
- ✅ BDD 測試通過

---

### **階段 6：限流與最後整合（1 天）**

#### Step 6.1：實作限流 Middleware

**需要實現：**

- `internal/middleware/rate_limit.go`
  - 基於 IP 或 userId 的限流
  - 可使用簡單的 in-memory map + mutex（或整合第三方庫）

**參考：**

- Java: Bucket4j 配置

#### Step 6.2：完善錯誤處理

**確保所有錯誤響應符合 swagger.yaml 的 `ErrorResponse` 格式：**

```json
{
	"error": "錯誤訊息"
}
```

#### Step 6.3：補充 Health Check

**確保：**

- `GET /healthz` 返回 `{"status": "ok"}`
- 可選：檢查資料庫連線狀態

#### Step 6.4：執行所有 BDD 測試

```bash
make bdd
```

**確保所有 28 個 feature files 通過：**

- ✅ auth/\*.feature
- ✅ journey/\*.feature
- ✅ mission/\*.feature
- ✅ progress/\*.feature
- ✅ order/\*.feature
- ✅ user/\*.feature

#### Step 6.5：程式碼品質檢查

**執行：**

```bash
# Lint
golangci-lint run

# Format
gofmt -w .

# Vet
go vet ./...
```

#### Step 6.6：撰寫 README.md

**包含：**

- 專案介紹
- 技術棧
- 快速開始指南
- API 文檔連結（swagger.yaml）
- 測試執行方式
- 部署說明

**完成標準：**

- ✅ 所有功能實現完成
- ✅ 所有 BDD 測試通過
- ✅ 程式碼格式化
- ✅ README 完整

---

### **階段 7：容器化與部署準備（選做，0.5 天）**

#### Step 7.1：撰寫 Dockerfile

**多階段構建：**

```dockerfile
# Build stage
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o server cmd/server/main.go

# Run stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/server .
COPY --from=builder /app/config ./config
EXPOSE 8080
CMD ["./server"]
```

#### Step 7.2：更新 docker-compose.yml

**新增 backend service：**

```yaml
services:
  backend:
    build: .
    ports:
      - '8080:8080'
    environment:
      DB_HOST: postgres
      DB_PORT: 5432
      DB_USER: postgres
      DB_PASSWORD: postgres
      DB_NAME: waterballsa
    depends_on:
      - postgres

  postgres:
    # ... (同前)
```

#### Step 7.3：測試 Docker 部署

```bash
docker-compose up --build
```

**完成標準：**

- ✅ Docker 容器成功構建和運行
- ✅ 應用可透過 Docker 訪問

---

## 📋 實作 BDD 測試的關鍵步驟

由於 BDD 測試是整個重寫的核心驗收標準，這裡詳細說明如何實現。

### godog Step Definitions 實作策略

#### 1. 建立可復用的 Common Steps（test/steps/common_steps.go）

**實現這些通用 steps：**

```go
// HTTP 請求相關
When(`^I send "([^"]*)" request to "([^"]*)"$`)
When(`^I set request body to:$`)
When(`^I set header "([^"]*)" to "([^"]*)"$`)

// HTTP 響應驗證
Then(`^the response status code should be (\d+)$`)
Then(`^the response body should contain field "([^"]*)"$`)
Then(`^the response body field "([^"]*)" should equal string "([^"]*)"$`)
Then(`^the response body field "([^"]*)" should equal number (\d+)$`)
```

**實作技巧：**

- 使用全域變數保存 HTTP request/response
- 使用 `encoding/json` + `gjson` 解析 JSON
- 使用 `httptest` 或直接呼叫 Gin router

#### 2. 建立資料庫設置 Steps（test/steps/database_steps.go）

**實現這些資料庫 steps：**

```go
Given(`^the database has a user:$`) // 插入測試用戶
Given(`^the database has a journey:$`) // 插入測試旅程
Given(`^the database has an order:$`) // 插入測試訂單
Given(`^the database is cleaned$`) // 清空資料庫
```

**實作技巧：**

- 直接使用 GORM 插入測試資料
- 每個 scenario 前清空相關表（使用 `@BeforeScenario` hook）

#### 3. 建立測試 Suite（test/steps/suite_test.go）

```go
package steps

import (
    "testing"
    "github.com/cucumber/godog"
)

func TestFeatures(t *testing.T) {
    suite := godog.TestSuite{
        ScenarioInitializer: InitializeScenario,
        Options: &godog.Options{
            Format:   "pretty",
            Paths:    []string{"../features"},
            TestingT: t,
        },
    }

    if suite.Run() != 0 {
        t.Fatal("non-zero status returned, failed to run feature tests")
    }
}

func InitializeScenario(ctx *godog.ScenarioContext) {
    // 註冊所有 step definitions
    RegisterCommonSteps(ctx)
    RegisterDatabaseSteps(ctx)
    RegisterAPISteps(ctx)
}
```

### 執行測試

```bash
# 執行所有 BDD 測試
go test ./test/steps -v

# 或使用 godog CLI
godog run test/features
```

---

## ⚠️ 常見陷阱與注意事項

### 1. GORM 與 JPA 的差異

| JPA 行為                   | GORM 等價                                     | 注意事項                |
| -------------------------- | --------------------------------------------- | ----------------------- |
| `@OneToMany` Lazy Loading  | 需手動 `Preload()`                            | GORM 默認不自動加載關聯 |
| `@Transactional`           | `db.Transaction(func(tx *gorm.DB) error {})`  | 手動管理事務            |
| `@Lock(PESSIMISTIC_WRITE)` | `Clauses(clause.Locking{Strength: "UPDATE"})` | 悲觀鎖語法不同          |
| Cascade operations         | 需手動設置 `constraint:OnDelete:CASCADE`      | GORM 不會自動級聯       |

### 2. 錯誤處理模式

**Go 的顯式錯誤處理：**

```go
// Java 風格（錯誤）
user, _ := userRepo.FindByID(id)  // 忽略錯誤 ❌

// Go 慣例（正確）
user, err := userRepo.FindByID(id)
if err != nil {
    return nil, fmt.Errorf("failed to find user: %w", err)
}
```

### 3. JSON 序列化

**注意 struct tags：**

```go
type User struct {
    ID       int64  `json:"id" gorm:"primaryKey"`           // 同時指定 JSON 和 GORM tags
    Username string `json:"username" gorm:"unique;not null"`
    Password string `json:"-" gorm:"not null"`              // json:"-" 表示不序列化到 JSON
}
```

### 4. 時間處理

**Go 的時間格式化：**

```go
// Java: new SimpleDateFormat("yyyy-MM-dd").format(date)
// Go: time.Now().Format("2006-01-02")  // Go 的魔法數字！

// Unix timestamp 轉換
millis := time.Now().UnixMilli()  // Java: System.currentTimeMillis()
```

### 5. Nullable 欄位

**使用指標或 sql.NullXXX：**

```go
type Order struct {
    PaidAt    *time.Time `json:"paidAt" gorm:"default:null"`  // 可為 null
    ExpiredAt *time.Time `json:"expiredAt"`
}
```

### 6. 分頁查詢

**GORM 分頁：**

```go
var orders []Order
db.Where("user_id = ?", userID).
   Limit(pageSize).
   Offset((page - 1) * pageSize).
   Find(&orders)
```

---

## 📊 預期時程表

| 階段 | 任務           | 預估時間 | 累計時間    |
| ---- | -------------- | -------- | ----------- |
| 0    | 專案初始化     | 1 天     | 1 天        |
| 1    | 基礎架構層     | 1-2 天   | 2-3 天      |
| 2    | 認證模塊       | 1 天     | 3-4 天      |
| 3    | 學習旅程模塊   | 2 天     | 5-6 天      |
| 4    | 進度追蹤模塊   | 1 天     | 6-7 天      |
| 5    | 訂單系統       | 2 天     | 8-9 天      |
| 6    | 限流與整合     | 1 天     | 9-10 天     |
| 7    | 容器化（選做） | 0.5 天   | 9.5-10.5 天 |

**樂觀估計：7 天**（有 AI 輔助 + Golang 基礎良好）
**保守估計：10-14 天**（包含學習曲線和除錯）

---

## ✅ 驗收標準

### 功能驗收

- [ ] 所有 28 個 .feature 文件的 BDD 測試通過
- [ ] 所有 API endpoints 符合 swagger.yaml 規格
- [ ] 資料庫 schema 與 Java 版本一致

### 技術驗收

- [ ] 使用 Gin 框架
- [ ] 使用 GORM 作為 ORM
- [ ] JWT 認證正確實現
- [ ] 悲觀鎖（訂單支付）正確實現
- [ ] 定時任務（訂單過期）正確運作
- [ ] 錯誤處理統一且符合規格

### 代碼品質

- [ ] 通過 `golangci-lint` 檢查
- [ ] 通過 `go vet` 檢查
- [ ] 程式碼格式化（`gofmt`）
- [ ] 有清晰的 README.md

### 部署就緒

- [ ] Docker 容器成功構建
- [ ] docker-compose 可啟動完整環境
- [ ] 環境變數配置完整

---

## 🎯 成功關鍵因素

1. **嚴格遵循 BDD 驅動開發**

   - 每完成一個模塊就跑對應的 .feature 測試
   - 測試通過才進入下一個模塊

2. **參考 Java 實現，但不盲目複製**

   - 理解業務邏輯
   - 用 Go 的慣用法重新實現

3. **保持簡單，避免過度設計**

   - 不追求完美的架構
   - 先求功能對齊，再考慮優化

4. **善用 AI 輔助**

   - 遇到 Go 語法問題立即查詢
   - 讓 AI 幫忙生成 GORM 查詢、DTO 轉換等重複代碼

5. **頻繁提交版本**
   - 每完成一個小功能就 commit
   - 方便回滾和追蹤進度

---

## 📚 參考資源

### 官方文檔

- [Gin Web Framework](https://gin-gonic.com/docs/)
- [GORM ORM](https://gorm.io/docs/)
- [golang-jwt](https://github.com/golang-jwt/jwt)
- [godog (Cucumber for Go)](https://github.com/cucumber/godog)
- [golang-migrate](https://github.com/golang-migrate/migrate)

### 學習資源

- [Effective Go](https://go.dev/doc/effective_go)
- [Go by Example](https://gobyexample.com/)
- [GORM Guides](https://gorm.io/docs/index.html)

### 專案參考

- 現有 Java 實現（waterballsa-backend）
- Swagger API 文檔（docs/api-docs/swagger.yaml）

---

## 🚀 下一步行動

1. **立即開始：建立新 GitHub repo**

   ```bash
   # 在 GitHub 建立 waterballsa-backend-go
   git clone https://github.com/YOUR_USERNAME/waterballsa-backend-go.git
   cd waterballsa-backend-go
   ```

2. **執行階段 0 的所有步驟**（預計 1 天）

3. **跑通第一個測試**（`/healthz` endpoint）

4. **開始階段 1，逐步推進**

---

**祝重寫順利！有任何問題隨時提問。** 🎉
