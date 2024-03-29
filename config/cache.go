package config

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

func InitCache() (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB}
	})

	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		return nil, fmt.Errorf("error connecting to redis")
	}
	return client, nil
}
