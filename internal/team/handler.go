// internal/team/handler.go
package team

import (
	"github.com/labstack/echo/v4"
	"github.com/sweatshop/sweatshop/internal/shared/response"
)

// Handler handles team HTTP requests
type Handler struct {
	service *Service
}

// NewHandler creates a new team handler
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// List handles GET /api/teams
func (h *Handler) List(c echo.Context) error {
	teams, err := h.service.List()
	if err != nil {
		return response.InternalError(c, "Failed to list teams")
	}
	return response.OK(c, teams)
}

// Get handles GET /api/teams/:id
func (h *Handler) Get(c echo.Context) error {
	id := c.Param("id")
	team, err := h.service.Get(id)
	if err != nil {
		return response.InternalError(c, "Failed to get team")
	}
	if team == nil {
		return response.NotFound(c, "Team not found")
	}
	return response.OK(c, team)
}

// Create handles POST /api/teams
func (h *Handler) Create(c echo.Context) error {
	var input CreateTeamInput
	if err := c.Bind(&input); err != nil {
		return response.BadRequest(c, "Invalid request body")
	}
	if input.Name == "" {
		return response.BadRequest(c, "Name is required")
	}

	team, err := h.service.Create(input)
	if err != nil {
		return response.InternalError(c, "Failed to create team")
	}
	return response.Created(c, team)
}

// Update handles PUT /api/teams/:id
func (h *Handler) Update(c echo.Context) error {
	id := c.Param("id")
	var input UpdateTeamInput
	if err := c.Bind(&input); err != nil {
		return response.BadRequest(c, "Invalid request body")
	}

	team, err := h.service.Update(id, input)
	if err != nil {
		return response.InternalError(c, "Failed to update team")
	}
	if team == nil {
		return response.NotFound(c, "Team not found")
	}
	return response.OK(c, team)
}

// Delete handles DELETE /api/teams/:id
func (h *Handler) Delete(c echo.Context) error {
	id := c.Param("id")
	if err := h.service.Delete(id); err != nil {
		return response.InternalError(c, "Failed to delete team")
	}
	return response.NoContent(c)
}

// RegisterRoutes registers team routes
func (h *Handler) RegisterRoutes(g *echo.Group) {
	g.GET("", h.List)
	g.GET("/:id", h.Get)
	g.POST("", h.Create)
	g.PUT("/:id", h.Update)
	g.DELETE("/:id", h.Delete)
}
