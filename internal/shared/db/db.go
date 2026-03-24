// internal/shared/db/db.go
package db

import (
	"database/sql"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

// Init initializes the database connection and creates tables
func Init(dataPath string) error {
	// Ensure data directory exists
	if err := os.MkdirAll(dataPath, 0755); err != nil {
		return err
	}

	dbPath := filepath.Join(dataPath, "sweatshop.db")

	var err error
	DB, err = sql.Open("sqlite", dbPath)
	if err != nil {
		return err
	}

	// Enable foreign keys
	if _, err := DB.Exec("PRAGMA foreign_keys = ON"); err != nil {
		return err
	}

	// Run migrations
	if err := runMigrations(); err != nil {
		return err
	}

	return nil
}

// Close closes the database connection
func Close() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}

func runMigrations() error {
	schema := `
	-- Settings (key-value store)
	CREATE TABLE IF NOT EXISTS settings (
		key TEXT PRIMARY KEY,
		value TEXT NOT NULL,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	-- Runtimes
	CREATE TABLE IF NOT EXISTS runtimes (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		command TEXT NOT NULL,
		args TEXT,
		supports_tools BOOLEAN DEFAULT false,
		supports_skills BOOLEAN DEFAULT false,
		supports_model BOOLEAN DEFAULT true,
		default_model TEXT,
		env_vars TEXT,
		is_built_in BOOLEAN DEFAULT false,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	-- Teams
	CREATE TABLE IF NOT EXISTS teams (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		description TEXT,
		lead_runtime_type TEXT DEFAULT 'claude-code',
		lead_runtime_model TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	-- Projects
	CREATE TABLE IF NOT EXISTS projects (
		id TEXT PRIMARY KEY,
		team_id TEXT NOT NULL,
		name TEXT NOT NULL,
		path TEXT NOT NULL,
		default_branch TEXT DEFAULT 'main',
		is_active BOOLEAN DEFAULT true,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (team_id) REFERENCES teams(id) ON DELETE CASCADE
	);

	-- Departments
	CREATE TABLE IF NOT EXISTS departments (
		id TEXT PRIMARY KEY,
		team_id TEXT NOT NULL,
		name TEXT NOT NULL,
		description TEXT,
		sort_order INTEGER DEFAULT 0,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (team_id) REFERENCES teams(id) ON DELETE CASCADE
	);

	-- Teammates
	CREATE TABLE IF NOT EXISTS teammates (
		id TEXT PRIMARY KEY,
		team_id TEXT NOT NULL,
		department_id TEXT NOT NULL,
		template_id TEXT NOT NULL,
		name TEXT NOT NULL,
		status TEXT DEFAULT 'idle',
		isolation_config TEXT,
		runtime_type TEXT,
		runtime_model TEXT,
		pid INTEGER,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		last_active_at DATETIME,
		FOREIGN KEY (team_id) REFERENCES teams(id) ON DELETE CASCADE,
		FOREIGN KEY (department_id) REFERENCES departments(id) ON DELETE SET NULL
	);

	-- Tasks
	CREATE TABLE IF NOT EXISTS tasks (
		id TEXT PRIMARY KEY,
		team_id TEXT NOT NULL,
		project_id TEXT,
		assigned_to TEXT,
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
	CREATE TABLE IF NOT EXISTS chat_sessions (
		id TEXT PRIMARY KEY,
		team_id TEXT NOT NULL,
		started_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		last_message_at DATETIME,
		message_count INTEGER DEFAULT 0,
		is_compacted BOOLEAN DEFAULT false,
		compaction_summary TEXT,
		FOREIGN KEY (team_id) REFERENCES teams(id) ON DELETE CASCADE
	);

	-- Messages
	CREATE TABLE IF NOT EXISTS messages (
		id TEXT PRIMARY KEY,
		team_id TEXT NOT NULL,
		session_id TEXT,
		sender_type TEXT NOT NULL,
		sender_id TEXT NOT NULL,
		content TEXT NOT NULL,
		message_type TEXT DEFAULT 'chat',
		importance INTEGER DEFAULT 5,
		is_summarized BOOLEAN DEFAULT false,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (team_id) REFERENCES teams(id) ON DELETE CASCADE,
		FOREIGN KEY (session_id) REFERENCES chat_sessions(id) ON DELETE SET NULL
	);

	-- Templates
	CREATE TABLE IF NOT EXISTS templates (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		department TEXT NOT NULL,
		description TEXT,
		icon TEXT,
		prompt TEXT NOT NULL,
		skills TEXT,
		tools TEXT,
		runtime_type TEXT DEFAULT 'claude-code',
		runtime_model TEXT,
		communication_peers TEXT,
		requires_lead BOOLEAN DEFAULT false,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	-- Skills
	CREATE TABLE IF NOT EXISTS skills (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		description TEXT,
		content TEXT NOT NULL,
		is_built_in BOOLEAN DEFAULT false,
		created_by TEXT,
		department TEXT,
		use_count INTEGER DEFAULT 0,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	-- Template-Skill Association
	CREATE TABLE IF NOT EXISTS template_skills (
		template_id TEXT NOT NULL,
		skill_id TEXT NOT NULL,
		added_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		PRIMARY KEY (template_id, skill_id),
		FOREIGN KEY (template_id) REFERENCES templates(id) ON DELETE CASCADE,
		FOREIGN KEY (skill_id) REFERENCES skills(id) ON DELETE CASCADE
	);

	-- Indexes
	CREATE INDEX IF NOT EXISTS idx_projects_team ON projects(team_id);
	CREATE INDEX IF NOT EXISTS idx_departments_team ON departments(team_id);
	CREATE INDEX IF NOT EXISTS idx_teammates_team ON teammates(team_id);
	CREATE INDEX IF NOT EXISTS idx_tasks_team ON tasks(team_id);
	CREATE INDEX IF NOT EXISTS idx_tasks_status ON tasks(status);
	CREATE INDEX IF NOT EXISTS idx_messages_team ON messages(team_id);
	CREATE INDEX IF NOT EXISTS idx_messages_session ON messages(session_id);
	`

	_, err := DB.Exec(schema)
	return err
}
