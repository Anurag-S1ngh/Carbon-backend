package redis

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisConfig struct {
	client *redis.Client
	logger *slog.Logger
}

func NewRedisConfig(redisURL string, logger *slog.Logger) (*RedisConfig, error) {
	client, err := RedisClient(redisURL)
	if err != nil {
		return nil, err
	}
	return &RedisConfig{
		client: client,
		logger: logger,
	}, nil
}

func RedisClient(redisURL string) (*redis.Client, error) {
	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, err
	}

	return redis.NewClient(opts), nil
}

func (r *RedisConfig) SetEx(key, value string, expiresInSeconds uint) error {
	ctx := context.Background()
	_, err := r.client.SetEx(ctx, key, value, time.Second*time.Duration(expiresInSeconds)).Result()
	if err != nil {
		errMsg := fmt.Sprintf("error while setting value of key:%s to redis", key)
		r.logger.Error(errMsg, "error", err)
		return err
	}

	return nil
}

func (r *RedisConfig) Get(key string) (string, error) {
	ctx := context.Background()
	value, err := r.client.Get(ctx, key).Result()
	if err != nil {
		errMsg := fmt.Sprintf("error while getting value of key:%s from redis", key)
		r.logger.Error(errMsg, "error", err)
		return "", err
	}

	return value, nil
}
