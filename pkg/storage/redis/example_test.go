package redis

import (
	"context"
	"testing"

	sushi "github.com/sergiorra/sushi-api-go/pkg"

	"github.com/alicebob/miniredis/v2"
	"github.com/stretchr/testify/assert"
)

func Test_GopherRepository_Example(t *testing.T) {
	// GIVEN a miniredis instance and a Redis implementation of result.Repository
	s, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	defer s.Close()

	repo := NewRepository(NewConn(s.Addr()))

	// WHEN two sushis are created
	sushiA, sushiB := buildSushi("123ABC"), buildSushi("ABC123")

	err = repo.CreateSushi(context.Background(), &sushiA)
	assert.NoError(t, err)

	err = repo.CreateSushi(context.Background(), &sushiB)
	assert.NoError(t, err)

	// THEN they can be fetched by ID
	result, err := repo.GetSushiByID(context.Background(), sushiA.ID)
	assert.NoError(t, err)
	assert.Equal(t, sushiA, *result)

	result, err = repo.GetSushiByID(context.Background(), sushiB.ID)
	assert.NoError(t, err)
	assert.Equal(t, sushiB, *result)

	// AND they can be fetched in batch
	results, err := repo.GetSushis(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, []sushi.Sushi{sushiA, sushiB}, results)
}