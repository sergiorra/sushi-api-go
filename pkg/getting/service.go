package getting

import (
	"context"
	sushi "github.com/sergiorra/sushi-api-go/pkg"
)

// Service provides getting operations
type Service interface {
	GetSushis(ctx context.Context) ([]sushi.Sushi, error)
	GetSushiByID(ctx context.Context, ID string) (*sushi.Sushi, error)
}

type service struct {
	repository sushi.Repository
}

// NewService creates a getting service with the necessary dependencies
func NewService(repository sushi.Repository) Service {
	return &service{repository}
}

// GetSushis returns all sushis
func (s *service) GetSushis(ctx context.Context) ([]sushi.Sushi, error) {
	return s.repository.GetSushis(ctx)
}

// GetSushiByID returns a sushi
func (s *service) GetSushiByID(ctx context.Context, ID string) (*sushi.Sushi, error) {
	return s.repository.GetSushiByID(ctx, ID)
}

