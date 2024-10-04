package cache

import (
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"os"
	"strconv"
)

var ctx = context.Background()

func NewRedisDB() (*redis.Client, error) {
	redisAddr := os.Getenv("REDIS_PORT")
	redisPassword := os.Getenv("REDIS_PASSWORD")
	redisDBStr := os.Getenv("REDIS_DB")

	redisDB, err := strconv.Atoi(redisDBStr)
	if err != nil {
		logrus.Fatalf("Invalid REDIS_DB value: %v", err)
	}

	client := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPassword,
		DB:       redisDB,
	})
	_, err = client.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %v", err)
	}

	return client, nil
}
