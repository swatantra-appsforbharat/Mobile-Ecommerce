package database

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

var RDB *redis.Client
var Ctx = context.Background()

func InitRedis() {
	RDB = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	// Test connection
	_, err := RDB.Ping(Ctx).Result()
	if err != nil {
		panic(fmt.Sprintf("❌ Failed to connect to Redis: %v", err))
	}

	fmt.Println("✅ Connected to Redis")
}
