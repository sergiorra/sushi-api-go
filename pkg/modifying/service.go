package modifying

import (
	sushi "github.com/sergiorra/sushi-api-go/pkg"
)

// Service provides modifying operations
type Service interface {
	ModifySushi(ID, ImageNumber, Name string, Ingredients []string) error
}

type service struct {
	repository sushi.Repository
}

// NewService creates a modifying service with the necessary dependencies
func NewService(repository sushi.Repository) Service {
	return &service{repository}
}

// ModifySushi modify a sushi data
func (s *service) ModifySushi(ID, ImageNumber, Name string, Ingredients []string) error {
	sushi := sushi.New(ID, ImageNumber, Name, Ingredients)
	return s.repository.UpdateSushi(ID, sushi)
}
