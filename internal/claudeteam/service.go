// internal/claudeteam/service.go
package claudeteam

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/sweatshop/sweatshop/pkg/logger"
)

// Service handles Claude Code team file operations
type Service struct {
	claudeDir string // ~/.claude path, injected for testability
}

// NewService creates a new Claude team service
func NewService() *Service {
	homeDir, _ := os.UserHomeDir()
	return &Service{
		claudeDir: filepath.Join(homeDir, ".claude"),
	}
}

// NewServiceWithDir creates a service with a custom claude directory (for testing)
func NewServiceWithDir(claudeDir string) *Service {
	return &Service{claudeDir: claudeDir}
}

// DiscoverTeams scans ~/.claude/teams/ and returns all team names
func (s *Service) DiscoverTeams() ([]string, error) {
	teamsDir := filepath.Join(s.claudeDir, "teams")

	entries, err := os.ReadDir(teamsDir)
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return nil, fmt.Errorf("failed to read teams directory: %w", err)
	}

	teams := make([]string, 0)
	for _, entry := range entries {
		if entry.IsDir() {
			configPath := filepath.Join(teamsDir, entry.Name(), "config.json")
			if _, err := os.Stat(configPath); err == nil {
				teams = append(teams, entry.Name())
			}
		}
	}

	sort.Strings(teams)
	return teams, nil
}

// ReadTeamConfig reads a team's config.json
func (s *Service) ReadTeamConfig(teamName string) (*TeamConfig, error) {
	configPath := filepath.Join(s.claudeDir, "teams", teamName, "config.json")

	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to read team config: %w", err)
	}

	var config TeamConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse team config: %w", err)
	}

	return &config, nil
}

// ReadInbox reads an agent's inbox messages
func (s *Service) ReadInbox(teamName, agentName string) ([]InboxMessage, error) {
	inboxPath := filepath.Join(s.claudeDir, "teams", teamName, "inboxes", agentName+".json")

	data, err := os.ReadFile(inboxPath)
	if err != nil {
		if os.IsNotExist(err) {
			return []InboxMessage{}, nil
		}
		return nil, fmt.Errorf("failed to read inbox: %w", err)
	}

	var messages []InboxMessage
	if err := json.Unmarshal(data, &messages); err != nil {
		return nil, fmt.Errorf("failed to parse inbox: %w", err)
	}

	// Sort by timestamp descending (newest first)
	sort.Slice(messages, func(i, j int) bool {
		return messages[i].Timestamp > messages[j].Timestamp
	})

	return messages, nil
}

// WriteMessage appends a message to an agent's inbox
func (s *Service) WriteMessage(teamName, to, from, message string) error {
	inboxDir := filepath.Join(s.claudeDir, "teams", teamName, "inboxes")
	inboxPath := filepath.Join(inboxDir, to+".json")

	// Ensure directory exists
	if err := os.MkdirAll(inboxDir, 0755); err != nil {
		return fmt.Errorf("failed to create inbox directory: %w", err)
	}

	// Read existing messages
	messages := make([]InboxMessage, 0)
	data, err := os.ReadFile(inboxPath)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to read existing inbox: %w", err)
	}
	if len(data) > 0 {
		if err := json.Unmarshal(data, &messages); err != nil {
			logger.Info.Printf("Warning: failed to parse existing inbox, starting fresh: %v", err)
			messages = []InboxMessage{}
		}
	}

	// Append new message
	newMsg := InboxMessage{
		From:      from,
		Text:      message,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Color:     "blue", // default color for dashboard messages
		Read:      false,
	}
	messages = append(messages, newMsg)

	// Write back with retry for file locking
	return s.writeInboxWithRetry(inboxPath, messages, 3)
}

// writeInboxWithRetry writes inbox with exponential backoff for file locking
func (s *Service) writeInboxWithRetry(path string, messages []InboxMessage, maxRetries int) error {
	var lastErr error
	for i := 0; i < maxRetries; i++ {
		data, err := json.MarshalIndent(messages, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal messages: %w", err)
		}

		if err := os.WriteFile(path, data, 0644); err != nil {
			lastErr = err
			if strings.Contains(err.Error(), "used by another process") {
				time.Sleep(time.Duration(100*(i+1)) * time.Millisecond)
				continue
			}
			return fmt.Errorf("failed to write inbox: %w", err)
		}
		return nil
	}
	return fmt.Errorf("failed to write inbox after %d retries: %w", maxRetries, lastErr)
}

// MarkMessageRead marks a message as read by updating the inbox file
func (s *Service) MarkMessageRead(teamName, agentName, timestamp string) error {
	inboxPath := filepath.Join(s.claudeDir, "teams", teamName, "inboxes", agentName+".json")

	data, err := os.ReadFile(inboxPath)
	if err != nil {
		return fmt.Errorf("failed to read inbox: %w", err)
	}

	var messages []InboxMessage
	if err := json.Unmarshal(data, &messages); err != nil {
		return fmt.Errorf("failed to parse inbox: %w", err)
	}

	// Find and mark the message
	found := false
	for i := range messages {
		if messages[i].Timestamp == timestamp {
			messages[i].Read = true
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("message not found: %s", timestamp)
	}

	return s.writeInboxWithRetry(inboxPath, messages, 3)
}

// ToTeamResponse converts TeamConfig to API response format
func (s *Service) ToTeamResponse(config *TeamConfig) *TeamResponse {
	members := make([]MemberResponse, len(config.Members))
	for i, m := range config.Members {
		members[i] = MemberResponse{
			AgentId:   m.AgentId,
			Name:      m.Name,
			AgentType: m.AgentType,
			Model:     m.Model,
			Status:    "unknown", // TODO: infer from inbox activity
		}
	}

	return &TeamResponse{
		Name:        config.Name,
		Description: config.Description,
		Members:     members,
	}
}
