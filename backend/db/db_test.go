package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/go-redis/redis/v8"
)

func testRedisConnectivity() error {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default database
	})

	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		return fmt.Errorf("failed to connect to Redis: %v", err)
	}

	return nil
}

func TestRedisConnectivity(t *testing.T) {
	err := testRedisConnectivity()
	if err != nil {
		t.Errorf("Redis connectivity test failed: %v", err)
	}
}
