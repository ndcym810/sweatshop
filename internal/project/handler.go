// internal/project/handler.go
package project

import (
	"github.com/labstack/echo/v4"
	"github.com/sweatshop/sweatshop/internal/shared/response"
)

// Handler handles project HTTP requests
type Handler struct {
	service *Service
}

// NewHandler creates a new project handler
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// List handles GET /api/teams/:teamId/projects
func (h *Handler) List(c echo.Context) error {
	teamID := c.Param("teamId")
	projects, err := h.service.ListByTeam(teamID)
	if err != nil {
		return response.InternalError(c, "Failed to list projects")
	}
	return response.OK(c, projects)
}

// Get handles GET /api/teams/:teamId/projects/:id
func (h *Handler) Get(c echo.Context) error {
	id := c.Param("id")
	project, err := h.service.Get(id)
	if err != nil {
		return response.InternalError(c, "Failed to get project")
	}
	if project == nil {
		return response.NotFound(c, "Project not found")
	}
	return response.OK(c, project)
}

// Create handles POST /api/teams/:teamId/projects
func (h *Handler) Create(c echo.Context) error {
	teamID := c.Param("teamId")
	var input CreateProjectInput
	if err := c.Bind(&input); err != nil {
		return response.BadRequest(c, "Invalid request body")
	}
	if input.Name == "" {
		return response.BadRequest(c, "Name is required")
	}

	project, err := h.service.Create(teamID, input)
	if err != nil {
		return response.InternalError(c, "Failed to create project")
	}
	return response.Created(c, project)
}

// Update handles PUT /api/teams/:teamId/projects/:id
func (h *Handler) Update(c echo.Context) error {
	id := c.Param("id")
	var input UpdateProjectInput
	if err := c.Bind(&input); err != nil {
		return response.BadRequest(c, "Invalid request body")
	}

	project, err := h.service.Update(id, input)
	if err != nil {
		return response.InternalError(c, "Failed to update project")
	}
	if project == nil {
		return response.NotFound(c, "Project not found")
	}
	return response.OK(c, project)
}

// Delete handles DELETE /api/teams/:teamId/projects/:id
func (h *Handler) Delete(c echo.Context) error {
	id := c.Param("id")
	if err := h.service.Delete(id); err != nil {
		return response.InternalError(c, "Failed to delete project")
	}
	return response.NoContent(c)
}

// RegisterRoutes registers project routes
func (h *Handler) RegisterRoutes(g *echo.Group) {
	g.GET("", h.List)
	g.GET("/:id", h.Get)
	g.POST("", h.Create)
	g.PUT("/:id", h.Update)
	g.DELETE("/:id", h.Delete)
}
