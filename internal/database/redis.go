package database

import (
	"context"
	"log"
	"os"

	"github.com/redis/go-redis/v9"
)

var Ctx = context.Background()

func CreateRedisClient() *redis.Client {
	addr := "redis:6379"

	if os.Getenv("INTEGRATION") != "" {
		addr = "localhost:6379"
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "",
	})

	_, err := rdb.Ping(Ctx).Result()
	if err != nil {
		log.Fatalf("Error Initializing redis: %v\n", err)
	}

	return rdb
}
