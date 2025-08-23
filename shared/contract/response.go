package contract

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// APIResponse is a generic wrapper for any single item or data payload.
type PaginatedResponse[T any] struct {
	Items      []T   `json:"items"`
	TotalCount int64 `json:"total_count"`
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
}

type APIResponse[T any] struct {
	Success bool   `json:"success"`
	Data    T      `json:"data,omitempty"`
	Error   string `json:"error,omitempty"`
}

type RegisterResponse struct {
	Message string `json:"message"`
}

type LoginResponse struct {
	UserID uint   `json:"user_identifier"`
	Role   string `json:"role"`
	Token  string `json:"token"`
}

type RegisterWithOfficeResponse struct {
	Message  string `json:"message"`
	OfficeID string `json:"office_id"`
	UserID   uint   `json:"user_id"`
}

type OfficeResponse struct {
	OfficeCode string `json:"officeCode"`
	Name       string `json:"name"`
	Address    string `json:"address"`
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

// PaginatedSuccess formats paginated result consistently
func PaginatedSuccess[T any](c echo.Context, items []T, total int64, page, pageSize int) error {
	resp := PaginatedResponse[T]{
		Items:      items,
		TotalCount: total,
		Page:       page,
		PageSize:   pageSize,
	}
	return c.JSON(http.StatusOK, APIResponse[PaginatedResponse[T]]{
		Success: true,
		Data:    resp,
	})
}

// SingleSuccess formats single object response
func SingleSuccess[T any](c echo.Context, data T) error {
	return c.JSON(http.StatusOK, APIResponse[T]{
		Success: true,
		Data:    data,
	})
}
