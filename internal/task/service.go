// internal/task/service.go
package task

import (
	"time"

	"github.com/sweatshop/sweatshop/pkg/uuid"
)

// Service provides task business logic
type Service struct{}

// NewService creates a new task service
func NewService() *Service {
	return &Service{}
}

// List returns all tasks for a team, optionally filtered by status
func (s *Service) List(teamID string, status string) ([]Task, error) {
	return GetAllByTeam(teamID, status)
}

// Get returns a task by ID
func (s *Service) Get(id string) (*Task, error) {
	return GetByID(id)
}

// Create creates a new task
func (s *Service) Create(teamID string, input CreateTaskInput) (*Task, error) {
	t := &Task{
		ID:          uuid.New(),
		TeamID:      teamID,
		ProjectID:   input.ProjectID,
		AssignedTo:  input.AssignedTo,
		Title:       input.Title,
		Description: input.Description,
		Status:      "pending",
		Priority:    "medium",
		CreatedAt:   time.Now(),
	}

	if input.Priority != "" {
		t.Priority = input.Priority
	}

	if err := Create(t); err != nil {
		return nil, err
	}
	return t, nil
}

// Update updates a task
func (s *Service) Update(id string, input UpdateTaskInput) (*Task, error) {
	t, err := GetByID(id)
	if err != nil {
		return nil, err
	}
	if t == nil {
		return nil, nil
	}

	if input.ProjectID != nil {
		t.ProjectID = input.ProjectID
	}
	if input.AssignedTo != nil {
		t.AssignedTo = input.AssignedTo
	}
	if input.Title != "" {
		t.Title = input.Title
	}
	if input.Description != "" {
		t.Description = input.Description
	}
	if input.Status != "" {
		t.Status = input.Status
	}
	if input.Priority != "" {
		t.Priority = input.Priority
	}
	if input.CompletedAt != nil {
		t.CompletedAt = input.CompletedAt
	}

	if err := Update(t); err != nil {
		return nil, err
	}
	return t, nil
}

// Delete deletes a task
func (s *Service) Delete(id string) error {
	return Delete(id)
}
