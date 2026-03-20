package ratelimiter

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"os"
	"time"
)

type Limiter interface {
	LimitByIP(userIP string) (bool, error)
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
		fmt.Printf("Failed to connect to redis server: %s\n", err)
		return nil, false
	}
	return rdb, true
}

func (r *rateLimiter) LimitByIP(ip string) (bool, error) {
	ctx := context.Background()

	key := fmt.Sprintf("ratelimiter:ip:%s", ip)
	count, err := r.Client.Incr(ctx, key).Result()
	if err != nil {
		return false, err
	}

	if count == 1 {
		r.Client.Expire(ctx, key, (time.Duration(r.Window))*time.Second)
	}
	if count > int64(r.Limit) {
		return false, nil
	}

	return true, nil
}
