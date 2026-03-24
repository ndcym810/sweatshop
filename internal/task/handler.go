// internal/task/handler.go
package task

import (
	"github.com/labstack/echo/v4"
	"github.com/sweatshop/sweatshop/internal/shared/response"
)

// Handler handles task HTTP requests
type Handler struct {
	service *Service
}

// NewHandler creates a new task handler
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// List handles GET /api/teams/:teamId/tasks
func (h *Handler) List(c echo.Context) error {
	teamID := c.Param("teamId")
	status := c.QueryParam("status")
	tasks, err := h.service.List(teamID, status)
	if err != nil {
		return response.InternalError(c, "Failed to list tasks")
	}
	return response.OK(c, tasks)
}

// Get handles GET /api/teams/:teamId/tasks/:id
func (h *Handler) Get(c echo.Context) error {
	id := c.Param("id")
	task, err := h.service.Get(id)
	if err != nil {
		return response.InternalError(c, "Failed to get task")
	}
	if task == nil {
		return response.NotFound(c, "Task not found")
	}
	return response.OK(c, task)
}

// Create handles POST /api/teams/:teamId/tasks
func (h *Handler) Create(c echo.Context) error {
	teamID := c.Param("teamId")
	var input CreateTaskInput
	if err := c.Bind(&input); err != nil {
		return response.BadRequest(c, "Invalid request body")
	}
	if input.Title == "" {
		return response.BadRequest(c, "Title is required")
	}

	task, err := h.service.Create(teamID, input)
	if err != nil {
		return response.InternalError(c, "Failed to create task")
	}
	return response.Created(c, task)
}

// Update handles PUT /api/teams/:teamId/tasks/:id
func (h *Handler) Update(c echo.Context) error {
	id := c.Param("id")
	var input UpdateTaskInput
	if err := c.Bind(&input); err != nil {
		return response.BadRequest(c, "Invalid request body")
	}

	task, err := h.service.Update(id, input)
	if err != nil {
		return response.InternalError(c, "Failed to update task")
	}
	if task == nil {
		return response.NotFound(c, "Task not found")
	}
	return response.OK(c, task)
}

// Delete handles DELETE /api/teams/:teamId/tasks/:id
func (h *Handler) Delete(c echo.Context) error {
	id := c.Param("id")
	if err := h.service.Delete(id); err != nil {
		return response.InternalError(c, "Failed to delete task")
	}
	return response.NoContent(c)
}

// RegisterRoutes registers task routes
func (h *Handler) RegisterRoutes(g *echo.Group) {
	g.GET("", h.List)
	g.GET("/:id", h.Get)
	g.POST("", h.Create)
	g.PUT("/:id", h.Update)
	g.DELETE("/:id", h.Delete)
}
