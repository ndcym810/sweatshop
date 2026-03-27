package claudeteam

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestService_DiscoverTeams(t *testing.T) {
	// Create temp directory for test
	tmpDir, err := os.MkdirTemp("", "claude-teams-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	teamsDir := filepath.Join(tmpDir, "teams")
	os.MkdirAll(teamsDir, 0755)

	svc := NewServiceWithDir(tmpDir)

	t.Run("empty directory returns empty list", func(t *testing.T) {
		teams, err := svc.DiscoverTeams()
		if err != nil {
			t.Errorf("DiscoverTeams() error = %v", err)
		}
		if len(teams) != 0 {
			t.Errorf("DiscoverTeams() = %v, want empty list", teams)
		}
	})

	t.Run("discovers teams with valid config.json", func(t *testing.T) {
		// Create team directory with config
		teamDir := filepath.Join(teamsDir, "test-team")
		os.MkdirAll(teamDir, 0755)
		config := TeamConfig{Name: "test-team"}
		configData, _ := json.Marshal(config)
		os.WriteFile(filepath.Join(teamDir, "config.json"), configData, 0644)

		teams, err := svc.DiscoverTeams()
		if err != nil {
			t.Errorf("DiscoverTeams() error = %v", err)
		}
		if len(teams) != 1 || teams[0] != "test-team" {
			t.Errorf("DiscoverTeams() = %v, want [test-team]", teams)
		}
	})

	t.Run("ignores directories without config.json", func(t *testing.T) {
		// Create directory without config
		os.MkdirAll(filepath.Join(teamsDir, "no-config-team"), 0755)

		teams, err := svc.DiscoverTeams()
		if err != nil {
			t.Errorf("DiscoverTeams() error = %v", err)
		}
		if len(teams) != 1 {
			t.Errorf("DiscoverTeams() = %v, want 1 team", teams)
		}
	})
}

func TestService_ReadTeamConfig(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "claude-teams-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	teamsDir := filepath.Join(tmpDir, "teams")
	os.MkdirAll(teamsDir, 0755)

	svc := NewServiceWithDir(tmpDir)

	t.Run("returns nil for non-existent team", func(t *testing.T) {
		config, err := svc.ReadTeamConfig("nonexistent")
		if err != nil {
			t.Errorf("ReadTeamConfig() error = %v", err)
		}
		if config != nil {
			t.Errorf("ReadTeamConfig() = %v, want nil", config)
		}
	})

	t.Run("reads valid config", func(t *testing.T) {
		teamDir := filepath.Join(teamsDir, "test-team")
		os.MkdirAll(teamDir, 0755)
		expected := TeamConfig{
			Name:          "test-team",
			Description:   "Test description",
			LeadAgentId:   "lead@test-team",
		}
		configData, _ := json.Marshal(expected)
		os.WriteFile(filepath.Join(teamDir, "config.json"), configData, 0644)

		config, err := svc.ReadTeamConfig("test-team")
		if err != nil {
			t.Errorf("ReadTeamConfig() error = %v", err)
		}
		if config == nil {
			t.Fatal("ReadTeamConfig() returned nil")
		}
		if config.Name != expected.Name {
			t.Errorf("Name = %v, want %v", config.Name, expected.Name)
		}
		if config.Description != expected.Description {
			t.Errorf("Description = %v, want %v", config.Description, expected.Description)
		}
	})
}

func TestService_ReadInbox(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "claude-teams-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	teamsDir := filepath.Join(tmpDir, "teams")
	os.MkdirAll(teamsDir, 0755)

	svc := NewServiceWithDir(tmpDir)

	t.Run("returns empty list for non-existent inbox", func(t *testing.T) {
		messages, err := svc.ReadInbox("nonexistent", "agent")
		if err != nil {
			t.Errorf("ReadInbox() error = %v", err)
		}
		if len(messages) != 0 {
			t.Errorf("ReadInbox() = %v, want empty list", messages)
		}
	})

	t.Run("reads messages sorted by timestamp desc", func(t *testing.T) {
		teamDir := filepath.Join(teamsDir, "test-team")
		inboxDir := filepath.Join(teamDir, "inboxes")
		os.MkdirAll(inboxDir, 0755)

		messages := []InboxMessage{
			{From: "lead", Text: "First", Timestamp: "2026-03-26T10:00:00Z"},
			{From: "lead", Text: "Second", Timestamp: "2026-03-26T11:00:00Z"},
		}
		data, _ := json.Marshal(messages)
		os.WriteFile(filepath.Join(inboxDir, "agent.json"), data, 0644)

		result, err := svc.ReadInbox("test-team", "agent")
		if err != nil {
			t.Errorf("ReadInbox() error = %v", err)
		}
		if len(result) != 2 {
			t.Fatalf("ReadInbox() returned %d messages, want 2", len(result))
		}
		// Should be sorted newest first
		if result[0].Timestamp < result[1].Timestamp {
			t.Error("Messages not sorted by timestamp descending")
		}
	})
}

func TestService_WriteMessage(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "claude-teams-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	svc := NewServiceWithDir(tmpDir)

	t.Run("creates inbox file if not exists", func(t *testing.T) {
		err := svc.WriteMessage("test-team", "agent", "sender", "test message")
		if err != nil {
			t.Errorf("WriteMessage() error = %v", err)
		}

		// Verify file was created
		inboxPath := filepath.Join(tmpDir, "teams", "test-team", "inboxes", "agent.json")
		if _, err := os.Stat(inboxPath); os.IsNotExist(err) {
			t.Error("Inbox file was not created")
		}

		// Verify message content
		data, _ := os.ReadFile(inboxPath)
		var messages []InboxMessage
		json.Unmarshal(data, &messages)
		if len(messages) != 1 {
			t.Fatalf("Expected 1 message, got %d", len(messages))
		}
		if messages[0].From != "sender" {
			t.Errorf("From = %v, want sender", messages[0].From)
		}
		if messages[0].Text != "test message" {
			t.Errorf("Text = %v, want 'test message'", messages[0].Text)
		}
	})

	t.Run("appends to existing inbox", func(t *testing.T) {
		// First message
		svc.WriteMessage("test-team", "agent", "sender1", "message1")
		// Second message
		svc.WriteMessage("test-team", "agent", "sender2", "message2")

		inboxPath := filepath.Join(tmpDir, "teams", "test-team", "inboxes", "agent.json")
		data, _ := os.ReadFile(inboxPath)
		var messages []InboxMessage
		json.Unmarshal(data, &messages)
		if len(messages) != 3 { // 1 from previous test + 2 from this test
			t.Errorf("Expected 3 messages, got %d", len(messages))
		}
	})
}

func TestService_MarkMessageRead(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "claude-teams-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	svc := NewServiceWithDir(tmpDir)

	t.Run("marks message as read", func(t *testing.T) {
		// Setup inbox with unread message
		inboxDir := filepath.Join(tmpDir, "teams", "test-team", "inboxes")
		os.MkdirAll(inboxDir, 0755)
		messages := []InboxMessage{
			{From: "lead", Timestamp: "2026-03-26T10:00:00Z", Read: false},
		}
		data, _ := json.Marshal(messages)
		os.WriteFile(filepath.Join(inboxDir, "agent.json"), data, 0644)

		err := svc.MarkMessageRead("test-team", "agent", "2026-03-26T10:00:00Z")
		if err != nil {
			t.Errorf("MarkMessageRead() error = %v", err)
		}

		// Verify message was marked read
		result, _ := svc.ReadInbox("test-team", "agent")
		if len(result) != 1 || !result[0].Read {
			t.Error("Message was not marked as read")
		}
	})

	t.Run("returns error for non-existent message", func(t *testing.T) {
		err := svc.MarkMessageRead("test-team", "agent", "nonexistent-timestamp")
		if err == nil {
			t.Error("Expected error for non-existent message")
		}
	})
}
