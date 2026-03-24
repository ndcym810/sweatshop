// internal/task/model.go
package task

import (
	"database/sql"
	"time"

	"github.com/sweatshop/sweatshop/internal/shared/db"
)

// Task represents a task entity
type Task struct {
	ID          string     `json:"id"`
	TeamID      string     `json:"teamId"`
	ProjectID   *string    `json:"projectId,omitempty"`
	AssignedTo  *string    `json:"assignedTo,omitempty"`
	Title       string     `json:"title"`
	Description string     `json:"description,omitempty"`
	Status      string     `json:"status"`
	Priority    string     `json:"priority"`
	CreatedAt   time.Time  `json:"createdAt"`
	CompletedAt *time.Time `json:"completedAt,omitempty"`
}

// CreateTaskInput represents input for creating a task
type CreateTaskInput struct {
	ProjectID   *string `json:"projectId"`
	AssignedTo  *string `json:"assignedTo"`
	Title       string  `json:"title" validate:"required"`
	Description string  `json:"description"`
	Priority    string  `json:"priority"`
}

// UpdateTaskInput represents input for updating a task
type UpdateTaskInput struct {
	ProjectID   *string    `json:"projectId"`
	AssignedTo  *string    `json:"assignedTo"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Status      string     `json:"status"`
	Priority    string     `json:"priority"`
	CompletedAt *time.Time `json:"completedAt"`
}

// GetAllByTeam retrieves all tasks for a team, optionally filtered by status
func GetAllByTeam(teamID string, status string) ([]Task, error) {
	var query string
	var args []interface{}

	if status != "" {
		query = `
			SELECT id, team_id, project_id, assigned_to, title, COALESCE(description, ''),
			       status, priority, created_at, completed_at
			FROM tasks WHERE team_id = ? AND status = ? ORDER BY created_at DESC
		`
		args = []interface{}{teamID, status}
	} else {
		query = `
			SELECT id, team_id, project_id, assigned_to, title, COALESCE(description, ''),
			       status, priority, created_at, completed_at
			FROM tasks WHERE team_id = ? ORDER BY created_at DESC
		`
		args = []interface{}{teamID}
	}

	rows, err := db.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var t Task
		var projectID, assignedTo sql.NullString
		var completedAt sql.NullTime
		if err := rows.Scan(&t.ID, &t.TeamID, &projectID, &assignedTo, &t.Title, &t.Description,
			&t.Status, &t.Priority, &t.CreatedAt, &completedAt); err != nil {
			return nil, err
		}
		if projectID.Valid {
			t.ProjectID = &projectID.String
		}
		if assignedTo.Valid {
			t.AssignedTo = &assignedTo.String
		}
		if completedAt.Valid {
			t.CompletedAt = &completedAt.Time
		}
		tasks = append(tasks, t)
	}
	return tasks, nil
}

// GetByID retrieves a task by ID
func GetByID(id string) (*Task, error) {
	var t Task
	var projectID, assignedTo sql.NullString
	var completedAt sql.NullTime
	err := db.DB.QueryRow(`
		SELECT id, team_id, project_id, assigned_to, title, COALESCE(description, ''),
		       status, priority, created_at, completed_at
		FROM tasks WHERE id = ?
	`, id).Scan(&t.ID, &t.TeamID, &projectID, &assignedTo, &t.Title, &t.Description,
		&t.Status, &t.Priority, &t.CreatedAt, &completedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if projectID.Valid {
		t.ProjectID = &projectID.String
	}
	if assignedTo.Valid {
		t.AssignedTo = &assignedTo.String
	}
	if completedAt.Valid {
		t.CompletedAt = &completedAt.Time
	}
	return &t, nil
}

// Create creates a new task
func Create(t *Task) error {
	var projectID, assignedTo interface{}
	var completedAt interface{}

	if t.ProjectID != nil {
		projectID = *t.ProjectID
	}
	if t.AssignedTo != nil {
		assignedTo = *t.AssignedTo
	}
	if t.CompletedAt != nil {
		completedAt = *t.CompletedAt
	}

	_, err := db.DB.Exec(`
		INSERT INTO tasks (id, team_id, project_id, assigned_to, title, description, status, priority, created_at, completed_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, t.ID, t.TeamID, projectID, assignedTo, t.Title, t.Description, t.Status, t.Priority, t.CreatedAt, completedAt)
	return err
}

// Update updates a task
func Update(t *Task) error {
	var projectID, assignedTo interface{}
	var completedAt interface{}

	if t.ProjectID != nil {
		projectID = *t.ProjectID
	}
	if t.AssignedTo != nil {
		assignedTo = *t.AssignedTo
	}
	if t.CompletedAt != nil {
		completedAt = *t.CompletedAt
	}

	_, err := db.DB.Exec(`
		UPDATE tasks SET project_id = ?, assigned_to = ?, title = ?, description = ?, status = ?, priority = ?, completed_at = ?
		WHERE id = ?
	`, projectID, assignedTo, t.Title, t.Description, t.Status, t.Priority, completedAt, t.ID)
	return err
}

// Delete deletes a task
func Delete(id string) error {
	_, err := db.DB.Exec("DELETE FROM tasks WHERE id = ?", id)
	return err
}
