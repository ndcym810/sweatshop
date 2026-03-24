// internal/app/router.go
package app

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// SetupRouter creates and configures the Echo router
func (a *App) SetupRouter() *echo.Echo {
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// API routes
	api := e.Group("/api")

	// Team routes
	a.TeamHandler.RegisterRoutes(api.Group("/teams"))

	// Nested routes under teams
	teams := api.Group("/teams/:teamId")
	a.ProjectHandler.RegisterRoutes(teams.Group("/projects"))
	a.DepartmentHandler.RegisterRoutes(teams.Group("/departments"))
	a.TaskHandler.RegisterRoutes(teams.Group("/tasks"))

	return e
}
