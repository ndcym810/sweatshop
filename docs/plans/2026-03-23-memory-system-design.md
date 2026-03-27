# Sweatshop - 3-Layer Memory System Design

**Version:** 1.0
**Date:** 2026-03-23
**Author:** Design Session

---

## Table of Contents

1. [Overview](#1-overview)
2. [Architecture](#2-architecture)
3. [Layer 1: Working Memory](#3-layer-1-working-memory)
4. [Layer 2: Entity Memory](#4-layer-2-entity-memory)
5. [Layer 3: Episodic Memory](#5-layer-3-episodic-memory)
6. [Memory Extraction](#6-memory-extraction)
7. [Memory Sync Strategy](#7-memory-sync-strategy)
8. [Memory Injection](#8-memory-injection)
9. [Implementation Phases](#9-implementation-phases)
10. [Database Schema](#10-database-schema)

---

## 1. Overview

### Problem Statement

Claude Code has a session-scoped memory system. While CLAUDE.md and memory files persist, the conversation context doesn't carry learnings from previous sessions well. For Sweatshop's multi-agent orchestration, we need:

1. **Persistent memory** that survives session boundaries
2. **Structured facts** about projects, teammates, and users
3. **Experiential memory** for patterns and learnings
4. **Cross-agent sharing** so Lead and teammates share knowledge

### Design Principles

| Principle | Description |
|-----------|-------------|
| **Lazy loading** | Only load what's needed for current context |
| **Size limits** | Strict token limits per memory tier |
| **Single responsibility** | Each memory entry = one concept |
| **Event-driven sync** | Critical memories sync immediately, others batch |
| **Search over scan** | Query specific memories, don't load all |

---

## 2. Architecture

### High-Level Architecture

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                              SWEATSHOP MEMORY SYSTEM                         │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  ┌─────────────────────────────────────────────────────────────────────┐    │
│  │ LAYER 1: WORKING MEMORY (Transient)                                 │    │
│  │ • Handled by Claude Code session context                            │    │
│  │ • Recent N conversation turns                                       │    │
│  │ • Auto-compacted by Claude Code                                     │    │
│  └─────────────────────────────────────────────────────────────────────┘    │
│                                      │                                       │
│                    Extract facts    │    Extract experiences                │
│                         ▼           │           ▼                           │
│  ┌─────────────────────────────┐    │    ┌─────────────────────────────┐   │
│  │ LAYER 2: ENTITY MEMORY      │    │    │ LAYER 3: EPISODIC MEMORY    │   │
│  │ Structured key-value store  │    │    │ Keyword/Vector search       │   │
│  │ (SQLite)                    │    │    │ (SQLite + embeddings later) │   │
│  │                             │    │    │                             │   │
│  │ project.language: "Go"      │    │    │ "Backend Dev 1 prefers      │   │
│  │ teammate.status: "busy"     │    │    │  explicit error messages"   │   │
│  │ user.style: "concise"       │    │    │                             │   │
│  └─────────────────────────────┘    │    └─────────────────────────────┘   │
│                    ▲                 │                  ▲                    │
│                    │    Query        │       Query      │                    │
│              Direct lookup           │     Keyword search                   │
│                    │                 │                  │                    │
│  ┌─────────────────────────────────────────────────────────────────────┐    │
│  │                        AGENT (Lead / Teammate)                      │    │
│  │  Injects relevant memories into context before responding           │    │
│  └─────────────────────────────────────────────────────────────────────┘    │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### Memory Types Comparison

| Layer | Purpose | Storage | Query Method | Persistence |
|-------|---------|---------|--------------|-------------|
| **Working** | Current conversation | Claude Code session | Direct access | Session only |
| **Entity** | Factual knowledge | SQLite (key-value) | Direct lookup | Permanent |
| **Episodic** | Experiences & patterns | SQLite (+ embeddings) | Keyword/semantic search | Permanent |

---

## 3. Layer 1: Working Memory

### Description

Working memory is the short-term context of the current conversation session. This is **already handled by Claude Code** - no custom implementation needed.

### Characteristics

| Aspect | Details |
|--------|---------|
| Storage | Claude Code session context |
| Capacity | Model's context window (~200K tokens) |
| Compaction | Auto-compacted by Claude Code when ~80% full |
| Persistence | Session-scoped, lost on session end |

### Sweatshop's Role

When important information emerges in working memory, **extract it** to Entity or Episodic memory:

```
Conversation → Memory Extractor → Entity/Episodic Memory
```

---

## 4. Layer 2: Entity Memory

### Description

Entity memory stores **structured, factual information** as key-value pairs. This is for things that can be clearly defined and directly looked up.

### What Goes Here

| Entity Type | Example Keys | Example Values |
|-------------|--------------|----------------|
| **Project** | language, framework, database, test_command | `"Go"`, `"Echo"`, `"SQLite"` |
| **Teammate** | status, current_task, strength, workload | `"busy"`, `"auth-feature"`, `["APIs", "security"]` |
| **User** | communication_style, preferred_language | `"concise"`, `"Go"` |
| **Team** | sprint_goal, active_project | `"Implement auth"`, `"proj-001"` |
| **Codebase** | entry_point, config_files, important_patterns | `"main.go"`, `["config.yaml"]` |

### Schema

```sql
CREATE TABLE entity_memory (
    id TEXT PRIMARY KEY,

    -- Entity identification
    scope TEXT NOT NULL,          -- 'global', 'team', 'project', 'agent'
    scope_id TEXT,                -- team_id, project_id, or agent_id
    entity_type TEXT NOT NULL,    -- 'project', 'teammate', 'user', 'codebase', 'task'
    entity_id TEXT,               -- UUID reference to actual entity

    -- Key-value pair
    key TEXT NOT NULL,            -- Attribute name
    value TEXT NOT NULL,          -- JSON value (string, number, array, object)

    -- Metadata
    confidence REAL DEFAULT 1.0,  -- How certain (0.0 - 1.0)
    source TEXT DEFAULT 'stated', -- 'stated', 'inferred', 'observed'
    ttl_days INTEGER,             -- Optional: auto-expire after N days

    -- Timestamps
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,

    -- Unique constraint
    UNIQUE(scope, scope_id, entity_type, entity_id, key)
);

-- Indexes
CREATE INDEX idx_entity_scope ON entity_memory(scope, scope_id);
CREATE INDEX idx_entity_type ON entity_memory(entity_type, entity_id);
CREATE INDEX idx_entity_key ON entity_memory(scope, scope_id, entity_type, key);
```

### Predefined Entity Keys

```yaml
project:
  - name: string
  - language: string
  - framework: string
  - database: string
  - test_framework: string
  - build_command: string
  - test_command: string
  - deploy_command: string
  - default_branch: string
  - code_style: string

teammate:
  - name: string
  - role: string
  - status: string              # idle, busy, error
  - current_task: string        # task_id
  - strength: [string]          # array of strengths
  - workload: number            # 0-100

user:
  - name: string
  - communication_style: string # concise, detailed, technical
  - preferred_language: string
  - timezone: string

team:
  - name: string
  - sprint_goal: string
  - active_project: string      # project_id

codebase:
  - entry_point: string
  - config_files: [string]
  - important_patterns: [string]
  - avoid_patterns: [string]
```

### Source Types

| Source | Description | Confidence |
|--------|-------------|------------|
| `stated` | Explicitly mentioned by user | 1.0 |
| `observed` | Seen in agent behavior | 0.8 |
| `inferred` | Deduced from context | 0.6 |

---

## 5. Layer 3: Episodic Memory

### Description

Episodic memory stores **experiences, patterns, and learnings** that aren't simple facts. These are retrieved via keyword search (Phase 2) or semantic search with embeddings (Phase 3).

### Episode Types

| Type | Description | Example |
|------|-------------|---------|
| `observation` | Noticed pattern or behavior | "Backend Dev 1 catches security issues others miss" |
| `lesson` | Something learned from experience | "SQLite locks with concurrent writes in tests" |
| `preference` | User/agent preference detected | "User prefers quick confirmations for routine tasks" |
| `decision` | Important decision made | "Decided JWT over session-based auth" |
| `feedback` | Code review or task feedback | "User said API responses should include request ID" |

### Context-Based Sharing

Memories are scoped by **context** (team, project, department) rather than by individual agent. This means:

- All agents working in the same context share relevant memories
- Backend teammates share "backend" department memories
- All teammates in a team share team-level decisions

```
┌─────────────────────────────────────────────────────────────┐
│                    TEAM CONTEXT                             │
│              (Shared by all agents in team)                 │
│  - Team decisions                                           │
│  - User preferences                                         │
│  - Sprint goals                                             │
└─────────────────────────────────────────────────────────────┘
                           │
              ┌────────────┼────────────┐
              ▼            ▼            ▼
┌──────────────────┐ ┌──────────────────┐ ┌──────────────────┐
│ Department:      │ │ Department:      │ │ Department:      │
│ backend          │ │ frontend         │ │ deployment       │
│                  │ │                  │ │                  │
│ - API patterns   │ │ - UI patterns    │ │ - CI/CD patterns │
│ - DB patterns    │ │ - CSS patterns   │ │ - Deploy configs │
│ - Auth patterns  │ │ - React patterns │ │ - Infra patterns │
└──────────────────┘ └──────────────────┘ └──────────────────┘
        │                    │                    │
        ▼                    ▼                    ▼
   All backend          All frontend         All devops
   agents share         agents share         agents share
```

### Schema

```sql
CREATE TABLE episodic_memory (
    id TEXT PRIMARY KEY,

    -- Scope (Context-based, not agent-based)
    scope TEXT NOT NULL,          -- 'global', 'team', 'project', 'department'
    scope_id TEXT,                -- team_id, project_id, or department_id
    created_by TEXT,              -- Which agent recorded this (for attribution only)

    -- Content
    episode_type TEXT NOT NULL,   -- 'observation', 'lesson', 'preference', 'decision', 'feedback'
    content TEXT NOT NULL,        -- The actual memory text

    -- Embedding (Phase 3)
    embedding BLOB,               -- Serialized float32 array (deferred to Phase 3)

    -- Metadata
    importance REAL DEFAULT 0.5,  -- 0.0 - 1.0, affects retrieval ranking
    tags TEXT,                    -- JSON array: ["security", "backend", "api"]

    -- Access tracking
    access_count INTEGER DEFAULT 0,
    last_accessed_at DATETIME,

    -- Timestamps
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    expires_at DATETIME,          -- Optional: auto-expire

    -- Soft delete
    is_active BOOLEAN DEFAULT TRUE
);

-- Indexes
CREATE INDEX idx_episodic_scope ON episodic_memory(scope, scope_id);
CREATE INDEX idx_episodic_type ON episodic_memory(episode_type);
CREATE INDEX idx_episodic_importance ON episodic_memory(importance DESC);
CREATE INDEX idx_episodic_active ON episodic_memory(is_active, created_at DESC);
```

### Importance Scoring

| Score Range | Description | Example |
|-------------|-------------|---------|
| 0.9 - 1.0 | Critical | Architectural decisions, user constraints |
| 0.7 - 0.8 | Important | Patterns that affect quality, team dynamics |
| 0.5 - 0.6 | Moderately useful | General observations, minor preferences |
| 0.3 - 0.4 | Low priority | Nice-to-know, may expire |
| 0.0 - 0.2 | Trivial | Almost noise, likely to be cleaned up |

### Retrieval (Without Embeddings - Phase 2)

```go
// Keyword-based retrieval for Phase 2
func (em *EpisodicMemory) Retrieve(query string, opts RetrieveOptions) ([]Episode, error) {
    keywords := extractKeywords(query)

    query := `
        SELECT id, scope, scope_id, agent_id, episode_type, content,
               importance, tags, access_count, created_at
        FROM episodic_memory
        WHERE is_active = TRUE
          AND (? IS NULL OR scope = ?)
          AND (? IS NULL OR scope_id = ?)
          AND (
              content LIKE '%' || ? || '%'
              OR tags LIKE '%' || ? || '%'
          )
        ORDER BY importance DESC
        LIMIT ?
    `

    // Match against keywords
    // ...
}
```

### Retrieval (With Embeddings - Phase 3)

```go
// Semantic search with embeddings (Phase 3)
func (em *EpisodicMemory) Retrieve(query string, opts RetrieveOptions) ([]Episode, error) {
    queryEmbedding := em.embedder.Embed(query)

    // Get candidates
    candidates := em.getCandidates(opts)

    // Rank by cosine similarity * importance
    for _, ep := range candidates {
        similarity := cosineSimilarity(queryEmbedding, ep.Embedding)
        score := similarity * ep.Importance
        // ...
    }

    // Return top N
}
```

---

## 6. Memory Extraction

### How Memories Are Created

```
┌─────────────────────────────────────────────────────────────────────┐
│                        CONVERSATION                                 │
│  User: "Backend Dev 1 is really good at catching security issues"  │
└─────────────────────────────────────────────────────────────────────┘
                                │
                                ▼
┌─────────────────────────────────────────────────────────────────────┐
│                    MEMORY EXTRACTOR                                 │
│  Analyzes conversation to identify:                                │
│  - Facts → Entity Memory                                           │
│  - Experiences → Episodic Memory                                   │
└─────────────────────────────────────────────────────────────────────┘
                                │
                    ┌───────────┴───────────┐
                    ▼                       ▼
        ┌───────────────────┐    ┌───────────────────┐
        │   ENTITY MEMORY   │    │  EPISODIC MEMORY  │
        │                   │    │                   │
        │ teammate.strength │    │ Type: observation │
        │ = ["security"]    │    │ "Backend Dev 1    │
        │ Source: stated    │    │  catches security │
        │ Confidence: 1.0   │    │  issues well"     │
        └───────────────────┘    └───────────────────┘
```

### Extraction Prompt Template

```markdown
You are a memory extraction system. Analyze the conversation and extract:

## 1. ENTITY MEMORY (Facts)
Extract factual, structured information that should be remembered.

Format:
```json
{
  "entities": [
    {"entity_type": "project", "entity_id": "proj-001", "key": "language", "value": "Go", "source": "stated", "confidence": 1.0},
    {"entity_type": "teammate", "entity_id": "agent-001", "key": "strength", "value": ["APIs", "security"], "source": "inferred", "confidence": 0.6}
  ]
}
```

Rules:
- Only extract definitive facts, not opinions
- Use "stated" if explicitly mentioned, "inferred" if deduced, "observed" if seen in behavior
- Confidence: 1.0 for stated, 0.8 for observed, 0.6 for inferred

## 2. EPISODIC MEMORY (Experiences)
Extract experiences, patterns, and learnings that aren't simple facts.

Format:
```json
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
```

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

Extract now. Output only valid JSON.
```

---

## 7. Memory Sync Strategy

### Event-Driven Sync

Different memory types have different sync urgency between agents:

| Memory Type | Sync Strategy | Why |
|-------------|---------------|-----|
| **Decisions** | Immediate | Affects everyone's work |
| **User preferences** | Immediate | Affects all interactions |
| **Lessons learned** | Immediate | Prevents repeat mistakes |
| **Teammate status** | Immediate | Orchestration needs it |
| **Project facts** | Immediate | Everyone needs same source of truth |
| **Observations** | Batch (5 min) | Nice to have, not critical |

### Sync Flow

```
┌─────────────────────────────────────────────────────────────┐
│                    LEAD AGENT                               │
├─────────────────────────────────────────────────────────────┤
│  Extracts memory from conversation                          │
│         │                                                   │
│         ▼                                                   │
│  ┌─────────────────┐                                        │
│  │ Classify by     │                                        │
│  │ sync urgency    │                                        │
│  └────────┬────────┘                                        │
│           │                                                 │
│     ┌─────┴─────┐                                           │
│     │           │                                           │
│     ▼           ▼                                           │
│  HIGH          LOW                                          │
│  Priority      Priority                                     │
│     │           │                                           │
│     │           └────▶ Queue for batch sync                 │
│     │                   (every 5 min)                       │
│     ▼                                                       │
│  Write immediately                                          │
│  to team memory                                             │
│     │                                                       │
│     ▼                                                       │
└─────┼───────────────────────────────────────────────────────┘
      │
      ▼
┌─────────────────────────────────────────────────────────────┐
│                    TEAM MEMORY                              │
│  (All teammates can read)                                   │
└─────────────────────────────────────────────────────────────┘
```

### High Priority Memories (Immediate Sync)

- `episode_type = decision`
- `episode_type = preference`
- `episode_type = lesson`
- `entity_type = teammate` + `key = status`
- `entity_type = project` (all keys)

### Low Priority Memories (Batch Sync)

- `episode_type = observation`
- `episode_type = feedback`

---

## 8. Memory Injection

### Before Agent Responds

Before an agent generates a response, relevant memories are injected into context:

```
┌─────────────────────────────────────────────────────────────────────┐
│                    MEMORY INJECTION                                 │
├─────────────────────────────────────────────────────────────────────┤
│                                                                     │
│  1. Get entity facts relevant to current context                    │
│     - Project facts (language, framework, etc.)                     │
│     - Teammate statuses                                             │
│     - User preferences                                              │
│                                                                     │
│  2. Query episodic memory for relevant experiences                  │
│     - Search by keywords/tags matching user query                   │
│     - Limit to top 5 results                                        │
│                                                                     │
│  3. Build context string                                            │
│     ## Known Facts                                                  │
│     - Project language: Go                                          │
│     - Project framework: Echo                                       │
│     - User prefers: concise responses                               │
│                                                                     │
│     ## Relevant Context                                             │
│     - [decision] Using JWT for authentication                       │
│     - [lesson] SQLite locks with concurrent writes in tests         │
│                                                                     │
│  4. Inject into agent's context window                              │
│                                                                     │
└─────────────────────────────────────────────────────────────────────┘
```

### Go Implementation

```go
// internal/memory/injector.go

type MemoryInjector struct {
    entityMem   *EntityMemory
    episodicMem *EpisodicMemory
}

func (mi *MemoryInjector) BuildMemoryContext(
    agentID string,
    teamID string,
    userQuery string,
) (string, error) {
    var context strings.Builder

    // 1. Entity facts
    context.WriteString("## Known Facts\n\n")

    projectFacts, _ := mi.entityMem.GetAllForEntity("team", teamID, "project", "*")
    for k, v := range projectFacts {
        context.WriteString(fmt.Sprintf("- %s: %v\n", k, v))
    }

    userPrefs, _ := mi.entityMem.GetAllForEntity("global", "", "user", "default")
    for k, v := range userPrefs {
        context.WriteString(fmt.Sprintf("- User %s: %v\n", k, v))
    }

    // 2. Episodic memories
    context.WriteString("\n## Relevant Context\n\n")

    episodes, _ := mi.episodicMem.Retrieve(userQuery, RetrieveOptions{
        Scope:  "team",
        ScopeID: teamID,
        Limit:  5,
    })

    for _, ep := range episodes {
        context.WriteString(fmt.Sprintf("- [%s] %s\n", ep.EpisodeType, ep.Content))
    }

    return context.String(), nil
}
```

---

## 9. Implementation Phases

### Phase Overview

| Phase | Focus | Features |
|-------|-------|----------|
| **1** | Entity Memory | Key-value storage, direct lookup, CRUD API |
| **2** | Episodic Memory (Basic) | Storage, keyword-based retrieval |
| **3** | Semantic Search | Add embeddings for episodic memory |
| **4** | Cross-Context Sharing | Memory sharing across team/project/department contexts |

### Phase 1: Entity Memory

**Goal:** Working entity memory system with CRUD operations

**Tasks:**
- Create database schema
- Implement `EntityMemory` Go package
- Implement CRUD handlers (REST API)
- Add extraction logic for entity memory
- Integrate with Lead agent

**Deliverable:** Can store and retrieve project facts, teammate info, user preferences

### Phase 2: Episodic Memory (Basic)

**Goal:** Working episodic memory with keyword retrieval

**Tasks:**
- Create database schema
- Implement `EpisodicMemory` Go package
- Implement keyword-based retrieval
- Add extraction logic for episodic memory
- Implement memory injection

**Deliverable:** Can store experiences and retrieve via keyword search

### Phase 3: Semantic Search

**Goal:** Embedding-based semantic search for episodic memory

**Tasks:**
- Add embedding column to schema
- Integrate embedding model (OpenAI or local)
- Implement semantic search
- Add embedding generation on episode creation

**Deliverable:** Can retrieve relevant memories via semantic similarity

### Phase 4: Cross-Context Sharing

**Goal:** Memory sharing across contexts (team, project, department)

**Tasks:**
- Implement event-driven sync logic
- Add batch sync for low-priority memories
- Integrate sync triggers in Lead agent
- Test cross-context memory access (department-level sharing)

**Deliverable:** Agents share relevant memories within their context scope

---

## 10. Database Schema

### Complete Schema

```sql
-- ============================================
-- SWEATSHOP MEMORY SYSTEM - DATABASE SCHEMA
-- ============================================

-- Layer 2: Entity Memory
CREATE TABLE entity_memory (
    id TEXT PRIMARY KEY,

    -- Entity identification
    scope TEXT NOT NULL,          -- 'global', 'team', 'project', 'agent'
    scope_id TEXT,                -- team_id, project_id, or agent_id
    entity_type TEXT NOT NULL,    -- 'project', 'teammate', 'user', 'codebase', 'task'
    entity_id TEXT,               -- UUID reference to actual entity

    -- Key-value pair
    key TEXT NOT NULL,            -- Attribute name
    value TEXT NOT NULL,          -- JSON value (string, number, array, object)

    -- Metadata
    confidence REAL DEFAULT 1.0,  -- How certain (0.0 - 1.0)
    source TEXT DEFAULT 'stated', -- 'stated', 'inferred', 'observed'
    ttl_days INTEGER,             -- Optional: auto-expire after N days

    -- Timestamps
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,

    -- Unique constraint
    UNIQUE(scope, scope_id, entity_type, entity_id, key)
);

-- Indexes for entity_memory
CREATE INDEX idx_entity_scope ON entity_memory(scope, scope_id);
CREATE INDEX idx_entity_type ON entity_memory(entity_type, entity_id);
CREATE INDEX idx_entity_key ON entity_memory(scope, scope_id, entity_type, key);

-- Layer 3: Episodic Memory (Context-based sharing)
CREATE TABLE episodic_memory (
    id TEXT PRIMARY KEY,

    -- Scope (Context-based, not agent-based)
    scope TEXT NOT NULL,          -- 'global', 'team', 'project', 'department'
    scope_id TEXT,                -- team_id, project_id, or department_id
    created_by TEXT,              -- Which agent recorded this (for attribution only)

    -- Content
    episode_type TEXT NOT NULL,   -- 'observation', 'lesson', 'preference', 'decision', 'feedback'
    content TEXT NOT NULL,        -- The actual memory text

    -- Embedding (Phase 3 - nullable until implemented)
    embedding BLOB,               -- Serialized float32 array

    -- Metadata
    importance REAL DEFAULT 0.5,  -- 0.0 - 1.0
    tags TEXT,                    -- JSON array

    -- Access tracking
    access_count INTEGER DEFAULT 0,
    last_accessed_at DATETIME,

    -- Timestamps
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    expires_at DATETIME,

    -- Soft delete
    is_active BOOLEAN DEFAULT TRUE
);

-- Indexes for episodic_memory
CREATE INDEX idx_episodic_scope ON episodic_memory(scope, scope_id);
CREATE INDEX idx_episodic_type ON episodic_memory(episode_type);
CREATE INDEX idx_episodic_importance ON episodic_memory(importance DESC);
CREATE INDEX idx_episodic_active ON episodic_memory(is_active, created_at DESC);

-- ============================================
-- SYNC QUEUE (for batch sync of low-priority memories)
-- ============================================

CREATE TABLE memory_sync_queue (
    id TEXT PRIMARY KEY,
    memory_type TEXT NOT NULL,    -- 'entity' or 'episodic'
    memory_id TEXT NOT NULL,      -- Reference to the memory entry
    sync_priority TEXT NOT NULL,  -- 'high', 'low'
    status TEXT DEFAULT 'pending', -- 'pending', 'synced', 'failed'
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    synced_at DATETIME
);

CREATE INDEX idx_sync_queue_status ON memory_sync_queue(status, sync_priority);
```

---

## Summary

| Aspect | Decision |
|--------|----------|
| **Architecture** | 3-layer (Working, Entity, Episodic) |
| **Working Memory** | Claude Code session (no custom implementation) |
| **Entity Memory** | SQLite key-value, direct lookup |
| **Episodic Memory** | SQLite + keyword search (Phase 2), embeddings (Phase 3) |
| **Sharing Model** | Context-based (team/project/department), not agent-based |
| **Extraction** | LLM-based prompt to classify and extract |
| **Sync Strategy** | Event-driven (immediate for critical, batch for others) |
| **Phases** | 4 phases: Entity → Episodic Basic → Semantic → Cross-Context Sharing |
