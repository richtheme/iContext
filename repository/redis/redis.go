package redis

import (
	"context"
	"errors"
	"iContext/internal/models"
	"strconv"

	"github.com/redis/go-redis/v9"
)

type RedisStorage struct {
	Client *redis.Client
}

func dial(addr string) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "",
		DB:       0,
	})
}

func RedisNew(redisAddr string) (*RedisStorage, error) {
	client := dial(redisAddr)
	if _, err := client.Ping(context.Background()).Result(); err != nil {
		return &RedisStorage{}, err
	}

	return &RedisStorage{
		Client: client,
	}, nil
}

func (c *RedisStorage) Get(ctx context.Context, key string) (models.RedisJSON, error) {
	value, err := c.Client.Get(ctx, key).Result()
	if err != nil {
		return models.RedisJSON{}, errors.New("no such key")
	}

	valueInt, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return models.RedisJSON{}, err
	}

	return models.RedisJSON{Key: key, Value: valueInt}, nil
}

func (c *RedisStorage) Set(ctx context.Context, u models.RedisJSON) error {
	return c.Client.Set(ctx, u.Key, u.Value, 0).Err()
}

func (c *RedisStorage) Close() error {
	return c.Client.Close()
}
