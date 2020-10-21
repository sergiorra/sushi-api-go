package getting

import (
	sushi "github.com/sergiorra/sushi-api-go/pkg"
)

// Service provides getting operations
type Service interface {
	GetSushis() ([]sushi.Sushi, error)
	GetSushiByID(ID string) (*sushi.Sushi, error)
}

type service struct {
	repository sushi.Repository
}

// NewService creates a getting service with the necessary dependencies
func NewService(repository sushi.Repository) Service {
	return &service{repository}
}

// GetSushis returns all sushis
func (s *service) GetSushis() ([]sushi.Sushi, error) {
	return s.repository.GetSushis()
}

// GetSushiByID returns a sushi
func (s *service) GetSushiByID(ID string) (*sushi.Sushi, error) {
	return s.repository.GetSushiByID(ID)
}

