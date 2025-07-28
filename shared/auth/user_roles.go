package auth

import (
	models "github.com/antoniusDoni/monorepo/shared/model"
	"gorm.io/gorm"
)

// LoadRolesForUser loads roles assigned to userID
func LoadRolesForUser(db *gorm.DB, userID uint) ([]models.Role, error) {
	var roles []models.Role

	// Join user_roles and roles to get all roles for a user
	err := db.Table("roles").
		Select("roles.*").
		Joins("inner join user_roles on user_roles.role_id = roles.id").
		Where("user_roles.user_id = ?", userID).
		Scan(&roles).Error

	if err != nil {
		return nil, err
	}
	return roles, nil
}
