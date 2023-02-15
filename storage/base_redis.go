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
	"time"

	"ohurlshortener/utils"

	"github.com/go-redis/redis/v8"
)

var (
	redisService = &RedisService{}
	ctx          = context.Background()
)

// RedisService Redis 服务
type RedisService struct {
	redisClient *redis.Client
}

// InitRedisService 初始化 Redis 服务
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

// RedisSet 设置 Redis 键值对
func RedisSet(key string, value interface{}, ttl time.Duration) error {
	return redisService.redisClient.Set(ctx, key, value, ttl).Err()
}

// RedisSet30m 设置 Redis 键值对，过期时间为 30 分钟
func RedisSet30m(key string, value interface{}) error {
	return RedisSet(key, value, 30*time.Minute)
}

// RedisSet4Ever 设置 Redis 键值对，永不过期
func RedisSet4Ever(key string, value interface{}) error {
	return RedisSet(key, value, redis.KeepTTL)
}

// RedisScan4Keys 获取 Redis 中所有以 prefix 开头的键
func RedisScan4Keys(prefix string) ([]string, error) {
	var keys []string
	sc := redisService.redisClient.Scan(ctx, 0, prefix, 0).Iterator()
	for sc.Next(ctx) {
		keys = append(keys, sc.Val())
	}
	return keys, nil
}

// RedisGetString 获取 Redis 中的字符串
func RedisGetString(key string) (string, error) {
	result, err := redisService.redisClient.Get(ctx, key).Result()
	if err == redis.Nil {
		return result, nil
	}
	return result, err
}

// RedisFlushDB 清空 Redis 中的所有键值对
func RedisFlushDB() error {
	return redisService.redisClient.FlushDB(ctx).Err()
}

// RedisDelete 删除 Redis 中的键值对
func RedisDelete(key ...string) error {
	if len(key) > 0 {
		return redisService.redisClient.Del(ctx, key...).Err()
	}
	return nil
}
