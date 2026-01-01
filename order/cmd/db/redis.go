package db

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"order/proto"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

var RedisClient *redis.Client

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

func SetCache(ctx context.Context, key string, data any, ttl time.Duration) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	if err := RedisClient.Set(ctx, key, jsonData, ttl).Err(); err != nil {
		return err
	}

	logrus.Infof("Cache set for key: %s (TTL: %v)", key, ttl)
	return nil
}

func GetCacheOrderList(ctx context.Context, key string) ([]*proto.Order, error) {
	data, err := RedisClient.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, errors.New(cachedMiss)
	} else if err != nil {
		return nil, err
	}

	var result []*proto.Order
	if err := json.Unmarshal([]byte(data), &result); err != nil {
		return nil, err
	}

	logrus.Debugf("Cache hit for key: %s", key)

	return result, nil
}

func DeleteCache(ctx context.Context, key string) error {

	if err := RedisClient.Del(ctx, key).Err(); err != nil {
		return err
	}

	return nil
}

func DeleteCacheByPattern(ctx context.Context, pattern string) error {
	var cursor uint64
	var keys []string

	for {
		var err error
		var scanKeys []string

		scanKeys, cursor, err = RedisClient.Scan(ctx, cursor, pattern, 100).Result()
		if err != nil {
			return err
		}
		keys = append(keys, scanKeys...)
		if cursor == 0 {
			break
		}
	}

	if err := RedisClient.Del(ctx, keys...).Err(); err != nil {
		return err
	}

	return nil
}

func RedisOrderKey(userID, page int) string {
	return fmt.Sprintf("orders:user%dpage%d", userID, page)
}
