package adding

import (
	sushi "github.com/sergiorra/sushi-api-go/pkg"
)

// Service provides adding operations.
type Service interface {
	AddSushi(ID, ImageNumber, Name string, Ingredients []string) error
}

type service struct {
	repository sushi.Repository
}

// NewService creates an adding service with the necessary dependencies
func NewService(repository sushi.Repository) Service {
	return &service{repository}
}

// AddSushi adds the given sushi to storage
func (s *service) AddSushi(ID, ImageNumber, Name string, Ingredients []string) error {
	sushi := sushi.New(ID, ImageNumber, Name, Ingredients)
	return s.repository.CreateSushi(sushi)
}
