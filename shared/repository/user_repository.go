package repository

import (
	model "github.com/antoniusDoni/monorepo/shared/model"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *model.User) error
	FindByUsername(username string) (*model.User, error)
	FindByID(id uint) (*model.User, error)
	GetRolesByUserID(userID uint) ([]model.Role, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user *model.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) FindByUsername(username string) (*model.User, error) {
	var user model.User
	err := r.db.Where("username = ?", username).
		Preload("Roles"). // preload roles association
		First(&user).Error
	return &user, err
}

func (r *userRepository) FindByID(id uint) (*model.User, error) {
	var user model.User
	err := r.db.Preload("Roles").First(&user, id).Error // preload roles association
	return &user, err
}

func (r *userRepository) GetRolesByUserID(userID uint) ([]model.Role, error) {
	var roles []model.Role
	err := r.db.Table("roles").
		Select("roles.*").
		Joins("inner join user_roles on user_roles.role_id = roles.id").
		Where("user_roles.user_id = ?", userID).
		Scan(&roles).Error
	return roles, err
}
