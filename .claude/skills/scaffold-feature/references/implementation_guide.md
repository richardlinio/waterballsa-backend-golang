# Implementation Guide: Patterns & Pitfalls

This guide provides common implementation patterns for this project, along with anti-patterns to avoid.

## Go Idiomatic Patterns

For general Go language conventions, including:

- **Slice Initialization**
- **Return Value Design**
- **Repository Error Handling**
- **Variable Naming Conventions**
- **Aggregate Structure Design**
- **Configuration Management**

Please refer to the `go-idiomatic` skill. It contains detailed explanations and code examples for ensuring your Go code is idiomatic and consistent with our project's standards.

You can activate it by typing `@go-idiomatic` in the chat.

---

The following sections cover patterns specific to scaffolding new features in this application.

## Authorization Checks

✅ **Do:** Perform authorization checks in the handler to protect endpoints. Retrieve the user from the context and verify their permissions against the requested resource.

```go
// Handler
user := c.MustGet("JWT_PAYLOAD").(*model.User)
requestedUserID, _ := strconv.ParseInt(c.Param("userId"), 10, 64)

if user.ID != requestedUserID {
    _ = c.Error(apperror.UnauthorizedAccess())
    return
}
```

❌ **Don't:** Skip authorization checks. This is a critical security vulnerability.

## DTO Assembly Pattern

✅ **Do:** Define converter functions in the DTO layer to assemble response DTOs. Extract message variables in handlers for better readability.

```go
// internal/dto/order.go
type PayOrderResponse struct {
    ID          int64   `json:"id"`
    OrderNumber string  `json:"orderNumber"`
    Status      string  `json:"status"`
    Price       float64 `json:"price"`
    PaidAt      int64   `json:"paidAt"`
    Message     string  `json:"message"`
}

// ToPayOrderResponse converts paid order data to PayOrderResponse DTO
func ToPayOrderResponse(order *model.Order, message string) PayOrderResponse {
    return PayOrderResponse{
        ID:          order.ID,
        OrderNumber: order.OrderNumber,
        Status:      order.Status,
        Price:       order.Price,
        PaidAt:      order.PaidAt.UnixMilli(),
        Message:     message,
    }
}
```

```go
// internal/handler/order.go
func (h *OrderHandler) PayOrder(c *gin.Context) {
    // ... validation and service call ...

    paidOrder, err := h.orderService.PayOrder(ctx, orderID, userID)
    if err != nil {
        _ = c.Error(err)
        return
    }

    message := "付款完成"
    response := dto.ToPayOrderResponse(paidOrder, message)
    c.JSON(http.StatusOK, response)
}
```

❌ **Don't:** Manually construct response DTOs in handlers. This violates separation of concerns.

```go
// Bad: Assembling DTO directly in handler
response := dto.PayOrderResponse{
    ID:          paidOrder.ID,
    OrderNumber: paidOrder.OrderNumber,
    Status:      paidOrder.Status,
    Price:       paidOrder.Price,
    PaidAt:      paidOrder.PaidAt.UnixMilli(),
    Message:     "付款完成",
}
c.JSON(http.StatusOK, response)
```

## UPSERT Pattern in SQL

✅ **Do:** Use the `INSERT ... ON CONFLICT DO UPDATE` statement for creating or updating records in a single database round-trip.

```sql
-- name: UpsertProgress :one
INSERT INTO user_mission_progress (user_id, mission_id, status, watch_position_seconds)
VALUES ($1, $2, $3, $4)
ON CONFLICT (user_id, mission_id)
DO UPDATE SET
  status = EXCLUDED.status,
  watch_position_seconds = EXCLUDED.watch_position_seconds,
  updated_at = NOW()
RETURNING *;
```

## SQLc Naming Conventions

✅ **Do:** Follow consistent naming conventions for SQLc queries.

- `GetXxxByYyy`: Single record query
- `ListXxxByYyy`: Multiple record query
- `CreateXxx`: Create a new record
- `UpdateXxx`: Update an existing record
- `UpsertXxx`: Create or update a record
- `DeleteXxx`: Delete a record

Remember to run `make sqlc` after adding or modifying queries.

## Router Registration

✅ **Do:** Register new handlers in `internal/router/router.go`.

```go
func Setup(
    engine *gin.Engine,
    cfg config.Config,
    fooHandler *handler.FooHandler, // Add new handler as a parameter
    pool *pgxpool.Pool,
    logger *slog.Logger,
) *gin.Engine {
    // ...
    protected := engine.Group("/")
    protected.Use(middleware.JWTAuth(cfg.JWT))
    {
        protected.GET("/foo/:id", fooHandler.GetFoo) // Register route
    }
    return engine
}
```

❌ **Don't:** Forget to register the route. The endpoint will not be available otherwise.
