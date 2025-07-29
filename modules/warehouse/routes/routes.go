package routes

import (
	"reflect"
	"strings"

	"github.com/antoniusDoni/monorepo/modules/warehouse/handler"
	"github.com/labstack/echo/v4"
)

// RegisterModules registers all provided modules automatically
func RegisterModules(e *echo.Echo, modules []handler.RouteRegistrar) {
	for _, module := range modules {
		groupName := deriveGroupName(module)
		g := e.Group(groupName)
		module.RegisterRoutes(g)
	}
}

// deriveGroupName derives route group path from handler struct name
func deriveGroupName(module handler.RouteRegistrar) string {
	t := reflect.TypeOf(module)
	name := t.Elem().Name()                         // e.g. "WarehouseHandler"
	baseName := strings.TrimSuffix(name, "Handler") // "Warehouse"
	return "/" + strings.ToLower(baseName) + "s"    // "/warehouses"
}
