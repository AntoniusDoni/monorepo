package handler

import (
	"net/http"

	"github.com/antoniusDoni/monorepo/shared/contract"
	"github.com/antoniusDoni/monorepo/shared/service"
	"github.com/labstack/echo/v4"
)

// AuthHandler handles user auth routes.
type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// RegisterRoutes implements RouteRegistrar interface.
// Registers routes under "/auths" group by default.
func (h *AuthHandler) RegisterRoutes(g *echo.Group) {
	g.POST("/register", h.Register)
	g.POST("/register-with-office", h.RegisterWithOffice)
	g.POST("/login", h.Login)
}

// Register godoc
// @Summary      Register a new user
// @Description  Register a new user with an existing office ID
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request  body      contract.RegisterRequest  true  "Registration request"
// @Success      201      {object}  contract.RegisterResponse
// @Failure      400      {object}  object
// @Router       /register [post]
func (h *AuthHandler) Register(c echo.Context) error {
	req := new(contract.RegisterRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}
	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	if err := h.authService.Register(req.Username, req.Password, req.Email, req.OfficeID); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, contract.RegisterResponse{Message: "User registered successfully"})
}

// RegisterWithOffice godoc
// @Summary      Register a new user with office creation
// @Description  Create a new office and register the first user for that office in a single operation
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request  body      contract.RegisterWithOfficeRequest  true  "Registration with office request"
// @Success      201      {object}  contract.RegisterWithOfficeResponse
// @Failure      400      {object}  object
// @Router       /register-with-office [post]
func (h *AuthHandler) RegisterWithOffice(c echo.Context) error {
	req := new(contract.RegisterWithOfficeRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}
	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	resp, err := h.authService.RegisterWithOffice(req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, resp)
}

// Login godoc
// @Summary      User login
// @Description  Authenticate user and return JWT token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request  body      contract.LoginRequest  true  "Login request"
// @Success      200      {object}  contract.LoginResponse
// @Failure      400      {object}  object
// @Failure      401      {object}  object
// @Router       /login [post]
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
