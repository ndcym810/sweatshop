// internal/claudeteam/model.go
package claudeteam

// TeamConfig represents ~/.claude/teams/{team-name}/config.json
type TeamConfig struct {
	Name          string   `json:"name"`
	Description   string   `json:"description"`
	CreatedAt     int64    `json:"createdAt"`
	LeadAgentId   string   `json:"leadAgentId"`
	LeadSessionId string   `json:"leadSessionId"`
	Members       []Member `json:"members"`
}

// Member represents a Claude Code agent in the team
type Member struct {
	AgentId          string   `json:"agentId"`
	Name             string   `json:"name"`
	AgentType        string   `json:"agentType"`
	Model            string   `json:"model"`
	Prompt           string   `json:"prompt"`
	Color            string   `json:"color"`
	PlanModeRequired bool     `json:"planModeRequired"`
	JoinedAt         int64    `json:"joinedAt"`
	TmuxPaneId       string   `json:"tmuxPaneId"`
	Cwd              string   `json:"cwd"`
	Subscriptions    []string `json:"subscriptions"`
	BackendType      string   `json:"backendType"`
}

// InboxMessage represents a message in ~/.claude/teams/{team}/inboxes/{agent}.json
type InboxMessage struct {
	From      string `json:"from"`
	Text      string `json:"text"`
	Summary   string `json:"summary"`
	Timestamp string `json:"timestamp"`
	Color     string `json:"color"`
	Read      bool   `json:"read"`
}

// SendMessageInput is the request body for POST /api/claude-teams/:name/message
type SendMessageInput struct {
	To      string `json:"to"`
	Message string `json:"message"`
}

// TeamResponse is the API response for a team
type TeamResponse struct {
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Members     []MemberResponse  `json:"members"`
}

// MemberResponse is a member in the API response
type MemberResponse struct {
	AgentId   string `json:"agentId"`
	Name      string `json:"name"`
	AgentType string `json:"agentType"`
	Model     string `json:"model"`
	Status    string `json:"status"` // inferred from inbox activity
}
