# Sweatshop Implementation Plans Index

This directory contains detailed, bite-sized implementation plans for the Sweatshop multi-agent orchestration system.

## Plan Files

### Core System (from sweatshop-design.md)

| Plan File | Phase | Description | Status |
|-----------|-------|-------------|--------|
| `2026-03-23-phase-1.1-ui-shell.md` | 1.1 | Frontend layout, React + Vite + TanStack Router | 📋 Ready |
| `2026-03-23-phase-1.2-backend-foundation.md` | 1.2 | Go server, SQLite, REST API for all entities | 📋 Ready |
| `2026-03-24-phase-1.3-template-management.md` | 1.3 | YAML templates, CRUD operations | 📋 Ready |

### Memory System (from memory-system-design.md)

| Plan File | Phase | Description | Status |
|-----------|-------|-------------|--------|
| `2026-03-24-memory-phase-1-entity-memory.md` | Mem 1 | Key-value entity storage, CRUD API | 📋 Ready |

### Design Documents

| Document | Description |
|----------|-------------|
| `2026-03-20-sweatshop-design.md` | Main system architecture and design |
| `2026-03-23-memory-system-design.md` | 3-layer memory system design |

---

## How to Execute Plans

Each plan follows this format:
1. **Task N** - Logical grouping of related steps
2. **Files** - Exact paths to create/modify
3. **Step X** - Individual action (2-5 minutes each)
4. **Commit** - Atomic commit after each task

### Execution Options

**Option A: Subagent-Driven (current session)**
```
Use superpowers:subagent-driven-development skill
- Fresh subagent per task
- Code review between tasks
- Fast iteration
```

**Option B: Parallel Session (separate)**
```
Open new session with superpowers:executing-plans skill
- Batch execution
- Checkpoints
- Resume capability
```

---

## Remaining Phases to Plan

These phases need implementation plans created:

### Core System
- [ ] Phase 2.1: Agent Spawning
- [ ] Phase 2.2: Lead Agent
- [ ] Phase 2.3: IPC System
- [ ] Phase 3.1: Full Team Workflow
- [ ] Phase 3.2: Memory System Integration
- [ ] Phase 4: Multi-Team Management

### Memory System
- [ ] Memory Phase 2: Episodic Memory (Basic)
- [ ] Memory Phase 3: Semantic Search
- [ ] Memory Phase 4: Cross-Context Sharing

---

## Quick Start

1. Start with **Phase 1.1** (UI Shell) to get visual progress
2. Then **Phase 1.2** (Backend) for API foundation
3. Then **Phase 1.3** (Templates) for agent configuration
4. Then **Memory Phase 1** for persistent facts

Or work on backend first:
1. Phase 1.2 (Backend)
2. Memory Phase 1 (Entity Memory)
3. Phase 1.1 (UI Shell)
4. Phase 1.3 (Templates)

---

## Architecture Overview

```
┌─────────────────────────────────────────────────────────┐
│                  Frontend (React + Vite)                 │
│  Teams UI | Departments | Lead Chat | Task Board        │
└─────────────────────────────────────────────────────────┘
                          │ REST API / WebSocket
                          ▼
┌─────────────────────────────────────────────────────────┐
│                  Backend (Go + Echo)                     │
│  Handlers | Services | Models                           │
└─────────────────────────────────────────────────────────┘
                          │
                          ▼
┌─────────────────────────────────────────────────────────┐
│                  Data Layer                              │
│  SQLite (entities, tasks, memory)                       │
└─────────────────────────────────────────────────────────┘
                          │
                          ▼
┌─────────────────────────────────────────────────────────┐
│                  Agent Layer                             │
│  Lead Agent | Teammates (Claude Code / Codex / etc.)    │
└─────────────────────────────────────────────────────────┘
```

---

## File Structure (Target)

```
sweatshop/
├── cmd/server/main.go           # Entry point
├── internal/
│   ├── app/                     # App initialization
│   ├── team/                    # Team feature
│   ├── project/                 # Project feature
│   ├── department/              # Department feature
│   ├── teammate/                # Teammate feature
│   ├── task/                    # Task feature
│   ├── template/                # Template feature
│   ├── memory/
│   │   ├── entity/              # Entity memory (Phase 1)
│   │   ├── episodic/            # Episodic memory (Phase 2)
│   │   └── extractor/           # Memory extraction
│   └── shared/
│       ├── db/                  # Database
│       ├── response/            # API helpers
│       └── middleware/          # HTTP middleware
├── pkg/
│   ├── uuid/                    # UUID generation
│   └── logger/                  # Logging
├── configs/templates/           # YAML templates
├── web/                         # Frontend
│   ├── src/
│   │   ├── app/                 # App setup
│   │   ├── features/            # Feature modules
│   │   ├── shared/              # Shared components
│   │   └── routes/              # TanStack routes
│   └── package.json
└── docs/plans/                  # This directory
```
