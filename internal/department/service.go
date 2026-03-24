// internal/department/service.go
package department

import (
	"time"

	"github.com/sweatshop/sweatshop/pkg/uuid"
)

// Service provides department business logic
type Service struct{}

// NewService creates a new department service
func NewService() *Service {
	return &Service{}
}

// List returns all departments for a team
func (s *Service) List(teamID string) ([]Department, error) {
	return GetAllByTeam(teamID)
}

// Get returns a department by ID
func (s *Service) Get(id string) (*Department, error) {
	return GetByID(id)
}

// Create creates a new department
func (s *Service) Create(teamID string, input CreateDepartmentInput) (*Department, error) {
	d := &Department{
		ID:          uuid.New(),
		TeamID:      teamID,
		Name:        input.Name,
		Description: input.Description,
		SortOrder:   input.SortOrder,
		CreatedAt:   time.Now(),
	}

	if err := Create(d); err != nil {
		return nil, err
	}
	return d, nil
}

// Update updates a department
func (s *Service) Update(id string, input UpdateDepartmentInput) (*Department, error) {
	d, err := GetByID(id)
	if err != nil {
		return nil, err
	}
	if d == nil {
		return nil, nil
	}

	if input.Name != "" {
		d.Name = input.Name
	}
	if input.Description != "" {
		d.Description = input.Description
	}
	if input.SortOrder != nil {
		d.SortOrder = *input.SortOrder
	}

	if err := Update(d); err != nil {
		return nil, err
	}
	return d, nil
}

// Delete deletes a department
func (s *Service) Delete(id string) error {
	return Delete(id)
}
