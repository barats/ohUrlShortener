package redis

import (
	"context"
	"ohurlshortener/utils"
	"time"

	"github.com/go-redis/redis/v8"
)

var (
	redisService = &RedisService{}
	ctx          = context.Background()
)

type RedisService struct {
	redisClient *redis.Client
}

func InitRedisService() (*RedisService, error) {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     utils.RedisConfig.Host,
		DB:       utils.RedisConfig.Database,
		Username: utils.RedisConfig.User,
		Password: utils.RedisConfig.Password,
	})

	_, err := redisClient.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}

	redisService.redisClient = redisClient

	return redisService, nil
}

func Set(key string, value interface{}, ttl time.Duration) error {
	return redisService.redisClient.Set(ctx, key, value, ttl).Err()
}

func Set30m(key string, value interface{}) error {
	return Set(key, value, 30*time.Minute)
}

//Set4Ever Needs redis version 6.0 or above
func Set4Ever(key string, value interface{}) error {
	return Set(key, value, redis.KeepTTL)
}

func GetString(key string) (string, error) {
	return redisService.redisClient.Get(ctx, key).Result()
}
