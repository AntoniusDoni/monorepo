package auth

import (
	"errors"
	"net/http"
	"strings"

	models "github.com/antoniusDoni/monorepo/shared/model"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type AuthMiddleware struct {
	JwtSecret []byte
	AuthMode  string // "jwt" or "token"
	DB        *gorm.DB
}

// NewAuthMiddleware creates a new AuthMiddleware instance
func NewAuthMiddleware(jwtSecret string, authMode string, db *gorm.DB) *AuthMiddleware {
	if jwtSecret == "" {
		panic("JWT secret cannot be empty")
	}
	return &AuthMiddleware{
		JwtSecret: []byte(jwtSecret),
		AuthMode:  authMode,
		DB:        db,
	}
}

// Handler returns the Echo middleware function
func (m *AuthMiddleware) Handler(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := strings.TrimSpace(c.Request().Header.Get("Authorization"))
		if authHeader == "" {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Missing Authorization header"})
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid Authorization header format"})
		}

		scheme := strings.ToLower(parts[0])
		tokenStr := parts[1]

		switch m.AuthMode {
		case "jwt":
			if scheme != "bearer" {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Authorization scheme must be Bearer for JWT mode"})
			}
			token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, errors.New("unexpected signing method")
				}
				return m.JwtSecret, nil
			})
			if err != nil || !token.Valid {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid or expired JWT token"})
			}
			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid JWT claims"})
			}
			if userID, ok := claims["user_id"].(float64); ok {
				c.Set("user_id", uint(userID))
			}
			return next(c)

		case "token":
			if scheme != "token" {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Authorization scheme must be Token for token mode"})
			}
			user, err := m.getUserByToken(tokenStr)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid or unknown token"})
			}
			c.Set("user_id", user.ID)
			c.Set("username", user.Username)
			return next(c)

		default:
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Invalid AUTH_MODE configuration"})
		}
	}
}

func (m *AuthMiddleware) getUserByToken(token string) (*models.User, error) {
	var user models.User
	result := m.DB.Where("api_token = ?", token).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}
