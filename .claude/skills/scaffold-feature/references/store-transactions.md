# Store Pattern - Transaction Implementation Reference

This guide provides detailed patterns for implementing atomic database transactions using the Store pattern. For the basic overview, see [CLAUDE.md](../../../../CLAUDE.md).

## Overview

The **Store pattern** provides a clean abstraction layer for database transactions that encapsulates multi-step atomic operations separately from business logic.

**Key Benefits:**

- Prevents TOCTOU (Time-of-Check-Time-of-Use) race conditions
- Ensures atomic multi-table operations
- Centralizes transaction management logic
- Simplifies service layer code

## Implementation Checklist

Follow these steps when implementing a new transaction:

### 1. Define Transaction Signature

**Location:** `internal/store/{domain}.go` (e.g., `auth.go`, `order.go`, `progress.go`)

Create three components:

```go
// Params: Input parameters for the transaction
type {Action}TxParams struct {
    // Include all required data
    UserID    int64
    ItemID    int64
    // ... other fields
}

// Result: Output data from the transaction
type {Action}TxResult struct {
    // Return domain models, not DTOs
    User  *model.User
    Order *model.Order
    // ... other fields
}

// Function: Transaction implementation
func (s *Store) {Action}Tx(ctx context.Context, arg {Action}TxParams) ({Action}TxResult, error)
```

### 2. Implement execTx Callback

All transaction logic must be inside the `execTx` callback:

```go
func (s *Store) ActionTx(ctx context.Context, arg ActionTxParams) (ActionTxResult, error) {
    var result ActionTxResult

    err := s.execTx(ctx, func(q *db.Queries) error {
        // 1. Create transaction-aware repositories INSIDE callback
        repo := repository.NewRepository(q)

        // 2. Execute transaction logic
        // 3. Populate result
        // 4. Return error or nil

        return nil
    })

    if err != nil {
        return ActionTxResult{}, fmt.Errorf("action transaction failed: %w", err)
    }

    return result, nil
}
```

### 3. Define Domain Errors

**Location:** `internal/store/errors.go`

Create descriptive domain-specific errors:

```go
var (
    // Group errors by domain
    ErrTokenUserMismatch       = errors.New("token user ID mismatch")
    ErrJourneyAlreadyPurchased = errors.New("journey already purchased")
    ErrMissionAlreadyDelivered = errors.New("mission already delivered")
)
```

These errors should describe **business rule violations**, not technical failures.

### 4. Implement TOCTOU Protection

For operations that check-then-act, re-validate inside the transaction:

```go
err := s.execTx(ctx, func(q *db.Queries) error {
    // 1. Get and lock the record
    record, err := repo.GetByID(ctx, id)
    if err != nil {
        return err
    }

    // 2. Re-validate conditions INSIDE transaction
    // Even if pre-transaction check passed, verify again
    if record.Status == "COMPLETED" {
        return ErrAlreadyCompleted
    }

    // 3. Proceed with operation
    // ...
})
```

### 5. Integrate with Service Layer

**Service Constructor:**

```go
func NewService(..., st *store.Store) *Service {
    return &Service{
        // ... other fields
        store: st,
    }
}
```

**Service Method:**

```go
func (s *Service) DoAction(ctx context.Context, req dto.Request) (*dto.Response, error) {
    // 1. Validate input (outside transaction)
    // 2. Perform non-DB operations (outside transaction)

    // 3. Call transaction
    result, err := s.store.ActionTx(ctx, store.ActionTxParams{
        UserID: req.UserID,
        // ... map request to params
    })
    if err != nil {
        // Handle errors
        if errors.Is(err, store.ErrSpecificError) {
            return nil, apperror.BusinessError()
        }
        return nil, apperror.DatabaseError(err)
    }

    // 4. Map result to DTO
    return &dto.Response{/* ... */}, nil
}
```

### 6. Update Dependency Injection

**CRITICAL:** Update BOTH files:

**Production:** `internal/app/app.go`

```go
// Initialize store (provides queries and transactions)
st := store.New(pool)

// Pass to services that need transactions
service := service.NewService(..., st)
```

**Testing:** `tests/testutil/server.go`

```go
// Initialize store
st := store.New(pool)

// Pass to services
service := service.NewService(..., st)
```

⚠️ **Forgetting `testutil/server.go` is a common mistake!**

## Anti-Patterns

### ❌ Creating Repository Outside execTx

**Wrong:**

```go
repo := repository.NewRepository(s.Queries) // Uses pool, not transaction!

err := s.execTx(ctx, func(q *db.Queries) error {
    return repo.Update(ctx, data) // Not in transaction!
})
```

**Correct:**

```go
err := s.execTx(ctx, func(q *db.Queries) error {
    repo := repository.NewRepository(q) // Uses transaction
    return repo.Update(ctx, data)
})
```

### ❌ Missing TOCTOU Re-validation

**Wrong:**

```go
// Service layer
if record.Status == "COMPLETED" {
    return apperror.AlreadyCompleted()
}
// Race condition: status might change before transaction!
result, err := s.store.UpdateTx(ctx, params)
```

**Correct:**

```go
// Transaction re-validates inside
err := s.execTx(ctx, func(q *db.Queries) error {
    record, err := repo.GetByID(ctx, id) // Locks record
    if record.Status == "COMPLETED" {
        return ErrAlreadyCompleted
    }
    // ... proceed with update
})
```

### ❌ Using Store for Simple Operations

**Wrong:**

```go
// Don't create transactions for single operations
func (s *Store) GetUserTx(ctx context.Context, id int64) (*model.User, error) {
    // Unnecessary transaction overhead
}
```

**Correct:**

```go
// Use Repository for simple queries
user, err := userRepo.GetByID(ctx, id)
```

### ❌ Breaking Error Chain

**Wrong:**

```go
if err != nil {
    return errors.New("failed to get user") // Lost original error!
}
```

**Correct:**

```go
if err != nil {
    return fmt.Errorf("failed to get user: %w", err) // Preserves error chain
}
```

### ❌ Forgetting Test DI Updates

**Wrong:**

```go
// Updated app.go but forgot testutil/server.go
// Tests will fail with nil pointer or missing dependency
```

**Correct:**

```go
// Update BOTH:
// - internal/app/app.go
// - tests/testutil/server.go
```

## When to Use Store vs Repository

### Use Store (Transactions) When:

✅ **Atomic multi-table updates**

- Create order + order items + update inventory
- Transfer funds between accounts

✅ **TOCTOU prevention required**

- Check balance then withdraw
- Check status then update
- Any check-then-act pattern

✅ **Complex validation with locks**

- Validate duplicate purchase with lock
- Verify uniqueness constraints

✅ **Multi-step operations must be atomic**

- Token rotation (revoke old + create new)
- Reward claiming (update progress + user experience)

### Use Repository When:

✅ **Single-table operations**

- Create user
- Update profile
- Delete record

✅ **Read-only queries**

- Get by ID
- List with filters
- Count operations

✅ **Simple CRUD**

- No race conditions
- No multi-table coordination
- No complex validation

## Real-World Examples

### Example 1: Token Rotation

**File:** [internal/store/auth.go](../../../../internal/store/auth.go)

**Problem:** Prevent concurrent logout from creating orphaned tokens.

**Solution:** Atomic operation that re-verifies token, revokes old, creates new.

### Example 2: Order Creation

**File:** [internal/store/order_create.go](../../../../internal/store/order_create.go)

**Problem:** Prevent duplicate purchases and handle existing unpaid orders.

**Solution:** Transaction that checks purchase history, validates existing orders, creates order atomically.

### Example 3: Mission Delivery

**File:** [internal/store/progress.go](../../../../internal/store/progress.go)

**Problem:** Prevent double-claiming rewards through concurrent requests.

**Solution:** Lock progress record, re-validate status, update atomically.

## Quick Reference

**Transaction Function Template:**

```go
type ActionTxParams struct { /* inputs */ }
type ActionTxResult struct { /* outputs */ }

func (s *Store) ActionTx(ctx context.Context, arg ActionTxParams) (ActionTxResult, error) {
    var result ActionTxResult

    err := s.execTx(ctx, func(q *db.Queries) error {
        repo := repository.NewRepo(q)

        // 1. Lock records
        // 2. Re-validate
        // 3. Execute operations
        // 4. Populate result

        return nil
    })

    if err != nil {
        return ActionTxResult{}, fmt.Errorf("action tx failed: %w", err)
    }
    return result, nil
}
```

**Service Integration Template:**

```go
result, err := s.store.ActionTx(ctx, store.ActionTxParams{...})
if err != nil {
    if errors.Is(err, store.ErrDomainError) {
        return nil, apperror.BusinessError()
    }
    return nil, apperror.DatabaseError(err)
}
```

## Summary

- ✅ Use Store for atomic multi-step operations
- ✅ Create repos INSIDE execTx callback
- ✅ Re-validate conditions inside transaction
- ✅ Define domain errors in store/errors.go
- ✅ Preserve error chains with fmt.Errorf
- ✅ Update BOTH app.go and testutil/server.go
- ❌ Don't use Store for simple single-table operations
- ❌ Don't create repos outside transaction
- ❌ Don't skip TOCTOU re-validation
