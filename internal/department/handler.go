// internal/department/handler.go
package department

import (
	"github.com/labstack/echo/v4"
	"github.com/sweatshop/sweatshop/internal/shared/response"
)

// Handler handles department HTTP requests
type Handler struct {
	service *Service
}

// NewHandler creates a new department handler
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// List handles GET /api/teams/:teamId/departments
func (h *Handler) List(c echo.Context) error {
	teamID := c.Param("teamId")
	departments, err := h.service.List(teamID)
	if err != nil {
		return response.InternalError(c, "Failed to list departments")
	}
	return response.OK(c, departments)
}

// Get handles GET /api/teams/:teamId/departments/:id
func (h *Handler) Get(c echo.Context) error {
	id := c.Param("id")
	department, err := h.service.Get(id)
	if err != nil {
		return response.InternalError(c, "Failed to get department")
	}
	if department == nil {
		return response.NotFound(c, "Department not found")
	}
	return response.OK(c, department)
}

// Create handles POST /api/teams/:teamId/departments
func (h *Handler) Create(c echo.Context) error {
	teamID := c.Param("teamId")
	var input CreateDepartmentInput
	if err := c.Bind(&input); err != nil {
		return response.BadRequest(c, "Invalid request body")
	}
	if input.Name == "" {
		return response.BadRequest(c, "Name is required")
	}

	department, err := h.service.Create(teamID, input)
	if err != nil {
		return response.InternalError(c, "Failed to create department")
	}
	return response.Created(c, department)
}

// Update handles PUT /api/teams/:teamId/departments/:id
func (h *Handler) Update(c echo.Context) error {
	id := c.Param("id")
	var input UpdateDepartmentInput
	if err := c.Bind(&input); err != nil {
		return response.BadRequest(c, "Invalid request body")
	}

	department, err := h.service.Update(id, input)
	if err != nil {
		return response.InternalError(c, "Failed to update department")
	}
	if department == nil {
		return response.NotFound(c, "Department not found")
	}
	return response.OK(c, department)
}

// Delete handles DELETE /api/teams/:teamId/departments/:id
func (h *Handler) Delete(c echo.Context) error {
	id := c.Param("id")
	if err := h.service.Delete(id); err != nil {
		return response.InternalError(c, "Failed to delete department")
	}
	return response.NoContent(c)
}

// RegisterRoutes registers department routes
func (h *Handler) RegisterRoutes(g *echo.Group) {
	g.GET("", h.List)
	g.GET("/:id", h.Get)
	g.POST("", h.Create)
	g.PUT("/:id", h.Update)
	g.DELETE("/:id", h.Delete)
}
