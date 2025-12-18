package ratelimit

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RateLimiter struct {
	client *redis.Client
	limit  int           // requests per window
	window time.Duration // time window
}

func NewRateLimiter(client *redis.Client, limit int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		client: client,
		limit:  limit,
		window: window,
	}
}

// Allow checks if request from IP is allowed
// Uses Fixed Window algorithm
func (rl *RateLimiter) Allow(ctx context.Context, ip string) (bool, error) {
	key := fmt.Sprintf("rate_limit:%s", ip)

	// Increment counter
	count, err := rl.client.Incr(ctx, key).Result()
	if err != nil {
		return false, err
	}

	// Set expiry on first request
	if count == 1 {
		rl.client.Expire(ctx, key, rl.window)
	}

	return count <= int64(rl.limit), nil
}

// Remaining returns remaining requests in current window
func (rl *RateLimiter) Remaining(ctx context.Context, ip string) (int, error) {
	key := fmt.Sprintf("rate_limit:%s", ip)
	count, err := rl.client.Get(ctx, key).Int()
	if err != nil {
		if err == redis.Nil {
			return rl.limit, nil
		}
		return 0, err
	}

	remaining := rl.limit - count
	if remaining < 0 {
		remaining = 0
	}
	return remaining, nil
}
