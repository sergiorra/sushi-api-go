package cockroach

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"

	sushi "github.com/sergiorra/sushi-api-go/pkg"
)

/*
	$ cockroach start-single-node --insecure
	--- open new terminal tab ---
	$ cockroach sql --insecure
	$ create database sushiapi;
	$ set database = sushiapi;
	$ CREATE TABLE sushis (
			id STRING(32),
			image_number STRING(100) NOT NULL,
			name STRING NULL,
			created_at TIMESTAMPTZ NOT NULL,
	   		updated_at TIMESTAMPTZ,
	   		PRIMARY KEY ("id")
		);
	--- open new terminal tab ---
	$ go run cmd/sushi-api/main.go -database cockroach
 */

type sushiRepository struct {
	db     *sql.DB
}

// NewRepository creates a cockroach repository with the necessary dependencies
func NewRepository(db *sql.DB) sushi.Repository {
	return sushiRepository{db: db}
}

func (r sushiRepository) CreateSushi(ctx context.Context, s *sushi.Sushi) error {
	fmt.Println("creating")
	sqlStm := `INSERT INTO sushis (id, image_number, name, created_at) 
				VALUES ($1, $2, $3, NOW())`
	_, err := r.db.Exec(sqlStm, s.ID, s.ImageNumber, s.Name)
	fmt.Println("err", err)
	if err != nil {
		return err
	}
	return nil
}

func (r sushiRepository) GetSushis(ctx context.Context) ([]sushi.Sushi, error) {
	sqlStm := `SELECT id, image_number, name, created_at, updated_at FROM sushis`
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
	sqlStm := `DELETE FROM sushis WHERE id=$1`
	_, err := r.db.Exec(sqlStm, ID)
	if err != nil {
		return err
	}
	return nil
}

func (r sushiRepository) UpdateSushi(ctx context.Context, ID string, s *sushi.Sushi) error {
	sqlStm := `UPDATE sushis SET image_number=$1, name=$2, updated_at=$3 WHERE id=$4`
	_, err := r.db.Exec(sqlStm, s.ImageNumber, s.Name, s.UpdatedAt, ID)
	if err != nil {
		return err
	}
	return nil
}

func (r sushiRepository) GetSushiByID(ctx context.Context, ID string) (*sushi.Sushi, error) {
	sqlStm := `SELECT id, image_number, name, created_at, updated_at FROM sushis WHERE id=$1`
	rows, err := r.db.Query(sqlStm, ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var s sushi.Sushi

	if rows.Next() {
		if err := rows.Scan(&s.ID, &s.ImageNumber, &s.Name, &s.CreatedAt, &s.UpdatedAt); err != nil {
			log.Println(err)
		}
	}
	return &s, nil
}
