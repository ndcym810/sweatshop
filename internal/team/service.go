// internal/team/service.go
package team

import (
	"time"

	"github.com/sweatshop/sweatshop/pkg/uuid"
)

// Service provides team business logic
type Service struct{}

// NewService creates a new team service
func NewService() *Service {
	return &Service{}
}

// List returns all teams
func (s *Service) List() ([]Team, error) {
	return GetAll()
}

// Get returns a team by ID
func (s *Service) Get(id string) (*Team, error) {
	return GetByID(id)
}

// Create creates a new team
func (s *Service) Create(input CreateTeamInput) (*Team, error) {
	now := time.Now()
	t := &Team{
		ID:               uuid.New(),
		Name:             input.Name,
		Description:      input.Description,
		LeadRuntimeType:  input.LeadRuntimeType,
		LeadRuntimeModel: input.LeadRuntimeModel,
		CreatedAt:        now,
		UpdatedAt:        now,
	}

	// Set default runtime type
	if t.LeadRuntimeType == "" {
		t.LeadRuntimeType = "claude-code"
	}

	if err := Create(t); err != nil {
		return nil, err
	}
	return t, nil
}

// Update updates a team
func (s *Service) Update(id string, input UpdateTeamInput) (*Team, error) {
	t, err := GetByID(id)
	if err != nil {
		return nil, err
	}
	if t == nil {
		return nil, nil
	}

	if input.Name != "" {
		t.Name = input.Name
	}
	if input.Description != "" {
		t.Description = input.Description
	}
	if input.LeadRuntimeType != "" {
		t.LeadRuntimeType = input.LeadRuntimeType
	}
	if input.LeadRuntimeModel != "" {
		t.LeadRuntimeModel = input.LeadRuntimeModel
	}
	t.UpdatedAt = time.Now()

	if err := Update(t); err != nil {
		return nil, err
	}
	return t, nil
}

// Delete deletes a team
func (s *Service) Delete(id string) error {
	return Delete(id)
}
