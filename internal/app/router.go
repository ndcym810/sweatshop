// internal/app/router.go
package app

import (
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for MVP
	},
}

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

	// Claude Code team routes (separate namespace)
	a.ClaudeTeamHandler.RegisterRoutes(api.Group("/claude-teams"))

	// Nested routes under teams
	teams := api.Group("/teams/:teamId")
	a.ProjectHandler.RegisterRoutes(teams.Group("/projects"))
	a.DepartmentHandler.RegisterRoutes(teams.Group("/departments"))
	a.TaskHandler.RegisterRoutes(teams.Group("/tasks"))

	// WebSocket endpoint
	e.GET("/ws", a.handleWebSocket)

	return e
}

// handleWebSocket handles WebSocket connections
func (a *App) handleWebSocket(c echo.Context) error {
	conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}

	a.WebSocketHub.Register(conn)
	defer a.WebSocketHub.Unregister(conn)

	// Keep connection alive, read messages (ping/pong handled by gorilla)
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			break
		}
	}
	return nil
}
