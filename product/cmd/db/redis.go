package db

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"product_service/proto"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

var RedisClient *redis.Client

type CachedProductList struct {
	Products     []*proto.Product
	TotalProduct int
	TotalPage    int
}

var cachedMiss = "cache miss"

// InitRedis initializes Redis connection
func InitRedis() error {
	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		redisURL = "localhost:6379"
	}

	RedisClient = redis.NewClient(&redis.Options{
		Addr:         redisURL,
		Password:     "", // no password for now
		DB:           0,  // default DB
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		PoolSize:     10,
		MinIdleConns: 5,
	})

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := RedisClient.Ping(ctx).Err(); err != nil {
		logrus.Errorf("Failed to connect to Redis: %v", err)
		return err
	}

	logrus.Info("Successfully connected to Redis")
	return nil
}

// SetCache sets a value in Redis with TTL
func SetCache(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	jsonData, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %v", err)
	}

	if err := RedisClient.Set(ctx, key, jsonData, ttl).Err(); err != nil {
		logrus.Warnf("Failed to set cache for key %s: %v", key, err)
		return err
	}

	logrus.Debugf("Cache set for key: %s (TTL: %v)", key, ttl)
	return nil
}

// GetCacheProduct gets a value from Redis and unmarshals it
func GetCacheProduct(ctx context.Context, key string) (*proto.Product, error) {
	val, err := RedisClient.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, fmt.Errorf(cachedMiss)
	} else if err != nil {
		logrus.Warnf("Failed to get cache for key %s: %v", key, err)
		return nil, err
	}

	var result *proto.Product
	if err := json.Unmarshal([]byte(val), &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal cache data: %v", err)
	}

	logrus.Debugf("Cache hit for key: %s", key)
	return result, nil
}

func GetCachedProductList(ctx context.Context, key string) (*CachedProductList, error) {
	val, err := RedisClient.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, errors.New(cachedMiss)
	}

	var result *CachedProductList
	if err := json.Unmarshal([]byte(val), &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal cache data: %v", err)
	}

	logrus.Debugf("Cache hit for key: %s", key)
	return result, nil
}

// DeleteCache deletes a key from Redis
func DeleteCache(ctx context.Context, key string) error {
	if err := RedisClient.Del(ctx, key).Err(); err != nil {
		logrus.Warnf("Failed to delete cache for key %s: %v", key, err)
		return err
	}

	logrus.Debugf("Cache deleted for key: %s", key)
	return nil
}

// DeleteCachePattern deletes all keys matching a pattern
func DeleteCachePattern(ctx context.Context, pattern string) error {
	var cursor uint64
	var keys []string

	// Scan for keys matching pattern
	for {
		var err error
		var scanKeys []string

		scanKeys, cursor, err = RedisClient.Scan(ctx, cursor, pattern, 100).Result()
		if err != nil {
			return fmt.Errorf("failed to scan keys: %v", err)
		}

		keys = append(keys, scanKeys...)

		if cursor == 0 {
			break
		}
	}

	// Delete all matching keys
	if len(keys) > 0 {
		if err := RedisClient.Del(ctx, keys...).Err(); err != nil {
			return fmt.Errorf("failed to delete keys: %v", err)
		}
		logrus.Debugf("Deleted %d cache keys matching pattern: %s", len(keys), pattern)
	}

	return nil
}

// GetCacheStats returns basic cache statistics
func GetCacheStats(ctx context.Context) map[string]interface{} {
	info := RedisClient.Info(ctx, "stats").Val()
	return map[string]interface{}{
		"info": info,
	}
}
