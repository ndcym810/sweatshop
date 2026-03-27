// internal/claudeteam/handler.go
package claudeteam

import (
	"github.com/labstack/echo/v4"
	"github.com/sweatshop/sweatshop/internal/shared/response"
)

// Handler handles Claude team HTTP requests
type Handler struct {
	service *Service
}

// NewHandler creates a new Claude team handler
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// List handles GET /api/claude-teams
func (h *Handler) List(c echo.Context) error {
	teams, err := h.service.DiscoverTeams()
	if err != nil {
		return response.InternalError(c, "Failed to discover teams")
	}
	return response.OK(c, teams)
}

// Get handles GET /api/claude-teams/:name
func (h *Handler) Get(c echo.Context) error {
	name := c.Param("name")

	config, err := h.service.ReadTeamConfig(name)
	if err != nil {
		return response.InternalError(c, "Failed to read team config")
	}
	if config == nil {
		return response.NotFound(c, "Team not found")
	}

	teamResp := h.service.ToTeamResponse(config)
	return response.OK(c, teamResp)
}

// GetInbox handles GET /api/claude-teams/:name/inbox/:agent
func (h *Handler) GetInbox(c echo.Context) error {
	teamName := c.Param("name")
	agentName := c.Param("agent")

	messages, err := h.service.ReadInbox(teamName, agentName)
	if err != nil {
		return response.InternalError(c, "Failed to read inbox")
	}

	return response.OK(c, messages)
}

// SendMessage handles POST /api/claude-teams/:name/message
func (h *Handler) SendMessage(c echo.Context) error {
	teamName := c.Param("name")

	var input SendMessageInput
	if err := c.Bind(&input); err != nil {
		return response.BadRequest(c, "Invalid request body")
	}
	if input.To == "" {
		return response.BadRequest(c, "Recipient 'to' is required")
	}
	if input.Message == "" {
		return response.BadRequest(c, "Message is required")
	}

	// Verify team exists
	config, err := h.service.ReadTeamConfig(teamName)
	if err != nil {
		return response.InternalError(c, "Failed to verify team")
	}
	if config == nil {
		return response.NotFound(c, "Team not found")
	}

	// Write message to inbox
	if err := h.service.WriteMessage(teamName, input.To, "dashboard", input.Message); err != nil {
		return response.InternalError(c, "Failed to send message")
	}

	return response.OK(c, map[string]string{"status": "sent"})
}

// MarkRead handles DELETE /api/claude-teams/:name/inbox/:agent/:timestamp
func (h *Handler) MarkRead(c echo.Context) error {
	teamName := c.Param("name")
	agentName := c.Param("agent")
	timestamp := c.Param("timestamp")

	if err := h.service.MarkMessageRead(teamName, agentName, timestamp); err != nil {
		return response.InternalError(c, "Failed to mark message as read")
	}

	return response.NoContent(c)
}

// RegisterRoutes registers Claude team routes
func (h *Handler) RegisterRoutes(g *echo.Group) {
	g.GET("", h.List)
	g.GET("/:name", h.Get)
	g.GET("/:name/inbox/:agent", h.GetInbox)
	g.POST("/:name/message", h.SendMessage)
	g.DELETE("/:name/inbox/:agent/:timestamp", h.MarkRead)
}
