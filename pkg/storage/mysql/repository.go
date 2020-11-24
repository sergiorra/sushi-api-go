package mysql

import (
	"context"
	"database/sql"
	"errors"
	"time"

	// sqlbuilder builds SQL string automatically given some arguments (like table, object,...)
	"github.com/huandu/go-sqlbuilder"
	_ "github.com/lib/pq"
	sushiapi "github.com/sergiorra/sushi-api-go/pkg"
)

type sushiRepository struct {
	table string
	db    *sql.DB
}

// NewRepository instances a MySQL implementation of the sushiapi.Repository
func NewRepository(table string, db *sql.DB) sushiapi.Repository {
	return sushiRepository{table: table, db: db}
}

// CreateSushi satisfies the sushiapi.Repository interface
func (r sushiRepository) CreateSushi(ctx context.Context, g *sushiapi.Sushi) error {
	insertBuilder := sqlbuilder.NewStruct(new(sqlSushi)).InsertInto(
		r.table,
		sqlSushi{
			ID:        		g.ID,
			ImageNumber:    g.ImageNumber,
			Name:     		g.Name,
			CreatedAt: 		g.CreatedAt,
			UpdatedAt: 		g.UpdatedAt,
		},
	)

	query, args := insertBuilder.Build()
	_, err := r.db.ExecContext(ctx, query, args...)
	return err
}

// GetSushis satisfies the sushiapi.Repository interface
func (r sushiRepository) GetSushis(ctx context.Context) ([]sushiapi.Sushi, error) {
	sqlSushiStruct := sqlbuilder.NewStruct(new(sqlSushi))

	selectBuilder := sqlSushiStruct.SelectFrom(r.table)
	query, args := selectBuilder.Build()

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	defer func() { _ = rows.Close() }()

	var sushis []sushiapi.Sushi
	for rows.Next() {
		sqlSushi := sqlSushi{}

		err := rows.Scan(sqlSushiStruct.Addr(&sqlSushi)...)
		if err != nil {
			return nil, err
		}

		sushis = append(sushis, sushiapi.Sushi{
			ID:        		sqlSushi.ID,
			ImageNumber:    sqlSushi.ImageNumber,
			Name:     		sqlSushi.Name,
			CreatedAt: 		sqlSushi.CreatedAt,
			UpdatedAt: 		sqlSushi.UpdatedAt,
		})
	}

	return sushis, nil
}

// DeleteSushi satisfies the sushiapi.Repository interface
func (r sushiRepository) DeleteSushi(ctx context.Context, ID string) error {
	deleteBuilder := sqlbuilder.NewStruct(new(sqlSushi)).DeleteFrom(r.table)
	query, args := deleteBuilder.Where(
		deleteBuilder.Equal("id", ID),
	).Build()

	_, err := r.db.ExecContext(ctx, query, args...)
	return err
}

// UpdateSushi satisfies the sushiapi.Repository interface
func (r sushiRepository) UpdateSushi(ctx context.Context, ID string, g *sushiapi.Sushi) error {
	updateBuilder := sqlbuilder.NewStruct(new(sqlSushi)).Update(
		r.table,
		sqlSushi{
			ID:        		g.ID,
			ImageNumber:    g.ImageNumber,
			Name:     		g.Name,
			CreatedAt: 		g.CreatedAt,
			UpdatedAt: 		g.UpdatedAt,
		},
	)

	query, args := updateBuilder.Where(
		updateBuilder.Equal("id", ID),
	).Build()

	result, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("not found")
	}

	return nil
}

// GetSushiByID satisfies the sushiapi.Repository interface
func (r sushiRepository) GetSushiByID(ctx context.Context, ID string) (*sushiapi.Sushi, error) {
	sqlSushiStruct := sqlbuilder.NewStruct(new(sqlSushi))

	selectBuilder := sqlSushiStruct.SelectFrom(r.table)

	query, args := selectBuilder.Where(
		selectBuilder.Equal("id", ID),
	).Build()

	row := r.db.QueryRowContext(ctx, query, args...)

	sqlSushi := sqlSushi{}

	err := row.Scan(sqlSushiStruct.Addr(&sqlSushi)...)
	if err != nil {
		return nil, err
	}

	return &sushiapi.Sushi{
		ID:        		sqlSushi.ID,
		ImageNumber:    sqlSushi.ImageNumber,
		Name:     		sqlSushi.Name,
		CreatedAt: 		sqlSushi.CreatedAt,
		UpdatedAt: 		sqlSushi.UpdatedAt,
	}, nil
}

type sqlSushi struct {
	ID        		string     `db:"id"`
	ImageNumber     string     `db:"image_number"`
	Name     		string     `db:"name"`
	CreatedAt 		*time.Time `db:"created_at"`
	UpdatedAt 		*time.Time `db:"updated_at"`
}