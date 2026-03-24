// internal/department/model.go
package department

import (
	"database/sql"
	"time"

	"github.com/sweatshop/sweatshop/internal/shared/db"
)

// Department represents a department entity
type Department struct {
	ID          string    `json:"id"`
	TeamID      string    `json:"teamId"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	SortOrder   int       `json:"sortOrder"`
	CreatedAt   time.Time `json:"createdAt"`
}

// CreateDepartmentInput represents input for creating a department
type CreateDepartmentInput struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
	SortOrder   int    `json:"sortOrder"`
}

// UpdateDepartmentInput represents input for updating a department
type UpdateDepartmentInput struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	SortOrder   *int   `json:"sortOrder"`
}

// GetAllByTeam retrieves all departments for a team
func GetAllByTeam(teamID string) ([]Department, error) {
	rows, err := db.DB.Query(`
		SELECT id, team_id, name, COALESCE(description, ''), sort_order, created_at
		FROM departments WHERE team_id = ? ORDER BY sort_order, created_at
	`, teamID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var departments []Department
	for rows.Next() {
		var d Department
		if err := rows.Scan(&d.ID, &d.TeamID, &d.Name, &d.Description, &d.SortOrder, &d.CreatedAt); err != nil {
			return nil, err
		}
		departments = append(departments, d)
	}
	return departments, nil
}

// GetByID retrieves a department by ID
func GetByID(id string) (*Department, error) {
	var d Department
	err := db.DB.QueryRow(`
		SELECT id, team_id, name, COALESCE(description, ''), sort_order, created_at
		FROM departments WHERE id = ?
	`, id).Scan(&d.ID, &d.TeamID, &d.Name, &d.Description, &d.SortOrder, &d.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &d, nil
}

// Create creates a new department
func Create(d *Department) error {
	_, err := db.DB.Exec(`
		INSERT INTO departments (id, team_id, name, description, sort_order, created_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`, d.ID, d.TeamID, d.Name, d.Description, d.SortOrder, d.CreatedAt)
	return err
}

// Update updates a department
func Update(d *Department) error {
	_, err := db.DB.Exec(`
		UPDATE departments SET name = ?, description = ?, sort_order = ?
		WHERE id = ?
	`, d.Name, d.Description, d.SortOrder, d.ID)
	return err
}

// Delete deletes a department
func Delete(id string) error {
	_, err := db.DB.Exec("DELETE FROM departments WHERE id = ?", id)
	return err
}
