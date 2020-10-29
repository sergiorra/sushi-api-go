package redis

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	sushiapi "github.com/sergiorra/sushi-api-go/pkg"

	"github.com/gomodule/redigo/redis"
	_ "github.com/lib/pq"
	"github.com/rafaeljusto/redigomock"
	"github.com/stretchr/testify/assert"

)

func Test_SushiRepository_CreateSushi_RepositoryError(t *testing.T) {
	sushi := buildSushi("01D3XZ38KDR")

	conn := redigomock.NewConn()
	conn.Command("SET", sushi.ID, sushiToJSONString(sushi)).ExpectError(errors.New("something failed"))

	repo := NewRepository(wrapRedisConn(conn))
	err := repo.CreateSushi(context.Background(), &sushi)

	assert.Error(t, err)
	assert.NoError(t, conn.ExpectationsWereMet())
}

func Test_SushiRepository_CreateSushi_Success(t *testing.T) {
	sushi := buildSushi("01D3XZ38KDR")

	conn := redigomock.NewConn()
	conn.Command("SET", sushi.ID, sushiToJSONString(sushi)).Expect("OK")

	repo := NewRepository(wrapRedisConn(conn))
	err := repo.CreateSushi(context.Background(), &sushi)

	assert.NoError(t, err)
	assert.NoError(t, conn.ExpectationsWereMet())
}

func Test_SushiRepository_GetSushis_RepositoryError(t *testing.T) {
	conn := redigomock.NewConn()
	conn.Command("KEYS", "*").ExpectError(errors.New("something failed"))

	repo := NewRepository(wrapRedisConn(conn))
	_, err := repo.GetSushis(context.Background())

	assert.Error(t, err)
	assert.NoError(t, conn.ExpectationsWereMet())
}

func Test_SushiRepository_GetSushis_NoRows(t *testing.T) {
	conn := redigomock.NewConn()
	conn.Command("KEYS", "*").Expect([]interface{}{})

	repo := NewRepository(wrapRedisConn(conn))
	sushis, err := repo.GetSushis(context.Background())

	assert.NoError(t, err)
	assert.NoError(t, conn.ExpectationsWereMet())
	assert.Len(t, sushis, 0)
}

func Test_SushiRepository_GetSushis_RowWithInvalidData(t *testing.T) {
	conn := redigomock.NewConn()
	conn.Command("KEYS", "*").Expect([]interface{}{"123", "456"})
	conn.Command("MGET", "123", "456").Expect([]interface{}{"invalid-data"})

	repo := NewRepository(wrapRedisConn(conn))
	_, err := repo.GetSushis(context.Background())

	assert.Error(t, err)
	assert.NoError(t, conn.ExpectationsWereMet())
}

func Test_SushiRepository_GetSushis_Succeeded(t *testing.T) {
	sushiA, sushiB := buildSushi("01D3XZ38KDR"), buildSushi("01D3XZ38TRE")
	expectedSushis := []sushiapi.Sushi{sushiA, sushiB}

	conn := redigomock.NewConn()
	conn.Command("KEYS", "*").Expect([]interface{}{sushiA.ID, sushiB.ID})
	conn.Command("MGET", sushiA.ID, sushiB.ID).Expect(
		[]interface{}{sushiToJSONString(sushiA), sushiToJSONString(sushiB)},
	)

	repo := NewRepository(wrapRedisConn(conn))
	sushis, err := repo.GetSushis(context.Background())

	assert.NoError(t, err)
	assert.NoError(t, conn.ExpectationsWereMet())
	assert.Equal(t, expectedSushis, sushis)
}

func Test_SushiRepository_DeleteSushi_RepositoryError(t *testing.T) {
	sushiID := "01D3XZ38KDR"

	conn := redigomock.NewConn()
	conn.Command("DEL", sushiID).ExpectError(errors.New("something failed"))

	repo := NewRepository(wrapRedisConn(conn))
	err := repo.DeleteSushi(context.Background(), sushiID)

	assert.Error(t, err)
	assert.NoError(t, conn.ExpectationsWereMet())
}

func Test_SushiRepository_DeleteSushi_Success(t *testing.T) {
	sushiID := "01D3XZ38KDR"

	conn := redigomock.NewConn()
	conn.Command("DEL", sushiID).Expect(1)

	repo := NewRepository(wrapRedisConn(conn))
	err := repo.DeleteSushi(context.Background(), sushiID)

	assert.NoError(t, err)
	assert.NoError(t, conn.ExpectationsWereMet())
}

func Test_SushiRepository_UpdateSushi_RepositoryError(t *testing.T) {
	sushi := buildSushi("01D3XZ38KDR")

	conn := redigomock.NewConn()
	conn.Command("SET", sushi.ID, sushiToJSONString(sushi), "XX").ExpectError(errors.New("something failed"))

	repo := NewRepository(wrapRedisConn(conn))
	err := repo.UpdateSushi(context.Background(), sushi.ID, &sushi)

	assert.Error(t, err)
	assert.NoError(t, conn.ExpectationsWereMet())
}

func Test_SushiRepository_UpdateSushi_NotFound(t *testing.T) {
	sushi := buildSushi("01D3XZ38KDR")

	conn := redigomock.NewConn()
	conn.Command("SET", sushi.ID, sushiToJSONString(sushi), "XX").Expect(nil)

	repo := NewRepository(wrapRedisConn(conn))
	err := repo.UpdateSushi(context.Background(), sushi.ID, &sushi)

	assert.Error(t, err)
	assert.NoError(t, conn.ExpectationsWereMet())
}

func Test_SushiRepository_UpdateSushi_Success(t *testing.T) {
	sushi := buildSushi("01D3XZ38KDR")

	conn := redigomock.NewConn()
	conn.Command("SET", sushi.ID, sushiToJSONString(sushi), "XX").Expect("OK")

	repo := NewRepository(wrapRedisConn(conn))
	err := repo.UpdateSushi(context.Background(), sushi.ID, &sushi)

	assert.NoError(t, err)
	assert.NoError(t, conn.ExpectationsWereMet())
}

func Test_SushiRepository_GetSushiByID_RepositoryError(t *testing.T) {
	sushiID := "01D3XZ38KDR"

	conn := redigomock.NewConn()
	conn.Command("GET", sushiID).ExpectError(errors.New("something failed"))

	repo := NewRepository(wrapRedisConn(conn))
	_, err := repo.GetSushiByID(context.Background(), sushiID)

	assert.Error(t, err)
	assert.NoError(t, conn.ExpectationsWereMet())
}

func Test_SushiRepository_GetSushiByID_NoRows(t *testing.T) {
	sushiID := "01D3XZ38KDR"

	conn := redigomock.NewConn()
	conn.Command("GET", sushiID).Expect(nil)

	repo := NewRepository(wrapRedisConn(conn))
	_, err := repo.GetSushiByID(context.Background(), sushiID)

	assert.Error(t, err)
	assert.NoError(t, conn.ExpectationsWereMet())
}

func Test_SushiRepository_GetSushiByID_RowWithInvalidData(t *testing.T) {
	sushiID := "01D3XZ38KDR"

	conn := redigomock.NewConn()
	conn.Command("GET", sushiID).Expect("invalid-data")

	repo := NewRepository(wrapRedisConn(conn))
	_, err := repo.GetSushiByID(context.Background(), sushiID)

	assert.Error(t, err)
	assert.NoError(t, conn.ExpectationsWereMet())
}

func Test_SushiRepository_GetSushiByID_Succeeded(t *testing.T) {
	sushiID := "01D3XZ38KDR"
	expectedSushi := buildSushi(sushiID)

	conn := redigomock.NewConn()
	conn.Command("GET", sushiID).Expect(sushiToJSONString(expectedSushi))

	repo := NewRepository(wrapRedisConn(conn))
	sushi, err := repo.GetSushiByID(context.Background(), sushiID)

	assert.NoError(t, err)
	assert.NoError(t, conn.ExpectationsWereMet())
	assert.Equal(t, &expectedSushi, sushi)
}

func buildSushi(ID string) sushiapi.Sushi {
	return sushiapi.Sushi{
		ID:    ID,
		ImageNumber: "3",
		Name:  "Salmon Roll",
		Ingredients:   []string {"Salmon"},
	}
}

func sushiToJSONString(sushi sushiapi.Sushi) string {
	bytes, _ := json.Marshal(&sushi)
	return string(bytes)
}

func wrapRedisConn(conn redis.Conn) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial:        func() (redis.Conn, error) { return conn, nil },
	}
}