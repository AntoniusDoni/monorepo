package handler

import (
	"reflect"
	"strings"

	"github.com/labstack/echo/v4"
)

// RegisterModules registers all provided modules by invoking their RegisterRoutes methods,
// grouping routes under pluralized names derived from the handler struct name.
func RegisterModules(e *echo.Echo, modules []RouteRegistrar) {
	for _, module := range modules {
		groupName := deriveGroupName(module)
		g := e.Group(groupName)
		module.RegisterRoutes(g)
	}
}

// deriveGroupName derives route group path from handler struct name.
// E.g. WarehouseHandler -> "/warehouses"
func deriveGroupName(module RouteRegistrar) string {
	t := reflect.TypeOf(module)
	name := t.Elem().Name()                         // e.g. "WarehouseHandler"
	baseName := strings.TrimSuffix(name, "Handler") // "Warehouse"
	return "/" + strings.ToLower(baseName) + "s"    // "/warehouses"
}
