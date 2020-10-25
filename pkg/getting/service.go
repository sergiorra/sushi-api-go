package getting

import (
	"context"

	sushi "github.com/sergiorra/sushi-api-go/pkg"
	"github.com/sergiorra/sushi-api-go/pkg/log"
)

// Service provides getting operations
type Service interface {
	GetSushis(ctx context.Context) ([]sushi.Sushi, error)
	GetSushiByID(ctx context.Context, ID string) *sushi.Sushi
}

type service struct {
	repository sushi.Repository
	logger     log.Logger
}

// NewService creates a getting service with the necessary dependencies
func NewService(repository sushi.Repository, logger log.Logger) Service {
	return &service{repository, logger}
}

// GetSushis returns all sushis
func (s *service) GetSushis(ctx context.Context) ([]sushi.Sushi, error) {
	return s.repository.GetSushis(ctx)
}

// GetSushiByID returns a sushi
func (s *service) GetSushiByID(ctx context.Context, ID string) *sushi.Sushi {
	g, err := s.repository.GetSushiByID(ctx, ID)

	if err != nil {
		s.logger.UnexpectedError(ctx, err)
		return nil
	}

	return g
}

