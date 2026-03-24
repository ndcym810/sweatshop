// internal/project/service.go
package project

import (
	"time"

	"github.com/sweatshop/sweatshop/pkg/uuid"
)

// Service provides project business logic
type Service struct{}

// NewService creates a new project service
func NewService() *Service {
	return &Service{}
}

// ListByTeam returns all projects for a team
func (s *Service) ListByTeam(teamID string) ([]Project, error) {
	return GetAllByTeam(teamID)
}

// Get returns a project by ID
func (s *Service) Get(id string) (*Project, error) {
	return GetByID(id)
}

// Create creates a new project
func (s *Service) Create(teamID string, input CreateProjectInput) (*Project, error) {
	now := time.Now()
	p := &Project{
		ID:            uuid.New(),
		TeamID:        teamID,
		Name:          input.Name,
		Path:          input.Path,
		DefaultBranch: input.DefaultBranch,
		IsActive:      input.IsActive,
		CreatedAt:     now,
	}

	// Set default branch if not provided
	if p.DefaultBranch == "" {
		p.DefaultBranch = "main"
	}

	if err := Create(p); err != nil {
		return nil, err
	}
	return p, nil
}

// Update updates a project
func (s *Service) Update(id string, input UpdateProjectInput) (*Project, error) {
	p, err := GetByID(id)
	if err != nil {
		return nil, err
	}
	if p == nil {
		return nil, nil
	}

	if input.Name != "" {
		p.Name = input.Name
	}
	if input.Path != "" {
		p.Path = input.Path
	}
	if input.DefaultBranch != "" {
		p.DefaultBranch = input.DefaultBranch
	}
	if input.IsActive != nil {
		p.IsActive = *input.IsActive
	}

	if err := Update(p); err != nil {
		return nil, err
	}
	return p, nil
}

// Delete deletes a project
func (s *Service) Delete(id string) error {
	return Delete(id)
}
