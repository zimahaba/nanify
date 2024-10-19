package main

import (
	"context"

	"github.com/redis/go-redis/v9"
)

func createCache() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     "redis-service:6379",
		Password: "",
		DB:       0,
	})

	if err := client.Ping(context.TODO()).Err(); err != nil {
		panic(err)
	}

	return client
}
