package handler

import "github.com/labstack/echo/v4"

// RouteRegistrar defines a contract for modules to register their routes.
type RouteRegistrar interface {
	RegisterRoutes(g *echo.Group)
}
