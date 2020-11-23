package cockroach

import (
	"context"
	"database/sql"
	"errors"

	_ "github.com/lib/pq"

	sushi "github.com/sergiorra/sushi-api-go/pkg"
)

type sushiRepository struct {
	db     *sql.DB
}

// NewRepository creates a cockroach repository with the necessary dependencies
func NewRepository(db *sql.DB) sushi.Repository {
	return sushiRepository{db: db}
}

func (r sushiRepository) CreateSushi(ctx context.Context, g *sushi.Sushi) error {
	return errors.New("method not implemented")
}

func (r sushiRepository) GetSushis(ctx context.Context) ([]sushi.Sushi, error) {
	return nil, errors.New("method not implemented")
}

func (r sushiRepository) DeleteSushi(ctx context.Context, ID string) error {
	return errors.New("method not implemented")
}

func (r sushiRepository) UpdateSushi(ctx context.Context, ID string, g *sushi.Sushi) error {
	return errors.New("method not implemented")
}

func (r sushiRepository) GetSushiByID(ctx context.Context, ID string) (*sushi.Sushi, error) {
	return nil, errors.New("method not implemented")
}
