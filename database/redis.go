package database

import (
	"blog-go/config"
	"blog-go/log"
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

var (
	RedisClient *redis.Client
	Rctx        = context.Background()
)

func GetRedisClient() (*redis.Client, context.Context) {
	if RedisClient == nil {
		log.Logger.Info("[Redis]redis连接初始化")

		config_ := config.ConfigContext.RedisConfig
		RedisClient = redis.NewClient(&redis.Options{
			Addr:     config_.Host + ":" + config_.Port,
			Password: config_.Password,
			DB:       config_.DB,
		})
		pong, err := RedisClient.Ping(Rctx).Result()
		if err != nil {
			log.Logger.Error("[Redis]redis连接失败", zap.String("error", err.Error()))
		} else {
			log.Logger.Info("[Redis]redis连接成功", zap.String("pong", pong))
		}
	}

	return RedisClient, Rctx
}

func Set(key, value string, expire time.Duration) {
	err := RedisClient.Set(Rctx, key, value, expire).Err()
	if err != nil {
		log.Logger.Error("[Redis]设置redis数据失败", zap.String("error", err.Error()))
	}
	log.Logger.Debug("[Redis]设置redis数据成功", zap.String("key", key), zap.String("value", value))
}
func Get(key string) string {
	result, err := RedisClient.Get(Rctx, key).Result()
	if err != nil {
		log.Logger.Error("[Redis]获取redis数据失败", zap.String("error", err.Error()))
		return ""
	}
	log.Logger.Debug("[Redis]获取redis数据成功", zap.String("key", key), zap.String("value", result))
	return result
}
