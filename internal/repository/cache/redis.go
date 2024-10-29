package cache

import (
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"golang.org/x/net/context"
	"strconv"
)

type Redis struct {
	cache *redis.Client
}

func (r *Redis) GetString(key string) (string, error) {
	result, err := r.cache.Get(context.Background(), key).Result()
	if err != nil {
		return "", err
	}

	return result, nil
}

func (r *Redis) GetInt(key string) (int, error) {
	strRes, err := r.GetString(key)
	res, err := strconv.Atoi(strRes)
	if err != nil {
		return 0, err
	}

	return res, nil
}

var ctx = context.Background()

func NewRedisCache(cache *redis.Client) *Redis {
	return &Redis{cache: cache}
}

func NewRedisDB() (*redis.Client, error) {
	redisHost := viper.GetString("REDIS_HOST")
	redisPort := viper.GetString("REDIS_PORT")
	redisPassword := viper.GetString("REDIS_PASSWORD")
	redisDBStr := viper.GetString("REDIS_DB")

	if redisDBStr == "" {
		redisDBStr = "0"
	}

	redisDB, err := strconv.Atoi(redisDBStr)
	if err != nil {
		return nil, fmt.Errorf("invalid REDIS_DB value: %v", err)
	}

	redisAddr := fmt.Sprintf("%s:%s", redisHost, redisPort)
	fmt.Println("REDIS_HOST:", redisAddr)
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
