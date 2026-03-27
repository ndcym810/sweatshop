// internal/claudeteam/watcher.go
package claudeteam

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/gorilla/websocket"
	"github.com/sweatshop/sweatshop/pkg/logger"
)

// EventType represents a WebSocket event type
type EventType string

const (
	EventTeamDiscovered EventType = "team:discovered"
	EventTeamUpdated    EventType = "team:updated"
	EventMessageNew     EventType = "message:new"
	EventMessageRead    EventType = "message:read"
)

// WSMessage is a WebSocket message sent to clients
type WSMessage struct {
	Event     EventType   `json:"event"`
	Timestamp string      `json:"timestamp"`
	Data      interface{} `json:"data"`
}

// Watcher watches Claude Code team files and broadcasts updates
type Watcher struct {
	service   *Service
	watcher   *fsnotify.Watcher
	hub       *WebSocketHub
	claudeDir string
	stopCh    chan struct{}
}

// WebSocketHub manages WebSocket connections
type WebSocketHub struct {
	clients    map[*websocket.Conn]bool
	broadcast  chan []byte
	register   chan *websocket.Conn
	unregister chan *websocket.Conn
	mu         sync.RWMutex
}

// NewWebSocketHub creates a new WebSocket hub
func NewWebSocketHub() *WebSocketHub {
	return &WebSocketHub{
		clients:    make(map[*websocket.Conn]bool),
		broadcast:  make(chan []byte, 256),
		register:   make(chan *websocket.Conn),
		unregister: make(chan *websocket.Conn),
	}
}

// Run starts the WebSocket hub
func (h *WebSocketHub) Run() {
	for {
		select {
		case conn := <-h.register:
			h.mu.Lock()
			h.clients[conn] = true
			h.mu.Unlock()
			logger.Info.Printf("WebSocket client connected. Total: %d", len(h.clients))

		case conn := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[conn]; ok {
				delete(h.clients, conn)
				conn.Close()
			}
			h.mu.Unlock()
			logger.Info.Printf("WebSocket client disconnected. Total: %d", len(h.clients))

		case message := <-h.broadcast:
			h.mu.RLock()
			for conn := range h.clients {
				if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
					conn.Close()
					delete(h.clients, conn)
				}
			}
			h.mu.RUnlock()
		}
	}
}

// Register adds a client to the hub
func (h *WebSocketHub) Register(conn *websocket.Conn) {
	h.register <- conn
}

// Unregister removes a client from the hub
func (h *WebSocketHub) Unregister(conn *websocket.Conn) {
	h.unregister <- conn
}

// Broadcast sends a message to all connected clients
func (h *WebSocketHub) Broadcast(msg WSMessage) {
	data, err := json.Marshal(msg)
	if err != nil {
		logger.Error.Printf("Failed to marshal WebSocket message: %v", err)
		return
	}
	h.broadcast <- data
}

// NewWatcher creates a new file watcher
func NewWatcher(service *Service, hub *WebSocketHub) (*Watcher, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, fmt.Errorf("failed to create fsnotify watcher: %w", err)
	}

	homeDir, _ := os.UserHomeDir()
	claudeDir := filepath.Join(homeDir, ".claude")

	return &Watcher{
		service:   service,
		watcher:   watcher,
		hub:       hub,
		claudeDir: claudeDir,
		stopCh:    make(chan struct{}),
	}, nil
}

// Start begins watching for file changes
func (w *Watcher) Start() error {
	teamsDir := filepath.Join(w.claudeDir, "teams")

	// Watch the teams directory
	if err := w.watcher.Add(teamsDir); err != nil {
		if os.IsNotExist(err) {
			logger.Info.Printf("Teams directory does not exist yet: %s", teamsDir)
			return nil // Will be created when user spawns a team
		}
		return fmt.Errorf("failed to watch teams directory: %w", err)
	}

	go w.eventLoop()
	logger.Info.Println("Started watching Claude Code teams directory")
	return nil
}

// Stop stops the watcher
func (w *Watcher) Stop() {
	close(w.stopCh)
	w.watcher.Close()
}

// eventLoop processes file system events
func (w *Watcher) eventLoop() {
	for {
		select {
		case <-w.stopCh:
			return
		case event, ok := <-w.watcher.Events:
			if !ok {
				return
			}
			w.handleEvent(event)
		case err, ok := <-w.watcher.Errors:
			if !ok {
				return
			}
			logger.Error.Printf("Watcher error: %v", err)
		}
	}
}

// handleEvent processes a single file system event
func (w *Watcher) handleEvent(event fsnotify.Event) {
	// Only care about write and create events
	if event.Op&fsnotify.Write == 0 && event.Op&fsnotify.Create == 0 {
		return
	}

	relPath, err := filepath.Rel(w.claudeDir, event.Name)
	if err != nil {
		return
	}

	// Check if it's a config.json change
	if filepath.Base(relPath) == "config.json" {
		teamName := filepath.Base(filepath.Dir(relPath))
		w.broadcastTeamUpdate(teamName)
		return
	}

	// Check if it's an inbox file change
	if filepath.Base(filepath.Dir(relPath)) == "inboxes" {
		agentName := filepath.Base(event.Name)
		agentName = agentName[:len(agentName)-5] // remove .json
		teamName := filepath.Base(filepath.Dir(filepath.Dir(relPath)))
		w.broadcastNewMessage(teamName, agentName)
		return
	}

	// New team directory created
	if filepath.Dir(relPath) == "teams" && event.Op&fsnotify.Create != 0 {
		teamName := filepath.Base(relPath)
		w.broadcastTeamDiscovered(teamName)
		return
	}
}

func (w *Watcher) broadcastTeamDiscovered(teamName string) {
	w.hub.Broadcast(WSMessage{
		Event: EventTeamDiscovered,
		Data:  map[string]string{"name": teamName},
	})
	logger.Info.Printf("Broadcasted team discovered: %s", teamName)
}

func (w *Watcher) broadcastTeamUpdate(teamName string) {
	config, err := w.service.ReadTeamConfig(teamName)
	if err != nil {
		logger.Error.Printf("Failed to read team config for broadcast: %v", err)
		return
	}

	w.hub.Broadcast(WSMessage{
		Event: EventTeamUpdated,
		Data:  w.service.ToTeamResponse(config),
	})
	logger.Info.Printf("Broadcasted team update: %s", teamName)
}

func (w *Watcher) broadcastNewMessage(teamName, agentName string) {
	w.hub.Broadcast(WSMessage{
		Event: EventMessageNew,
		Data: map[string]string{
			"team":  teamName,
			"agent": agentName,
		},
	})
	logger.Info.Printf("Broadcasted new message for %s/%s", teamName, agentName)
}
