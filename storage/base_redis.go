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
	"log"
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
	redisClient        *redis.Client
	redisClusterClient *redis.ClusterClient
	clusterMode        bool
}

// InitRedisService 初始化 Redis 服务
func InitRedisService() (*RedisService, error) {

	redisClusterClient := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:    utils.RedisClusterConfig.Hosts,
		Username: utils.RedisClusterConfig.User,
		Password: utils.RedisClusterConfig.Password,
		PoolSize: utils.RedisClusterConfig.PoolSize,
	})

	_, err := redisClusterClient.Ping(ctx).Result()
	if err != nil {
		log.Println(err)
		log.Println("Failed to connect to Redis cluster. Will try to connect to single node.")
		// If there's any error while connecting to Redis cluster,
		// then try to connect to single Redis node.
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
	} else {
		redisService.clusterMode = true
		redisService.redisClusterClient = redisClusterClient
	}

	return redisService, nil
}

// RedisSet 设置 Redis 键值对
func RedisSet(key string, value interface{}, ttl time.Duration) error {
	if redisService.clusterMode {
		return redisService.redisClusterClient.Set(ctx, key, value, ttl).Err()
	}
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
	var sc *redis.ScanIterator
	if redisService.clusterMode {
		err := redisService.redisClusterClient.ForEachMaster(ctx, func(ctx context.Context, client *redis.Client) error {
			sc = client.Scan(ctx, 0, prefix, 0).Iterator()
			for sc.Next(ctx) {
				keys = append(keys, sc.Val())
			}
			return sc.Err()
		})
		if err != nil {
			return keys, err
		}
	} else {
		sc = redisService.redisClient.Scan(ctx, 0, prefix, 0).Iterator()
	}
	for sc.Next(ctx) {
		keys = append(keys, sc.Val())
	}
	return keys, sc.Err()
}

// RedisGetString 获取 Redis 中的字符串
func RedisGetString(key string) (string, error) {
	var result string
	var err error
	if redisService.clusterMode {
		result, err = redisService.redisClusterClient.Get(ctx, key).Result()
	} else {
		result, err = redisService.redisClient.Get(ctx, key).Result()
	}
	if err == redis.Nil {
		return result, nil
	}
	return result, err
}

// RedisFlushDB 清空 Redis 中的所有键值对
func RedisFlushDB() error {
	if redisService.clusterMode {
		return redisService.redisClusterClient.ForEachMaster(ctx, func(ctx context.Context, client *redis.Client) error {
			return client.FlushDB(ctx).Err()
		})
	}
	return redisService.redisClient.FlushDB(ctx).Err()
}

// RedisDelete 删除 Redis 中的键值对
func RedisDelete(key ...string) error {
	if len(key) > 0 {
		if redisService.clusterMode {
			// Apparently you can NOT delete multiple keys in a single request,since keys are distributed across multiple nodes.
			var err error
			for _, k := range key {
				err = redisService.redisClusterClient.Del(ctx, k).Err()
			}
			return err
		} else {
			return redisService.redisClient.Del(ctx, key...).Err()
		}
	}
	return nil
}
