package module

import (
	"gorm.io/gorm"

	"github.com/antoniusDoni/monorepo/shared/repository"
	"github.com/antoniusDoni/monorepo/shared/service"
)

// ModuleContext holds shared dependencies for modules.
type ModuleContext struct {
	DB          *gorm.DB
	UserRepo    repository.UserRepository
	AuthService *service.AuthService
}
