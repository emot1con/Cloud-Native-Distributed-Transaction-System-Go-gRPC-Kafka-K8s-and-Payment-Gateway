package middleware

import (
	"context"
	_ "embed"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

//go:embed rate_limiter.lua
var luaScript string

var RedisClient *redis.Client

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

// RateLimiterConfig holds rate limiter configuration
type RateLimiterConfig struct {
	Capacity   int     // Maximum number of tokens (max requests)
	RefillRate float64 // Tokens refilled per second
}

func DefaultRateLimiterConfig() *RateLimiterConfig {
	return &RateLimiterConfig{
		Capacity:   20,
		RefillRate: 10.0,
	}
}

func StrictRateLimiterConfig() *RateLimiterConfig {
	return &RateLimiterConfig{
		Capacity:   5,
		RefillRate: 1.0,
	}
}

// RateLimiterMiddleware implements Token Bucket algorithm using Redis + Lua script
func RateLimiterMiddleware(config *RateLimiterConfig) gin.HandlerFunc {
	// Load Lua script into Redis (SHA hash will be cached)
	ctx := context.Background()
	scriptSHA, err := RedisClient.ScriptLoad(ctx, luaScript).Result()
	if err != nil {
		logrus.WithError(err).Fatal("Failed to load Lua script for rate limiter")
	}

	return func(c *gin.Context) {
		// Get user ID from JWT context
		userID, ok := c.Request.Context().Value(UserKey).(int)
		if !ok {
			c.JSON(401, gin.H{"error": "User ID not found"})
			return
		}

		key := UserRateLimiterKey(userID)
		now := time.Now().Unix()
		result, err := RedisClient.EvalSha(ctx, scriptSHA, []string{key},
			config.Capacity,
			config.RefillRate,
			now,
		).Result()

		if err != nil {
			logrus.WithError(err).Error("Failed to execute rate limiter Lua script")
			// Fail open: allow request if Redis fails
			c.Next()
			return
		}

		// Check if request is allowed
		allowed := result.(int64)
		if allowed == 0 {
			// Rate limit exceeded
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":       "Rate limit exceeded",
				"message":     fmt.Sprintf("Maximum %d requests per second allowed", int(config.RefillRate)),
				"retry_after": fmt.Sprintf("%.1f seconds", 1.0/config.RefillRate),
			})
			logrus.Infof("Rate-Limit Exceed")
			c.Abort()
			return
		}

		logrus.Infof("Rate Limit Pass")

		// Request allowed, continue
		c.Next()
	}
}

// Build cache key for user rate limiting
func UserRateLimiterKey(userID int) string {
	return fmt.Sprintf("rate_limiter:user:%d", userID)
}
