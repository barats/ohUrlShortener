package redis

import (
	"context"
	"ohurlshortener/utils"
	"time"

	oredis "github.com/go-redis/redis/v8"
)

var (
	redisService = &RedisService{}
	ctx          = context.Background()
)

type RedisService struct {
	redisClient *oredis.Client
}

func InitRedisService() (*RedisService, error) {
	redisClient := oredis.NewClient(&oredis.Options{
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
	return Set(key, value, oredis.KeepTTL)
}

func GetString(key string) (string, error) {
	result, err := redisService.redisClient.Get(ctx, key).Result()
	if err == oredis.Nil {
		return result, nil
	}
	return result, err
}

func FlushDB() error {
	return redisService.redisClient.FlushDB(ctx).Err()
}

func PushToList(key string, value interface{}) error {
	return redisService.redisClient.LPush(ctx, key, value).Err()
	//	return redisService.redisClient.Expire(ctx, key, oredis.KeepTTL).Err()
}

func GetAllFromList(key string) ([]string, error) {
	return redisService.redisClient.LRange(ctx, key, 0, -1).Result()
}

func Expire(key string) error {
	return redisService.redisClient.Expire(ctx, key, 1*time.Second).Err()
}
