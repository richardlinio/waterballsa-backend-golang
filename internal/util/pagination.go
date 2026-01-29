package util

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/richardlinio/waterballsa-backend-golang/internal/apperror"
)

const (
	DefaultPage  = int32(1)
	DefaultLimit = int32(20)
	MaxLimit     = int32(100)
)

// PaginationParams holds pagination parameters
type PaginationParams struct {
	Page   int32
	Limit  int32
	Offset int32
}

// ParsePaginationParams parses pagination parameters from query string
// Returns error if page or limit parameters are invalid
func ParsePaginationParams(c *gin.Context) (*PaginationParams, error) {
	page := DefaultPage
	limit := DefaultLimit

	// Parse page parameter
	if pageStr := c.Query("page"); pageStr != "" {
		pageInt, err := strconv.ParseInt(pageStr, 10, 32)
		if err != nil {
			return nil, apperror.ValidationFailed()
		}
		if pageInt < 1 {
			return nil, apperror.ValidationFailed()
		}
		page = int32(pageInt)
	}

	// Parse limit parameter
	if limitStr := c.Query("limit"); limitStr != "" {
		limitInt, err := strconv.ParseInt(limitStr, 10, 32)
		if err != nil {
			return nil, apperror.ValidationFailed()
		}
		if limitInt < 1 {
			return nil, apperror.ValidationFailed()
		}
		if limitInt > int64(MaxLimit) {
			return nil, apperror.ValidationFailed()
		}
		limit = int32(limitInt)
	}

	// Calculate offset
	offset := (page - 1) * limit

	return &PaginationParams{
		Page:   page,
		Limit:  limit,
		Offset: offset,
	}, nil
}
