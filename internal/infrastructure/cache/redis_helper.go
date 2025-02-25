package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type Redis struct {
	redisClient *redis.Client
	rctx        context.Context
}

var Rdb *Redis

func NewRedis(host, port, username, password string, db int) *Redis {
	return &Redis{
		rctx: context.Background(),
		redisClient: redis.NewClient(&redis.Options{
			Addr:     host + ":" + port,
			Password: password,
			DB:       db,
		}),
	}
}

func Setup(host, port, username, password string, db int) {
	Rdb = NewRedis(host, port, username, password, db)
}

func GetRedis() (*Redis, error) {
	if Rdb == nil {
		return nil, fmt.Errorf("redis未初始化")
	}
	return Rdb, nil
}

func (r *Redis) Set(key, value string, expire time.Duration) error {
	return r.redisClient.Set(r.rctx, key, value, expire).Err()
}
func (r *Redis) Get(key string) (string, error) {
	return r.redisClient.Get(r.rctx, key).Result()
}

func (r *Redis) Del(key string) error {
	return r.redisClient.Del(r.rctx, key).Err()
}
