package cockroach

import (
	"context"
	"database/sql"
	"errors"
	"log"

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

func (r sushiRepository) CreateSushi(ctx context.Context, s *sushi.Sushi) error {
	sqlStm := `INSERT INTO sushis (id, image_number, name, created_at) 
				VALUES ($1, $2, $3, NOW())`
	_, err := r.db.Exec(sqlStm, s.ID, s.ImageNumber, s.Name)
	if err != nil {
		return err
	}
	return nil
}

func (r sushiRepository) GetSushis(ctx context.Context) ([]sushi.Sushi, error) {
	sqlStm := `SELECT id, image_number, name, created_at, updated_at FROM gophers`
	rows, err := r.db.Query(sqlStm)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var sushis []sushi.Sushi

	for rows.Next() {
		var s sushi.Sushi
		if err := rows.Scan(&s.ID, &s.ImageNumber, &s.Name, &s.CreatedAt, &s.UpdatedAt); err != nil {
			log.Println(err)
			continue
		}
		sushis = append(sushis, s)
	}
	return sushis, nil
}

func (r sushiRepository) DeleteSushi(ctx context.Context, ID string) error {
	return errors.New("method not implemented")
}

func (r sushiRepository) UpdateSushi(ctx context.Context, ID string, s *sushi.Sushi) error {
	return errors.New("method not implemented")
}

func (r sushiRepository) GetSushiByID(ctx context.Context, ID string) (*sushi.Sushi, error) {
	return nil, errors.New("method not implemented")
}
