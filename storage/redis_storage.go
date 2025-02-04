package storage

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisStorage struct {
	client *redis.Client
}

func NewRedisStorage(addr, password string, db int) (*RedisStorage, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	// Testar conexão
	ctx := context.Background()
	_, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}

	return &RedisStorage{client: client}, nil
}

func (rs *RedisStorage) Increment(key string) (int64, error) {
	ctx := context.Background()
	countKey := "count:" + key

	// Aumentar TTL para 5 segundos para dar tempo dos testes executarem
	pipe := rs.client.Pipeline()
	incr := pipe.Incr(ctx, countKey)
	pipe.Expire(ctx, countKey, 5*time.Second)
	_, err := pipe.Exec(ctx)

	if err != nil {
		return 0, err
	}

	return incr.Val(), nil
}

func (rs *RedisStorage) Block(key string, seconds int64) error {
	ctx := context.Background()
	blockKey := "blocked:" + key
	return rs.client.Set(ctx, blockKey, "1", time.Duration(seconds)*time.Second).Err()
}

func (rs *RedisStorage) IsBlocked(key string) (bool, error) {
	ctx := context.Background()
	blockKey := "blocked:" + key
	exists, err := rs.client.Exists(ctx, blockKey).Result()
	if err != nil {
		return false, err
	}
	return exists == 1, nil
}

func (rs *RedisStorage) Reset(key string) error {
	ctx := context.Background()
	countKey := "count:" + key
	blockKey := "blocked:" + key

	pipe := rs.client.Pipeline()
	pipe.Del(ctx, countKey)
	pipe.Del(ctx, blockKey)
	_, err := pipe.Exec(ctx)
	return err
}
