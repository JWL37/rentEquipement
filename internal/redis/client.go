package redis

import (
	"context"
	"rentEquipement/internal/config"

	"github.com/go-redis/redis/v8"
)

type RedisClient struct {
	Client *redis.Client
}

func NewRedisClient(cfg *config.Config) *RedisClient {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})

	return &RedisClient{
		Client: client,
	}
}

func (rc *RedisClient) Ping(ctx context.Context) error {
	return rc.Client.Ping(ctx).Err()
}
