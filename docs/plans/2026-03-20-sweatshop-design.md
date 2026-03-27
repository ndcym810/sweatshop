# Sweatshop - Multi-Agent Orchestration System Design

**Version:** 1.0
**Date:** 2026-03-20
**Author:** Design Session

---

## Table of Contents

1. [Overview & Vision](#1-overview--vision)
2. [Architecture](#2-architecture)
3. [Tech Stack](#3-tech-stack)
4. [Data Model & Storage](#4-data-model--storage)
5. [API Design](#5-api-design)
6. [Frontend UI Layout](#6-frontend-ui-layout)
7. [Memory System](#7-memory-system)
8. [Project Folder Structure](#8-project-folder-structure)
9. [Implementation Phases](#9-implementation-phases)
10. [Multi-Runtime Configuration](#10-multi-runtime-configuration)
11. [Settings System](#11-settings-system)
12. [Future Considerations](#12-future-considerations)

---

## 1. Overview & Vision

### Project Name: Sweatshop

### Vision

A multi-agent orchestration system that manages AI agent teams like a company. Each team is cross-functional (dev, marketing, deployment, research) and works on one or more projects. Teams are coordinated by a Lead agent that allocates tasks and manages teammates.

### Core Concepts

| Concept | Description | CRUD |
|---------|-------------|------|
| **Team** | A cross-functional group working on one or more projects. Includes Lead + multiple teammates from different departments. | Create, Read, Update, Delete |
| **Lead** | Orchestration agent that manages a team, allocates tasks, coordinates communication. | Auto-created with team |
| **Teammate** | A specialized agent (dev, marketing, deployment, etc.) that executes tasks. | Create, Read, Update, Delete |
| **Department** | Agent type category within a team (development, marketing, deployment, research). | Create, Read, Update, Delete |
| **Project** | The work product a team is building (app, service, etc.). A team can handle multiple projects. | Create, Read, Update, Delete |
| **Task** | Unit of work assigned by Lead to teammates. | Create, Read, Update, Delete |
| **Template** | Reusable agent configuration (prompt, skills, tools, model, communication rules). | Create, Read, Update, Delete |
| **Skill** | Custom capability that can be added to agent templates. | Create, Read, Update, Delete |
| **Runtime** | Agent tool configuration (Claude Code, Codex, Gemini CLI, etc.). | Create, Read, Update, Delete |

### User Flow

1. User opens Sweatshop → sees list of teams
2. User can create/delete teams
3. User selects a team → sees Projects panel + Tasks/Departments/Lead tabs
4. User selects a project to focus on → tasks filter to that project
5. User can add/remove projects for the team
6. User can manage departments (add agent types available for this team)
7. User can spawn/remove teammates from available templates
8. User chats with Lead → Lead delegates work across teammates
9. User monitors progress via task board and teammate views

---

## 2. Architecture

### High-Level Architecture

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                         Frontend Layer (React + Vite + TanStack Start)      │
│        (Teams UI | Department UI | Lead Chat | Task Board | File Browser)   │
└─────────────────────────────────────────────────────────────────────────────┘
                                      │ REST API / WebSocket
                                      ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                         Backend Layer (Go + Echo)                            │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐        │
│  │  REST API   │  │  WebSocket  │  │File Watcher │  │  Process    │        │
│  │  Handlers   │  │   Handler   │  │  Service    │  │  Spawner    │        │
│  └─────────────┘  └─────────────┘  └─────────────┘  └─────────────┘        │
└─────────────────────────────────────────────────────────────────────────────┘
                                      │
                                      ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                           Data Layer                                         │
│  ┌─────────────────────────────┐  ┌─────────────────────────────────────┐  │
│  │   SQLite Database           │  │   ~/.sweatshop/ File System         │  │
│  │   - Teams                   │  │   - Agent memory                    │  │
│  │   - Projects                │  │   - IPC messages                    │  │
│  │   - Departments             │  │   - Task artifacts                  │  │
│  │   - Teammates               │  │   - Templates                       │  │
│  │   - Tasks                   │  │   - Session state                   │  │
│  └─────────────────────────────┘  └─────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────────────────────┘
                                      │
                                      ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                    Agent Layer (Claude Code / Codex / Gemini / Others)       │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐        │
│  │    Lead     │  │ Frontend Dev│  │  Marketing  │  │  DevOps     │  ...   │
│  │   Agent     │  │   Agent     │  │   Agent     │  │   Agent     │        │
│  └─────────────┘  └─────────────┘  └─────────────┘  └─────────────┘        │
│        │              │              │              │                       │
│        └──────────────┴──────────────┴──────────────┘                       │
│                         File-based IPC (messages/)                          │
└─────────────────────────────────────────────────────────────────────────────┘
```

### Communication Flow

```
Frontend                    Backend                    Agent
   │                          │                         │
   │──── REST: GET /api/teams ────│                         │
   │◄── JSON response ────────│                         │
   │                          │                         │
   │──── REST: POST spawn ────│──── Process spawn ─────►│
   │◄── 201 Created ──────────│                         │
   │                          │                         │
   │                          │◄── Write to IPC file ───│
   │◄── WS: agent update ──────│                         │
   │                          │                         │
   │──── WS: send message ────│──── Read IPC file ─────►│
   │                          │◄── Write response ─────│
   │◄── WS: agent response ───│                         │
```

---

## 3. Tech Stack

### Frontend

| Layer | Technology | Version | Purpose |
|-------|------------|---------|---------|
| Framework | React | 19.x | UI components |
| Build Tool | Vite | 6.x | Fast bundling, HMR |
| Router | TanStack Start | 1.x | File-based routing, SSR-ready |
| State | Zustand | 5.x | Global state management |
| Server State | TanStack Query | 5.x | Server state, caching |
| Styling | Tailwind CSS | 4.x | Utility-first CSS |
| Components | shadcn/ui | Latest | Pre-built accessible components |
| Icons | Lucide React | Latest | Consistent icon set |
| TypeScript | TypeScript | 5.x | Type safety |

### Backend

| Layer | Technology | Version | Purpose |
|-------|------------|---------|---------|
| Language | Go | 1.26.x | Performance, concurrency, single binary |
| Web Framework | Echo | 4.x | REST API, WebSocket, middleware |
| Database | SQLite (modernc.org/sqlite) | 1.x | Embedded SQL, no CGO dependency |
| File Watching | fsnotify | 1.x | Cross-platform file system events |
| Process Management | os/exec | stdlib | Spawn and manage agent processes |
| Config | YAML (gopkg.in/yaml.v3) | 3.x | Templates, team configs |
| UUID | github.com/google/uuid | 1.x | UUID generation |

### Development Tools

| Tool | Version | Purpose |
|------|---------|---------|
| Air | Latest | Go backend hot reload |
| Node.js | 22.x LTS | Runtime for frontend tooling |
| pnpm | 10.x | Fast package manager |
| ESLint | 9.x | Code linting |
| Prettier | 3.x | Code formatting |
| golangci-lint | Latest | Go linting |

---

## 4. Data Model & Storage

### Database Schema (SQLite)

```sql
-- Settings (key-value store)
CREATE TABLE settings (
    key TEXT PRIMARY KEY,
    value TEXT NOT NULL,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Runtimes
CREATE TABLE runtimes (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    command TEXT NOT NULL,
    args TEXT,                    -- JSON array
    supports_tools BOOLEAN DEFAULT false,
    supports_skills BOOLEAN DEFAULT false,
    supports_model BOOLEAN DEFAULT true,
    default_model TEXT,
    env_vars TEXT,                -- JSON object
    is_built_in BOOLEAN DEFAULT false,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Teams
CREATE TABLE teams (
    id TEXT PRIMARY KEY,          -- UUID
    name TEXT NOT NULL,
    description TEXT,
    lead_runtime_type TEXT DEFAULT 'claude-code',
    lead_runtime_model TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Projects (team can have multiple)
CREATE TABLE projects (
    id TEXT PRIMARY KEY,          -- UUID
    team_id TEXT NOT NULL,        -- UUID
    name TEXT NOT NULL,
    path TEXT NOT NULL,
    default_branch TEXT DEFAULT 'main',
    is_active BOOLEAN DEFAULT true,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (team_id) REFERENCES teams(id) ON DELETE CASCADE
);

-- Departments (per team)
CREATE TABLE departments (
    id TEXT PRIMARY KEY,          -- UUID
    team_id TEXT NOT NULL,        -- UUID
    name TEXT NOT NULL,
    description TEXT,
    sort_order INTEGER DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (team_id) REFERENCES teams(id) ON DELETE CASCADE
);

-- Teammates (agents)
CREATE TABLE teammates (
    id TEXT PRIMARY KEY,          -- UUID
    team_id TEXT NOT NULL,        -- UUID
    department_id TEXT NOT NULL,  -- UUID
    template_id TEXT NOT NULL,    -- UUID
    name TEXT NOT NULL,
    status TEXT DEFAULT 'idle',   -- idle, running, stopped, error
    isolation_config TEXT,        -- JSON: {"type": "process|docker", "mounts": [], ...}
    runtime_type TEXT,            -- Override from template
    runtime_model TEXT,           -- Override from template
    pid INTEGER,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    last_active_at DATETIME,
    FOREIGN KEY (team_id) REFERENCES teams(id) ON DELETE CASCADE,
    FOREIGN KEY (department_id) REFERENCES departments(id) ON DELETE SET NULL
);

-- Tasks
CREATE TABLE tasks (
    id TEXT PRIMARY KEY,          -- UUID
    team_id TEXT NOT NULL,        -- UUID
    project_id TEXT,              -- UUID (nullable)
    assigned_to TEXT,             -- UUID (teammate_id, nullable)
    title TEXT NOT NULL,
    description TEXT,
    status TEXT DEFAULT 'pending',
    priority TEXT DEFAULT 'medium',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    completed_at DATETIME,
    FOREIGN KEY (team_id) REFERENCES teams(id) ON DELETE CASCADE,
    FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE SET NULL,
    FOREIGN KEY (assigned_to) REFERENCES teammates(id) ON DELETE SET NULL
);

-- Chat Sessions
CREATE TABLE chat_sessions (
    id TEXT PRIMARY KEY,
    team_id TEXT NOT NULL,
    started_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    last_message_at DATETIME,
    message_count INTEGER DEFAULT 0,
    is_compacted BOOLEAN DEFAULT false,
    compaction_summary TEXT,
    FOREIGN KEY (team_id) REFERENCES teams(id) ON DELETE CASCADE
);

-- Messages (chat history)
CREATE TABLE messages (
    id TEXT PRIMARY KEY,          -- UUID
    team_id TEXT NOT NULL,        -- UUID
    session_id TEXT,              -- UUID
    sender_type TEXT NOT NULL,    -- 'user', 'lead', 'teammate'
    sender_id TEXT NOT NULL,      -- UUID
    content TEXT NOT NULL,
    message_type TEXT DEFAULT 'chat',
    importance INTEGER DEFAULT 5, -- 1-10
    is_summarized BOOLEAN DEFAULT false,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (team_id) REFERENCES teams(id) ON DELETE CASCADE,
    FOREIGN KEY (session_id) REFERENCES chat_sessions(id) ON DELETE SET NULL
);

-- Templates
CREATE TABLE templates (
    id TEXT PRIMARY KEY,          -- UUID
    name TEXT NOT NULL,
    department TEXT NOT NULL,
    description TEXT,
    icon TEXT,
    prompt TEXT NOT NULL,
    skills TEXT,                  -- JSON array
    tools TEXT,                   -- JSON array
    runtime_type TEXT DEFAULT 'claude-code',
    runtime_model TEXT,
    communication_peers TEXT,     -- JSON array of template_ids
    requires_lead BOOLEAN DEFAULT false,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Skills Library
CREATE TABLE skills (
    id TEXT PRIMARY KEY,          -- UUID
    name TEXT NOT NULL,
    description TEXT,
    content TEXT NOT NULL,
    is_built_in BOOLEAN DEFAULT false,
    created_by TEXT,              -- 'system', 'user', or teammate_id
    department TEXT,
    use_count INTEGER DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Template-Skill Association
CREATE TABLE template_skills (
    template_id TEXT NOT NULL,
    skill_id TEXT NOT NULL,
    added_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (template_id, skill_id),
    FOREIGN KEY (template_id) REFERENCES templates(id) ON DELETE CASCADE,
    FOREIGN KEY (skill_id) REFERENCES skills(id) ON DELETE CASCADE
);

-- Agent Memory (for evolution - future use)
CREATE TABLE agent_memory (
    id TEXT PRIMARY KEY,
    teammate_id TEXT NOT NULL,
    memory_type TEXT NOT NULL,
    category TEXT,
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    source TEXT,
    importance INTEGER DEFAULT 5,
    use_count INTEGER DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    last_used_at DATETIME,
    FOREIGN KEY (teammate_id) REFERENCES teammates(id) ON DELETE CASCADE
);
```

### File System Structure

```
~/.sweatshop/
├── sweatshop.db                 # SQLite database
├── logs/
│   └── sweatshop.log            # System logs
├── global/
│   ├── rules/                   # Global rules/policies
│   ├── templates/               # YAML template files
│   │   ├── frontend-dev.yaml
│   │   ├── backend-dev.yaml
│   │   ├── marketing.yaml
│   │   └── deployment.yaml
│   └── skills/                  # Pre-built skill files
├── teams/
│   └── {team-id}/
│       ├── config.yaml          # Team settings
│       ├── projects.json        # Project paths
│       ├── lead/
│       │   ├── memory/
│       │   │   ├── context.md       # Working memory
│       │   │   └── knowledge.md     # Long-term memory
│       │   ├── session.json         # Current session state
│       │   ├── summaries/           # Compacted summaries
│       │   │   └── 2026-03-20.md
│       │   └── decisions/           # Important decisions log
│       ├── teammates/
│       │   └── {teammate-id}/
│       │       ├── state.json       # Current state
│       │       ├── memory/          # Agent memory
│       │       ├── worktree/        # Git worktree (if applicable)
│       │       └── logs/            # Agent output logs
│       └── shared/
│           ├── messages/            # Inter-agent messages (IPC)
│           └── artifacts/           # Shared outputs
```

### Isolation Config JSON Schema

```json
// Process-based (default)
{
    "type": "process"
}

// Docker-based
{
    "type": "docker",
    "image": "claude-agent:latest",
    "mounts": [
        {"host": "/projects/myapp", "container": "/workspace", "readonly": false}
    ],
    "env_vars": {
        "NODE_ENV": "development"
    },
    "network_access": false,
    "resource_limits": {
        "memory": "2g",
        "cpu": "1.0"
    }
}
```

---

## 5. API Design

### REST API Endpoints

#### Teams

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/teams` | List all teams |
| POST | `/api/teams` | Create a new team |
| GET | `/api/teams/:id` | Get team details |
| PUT | `/api/teams/:id` | Update team |
| DELETE | `/api/teams/:id` | Delete team |

#### Projects

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/teams/:teamId/projects` | List team's projects |
| POST | `/api/teams/:teamId/projects` | Add project to team |
| PUT | `/api/teams/:teamId/projects/:id` | Update project |
| DELETE | `/api/teams/:teamId/projects/:id` | Remove project from team |

#### Departments

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/teams/:teamId/departments` | List team's departments |
| POST | `/api/teams/:teamId/departments` | Create department |
| PUT | `/api/teams/:teamId/departments/:id` | Update department |
| DELETE | `/api/teams/:teamId/departments/:id` | Delete department |

#### Teammates

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/teams/:teamId/teammates` | List teammates (optional: `?departmentId=xxx`) |
| POST | `/api/teams/:teamId/teammates` | Spawn teammate (departmentId in body) |
| GET | `/api/teams/:teamId/teammates/:id` | Get teammate details |
| PUT | `/api/teams/:teamId/teammates/:id` | Update teammate |
| DELETE | `/api/teams/:teamId/teammates/:id` | Stop and remove teammate |
| POST | `/api/teams/:teamId/teammates/:id/start` | Start teammate |
| POST | `/api/teams/:teamId/teammates/:id/stop` | Stop teammate |

#### Tasks

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/teams/:teamId/tasks` | List tasks (filter: `?projectId=&status=`) |
| POST | `/api/teams/:teamId/tasks` | Create task |
| GET | `/api/teams/:teamId/tasks/:id` | Get task details |
| PUT | `/api/teams/:teamId/tasks/:id` | Update task |
| DELETE | `/api/teams/:teamId/tasks/:id` | Delete task |

#### Messages

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/teams/:teamId/messages` | Get chat history |
| POST | `/api/teams/:teamId/messages` | Send message to Lead |

#### Templates

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/templates` | List all templates |
| POST | `/api/templates` | Create template |
| GET | `/api/templates/:id` | Get template |
| PUT | `/api/templates/:id` | Update template |
| DELETE | `/api/templates/:id` | Delete template |

#### Skills

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/skills` | List all skills |
| GET | `/api/skills/prebuilt` | List pre-built skills |
| POST | `/api/skills` | Create skill |
| POST | `/api/skills/import` | Import skills (file upload) |
| GET | `/api/skills/export` | Export all skills |
| GET | `/api/skills/:id` | Get skill |
| PUT | `/api/skills/:id` | Update skill |
| DELETE | `/api/skills/:id` | Delete skill |

#### Runtimes

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/runtimes` | List all runtimes |
| POST | `/api/runtimes` | Create runtime |
| GET | `/api/runtimes/:id` | Get runtime |
| PUT | `/api/runtimes/:id` | Update runtime |
| DELETE | `/api/runtimes/:id` | Delete runtime |
| POST | `/api/runtimes/:id/test` | Test runtime availability |

#### Settings

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/settings` | Get all settings |
| PUT | `/api/settings` | Update settings |

#### Files

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/teams/:teamId/files?path=` | Browse project files |
| GET | `/api/teams/:teamId/files/content?path=` | Get file content |

### WebSocket Events

| Direction | Event | Description |
|-----------|-------|-------------|
| Client → Server | `subscribe:team` | Subscribe to team updates |
| Client → Server | `unsubscribe:team` | Unsubscribe |
| Client → Server | `send:message` | Send chat to Lead |
| Server → Client | `teammate:status` | Status change |
| Server → Client | `teammate:log` | Agent output |
| Server → Client | `task:updated` | Task change |
| Server → Client | `message:new` | New chat message |
| Server → Client | `file:changed` | File change detected |

---

## 6. Frontend UI Layout

### Main Layout Structure

```
┌─────────────────────────────────────────────────────────────────────────────┐
│  Header: [Sweatshop Logo]                         [Settings] [Theme] [User]│
├──────────┬──────────────────────────┬─────────────────────────┬────────────┤
│          │                          │                         │            │
│  Teams   │   Sidebar Panel 2        │       Main Area         │   File     │
│  Sidebar │                          │                         │   Browser  │
│          │   ┌──────────────────┐   │                         │            │
│  [+] New │   │ Project Selector │   │  Content based on       │   /src/    │
│   Team   │   └──────────────────┘   │  active tab:            │   /tests/  │
│          │                          │  - Tasks → Task detail  │   /docs/   │
│  • Team A│   [Tasks | Departments   │  - Departments →        │            │
│  • Team B│    | Lead]               │    Teammate detail      │   main.go  │
│  • Team C│                          │  - Lead → Chat view     │   go.mod   │
│          │   Tab-specific info      │                         │            │
│          │   in sidebar             │                         │            │
│          │                          │                         │            │
│  ← Hide  │                          │                         │            │
├──────────┴──────────────────────────┴─────────────────────────┴────────────┤
│  Footer: [Status Bar - Server connection, active agents count]             │
└─────────────────────────────────────────────────────────────────────────────┘
```

### Tab Contents

| Tab | Sidebar Shows | Main Area Shows |
|-----|---------------|-----------------|
| **Tasks** | Task list + filters | Selected task detail OR task board view |
| **Departments** | Collapsible department list with teammates | Selected teammate detail (logs, files, memory) |
| **Lead** | Lead status, team overview, quick actions | Full chat interface with message input |

### Departments Tab (Sidebar)

```
┌──────────────────────────────┐
│  Departments              [+]│  ← Add department
├──────────────────────────────┤
│  ▼ Development (2)        [+]│  ← Collapse + spawn
│     🎨 Frontend Dev 1     ●  │  ● = running, ○ = idle
│     ⚙️ Backend Dev 1      ○  │
│                              │
│  ▶ Marketing (1)             │  ← Collapsed
│                              │
│  ▼ Deployment (1)         [+]│
│     🚀 DevOps 1           ●  │
│                              │
│  ▶ Research (0)              │
│     (No teammates)           │
└──────────────────────────────┘
```

### Lead Tab Views

#### Sidebar (when Lead tab active)
```
┌──────────────────────────────┐
│  Lead                     ●  │  ● = online status
├──────────────────────────────┤
│  Name: Team Alpha Lead       │
│  Model: Claude Opus 4.6      │
│  Status: Managing 3 agents   │
├──────────────────────────────┤
│  Team Overview:              │
│  • 3 agents active           │
│  • 5 tasks in progress       │
│  • 12 tasks completed today  │
├──────────────────────────────┤
│  Quick Actions:              │
│  [Clear Chat]                │
│  [Export Conversation]       │
│  [View Full Report]          │
└──────────────────────────────┘
```

#### Main Area (when Lead tab active)
```
┌─────────────────────────────────────────────────────────────────┐
│  Chat with Lead                                           [⋯]  │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  ┌─────────────────────────────────────────────────────────┐   │
│  │ Lead                                           10:30 AM│   │
│  │ Good morning! I've reviewed the sprint backlog.        │   │
│  │ Currently 3 tasks in progress:                         │   │
│  │ • Login API (Backend Dev 1) - 60% done                 │   │
│  └─────────────────────────────────────────────────────────┘   │
│                                                                 │
│  ┌─────────────────────────────────────────────────────────┐   │
│  │ You                                             10:32 AM│   │
│  │ Unblock DevOps 1 with the staging credentials          │   │
│  └─────────────────────────────────────────────────────────┘   │
│                                                                 │
├─────────────────────────────────────────────────────────────────┤
│  ┌───────────────────────────────────────────────┐ [Send] [→]  │
│  │ Type a message to Lead...                     │              │
│  └───────────────────────────────────────────────┘              │
└─────────────────────────────────────────────────────────────────┘
```

---

## 7. Memory System

### Memory Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                        Lead Agent                               │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  ┌─────────────────────────────────────────────────────────┐   │
│  │              Session Context (Transient)                 │   │
│  │  - Recent messages (last N turns)                       │   │
│  │  - Current task being discussed                         │   │
│  │  - Active teammate references                           │   │
│  └─────────────────────────────────────────────────────────┘   │
│                           │                                     │
│                           ▼                                     │
│  ┌─────────────────────────────────────────────────────────┐   │
│  │              Working Memory (Active Context)             │   │
│  │  - Current sprint goals                                  │   │
│  │  - Active tasks and assignments                         │   │
│  │  - Recent decisions and outcomes                        │   │
│  │  - Loaded on each interaction                           │   │
│  └─────────────────────────────────────────────────────────┘   │
│                           │                                     │
│                           ▼                                     │
│  ┌─────────────────────────────────────────────────────────┐   │
│  │              Long-term Memory (Persistent)               │   │
│  │  - Project knowledge and patterns                       │   │
│  │  - Team capabilities and preferences                    │   │
│  │  - Historical decisions and rationale                   │   │
│  │  - Lessons learned                                      │   │
│  └─────────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────────┘
```

### Memory Files

#### context.md (Working Memory)
```markdown
# Team Alpha - Current Context

## Sprint Goal
Implement user authentication for E-commerce Platform

## Active Tasks
- [ ] Implement login API (Backend Dev 1) - 60%
- [ ] Design checkout page (Frontend Dev 1) - 40%
- [x] Setup CI/CD pipeline (DevOps 1) - Completed today

## Team Status
- Backend Dev 1: Working on login API
- Frontend Dev 1: Designing checkout UI
- DevOps 1: Available, CI/CD done

## Recent Decisions
- 2026-03-20: Use JWT for authentication
- 2026-03-20: PostgreSQL for user data

## Pending Questions
- None currently
```

#### knowledge.md (Long-term Memory)
```markdown
# Team Alpha - Knowledge Base

## Project Overview
E-commerce Platform - Main product for selling digital goods

## Tech Stack
- Frontend: React 19, Vite 6, TanStack Start
- Backend: Go 1.26, Echo
- Database: SQLite (dev), PostgreSQL (prod)

## Team Capabilities
- Development: Strong in React, Go, SQL
- Deployment: Experienced with Docker, CI/CD

## Coding Standards
- Use TypeScript for all frontend code
- Follow Go standard formatting
- Write tests for all API endpoints
```

### Compaction Strategy

| Trigger | Action |
|---------|--------|
| **Message count > 50** | Summarize messages older than 24 hours |
| **User clicks "Summarize"** | Manual compaction |
| **Session ends** | Create daily summary |
| **Token estimate > 80% limit** | Emergency compaction |

**What gets kept vs compacted:**
- **Kept:** High importance messages, decisions, task assignments
- **Compacted:** Casual conversation, status updates, redundant info

---

## 8. Project Folder Structure

```
sweatshop/
├── cmd/
│   └── server/
│       └── main.go                    # Application entry point
│
├── internal/
│   ├── app/                           # App initialization
│   │   ├── app.go                     # App struct & dependencies
│   │   └── router.go                  # Echo router setup
│   │
│   ├── team/
│   │   ├── handler.go                 # HTTP handlers
│   │   ├── service.go                 # Business logic
│   │   └── model.go                   # Entity + database operations
│   │
│   ├── project/
│   │   ├── handler.go
│   │   ├── service.go
│   │   └── model.go
│   │
│   ├── department/
│   │   ├── handler.go
│   │   ├── service.go
│   │   └── model.go
│   │
│   ├── teammate/
│   │   ├── handler.go
│   │   ├── service.go
│   │   └── model.go
│   │
│   ├── task/
│   │   ├── handler.go
│   │   ├── service.go
│   │   └── model.go
│   │
│   ├── message/
│   │   ├── handler.go
│   │   ├── service.go
│   │   └── model.go
│   │
│   ├── template/
│   │   ├── handler.go
│   │   ├── service.go
│   │   └── model.go
│   │
│   ├── skill/
│   │   ├── handler.go
│   │   ├── service.go
│   │   └── model.go
│   │
│   ├── runtime/
│   │   ├── handler.go
│   │   ├── service.go
│   │   └── model.go
│   │
│   ├── setting/
│   │   ├── handler.go
│   │   ├── service.go
│   │   └── model.go
│   │
│   ├── lead/
│   │   ├── handler.go
│   │   ├── service.go               # Lead orchestration logic
│   │   └── memory.go                 # Memory management
│   │
│   ├── file/
│   │   ├── handler.go
│   │   └── service.go               # File browsing logic
│   │
│   ├── websocket/
│   │   ├── hub.go
│   │   ├── client.go
│   │   └── handler.go
│   │
│   ├── agent/
│   │   ├── spawner.go               # Process spawning
│   │   ├── docker.go                # Docker isolation
│   │   └── worktree.go              # Git worktree
│   │
│   ├── ipc/
│   │   ├── watcher.go
│   │   └── protocol.go
│   │
│   ├── shared/
│   │   ├── db/
│   │   │   ├── db.go                # DB connection
│   │   │   └── migrations.sql       # All migrations
│   │   ├── middleware/
│   │   │   ├── cors.go
│   │   │   ├── logger.go
│   │   │   └── recover.go
│   │   ├── response/
│   │   │   └── response.go          # Common response helpers
│   │   └── config/
│   │       └── config.go
│   │
│   └── router/
│       └── router.go                # Route registration
│
├── pkg/
│   ├── uuid/
│   │   └── uuid.go
│   └── logger/
│       └── logger.go
│
├── web/                              # Frontend
│   ├── src/
│   │   ├── app/
│   │   │   ├── App.tsx
│   │   │   ├── main.tsx
│   │   │   ├── router.tsx
│   │   │   └── providers.tsx
│   │   │
│   │   ├── features/
│   │   │   ├── team/
│   │   │   │   ├── components/
│   │   │   │   ├── hooks/
│   │   │   │   ├── api/
│   │   │   │   ├── types/
│   │   │   │   └── store/
│   │   │   ├── project/
│   │   │   ├── department/
│   │   │   ├── teammate/
│   │   │   ├── task/
│   │   │   ├── lead/
│   │   │   ├── template/
│   │   │   ├── skill/
│   │   │   ├── runtime/
│   │   │   ├── setting/
│   │   │   └── file/
│   │   │
│   │   ├── shared/
│   │   │   ├── components/
│   │   │   │   ├── layout/
│   │   │   │   └── ui/              # shadcn/ui components
│   │   │   ├── hooks/
│   │   │   ├── lib/
│   │   │   ├── types/
│   │   │   └── stores/
│   │   │
│   │   ├── routes/
│   │   │   ├── __root.tsx
│   │   │   ├── index.tsx
│   │   │   ├── teams/
│   │   │   │   ├── index.tsx
│   │   │   │   └── $teamId/
│   │   │   │       ├── index.tsx
│   │   │   │       ├── tasks.tsx
│   │   │   │       ├── departments.tsx
│   │   │   │       └── lead.tsx
│   │   │   └── settings/
│   │   │       ├── index.tsx
│   │   │       ├── general.tsx
│   │   │       ├── runtimes.tsx
│   │   │       ├── skills.tsx
│   │   │       └── templates.tsx
│   │   │
│   │   └── index.css
│   │
│   ├── public/
│   ├── package.json
│   ├── tsconfig.json
│   ├── vite.config.ts
│   ├── tailwind.config.ts
│   └── tanstack-start.config.ts
│
├── configs/
│   └── templates/
│       ├── frontend-dev.yaml
│       ├── backend-dev.yaml
│       ├── marketing.yaml
│       └── deployment.yaml
│
├── scripts/
│   ├── build.sh
│   └── dev.sh
│
├── go.mod
├── go.sum
├── Makefile
├── .air.toml
├── .gitignore
└── README.md
```

---

## 9. Implementation Phases

### Phase Overview

| Phase | Focus | Deliverable |
|-------|-------|-------------|
| **1.1** | UI Shell | Frontend layout, mock data |
| **1.2** | Backend Foundation | Go server, SQLite, REST API |
| **1.3** | Template Management | Template CRUD, YAML loading |
| **2.1** | Agent Spawning | Process spawning, lifecycle |
| **2.2** | Lead Agent | Lead orchestration, chat |
| **2.3** | IPC System | File-based communication |
| **3.1** | Full Team Workflow | Task board, file browser |
| **3.2** | Memory System | Context, compaction |
| **4** | Multi-Team | Multiple teams, reporting |

### Phase 1.1: UI Shell

**Goal:** Working frontend layout with mock data

**Tasks:**
- Initialize Vite + React 19 + TanStack Start project
- Setup Tailwind CSS 4 and shadcn/ui
- Setup Zustand stores
- Create main layout (4-panel structure)
- Build TeamsSidebar component (collapsible)
- Build SecondarySidebar (tabs: Tasks, Departments, Lead)
- Build MainArea component
- Build FileBrowser component
- Add mock data for teams, teammates, tasks
- Implement routing (/, /teams/:id, /settings)

**Deliverable:** Navigable UI with mock data, no backend connections

### Phase 1.2: Backend Foundation

**Goal:** Working Go server with REST API

**Tasks:**
- Initialize Go 1.26 project with Echo
- Setup SQLite database with migrations
- Implement shared/db package
- Implement all features (model, service, handler)
- Setup router and middleware
- Create API client in frontend
- Connect frontend to backend

**Deliverable:** Full CRUD via REST API

### Phase 1.3: Template Management

**Goal:** Template system with YAML support

**Tasks:**
- Create default template YAML files
- Implement template loader
- Build Template browser UI
- Build Template form

**Deliverable:** Templates can be loaded from YAML, created/edited via UI

### Phase 2.1: Agent Spawning

**Goal:** Spawn and manage agent processes

**Tasks:**
- Implement process spawner
- Implement process manager
- Add isolation config support
- Implement git worktree creation
- Add Docker isolation support
- Track process status
- Update UI with spawn/stop buttons

**Deliverable:** Can spawn teammates that run agent processes

### Phase 2.2: Lead Agent

**Goal:** Lead agent that orchestrates team

**Tasks:**
- Implement Lead service
- Create Lead template
- Spawn Lead with team creation
- Implement chat with Lead
- Build Lead chat UI
- Add message persistence

**Deliverable:** Can chat with Lead agent

### Phase 2.3: IPC System

**Goal:** File-based inter-agent communication

**Tasks:**
- Design IPC message format
- Implement file watcher
- Implement message queue
- Add WebSocket events for IPC

**Deliverable:** Agents can communicate via file-based IPC

### Phase 3.1: Full Team Workflow

**Goal:** Complete task board and file browser

**Tasks:**
- Build TaskBoard UI (Kanban view)
- Implement file browser API
- Build FileTree and FileViewer components
- Connect tasks to projects

**Deliverable:** Full task management and file browsing

### Phase 3.2: Memory System

**Goal:** Lead memory management and compaction

**Tasks:**
- Implement memory file structure
- Implement context.md generation
- Implement compaction service
- Add chat session management

**Deliverable:** Lead has persistent memory, can handle long conversations

### Phase 4: Multi-Team Management

**Goal:** Support multiple teams with reporting

**Tasks:**
- Add team switching UI
- Implement cross-team isolation
- Add team-level reporting
- Performance optimization

**Deliverable:** Can manage multiple independent teams

---

## 10. Multi-Runtime Configuration

### Concept

Each agent (Lead or Teammate) can be configured to use a different LLM/agent tool:

| Agent | Tool | Reason |
|-------|------|--------|
| Lead | Codex | Good at orchestration, planning |
| Frontend Dev | Claude Code | Best coding capabilities |
| Designer | Gemini | Good at visual/design tasks |
| Research | Claude Opus | Deep analysis |

### Template Configuration

```yaml
# templates/frontend-dev.yaml
id: frontend-dev
name: Frontend Developer
department: development
description: Specialized in React, CSS, and UI components
icon: 🎨

# Agent tool/runtime selection
runtime:
  type: claude-code           # claude-code | opencode | codex | gemini-cli
  model: claude-sonnet-4-6    # Model override (optional)

# Skills and tools
skills:
  - frontend-design
tools:
  - Read
  - Write
  - Edit
  - Glob
  - Grep
  - Bash
prompt: |
  You are a frontend developer specialized in React...

# Communication config
communication:
  peers: [frontend-dev, backend-dev]
  requires_lead: false
```

### Default Runtimes

| Runtime | Command | Supports Tools | Supports Skills |
|---------|---------|----------------|-----------------|
| Claude Code | `claude --print` | ✓ | ✓ |
| Codex | `codex --quiet` | ✓ | ✗ |
| Gemini CLI | `gemini --output json` | ✗ | ✗ |
| OpenCode | `opencode --non-interactive` | ✓ | ✗ |

---

## 11. Settings System

### Settings Sections

| Section | Description |
|---------|-------------|
| **General** | App preferences, storage path, theme, notifications |
| **Runtimes** | Add/edit/delete agent runtimes |
| **Skills** | Manage pre-built and custom skills, import/export |
| **Templates** | Manage agent templates |
| **About** | Version, license, documentation links |

### General Settings

- Theme (Dark/Light/System)
- Language
- Data Path
- Default Isolation (Process/Docker)
- Default Runtime
- Max Concurrent Agents
- Notifications (task completed, agent error, lead message)

### Runtime Settings

- List all configured runtimes
- Add custom runtimes
- Edit runtime configuration (command, args, capabilities)
- Test runtime availability
- Delete custom runtimes (built-in cannot be deleted)

### Skills Settings

- View pre-built skills (read-only)
- Create custom skills
- Import skills from file/URL
- Export skills
- Delete custom skills

---

## 12. Future Considerations

### Multi-Runtime Support (Post-Phase 4)

```
~/.sweatshop/
├── global/
│   ├── templates/
│   │   ├── claude/
│   │   ├── opencode/
│   │   └── codex/
│   └── adapters/
│       ├── claude.go
│       ├── opencode.go
│       └── codex.go
```

### Mobile App

- React Native app using same REST API
- Focus on monitoring and chat
- Push notifications

### Multi-User Support

- Authentication system
- Role-based permissions
- Team membership
- PostgreSQL migration

### Cloud Deployment

- Deploy backend to cloud
- Secure remote access
- Team collaboration across locations

---

## Success Criteria

### Phase 1 Complete When:
- UI renders with 4-panel layout
- Can create/view/delete teams
- Can create/view/delete departments
- Can view templates
- All REST APIs working
- Frontend connected to backend

### Phase 2 Complete When:
- Can spawn teammate (agent process starts)
- Can stop teammate
- Can see teammate status
- Can chat with Lead
- Lead responds to messages
- Messages persist across restarts

### Phase 3 Complete When:
- Task board shows all tasks
- Can create/assign/complete tasks
- File browser shows project files
- Can view file contents
- Lead has working memory system

---

## Summary

| Aspect | Decision |
|--------|----------|
| **Architecture** | Local server (Go) + Web UI (React) |
| **Backend** | Go 1.26 + Echo + SQLite |
| **Frontend** | React 19 + Vite 6 + TanStack Start + Zustand |
| **Structure** | Feature-based (simplified layered) |
| **API** | REST + WebSocket |
| **Agent IPC** | File-based |
| **Agent Isolation** | Configurable (process/docker) |
| **Memory** | File-based with compaction |
| **Runtime** | Multi-runtime support per agent |
| **Phases** | 4 phases with sub-phases |
