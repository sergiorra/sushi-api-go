package removing

import (
	gopher "github.com/sergiorra/sushi-api-go/pkg"
)

// Service provides removing operations.
type Service interface {
	RemoveSushi(ID string) error
}

type service struct {
	repository gopher.Repository
}

// NewService creates a removing service with the necessary dependencies
func NewService(repository gopher.Repository) Service {
	return &service{repository}
}

// RemoveSushi remove sushi from the storage
func (s *service) RemoveSushi(ID string) error {
	return s.repository.DeleteSushi(ID)
}
