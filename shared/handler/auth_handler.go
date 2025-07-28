package handler

import (
	"net/http"

	"github.com/antoniusDoni/monorepo/shared/contract"
	"github.com/antoniusDoni/monorepo/shared/service"
	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Register(c echo.Context) error {
	req := new(contract.RegisterRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}
	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	if err := h.authService.Register(req.Username, req.Password); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, contract.RegisterResponse{Message: "User registered successfully"})
}
func (h *AuthHandler) Login(c echo.Context) error {
	req := new(contract.LoginRequest)
	if err := c.Bind(req); err != nil {
		return contract.MessageResponse(c, http.StatusBadRequest, "invalid request")
	}

	if err := c.Validate(req); err != nil {
		return contract.ErrorResponse(c, http.StatusBadRequest, err)
	}

	resp, err := h.authService.Login(req.Username, req.Password)
	if err != nil {
		return contract.ErrorResponse(c, http.StatusUnauthorized, err)
	}

	return contract.SuccessResponse(c, resp)

}
