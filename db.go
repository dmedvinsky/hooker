package main

import (
	"gopkg.in/redis.v5"
)

var Redis *redis.Client

func ConnectDB() {
	Redis = redis.NewClient(&redis.Options{
		Addr:     Config.RedisAddr,
		DB:       Config.RedisDB,
	})
}
