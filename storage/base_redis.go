// Copyright (c) [2022] [巴拉迪维 BaratSemet]
// [ohUrlShortener] is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
// 				 http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package storage

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
		PoolSize: utils.RedisConfig.PoolSize,
	})
	_, err := redisClient.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}
	redisService.redisClient = redisClient
	return redisService, nil
}

func RedisSet(key string, value interface{}, ttl time.Duration) error {
	return redisService.redisClient.Set(ctx, key, value, ttl).Err()
}

func RedisSet30m(key string, value interface{}) error {
	return RedisSet(key, value, 30*time.Minute)
}

func RedisSet4Ever(key string, value interface{}) error {
	return RedisSet(key, value, redis.KeepTTL)
}

func RedisScan4Keys(prefix string) ([]string, error) {
	keys := []string{}
	sc := redisService.redisClient.Scan(ctx, 0, prefix, 0).Iterator()
	for sc.Next(ctx) {
		keys = append(keys, sc.Val())
	}
	return keys, nil
}

func RedisGetString(key string) (string, error) {
	result, err := redisService.redisClient.Get(ctx, key).Result()
	if err == redis.Nil {
		return result, nil
	}
	return result, err
}

func RedisFlushDB() error {
	return redisService.redisClient.FlushDB(ctx).Err()
}

func RedisDelete(key ...string) error {
	if len(key) > 0 {
		return redisService.redisClient.Del(ctx, key...).Err()
	}
	return nil
}
