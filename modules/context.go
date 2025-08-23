package modules

import (
	"gorm.io/gorm"

	"github.com/antoniusDoni/monorepo/shared/repository"
	"github.com/antoniusDoni/monorepo/shared/service"
	wrepo "github.com/antoniusDoni/monorepo/modules/warehouse/repository"
)

// ModuleContext holds shared dependencies for modules.
type ModuleContext struct {
	DB          *gorm.DB
	UserRepo    repository.UserRepository
	AuthService *service.AuthService
	OfficeRepo  wrepo.OfficeRepository
}
