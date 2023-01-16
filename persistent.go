package main

import "github.com/go-redis/redis/v9"

func NewRedis() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
}
