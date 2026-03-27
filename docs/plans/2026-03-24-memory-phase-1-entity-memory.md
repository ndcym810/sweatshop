# Sweatshop Memory Phase 1: Entity Memory Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Build working entity memory system with key-value storage and CRUD operations for structured facts.

**Architecture:** Entity memory stores structured facts as key-value pairs scoped by context (global/team/project). SQLite backend with direct lookup queries. Go package with REST API. Integrated with Lead agent for memory extraction.

**Tech Stack:** Go 1.26, Echo, SQLite, JSON for values

---

## Task 1: Create Entity Memory Database Schema

**Files:**
- Modify: `internal/shared/db/db.go`

**Step 1: Add entity_memory table to schema**

Add to the `runMigrations()` function in `internal/shared/db/db.go`:

```go
// Add after existing tables

// Entity Memory (Layer 2)
CREATE TABLE IF NOT EXISTS entity_memory (
    id TEXT PRIMARY KEY,

    -- Entity identification
    scope TEXT NOT NULL,
    scope_id TEXT,
    entity_type TEXT NOT NULL,
    entity_id TEXT,

    -- Key-value pair
    key TEXT NOT NULL,
    value TEXT NOT NULL,

    -- Metadata
    confidence REAL DEFAULT 1.0,
    source TEXT DEFAULT 'stated',
    ttl_days INTEGER,

    -- Timestamps
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,

    -- Unique constraint
    UNIQUE(scope, scope_id, entity_type, entity_id, key)
);

-- Indexes for entity_memory
CREATE INDEX IF NOT EXISTS idx_entity_scope ON entity_memory(scope, scope_id);
CREATE INDEX IF NOT EXISTS idx_entity_type ON entity_memory(entity_type, entity_id);
CREATE INDEX IF NOT EXISTS idx_entity_key ON entity_memory(scope, scope_id, entity_type, key);
```

**Step 2: Verify database migration**

Run:
```bash
rm -f data/sweatshop.db
go build -o bin/sweatshop ./cmd/server
./bin/sweatshop -port 8081 &
sleep 2
sqlite3 data/sweatshop.db ".schema entity_memory"
pkill -f sweatshop
```
Expected: Table schema displayed

**Step 3: Commit**

```bash
git add internal/shared/db/db.go
git commit -m "feat: add entity_memory table schema

Co-Authored-By: Claude Opus 4.6 <noreply@anthropic.com>"
```

---

## Task 2: Create Entity Memory Model

**Files:**
- Create: `internal/memory/entity/model.go`

**Step 1: Create entity memory model**

```go
// internal/memory/entity/model.go
package entity

import (
	"database/sql"
	"time"

	"github.com/sweatshop/sweatshop/internal/shared/db"
)

// Scope defines the visibility scope of an entity memory
type Scope string

const (
	ScopeGlobal  Scope = "global"
	ScopeTeam    Scope = "team"
	ScopeProject Scope = "project"
	ScopeAgent   Scope = "agent"
)

// Source indicates how the memory was obtained
type Source string

const (
	SourceStated   Source = "stated"   // Explicitly mentioned
	SourceObserved Source = "observed" // Seen in behavior
	SourceInferred Source = "inferred" // Deduced from context
)

// EntityType defines what kind of entity the memory is about
type EntityType string

const (
	EntityTypeProject   EntityType = "project"
	EntityTypeTeammate  EntityType = "teammate"
	EntityTypeUser      EntityType = "user"
	EntityTypeTeam      EntityType = "team"
	EntityTypeCodebase  EntityType = "codebase"
	EntityTypeTask      EntityType = "task"
)

// Memory represents a single entity memory entry
type Memory struct {
	ID          string    `json:"id"`
	Scope       Scope     `json:"scope"`
	ScopeID     string    `json:"scopeId,omitempty"`
	EntityType  EntityType `json:"entityType"`
	EntityID    string    `json:"entityId,omitempty"`
	Key         string    `json:"key"`
	Value       string    `json:"value"` // JSON-encoded value
	Confidence  float64   `json:"confidence"`
	Source      Source    `json:"source"`
	TTLDays     *int      `json:"ttlDays,omitempty"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// CreateInput represents input for creating an entity memory
type CreateInput struct {
	Scope      Scope      `json:"scope" validate:"required"`
	ScopeID    string     `json:"scopeId"`
	EntityType EntityType `json:"entityType" validate:"required"`
	EntityID   string     `json:"entityId"`
	Key        string     `json:"key" validate:"required"`
	Value      string     `json:"value" validate:"required"`
	Confidence float64    `json:"confidence"`
	Source     Source     `json:"source"`
	TTLDays    *int       `json:"ttlDays"`
}

// QueryInput represents parameters for querying entity memories
type QueryInput struct {
	Scope      Scope
	ScopeID    string
	EntityType EntityType
	EntityID   string
	Key        string
}

// Get retrieves a specific entity memory
func Get(scope Scope, scopeID string, entityType EntityType, entityID string, key string) (*Memory, error) {
	var m Memory
	var scopeIDNull, entityIDNull sql.NullString
	var ttlDaysNull sql.NullInt64

	err := db.DB.QueryRow(`
		SELECT id, scope, scope_id, entity_type, entity_id, key, value, confidence, source, ttl_days, created_at, updated_at
		FROM entity_memory
		WHERE scope = ? AND COALESCE(scope_id, '') = ?
		  AND entity_type = ? AND COALESCE(entity_id, '') = ?
		  AND key = ?
	`, string(scope), scopeID, string(entityType), entityID, key).Scan(
		&m.ID, &m.Scope, &scopeIDNull, &m.EntityType, &entityIDNull,
		&m.Key, &m.Value, &m.Confidence, &m.Source, &ttlDaysNull,
		&m.CreatedAt, &m.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if scopeIDNull.Valid {
		m.ScopeID = scopeIDNull.String
	}
	if entityIDNull.Valid {
		m.EntityID = entityIDNull.String
	}
	if ttlDaysNull.Valid {
		days := int(ttlDaysNull.Int64)
		m.TTLDays = &days
	}

	return &m, nil
}

// GetAll retrieves all entity memories matching a query
func GetAll(input QueryInput) ([]Memory, error) {
	query := `SELECT id, scope, scope_id, entity_type, entity_id, key, value, confidence, source, ttl_days, created_at, updated_at
			  FROM entity_memory WHERE 1=1`
	args := []interface{}{}

	if input.Scope != "" {
		query += " AND scope = ?"
		args = append(args, string(input.Scope))
	}
	if input.ScopeID != "" {
		query += " AND scope_id = ?"
		args = append(args, input.ScopeID)
	}
	if input.EntityType != "" {
		query += " AND entity_type = ?"
		args = append(args, string(input.EntityType))
	}
	if input.EntityID != "" {
		query += " AND entity_id = ?"
		args = append(args, input.EntityID)
	}
	if input.Key != "" {
		query += " AND key = ?"
		args = append(args, input.Key)
	}

	query += " ORDER BY updated_at DESC"

	rows, err := db.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var memories []Memory
	for rows.Next() {
		var m Memory
		var scopeIDNull, entityIDNull sql.NullString
		var ttlDaysNull sql.NullInt64

		if err := rows.Scan(
			&m.ID, &m.Scope, &scopeIDNull, &m.EntityType, &entityIDNull,
			&m.Key, &m.Value, &m.Confidence, &m.Source, &ttlDaysNull,
			&m.CreatedAt, &m.UpdatedAt,
		); err != nil {
			return nil, err
		}

		if scopeIDNull.Valid {
			m.ScopeID = scopeIDNull.String
		}
		if entityIDNull.Valid {
			m.EntityID = entityIDNull.String
		}
		if ttlDaysNull.Valid {
			days := int(ttlDaysNull.Int64)
			m.TTLDays = &days
		}
		memories = append(memories, m)
	}

	return memories, nil
}

// Create creates a new entity memory
func Create(m *Memory) error {
	var scopeID, entityID interface{}
	var ttlDays interface{}

	if m.ScopeID != "" {
		scopeID = m.ScopeID
	}
	if m.EntityID != "" {
		entityID = m.EntityID
	}
	if m.TTLDays != nil {
		ttlDays = *m.TTLDays
	}

	_, err := db.DB.Exec(`
		INSERT INTO entity_memory (id, scope, scope_id, entity_type, entity_id, key, value, confidence, source, ttl_days, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, m.ID, string(m.Scope), scopeID, string(m.EntityType), entityID,
		m.Key, m.Value, m.Confidence, string(m.Source), ttlDays,
		m.CreatedAt, m.UpdatedAt)

	return err
}

// Upsert creates or updates an entity memory
func Upsert(m *Memory) error {
	var scopeID, entityID interface{}
	var ttlDays interface{}

	if m.ScopeID != "" {
		scopeID = m.ScopeID
	}
	if m.EntityID != "" {
		entityID = m.EntityID
	}
	if m.TTLDays != nil {
		ttlDays = *m.TTLDays
	}

	_, err := db.DB.Exec(`
		INSERT INTO entity_memory (id, scope, scope_id, entity_type, entity_id, key, value, confidence, source, ttl_days, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(scope, scope_id, entity_type, entity_id, key) DO UPDATE SET
			value = excluded.value,
			confidence = excluded.confidence,
			source = excluded.source,
			ttl_days = excluded.ttl_days,
			updated_at = excluded.updated_at
	`, m.ID, string(m.Scope), scopeID, string(m.EntityType), entityID,
		m.Key, m.Value, m.Confidence, string(m.Source), ttlDays,
		m.CreatedAt, m.UpdatedAt)

	return err
}

// Delete deletes an entity memory
func Delete(scope Scope, scopeID string, entityType EntityType, entityID string, key string) error {
	_, err := db.DB.Exec(`
		DELETE FROM entity_memory
		WHERE scope = ? AND COALESCE(scope_id, '') = ?
		  AND entity_type = ? AND COALESCE(entity_id, '') = ?
		  AND key = ?
	`, string(scope), scopeID, string(entityType), entityID, key)
	return err
}

// DeleteByID deletes an entity memory by ID
func DeleteByID(id string) error {
	_, err := db.DB.Exec("DELETE FROM entity_memory WHERE id = ?", id)
	return err
}
```

**Step 2: Verify compilation**

Run:
```bash
mkdir -p internal/memory/entity
go build ./internal/memory/entity/...
```
Expected: No errors

**Step 3: Commit**

```bash
git add internal/memory/entity/model.go
git commit -m "feat: add entity memory model with CRUD operations

Co-Authored-By: Claude Opus 4.6 <noreply@anthropic.com>"
```

---

## Task 3: Create Entity Memory Service

**Files:**
- Create: `internal/memory/entity/service.go`

**Step 1: Create entity memory service**

```go
// internal/memory/entity/service.go
package entity

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/sweatshop/sweatshop/pkg/uuid"
)

// Service provides entity memory business logic
type Service struct{}

// NewService creates a new entity memory service
func NewService() *Service {
	return &Service{}
}

// Get retrieves a specific entity memory
func (s *Service) Get(scope Scope, scopeID string, entityType EntityType, entityID string, key string) (*Memory, error) {
	return Get(scope, scopeID, entityType, entityID, key)
}

// Query retrieves entity memories matching criteria
func (s *Service) Query(input QueryInput) ([]Memory, error) {
	return GetAll(input)
}

// Set creates or updates an entity memory
func (s *Service) Set(input CreateInput) (*Memory, error) {
	// Set defaults
	if input.Confidence == 0 {
		input.Confidence = 1.0
	}
	if input.Source == "" {
		input.Source = SourceStated
	}

	now := time.Now()
	m := &Memory{
		ID:         uuid.New(),
		Scope:      input.Scope,
		ScopeID:    input.ScopeID,
		EntityType: input.EntityType,
		EntityID:   input.EntityID,
		Key:        input.Key,
		Value:      input.Value,
		Confidence: input.Confidence,
		Source:     input.Source,
		TTLDays:    input.TTLDays,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	if err := Upsert(m); err != nil {
		return nil, err
	}

	// Reload to get actual stored values
	return Get(input.Scope, input.ScopeID, input.EntityType, input.EntityID, input.Key)
}

// SetTyped creates or updates an entity memory with a typed value (auto-encodes to JSON)
func (s *Service) SetTyped(scope Scope, scopeID string, entityType EntityType, entityID string, key string, value interface{}, source Source) error {
	jsonValue, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to encode value: %w", err)
	}

	input := CreateInput{
		Scope:      scope,
		ScopeID:    scopeID,
		EntityType: entityType,
		EntityID:   entityID,
		Key:        key,
		Value:      string(jsonValue),
		Source:     source,
	}

	_, err = s.Set(input)
	return err
}

// GetTyped retrieves an entity memory and decodes the JSON value
func (s *Service) GetTyped(scope Scope, scopeID string, entityType EntityType, entityID string, key string, dest interface{}) error {
	m, err := Get(scope, scopeID, entityType, entityID, key)
	if err != nil {
		return err
	}
	if m == nil {
		return fmt.Errorf("memory not found")
	}

	return json.Unmarshal([]byte(m.Value), dest)
}

// Delete removes an entity memory
func (s *Service) Delete(scope Scope, scopeID string, entityType EntityType, entityID string, key string) error {
	return Delete(scope, scopeID, entityType, entityID, key)
}

// DeleteByID removes an entity memory by ID
func (s *Service) DeleteByID(id string) error {
	return DeleteByID(id)
}

// GetAllForEntity retrieves all memories for a specific entity
func (s *Service) GetAllForEntity(scope Scope, scopeID string, entityType EntityType, entityID string) (map[string]string, error) {
	memories, err := GetAll(QueryInput{
		Scope:      scope,
		ScopeID:    scopeID,
		EntityType: entityType,
		EntityID:   entityID,
	})
	if err != nil {
		return nil, err
	}

	result := make(map[string]string)
	for _, m := range memories {
		result[m.Key] = m.Value
	}
	return result, nil
}

// SetBatch creates or updates multiple entity memories
func (s *Service) SetBatch(inputs []CreateInput) error {
	for _, input := range inputs {
		if _, err := s.Set(input); err != nil {
			return err
		}
	}
	return nil
}
```

**Step 2: Verify compilation**

Run:
```bash
go build ./internal/memory/entity/...
```
Expected: No errors

**Step 3: Commit**

```bash
git add internal/memory/entity/service.go
git commit -m "feat: add entity memory service with typed helpers

Co-Authored-By: Claude Opus 4.6 <noreply@anthropic.com>"
```

---

## Task 4: Create Entity Memory Handler

**Files:**
- Create: `internal/memory/entity/handler.go`

**Step 1: Create entity memory handler**

```go
// internal/memory/entity/handler.go
package entity

import (
	"github.com/labstack/echo/v4"
	"github.com/sweatshop/sweatshop/internal/shared/response"
)

// Handler handles entity memory HTTP requests
type Handler struct {
	service *Service
}

// NewHandler creates a new entity memory handler
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// List handles GET /api/memory/entity
// Query params: scope, scopeId, entityType, entityId, key
func (h *Handler) List(c echo.Context) error {
	var input QueryInput

	if scope := c.QueryParam("scope"); scope != "" {
		input.Scope = Scope(scope)
	}
	input.ScopeID = c.QueryParam("scopeId")
	if entityType := c.QueryParam("entityType"); entityType != "" {
		input.EntityType = EntityType(entityType)
	}
	input.EntityID = c.QueryParam("entityId")
	input.Key = c.QueryParam("key")

	memories, err := h.service.Query(input)
	if err != nil {
		return response.InternalError(c, "Failed to query entity memories")
	}
	return response.OK(c, memories)
}

// Get handles GET /api/memory/entity/:id
func (h *Handler) Get(c echo.Context) error {
	id := c.Param("id")

	// Query by ID - we need to get by exact key
	// For simplicity, use the query endpoint with specific params
	return response.NotFound(c, "Use query parameters to retrieve specific memories")
}

// Set handles POST /api/memory/entity
func (h *Handler) Set(c echo.Context) error {
	var input CreateInput
	if err := c.Bind(&input); err != nil {
		return response.BadRequest(c, "Invalid request body")
	}

	// Validate required fields
	if input.Scope == "" || input.EntityType == "" || input.Key == "" || input.Value == "" {
		return response.BadRequest(c, "scope, entityType, key, and value are required")
	}

	memory, err := h.service.Set(input)
	if err != nil {
		return response.InternalError(c, "Failed to set entity memory")
	}
	return response.OK(c, memory)
}

// SetBatch handles POST /api/memory/entity/batch
func (h *Handler) SetBatch(c echo.Context) error {
	var inputs []CreateInput
	if err := c.Bind(&inputs); err != nil {
		return response.BadRequest(c, "Invalid request body")
	}

	if err := h.service.SetBatch(inputs); err != nil {
		return response.InternalError(c, "Failed to set entity memories")
	}
	return response.OK(c, map[string]string{"status": "ok"})
}

// Delete handles DELETE /api/memory/entity
// Query params: scope, scopeId, entityType, entityId, key
func (h *Handler) Delete(c echo.Context) error {
	scope := Scope(c.QueryParam("scope"))
	scopeID := c.QueryParam("scopeId")
	entityType := EntityType(c.QueryParam("entityType"))
	entityID := c.QueryParam("entityId")
	key := c.QueryParam("key")

	if scope == "" || entityType == "" || key == "" {
		return response.BadRequest(c, "scope, entityType, and key are required")
	}

	if err := h.service.Delete(scope, scopeID, entityType, entityID, key); err != nil {
		return response.InternalError(c, "Failed to delete entity memory")
	}
	return response.NoContent(c)
}

// DeleteByID handles DELETE /api/memory/entity/:id
func (h *Handler) DeleteByID(c echo.Context) error {
	id := c.Param("id")

	if err := h.service.DeleteByID(id); err != nil {
		return response.InternalError(c, "Failed to delete entity memory")
	}
	return response.NoContent(c)
}

// GetEntity handles GET /api/memory/entity/entity
// Returns all memories for a specific entity as key-value map
func (h *Handler) GetEntity(c echo.Context) error {
	scope := Scope(c.QueryParam("scope"))
	scopeID := c.QueryParam("scopeId")
	entityType := EntityType(c.QueryParam("entityType"))
	entityID := c.QueryParam("entityId")

	if scope == "" || entityType == "" {
		return response.BadRequest(c, "scope and entityType are required")
	}

	memories, err := h.service.GetAllForEntity(scope, scopeID, entityType, entityID)
	if err != nil {
		return response.InternalError(c, "Failed to get entity memories")
	}
	return response.OK(c, memories)
}

// RegisterRoutes registers entity memory routes
func (h *Handler) RegisterRoutes(g *echo.Group) {
	g.GET("", h.List)
	g.GET("/entity", h.GetEntity)
	g.POST("", h.Set)
	g.POST("/batch", h.SetBatch)
	g.DELETE("", h.Delete)
	g.DELETE("/:id", h.DeleteByID)
}
```

**Step 2: Verify compilation**

Run:
```bash
go build ./internal/memory/entity/...
```
Expected: No errors

**Step 3: Commit**

```bash
git add internal/memory/entity/handler.go
git commit -m "feat: add entity memory REST API handler

Co-Authored-By: Claude Opus 4.6 <noreply@anthropic.com>"
```

---

## Task 5: Integrate Entity Memory into App

**Files:**
- Modify: `internal/app/app.go`
- Modify: `internal/app/router.go`

**Step 1: Update app.go**

```go
// internal/app/app.go
package app

import (
	"github.com/sweatshop/sweatshop/internal/department"
	"github.com/sweatshop/sweatshop/internal/memory/entity"
	"github.com/sweatshop/sweatshop/internal/project"
	"github.com/sweatshop/sweatshop/internal/task"
	"github.com/sweatshop/sweatshop/internal/team"
	"github.com/sweatshop/sweatshop/internal/template"
	"github.com/sweatshop/sweatshop/pkg/logger"
)

// App holds application dependencies
type App struct {
	TeamHandler       *team.Handler
	ProjectHandler    *project.Handler
	DepartmentHandler *department.Handler
	TaskHandler       *task.Handler
	TemplateHandler   *template.Handler
	EntityHandler     *entity.Handler

	templateService *template.Service
	EntityService   *entity.Service
}

// New creates a new App instance
func New(templateDir string) *App {
	// Initialize services
	teamSvc := team.NewService()
	projectSvc := project.NewService()
	departmentSvc := department.NewService()
	taskSvc := task.NewService()
	templateSvc := template.NewService(templateDir)
	entitySvc := entity.NewService()

	// Initialize handlers
	app := &App{
		TeamHandler:       team.NewHandler(teamSvc),
		ProjectHandler:    project.NewHandler(projectSvc),
		DepartmentHandler: department.NewHandler(departmentSvc),
		TaskHandler:       task.NewHandler(taskSvc),
		TemplateHandler:   template.NewHandler(templateSvc),
		EntityHandler:     entity.NewHandler(entitySvc),
		templateService:   templateSvc,
		EntityService:     entitySvc,
	}

	logger.Info.Println("Application initialized")
	return app
}

// LoadDefaultTemplates loads default templates from YAML files
func (a *App) LoadDefaultTemplates() error {
	return a.templateService.LoadDefaults()
}
```

**Step 2: Update router.go**

```go
// internal/app/router.go
package app

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// SetupRouter creates and configures the Echo router
func (a *App) SetupRouter() *echo.Echo {
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// API routes
	api := e.Group("/api")

	// Memory routes
	memory := api.Group("/memory")
	a.EntityHandler.RegisterRoutes(memory.Group("/entity"))

	// Template routes (top-level)
	a.TemplateHandler.RegisterRoutes(api.Group("/templates"))

	// Team routes
	a.TeamHandler.RegisterRoutes(api.Group("/teams"))

	// Nested routes under teams
	teams := api.Group("/teams/:teamId")
	a.ProjectHandler.RegisterRoutes(teams.Group("/projects"))
	a.DepartmentHandler.RegisterRoutes(teams.Group("/departments"))
	a.TaskHandler.RegisterRoutes(teams.Group("/tasks"))

	return e
}
```

**Step 3: Verify build**

Run:
```bash
go build -o bin/sweatshop ./cmd/server
```
Expected: No errors

**Step 4: Commit**

```bash
git add internal/app/
git commit -m "feat: integrate entity memory into application

Co-Authored-By: Claude Opus 4.6 <noreply@anthropic.com>"
```

---

## Task 6: Test Entity Memory API

**Step 1: Start the server**

Run:
```bash
rm -f data/sweatshop.db
./bin/sweatshop -port 8080 &
sleep 2
```

**Step 2: Test creating entity memory**

Run:
```bash
curl -X POST http://localhost:8080/api/memory/entity \
  -H "Content-Type: application/json" \
  -d '{
    "scope": "team",
    "scopeId": "team-1",
    "entityType": "project",
    "entityId": "proj-1",
    "key": "language",
    "value": "\"Go\"",
    "source": "stated"
  }'
```
Expected: JSON response with created memory

**Step 3: Test querying entity memories**

Run:
```bash
curl "http://localhost:8080/api/memory/entity?scope=team&scopeId=team-1&entityType=project&entityId=proj-1"
```
Expected: JSON array with the created memory

**Step 4: Test getting all for entity**

Run:
```bash
curl "http://localhost:8080/api/memory/entity/entity?scope=team&scopeId=team-1&entityType=project&entityId=proj-1"
```
Expected: JSON object with key-value pairs

**Step 5: Test batch create**

Run:
```bash
curl -X POST http://localhost:8080/api/memory/entity/batch \
  -H "Content-Type: application/json" \
  -d '[
    {"scope": "team", "scopeId": "team-1", "entityType": "project", "entityId": "proj-1", "key": "framework", "value": "\"Echo\""},
    {"scope": "team", "scopeId": "team-1", "entityType": "project", "entityId": "proj-1", "key": "database", "value": "\"SQLite\""}
  ]'
```
Expected: `{"status":"ok"}`

**Step 6: Verify all memories exist**

Run:
```bash
curl "http://localhost:8080/api/memory/entity/entity?scope=team&scopeId=team-1&entityType=project&entityId=proj-1"
```
Expected: `{"data":{"database":"\"SQLite\"","framework":"\"Echo\"","language":"\"Go\""}}`

**Step 7: Cleanup**

Run:
```bash
pkill -f sweatshop || true
```

**Step 8: Commit**

```bash
git add -A
git commit -m "test: verify entity memory API endpoints

Co-Authored-By: Claude Opus 4.6 <noreply@anthropic.com>"
```

---

## Task 7: Add Memory Extraction Prompt

**Files:**
- Create: `internal/memory/extractor/prompt.go`
- Create: `internal/memory/extractor/extractor.go`

**Step 1: Create extraction prompt**

```go
// internal/memory/extractor/prompt.go
package extractor

// ExtractionPrompt is the system prompt for memory extraction
const ExtractionPrompt = `You are a memory extraction system. Analyze the conversation and extract:

## 1. ENTITY MEMORY (Facts)
Extract factual, structured information that should be remembered.

Format:
` + "```json" + `
{
  "entities": [
    {"entity_type": "project", "entity_id": "proj-001", "key": "language", "value": "Go", "source": "stated", "confidence": 1.0},
    {"entity_type": "teammate", "entity_id": "agent-001", "key": "strength", "value": ["APIs", "security"], "source": "inferred", "confidence": 0.6}
  ]
}
` + "```" + `

Rules:
- Only extract definitive facts, not opinions
- Use "stated" if explicitly mentioned, "inferred" if deduced, "observed" if seen in behavior
- Confidence: 1.0 for stated, 0.8 for observed, 0.6 for inferred
- entity_type must be one of: project, teammate, user, team, codebase, task

## 2. EPISODIC MEMORY (Experiences)
Extract experiences, patterns, and learnings that aren't simple facts.

Format:
` + "```json" + `
{
  "episodes": [
    {
      "episode_type": "observation",
      "content": "Backend Dev 1 consistently catches security issues during code review",
      "importance": 0.8,
      "tags": ["security", "code-review", "backend-dev-1"]
    }
  ]
}
` + "```" + `

Episode Types:
- observation: A pattern you noticed
- lesson: Something learned from experience
- preference: A user or agent preference
- decision: An important decision made
- feedback: Code review or task feedback

Importance scoring:
- 0.9-1.0: Critical, affects future decisions significantly
- 0.7-0.8: Important, useful context
- 0.5-0.6: Moderately useful
- 0.3-0.4: Low priority, nice to have

## CONVERSATION TO ANALYZE:
{conversation}

Extract now. Output only valid JSON with both "entities" and "episodes" arrays (can be empty).`
```

**Step 2: Create extractor service stub**

```go
// internal/memory/extractor/extractor.go
package extractor

import (
	"encoding/json"
)

// ExtractionResult represents the result of memory extraction
type ExtractionResult struct {
	Entities  []EntityExtraction  `json:"entities"`
	Episodes  []EpisodeExtraction `json:"episodes"`
}

// EntityExtraction represents an extracted entity memory
type EntityExtraction struct {
	EntityType string      `json:"entity_type"`
	EntityID   string      `json:"entity_id"`
	Key        string      `json:"key"`
	Value      interface{} `json:"value"`
	Source     string      `json:"source"`
	Confidence float64     `json:"confidence"`
}

// EpisodeExtraction represents an extracted episodic memory
type EpisodeExtraction struct {
	EpisodeType string   `json:"episode_type"`
	Content     string   `json:"content"`
	Importance  float64  `json:"importance"`
	Tags        []string `json:"tags"`
}

// Extractor handles memory extraction from conversations
type Extractor struct {
	// Will be integrated with LLM for actual extraction
}

// NewExtractor creates a new extractor
func NewExtractor() *Extractor {
	return &Extractor{}
}

// GetPrompt returns the extraction prompt with conversation filled in
func (e *Extractor) GetPrompt(conversation string) string {
	// Simple string replacement for now
	// In production, use proper templating
	result := ExtractionPrompt
	// Replace {conversation} placeholder
	result = result[:len(result)-len("{conversation}\n\nExtract now. Output only valid JSON with both \"entities\" and \"episodes\" arrays (can be empty).")]
	result += conversation + "\n\nExtract now. Output only valid JSON with both \"entities\" and \"episodes\" arrays (can be empty)."
	return result
}

// ParseResult parses extraction result from LLM response
func (e *Extractor) ParseResult(response string) (*ExtractionResult, error) {
	var result ExtractionResult
	if err := json.Unmarshal([]byte(response), &result); err != nil {
		return nil, err
	}
	return &result, nil
}
```

**Step 3: Verify compilation**

Run:
```bash
mkdir -p internal/memory/extractor
go build ./internal/memory/extractor/...
```
Expected: No errors

**Step 4: Commit**

```bash
git add internal/memory/extractor/
git commit -m "feat: add memory extraction prompt and parser

Co-Authored-By: Claude Opus 4.6 <noreply@anthropic.com>"
```

---

## Verification

Run these commands to verify Memory Phase 1 is complete:

```bash
# Build
go build -o bin/sweatshop ./cmd/server
# Expected: No errors

# Start and test
rm -f data/sweatshop.db
./bin/sweatshop &
sleep 2

# Create memory
curl -X POST http://localhost:8080/api/memory/entity \
  -H "Content-Type: application/json" \
  -d '{"scope":"global","entityType":"user","key":"name","value":"\"Test User\""}'

# Query memory
curl "http://localhost:8080/api/memory/entity?scope=global&entityType=user&key=name"

# Expected: JSON with the memory

pkill -f sweatshop
```

---

## Summary

| Component | Status |
|-----------|--------|
| Database Schema | ✅ entity_memory table |
| Model | ✅ CRUD operations |
| Service | ✅ Business logic |
| Handler | ✅ REST API |
| App Integration | ✅ Routes registered |
| Extraction Prompt | ✅ Template ready |
