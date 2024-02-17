package main

import (
	"context"
	"os"
	"testing"

	"github.com/asirago/shorturl/internal/database"
)

func TestMain(m *testing.M) {
	if os.Getenv("INTEGRATION") != "" {
		rdb := database.CreateRedisClient()

		exit := m.Run()
		rdb.FlushAll(context.Background()).Result()
		os.Exit(exit)

	}
}
