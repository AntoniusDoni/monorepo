package modules

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/antoniusDoni/monorepo/modules/warehouse"
	"github.com/labstack/echo/v4"
)

// ModuleRegistrar defines the interface for module registration
type ModuleRegistrar interface {
	RegisterRoutes(apiGroup *echo.Group) error
	GetName() string
}

// ModuleRegistry manages all available modules
type ModuleRegistry struct {
	modules map[string]ModuleRegistrar
}

// NewModuleRegistry creates a new module registry
func NewModuleRegistry() *ModuleRegistry {
	return &ModuleRegistry{
		modules: make(map[string]ModuleRegistrar),
	}
}

// RegisterModule adds a module to the registry
func (r *ModuleRegistry) RegisterModule(name string, module ModuleRegistrar) {
	r.modules[name] = module
	log.Printf("Module registered: %s", name)
}

// InitializeModules initializes and registers routes for enabled modules
func (r *ModuleRegistry) InitializeModules(apiGroup *echo.Group, services *ModuleContext) error {
	enabledModules := getEnabledModules()

	if len(enabledModules) == 0 {
		log.Println("No modules enabled")
		return nil
	}

	log.Printf("Enabled modules: %v", getEnabledModuleNames(enabledModules))

	for moduleName := range enabledModules {
		if module, exists := r.modules[moduleName]; exists {
			if err := module.RegisterRoutes(apiGroup); err != nil {
				return fmt.Errorf("failed to initialize module %s: %w", moduleName, err)
			}
			log.Printf("Module %s initialized successfully", moduleName)
		} else {
			log.Printf("Warning: Module %s is enabled but not registered", moduleName)
		}
	}

	return nil
}

// getEnabledModules parses enabled modules from environment
func getEnabledModules() map[string]bool {
	enabled := make(map[string]bool)

	enabledModulesEnv := os.Getenv("ENABLE_MODULES")
	if enabledModulesEnv == "" {
		return enabled
	}

	for _, m := range strings.Split(enabledModulesEnv, ",") {
		module := strings.TrimSpace(m)
		if module != "" {
			enabled[module] = true
		}
	}

	return enabled
}

// getEnabledModuleNames returns a slice of enabled module names for logging
func getEnabledModuleNames(enabled map[string]bool) []string {
	var names []string
	for name := range enabled {
		names = append(names, name)
	}
	return names
}

// WarehouseModule implements ModuleRegistrar for the warehouse module
type WarehouseModule struct {
	deps *warehouse.ModuleDependencies
}

// NewWarehouseModule creates a new warehouse module instance
func NewWarehouseModule(services *ModuleContext) *WarehouseModule {
	return &WarehouseModule{
		deps: &warehouse.ModuleDependencies{
			DB:         services.DB,
			OfficeRepo: services.OfficeRepo,
		},
	}
}

// RegisterRoutes implements ModuleRegistrar interface
func (w *WarehouseModule) RegisterRoutes(apiGroup *echo.Group) error {
	return warehouse.RegisterRoutes(apiGroup, w.deps)
}

// GetName implements ModuleRegistrar interface
func (w *WarehouseModule) GetName() string {
	return "warehouse"
}

// SetupModules initializes and returns a configured module registry
func SetupModules(services *ModuleContext) *ModuleRegistry {
	registry := NewModuleRegistry()

	// Register warehouse module
	warehouseModule := NewWarehouseModule(services)
	registry.RegisterModule("warehouse", warehouseModule)

	// Future modules can be registered here
	// hrModule := NewHRModule(services)
	// registry.RegisterModule("hr", hrModule)

	// financeModule := NewFinanceModule(services)
	// registry.RegisterModule("finance", financeModule)

	return registry
}
