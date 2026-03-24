// internal/team/model.go
package team

import (
	"database/sql"
	"time"

	"github.com/sweatshop/sweatshop/internal/shared/db"
)

// Team represents a team entity
type Team struct {
	ID               string    `json:"id"`
	Name             string    `json:"name"`
	Description      string    `json:"description,omitempty"`
	LeadRuntimeType  string    `json:"leadRuntimeType"`
	LeadRuntimeModel string    `json:"leadRuntimeModel,omitempty"`
	CreatedAt        time.Time `json:"createdAt"`
	UpdatedAt        time.Time `json:"updatedAt"`
}

// CreateTeamInput represents input for creating a team
type CreateTeamInput struct {
	Name             string `json:"name" validate:"required"`
	Description      string `json:"description"`
	LeadRuntimeType  string `json:"leadRuntimeType"`
	LeadRuntimeModel string `json:"leadRuntimeModel"`
}

// UpdateTeamInput represents input for updating a team
type UpdateTeamInput struct {
	Name             string `json:"name"`
	Description      string `json:"description"`
	LeadRuntimeType  string `json:"leadRuntimeType"`
	LeadRuntimeModel string `json:"leadRuntimeModel"`
}

// GetAll retrieves all teams
func GetAll() ([]Team, error) {
	rows, err := db.DB.Query(`
		SELECT id, name, COALESCE(description, ''), COALESCE(lead_runtime_type, 'claude-code'),
		       COALESCE(lead_runtime_model, ''), created_at, updated_at
		FROM teams ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var teams []Team
	for rows.Next() {
		var t Team
		if err := rows.Scan(&t.ID, &t.Name, &t.Description, &t.LeadRuntimeType,
			&t.LeadRuntimeModel, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, err
		}
		teams = append(teams, t)
	}
	return teams, nil
}

// GetByID retrieves a team by ID
func GetByID(id string) (*Team, error) {
	var t Team
	err := db.DB.QueryRow(`
		SELECT id, name, COALESCE(description, ''), COALESCE(lead_runtime_type, 'claude-code'),
		       COALESCE(lead_runtime_model, ''), created_at, updated_at
		FROM teams WHERE id = ?
	`, id).Scan(&t.ID, &t.Name, &t.Description, &t.LeadRuntimeType,
		&t.LeadRuntimeModel, &t.CreatedAt, &t.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &t, nil
}

// Create creates a new team
func Create(t *Team) error {
	_, err := db.DB.Exec(`
		INSERT INTO teams (id, name, description, lead_runtime_type, lead_runtime_model, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`, t.ID, t.Name, t.Description, t.LeadRuntimeType, t.LeadRuntimeModel, t.CreatedAt, t.UpdatedAt)
	return err
}

// Update updates a team
func Update(t *Team) error {
	_, err := db.DB.Exec(`
		UPDATE teams SET name = ?, description = ?, lead_runtime_type = ?, lead_runtime_model = ?, updated_at = ?
		WHERE id = ?
	`, t.Name, t.Description, t.LeadRuntimeType, t.LeadRuntimeModel, t.UpdatedAt, t.ID)
	return err
}

// Delete deletes a team
func Delete(id string) error {
	_, err := db.DB.Exec("DELETE FROM teams WHERE id = ?", id)
	return err
}
