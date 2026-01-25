# 依賴注入範例 (Dependency Injection)

本文件說明如何在不同環境中設定依賴注入。

## 生產環境 (`internal/app/app.go`)

```go
// Repositories
fooRepository := repository.NewFooRepository(pool)

// Services
fooService := service.NewFooService(fooRepository, logger, cfg.Server.RequestTimeout)

// Handlers
fooHandler := handler.NewFooHandler(fooService, logger, cfg.Server.RequestTimeout)

// Router
engine := router.Setup(engine, cfg, fooHandler, pool, logger)
```

## 測試環境 (`tests/testutil/server.go`)

**重要:** 新增功能時必須同步更新此檔案，否則測試編譯會失敗！

```go
// 1. Repository 初始化 (第 93-99 行附近)
fooRepository := repository.NewFooRepository(queries)

// 2. Service 初始化 (第 104-108 行附近)
fooService := service.NewFooService(fooRepository)

// 3. Handler 初始化 (第 122-126 行附近)
fooHandler := handler.NewFooHandler(fooService, log, cfg.Server.RequestTimeout)

// 4. Router Setup (第 135 行附近)
r := router.NewRouter(
    ginEngine,
    cfg.CORS,
    cfg.RateLimit,
    cfg.JWT,
    log,
    healthHandler,
    authHandler,
    journeyHandler,
    missionHandler,
    progressHandler,
    fooHandler, // 新增此參數
    jwtMiddleware,
    blacklistChecker,
)
```
