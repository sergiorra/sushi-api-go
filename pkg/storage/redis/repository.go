package redis

import (
	"context"
	"encoding/json"
	"errors"

	sushiapi "github.com/sergiorra/sushi-api-go/pkg"

	"github.com/gomodule/redigo/redis"
	_ "github.com/lib/pq"
)

const (
	onlyIfExists = "XX"
)

type sushiRepository struct {
	pool *redis.Pool
}

// NewRepository instances a Redis implementation of the sushiapi.Repository
func NewRepository(pool *redis.Pool) sushiapi.Repository {
	return sushiRepository{
		pool: pool,
	}
}

// CreateSushi satisfies the sushiapi.Repository interface
func (s sushiRepository) CreateSushi(ctx context.Context, sushi *sushiapi.Sushi) error {
	bytes, err := json.Marshal(sushi)
	if err != nil {
		return err
	}

	conn, err := s.pool.GetContext(ctx)
	if err != nil {
		return err
	}

	_, err = conn.Do("SET", sushi.ID, string(bytes))
	return err
}

// GetSushis satisfies the sushiapi.Repository interface
func (s sushiRepository) GetSushis(ctx context.Context) ([]sushiapi.Sushi, error) {
	conn, err := s.pool.GetContext(ctx)
	if err != nil {
		return nil, err
	}

	keys, err := redis.Strings(conn.Do("KEYS", "*"))
	if err != nil {
		return nil, err
	}

	if len(keys) == 0 {
		return []sushiapi.Sushi{}, nil
	}

	args := make([]interface{}, 0, len(keys))
	for _, key := range keys {
		args = append(args, key)
	}

	results, err := redis.Strings(conn.Do("MGET", args...))
	if err != nil {
		return nil, err
	}

	sushis := make([]sushiapi.Sushi, 0, len(results))
	for _, result := range results {
		sushi := sushiapi.Sushi{}

		err := json.Unmarshal([]byte(result), &sushi)
		if err != nil {
			return nil, err
		}

		sushis = append(sushis, sushi)
	}
	return sushis, nil
}

// GetSushiByID satisfies the sushiapi.Repository interface
func (s sushiRepository) GetSushiByID(ctx context.Context, ID string) (*sushiapi.Sushi, error) {
	conn, err := s.pool.GetContext(ctx)
	if err != nil {
		return nil, err
	}

	result, err := redis.String(conn.Do("GET", ID))
	if err != nil {
		return nil, err
	}

	if result == "" {
		return nil, errors.New("not found")
	}

	sushi := &sushiapi.Sushi{}
	err = json.Unmarshal([]byte(result), sushi)

	return sushi, err
}

// DeleteSushi satisfies the sushiapi.Repository interface
func (s sushiRepository) DeleteSushi(ctx context.Context, ID string) error {
	conn, err := s.pool.GetContext(ctx)
	if err != nil {
		return err
	}

	_, err = conn.Do("DEL", ID)
	return err
}

// UpdateSushi satisfies the sushiapi.Repository interface
func (s sushiRepository) UpdateSushi(ctx context.Context, ID string, sushi *sushiapi.Sushi) error {
	bytes, err := json.Marshal(sushi)
	if err != nil {
		return err
	}

	conn, err := s.pool.GetContext(ctx)
	if err != nil {
		return err
	}

	result, err := conn.Do("SET", ID, string(bytes), onlyIfExists)
	if result == nil {
		return errors.New("not found")
	}
	return err
}
