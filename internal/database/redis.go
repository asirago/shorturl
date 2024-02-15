package database

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

func CreateRedisClient() *redis.Client {

	rdb := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "",
	})

	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("Error Initializing redis: %v\n", err)
	}

	return rdb
}
