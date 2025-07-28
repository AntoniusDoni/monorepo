package contract

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// APIResponse is a generic wrapper for any single item or data payload.
type APIResponse[T any] struct {
	Success bool   `json:"success"`
	Data    T      `json:"data,omitempty"`
	Error   string `json:"error,omitempty"`
}

// PaginationResponse wraps a paginated list of any type T.
type PaginationResponse[T any] struct {
	TotalRecords int64 `json:"total_records"`
	Page         int   `json:"page"`
	PageSize     int   `json:"page_size"`
	Items        []T   `json:"items"`
}

type RegisterResponse struct {
	Message string `json:"message"`
}

type LoginResponse struct {
	UserID uint   `json:"user_identifier"`
	Role   string `json:"role"`
	Token  string `json:"token"`
}

func ErrorResponse(c echo.Context, statusCode int, err error) error {
	return c.JSON(statusCode, map[string]string{
		"error": err.Error(),
	})
}

// MessageResponse formats error string directly
func MessageResponse(c echo.Context, statusCode int, msg string) error {
	return c.JSON(statusCode, map[string]string{
		"error": msg,
	})
}

// SuccessResponse formats a successful response
func SuccessResponse(c echo.Context, data any) error {
	return c.JSON(http.StatusOK, data)
}
