package mysql

import (
	"context"
	"errors"
	"testing"
	"time"

	sushiapi "github.com/sergiorra/sushi-api-go/pkg"

	// sqlmock simulates any sql driver behavior in tests, without needing a real database connection
	"github.com/DATA-DOG/go-sqlmock"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func Test_SushiRepository_CreateSushi_RepositoryError(t *testing.T) {
	sushi := buildSushi()

	db, sqlMock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		assert.NoError(t, err)
	}

	sqlMock.ExpectExec(
		"INSERT INTO sushis (id, image_number, name, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)").
		WithArgs(sushi.ID, sushi.ImageNumber, sushi.Name, sushi.CreatedAt, sushi.UpdatedAt).
		WillReturnError(errors.New("database failed"))

	repo := NewRepository("sushis", db)
	err = repo.CreateSushi(context.Background(), &sushi)

	assert.Error(t, err)
	assert.NoError(t, sqlMock.ExpectationsWereMet())
}

func Test_SushiRepository_CreateSushi_Success(t *testing.T) {
	sushi := buildSushi()

	db, sqlMock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		assert.NoError(t, err)
	}

	sqlMock.ExpectExec(
		"INSERT INTO sushis (id, image_number, name, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)").
		WithArgs(sushi.ID, sushi.ImageNumber, sushi.Name, sushi.CreatedAt, sushi.UpdatedAt).
		WillReturnResult(sqlmock.NewResult(1, 1))

	repo := NewRepository("sushis", db)
	err = repo.CreateSushi(context.Background(), &sushi)

	assert.NoError(t, err)
	assert.NoError(t, sqlMock.ExpectationsWereMet())
}

func Test_SushiRepository_GetSushis_RepositoryError(t *testing.T) {
	db, sqlMock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		assert.NoError(t, err)
	}

	sqlMock.ExpectQuery(
		"SELECT sushis.id, sushis.image_number, sushis.name, sushis.created_at, sushis.updated_at FROM sushis").
		WillReturnError(errors.New("something-failed"))

	repo := NewRepository("sushis", db)
	_, err = repo.GetSushis(context.Background())

	assert.Error(t, err)
	assert.NoError(t, sqlMock.ExpectationsWereMet())
}

func Test_SushiRepository_GetSushis_NoRows(t *testing.T) {
	db, sqlMock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		assert.NoError(t, err)
	}

	sqlMock.ExpectQuery(
		"SELECT sushis.id, sushis.image_number, sushis.name, sushis.created_at, sushis.updated_at FROM sushis").
		WillReturnRows(sqlmock.NewRows(
			[]string{"id", "image_number", "name", "created_at", "updated_at"}),
		)

	repo := NewRepository("sushis", db)
	sushis, err := repo.GetSushis(context.Background())

	assert.NoError(t, err)
	assert.NoError(t, sqlMock.ExpectationsWereMet())

	assert.Len(t, sushis, 0)
}

func Test_SushiRepository_GetSushis_RowWithInvalidData(t *testing.T) {
	db, sqlMock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		assert.NoError(t, err)
	}

	sqlMock.ExpectQuery(
		"SELECT sushis.id, sushis.image_number, sushis.name, sushis.created_at, sushis.updated_at FROM sushis").
		WillReturnRows(sqlmock.NewRows(
			[]string{"id", "image_number", "name", "created_at", "updated_at"}).
			AddRow(nil, nil, nil, nil, nil), // This is a row failure as the data type is wrong
		)

	repo := NewRepository("sushis", db)
	_, err = repo.GetSushis(context.Background())

	assert.Error(t, err)
	assert.NoError(t, sqlMock.ExpectationsWereMet())
}

func Test_SushiRepository_GetSushis_Succeeded(t *testing.T) {
	expectedSushis := []sushiapi.Sushi{
		buildSushi(),
		buildSushi(),
	}

	db, sqlMock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		assert.NoError(t, err)
	}

	sqlMock.ExpectQuery(
		"SELECT sushis.id, sushis.image_number, sushis.name, sushis.created_at, sushis.updated_at FROM sushis").
		WillReturnRows(sqlmock.NewRows(
			[]string{"id", "image_number", "name", "created_at", "updated_at"}).
			AddRow(expectedSushis[0].ID, expectedSushis[0].ImageNumber, expectedSushis[0].Name, expectedSushis[0].CreatedAt, expectedSushis[0].UpdatedAt).
			AddRow(expectedSushis[1].ID, expectedSushis[1].ImageNumber, expectedSushis[1].Name, expectedSushis[1].CreatedAt, expectedSushis[1].UpdatedAt),
		)

	repo := NewRepository("sushis", db)
	sushis, err := repo.GetSushis(context.Background())

	assert.NoError(t, err)
	assert.NoError(t, sqlMock.ExpectationsWereMet())
	assert.Equal(t, expectedSushis, sushis)
}

func Test_SushiRepository_DeleteSushi_RepositoryError(t *testing.T) {
	sushiID := "1"

	db, sqlMock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		assert.NoError(t, err)
	}

	sqlMock.ExpectExec(
		"DELETE FROM sushis WHERE id = ?").
		WithArgs(sushiID).
		WillReturnError(errors.New("database failed"))

	repo := NewRepository("sushis", db)
	err = repo.DeleteSushi(context.Background(), sushiID)

	assert.Error(t, err)
	assert.NoError(t, sqlMock.ExpectationsWereMet())
}

func Test_SushiRepository_DeleteSushi_Success(t *testing.T) {
	sushiID := "1"

	db, sqlMock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		assert.NoError(t, err)
	}

	sqlMock.ExpectExec(
		"DELETE FROM sushis WHERE id = ?").
		WithArgs(sushiID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	repo := NewRepository("sushis", db)
	err = repo.DeleteSushi(context.Background(), sushiID)

	assert.NoError(t, err)
	assert.NoError(t, sqlMock.ExpectationsWereMet())
}

func Test_SushiRepository_UpdateSushi_RepositoryError(t *testing.T) {
	sushi := buildSushi()

	db, sqlMock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		assert.NoError(t, err)
	}

	sqlMock.ExpectExec(
		"UPDATE sushis SET id = ?, image_number = ?, name = ?, created_at = ?, updated_at = ? WHERE id = ?").
		WithArgs(sushi.ID, sushi.ImageNumber, sushi.Name, sushi.CreatedAt, sushi.UpdatedAt, sushi.ID).
		WillReturnError(errors.New("database failed"))

	repo := NewRepository("sushis", db)
	err = repo.UpdateSushi(context.Background(), sushi.ID, &sushi)

	assert.Error(t, err)
	assert.NoError(t, sqlMock.ExpectationsWereMet())
}

func Test_SushiRepository_UpdateSushi_NotFound(t *testing.T) {
	sushi := buildSushi()

	db, sqlMock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		assert.NoError(t, err)
	}

	sqlMock.ExpectExec(
		"UPDATE sushis SET id = ?, image_number = ?, name = ?, created_at = ?, updated_at = ? WHERE id = ?").
		WithArgs(sushi.ID, sushi.ImageNumber, sushi.Name, sushi.CreatedAt, sushi.UpdatedAt, sushi.ID).
		WillReturnResult(sqlmock.NewResult(0, 0))

	repo := NewRepository("sushis", db)
	err = repo.UpdateSushi(context.Background(), sushi.ID, &sushi)

	assert.Error(t, err)
	assert.NoError(t, sqlMock.ExpectationsWereMet())
}

func Test_SushiRepository_UpdateSushi_Success(t *testing.T) {
	sushi := buildSushi()

	db, sqlMock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		assert.NoError(t, err)
	}

	sqlMock.ExpectExec(
		"UPDATE sushis SET id = ?, image_number = ?, name = ?, created_at = ?, updated_at = ? WHERE id = ?").
		WithArgs(sushi.ID, sushi.ImageNumber, sushi.Name, sushi.CreatedAt, sushi.UpdatedAt, sushi.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	repo := NewRepository("sushis", db)
	err = repo.UpdateSushi(context.Background(), sushi.ID, &sushi)

	assert.NoError(t, err)
	assert.NoError(t, sqlMock.ExpectationsWereMet())
}

func Test_SushiRepository_GetSushiByID_RepositoryError(t *testing.T) {
	sushiID := "1"

	db, sqlMock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		assert.NoError(t, err)
	}

	sqlMock.ExpectQuery(
		"SELECT sushis.id, sushis.image_number, sushis.name, sushis.created_at, sushis.updated_at FROM sushis WHERE id = ?").
		WithArgs(sushiID).
		WillReturnError(errors.New("something-failed"))

	repo := NewRepository("sushis", db)
	_, err = repo.GetSushiByID(context.Background(), sushiID)

	assert.Error(t, err)
	assert.NoError(t, sqlMock.ExpectationsWereMet())
}

func Test_SushiRepository_GetSushiByID_NoRows(t *testing.T) {
	sushiID := "1"

	db, sqlMock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		assert.NoError(t, err)
	}

	sqlMock.ExpectQuery(
		"SELECT sushis.id, sushis.image_number, sushis.name, sushis.created_at, sushis.updated_at FROM sushis WHERE id = ?").
		WithArgs(sushiID).
		WillReturnRows(sqlmock.NewRows(
			[]string{"id", "image_number", "name", "created_at", "updated_at"}),
		)

	repo := NewRepository("sushis", db)
	_, err = repo.GetSushiByID(context.Background(), sushiID)

	assert.Error(t, err)
	assert.NoError(t, sqlMock.ExpectationsWereMet())
}

func Test_SushiRepository_GetSushiByID_RowWithInvalidData(t *testing.T) {
	sushiID := "1"

	db, sqlMock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		assert.NoError(t, err)
	}

	sqlMock.ExpectQuery(
		"SELECT sushis.id, sushis.image_number, sushis.name, sushis.created_at, sushis.updated_at FROM sushis WHERE id = ?").
		WillReturnRows(sqlmock.NewRows(
			[]string{"id", "image_number", "name", "created_at", "updated_at"}).
			AddRow(nil, nil, nil, nil, nil), // This is a row failure as the data type is wrong
		)

	repo := NewRepository("sushis", db)
	_, err = repo.GetSushiByID(context.Background(), sushiID)

	assert.Error(t, err)
	assert.NoError(t, sqlMock.ExpectationsWereMet())
}

func Test_SushiRepository_GetSushiByID_Succeeded(t *testing.T) {
	expectedSushi := buildSushi()

	db, sqlMock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		assert.NoError(t, err)
	}

	sqlMock.ExpectQuery(
		"SELECT sushis.id, sushis.image_number, sushis.name, sushis.created_at, sushis.updated_at FROM sushis WHERE id = ?",
	).
		WithArgs(expectedSushi.ID).
		WillReturnRows(sqlmock.NewRows(
			[]string{"id", "image_number", "name", "created_at", "updated_at"}).
			AddRow(expectedSushi.ID, expectedSushi.ImageNumber, expectedSushi.Name, expectedSushi.CreatedAt, expectedSushi.UpdatedAt),
		)

	repo := NewRepository("sushis", db)
	sushi, err := repo.GetSushiByID(context.Background(), expectedSushi.ID)

	assert.NoError(t, err)
	assert.NoError(t, sqlMock.ExpectationsWereMet())
	assert.Equal(t, &expectedSushi, sushi)
}

func buildSushi() sushiapi.Sushi {
	now := time.Now()
	return sushiapi.Sushi{
		ID:        		"123ABC",
		ImageNumber: 	"Test_image",
		Name:      		"Test_name",
		CreatedAt: 		&now,
		UpdatedAt: 		&now,
	}
}
