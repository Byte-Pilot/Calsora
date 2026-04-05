package ratelimiter

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"log"
	"os"
	"time"
)

type Limiter interface {
	LimitByIP(userIP string, limit, window int) (bool, error)
	//LimitByUserID(userID int) (bool, error)
}

type rateLimiter struct {
	Client *redis.Client
	Limit  int
	Window int
}

func NewRateLimiter(client *redis.Client) *rateLimiter {
	return &rateLimiter{
		Client: client,
	}
}

func ConnectRedis() (*redis.Client, bool) {
	addr := os.Getenv("REDIS_ADDR")
	if addr == "" {
		addr = "localhost:6379"
	}

	rdb := redis.NewClient(&redis.Options{
		Addr: addr,
	})

	pingCtx, pingCancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer pingCancel()

	if err := rdb.Ping(pingCtx).Err(); err != nil {
		log.Fatalf("Failed to connect to redis server: %s\n", err)
		return nil, false
	}
	return rdb, true
}

func (r *rateLimiter) LimitByIP(userIP string, limit, window int) (bool, error) {
	ctx := context.Background()

	key := fmt.Sprintf("ratelimiter:ip:%s", userIP)
	count, err := r.Client.Incr(ctx, key).Result()
	if err != nil {
		return false, err
	}

	if count == 1 {
		r.Client.Expire(ctx, key, (time.Duration(window))*time.Second)
	}
	if count > int64(limit) {
		return false, nil
	}

	return true, nil
}
