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

type contextKey string

const (
	ContextKeyUserID   contextKey = "user_id"
	ContextKeyUsername contextKey = "username"
	ContextKeyRoles    contextKey = "roles"
)

// Middleware is the Echo middleware function for authentication and role loading
func (a *AuthMiddleware) Middleware(next echo.HandlerFunc) echo.HandlerFunc {
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

		switch a.AuthMode {
		case "jwt":
			if scheme != "bearer" {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Authorization scheme must be Bearer for JWT mode"})
			}

			token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, errors.New("unexpected signing method")
				}
				return a.JwtSecret, nil
			})

			if err != nil || !token.Valid {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid or expired JWT token"})
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid JWT claims"})
			}

			userIDFloat, ok := claims["user_id"].(float64)
			if !ok {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Missing user_id claim"})
			}

			userID := uint(userIDFloat)
			c.Set(string(ContextKeyUserID), userID)

			roles, err := LoadRolesForUser(a.DB, userID)
			if err == nil {
				c.Set(string(ContextKeyRoles), roles)
			}

		case "token":
			if scheme != "token" {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Authorization scheme must be Token for token mode"})
			}

			var user models.User
			if err := a.DB.Where("api_token = ?", tokenStr).First(&user).Error; err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid or unknown token"})
			}

			c.Set(string(ContextKeyUserID), user.ID)
			c.Set(string(ContextKeyUsername), user.Username)

			roles, err := LoadRolesForUser(a.DB, user.ID)
			if err == nil {
				c.Set(string(ContextKeyRoles), roles)
			}

		default:
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Invalid AUTH_MODE configuration"})
		}

		return next(c)
	}
}

// LoadRolesForUser loads roles assigned to userID
func LoadRolesForUser(db *gorm.DB, userID uint) ([]models.Role, error) {
	var roles []models.Role
	err := db.Table("roles").
		Select("roles.*").
		Joins("inner join user_roles on user_roles.role_id = roles.id").
		Where("user_roles.user_id = ?", userID).
		Scan(&roles).Error
	return roles, err
}

// RoleMiddleware enforces required role presence in user roles set in context
func RoleMiddleware(requiredRole string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			rolesAny := c.Get(string(ContextKeyRoles))
			roles, ok := rolesAny.([]models.Role)
			if !ok || !hasRole(roles, requiredRole) {
				return c.JSON(http.StatusForbidden, map[string]string{"error": "Insufficient role permission"})
			}
			return next(c)
		}
	}
}

func hasRole(roles []models.Role, required string) bool {
	for _, r := range roles {
		if r.Name == required {
			return true
		}
	}
	return false
}
