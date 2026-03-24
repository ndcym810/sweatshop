// internal/project/model.go
package project

import (
	"database/sql"
	"time"

	"github.com/sweatshop/sweatshop/internal/shared/db"
)

// Project represents a project entity
type Project struct {
	ID            string    `json:"id"`
	TeamID        string    `json:"teamId"`
	Name          string    `json:"name"`
	Path          string    `json:"path,omitempty"`
	DefaultBranch string    `json:"defaultBranch,omitempty"`
	IsActive      bool      `json:"isActive"`
	CreatedAt     time.Time `json:"createdAt"`
}

// CreateProjectInput represents input for creating a project
type CreateProjectInput struct {
	Name          string `json:"name" validate:"required"`
	Path          string `json:"path"`
	DefaultBranch string `json:"defaultBranch"`
	IsActive      bool   `json:"isActive"`
}

// UpdateProjectInput represents input for updating a project
type UpdateProjectInput struct {
	Name          string `json:"name"`
	Path          string `json:"path"`
	DefaultBranch string `json:"defaultBranch"`
	IsActive      *bool  `json:"isActive"`
}

// GetAllByTeam retrieves all projects for a given team
func GetAllByTeam(teamID string) ([]Project, error) {
	rows, err := db.DB.Query(`
		SELECT id, team_id, name, COALESCE(path, ''), COALESCE(default_branch, ''), is_active, created_at
		FROM projects WHERE team_id = ? ORDER BY created_at DESC
	`, teamID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []Project
	for rows.Next() {
		var p Project
		if err := rows.Scan(&p.ID, &p.TeamID, &p.Name, &p.Path, &p.DefaultBranch, &p.IsActive, &p.CreatedAt); err != nil {
			return nil, err
		}
		projects = append(projects, p)
	}
	return projects, nil
}

// GetByID retrieves a project by ID
func GetByID(id string) (*Project, error) {
	var p Project
	err := db.DB.QueryRow(`
		SELECT id, team_id, name, COALESCE(path, ''), COALESCE(default_branch, ''), is_active, created_at
		FROM projects WHERE id = ?
	`, id).Scan(&p.ID, &p.TeamID, &p.Name, &p.Path, &p.DefaultBranch, &p.IsActive, &p.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &p, nil
}

// Create creates a new project
func Create(p *Project) error {
	_, err := db.DB.Exec(`
		INSERT INTO projects (id, team_id, name, path, default_branch, is_active, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`, p.ID, p.TeamID, p.Name, p.Path, p.DefaultBranch, p.IsActive, p.CreatedAt)
	return err
}

// Update updates a project
func Update(p *Project) error {
	_, err := db.DB.Exec(`
		UPDATE projects SET name = ?, path = ?, default_branch = ?, is_active = ?
		WHERE id = ?
	`, p.Name, p.Path, p.DefaultBranch, p.IsActive, p.ID)
	return err
}

// Delete deletes a project
func Delete(id string) error {
	_, err := db.DB.Exec("DELETE FROM projects WHERE id = ?", id)
	return err
}
