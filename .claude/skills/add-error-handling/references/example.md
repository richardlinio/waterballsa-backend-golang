# 完整範例：新增 "Journey Not Found" 錯誤

### 1. internal/apperror/codes.go

```go
const (
    // ... other codes
    // 404 Not Found
    CodeUserNotFound    = "ERR_USER_NOT_FOUND"
    CodeJourneyNotFound = "ERR_JOURNEY_NOT_FOUND" // New
    CodeMissionNotFound = "ERR_MISSION_NOT_FOUND"
    // ... other codes
)
```

### 2. internal/apperror/messages.go

```go
var errorMessages = map[string]string{
    // ... other messages
    CodeUserNotFound:    "使用者不存在",
    CodeJourneyNotFound: "旅程不存在", // New
    CodeMissionNotFound: "任務不存在",
    // ... other messages
}
```

### 3. internal/apperror/status.go

```go
var httpStatusMap = map[string]int{
    // ... other statuses
    CodeUserNotFound:    http.StatusNotFound,
    CodeJourneyNotFound: http.StatusNotFound, // New
    CodeMissionNotFound: http.StatusNotFound,
    // ... other statuses
}
```

### 4. internal/apperror/error.go

```go
// 404 Not Found
func UserNotFound() *AppError {
    return New(CodeUserNotFound)
}

func JourneyNotFound() *AppError { // New
    return New(CodeJourneyNotFound)
}

func MissionNotFound() *AppError {
    return New(CodeMissionNotFound)
}
```

### 5. internal/service/journey.go

```go
import (
    "context"
    "errors"

    "github.com/waterballsa/waterballsa-backend/internal/apperror"
    "github.com/waterballsa/waterballsa-backend/internal/model"
    "github.com/waterballsa/waterballsa-backend/internal/repository"
)

// ... other service methods

func (s *Service) GetJourney(ctx context.Context, journeyID int64) (*model.Journey, error) {
    journey, err := s.repository.GetJourneyByID(ctx, journeyID)
    if err != nil {
        if errors.Is(err, repository.ErrResourceNotFound) {
            return nil, apperror.JourneyNotFound() // Use the new error
        }
        return nil, apperror.DatabaseError(err) // Wrap other errors
    }
    return journey, nil
}
```
