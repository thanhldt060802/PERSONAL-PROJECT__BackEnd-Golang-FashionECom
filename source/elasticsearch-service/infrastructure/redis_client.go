package infrastructure

import (
	"context"
	"fmt"
	"log"
	"thanhldt060802/config"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

func InitRedisClient() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", config.AppConfig.RedisHost, config.AppConfig.RedisPort),
		Password: config.AppConfig.RedisPassword,
		DB:       0,
	})

	if _, err := RedisClient.Ping(context.Background()).Result(); err != nil {
		log.Fatal("Ping to Redis failed: ", err)
	}

	log.Println("Connect to Redis successful")
}
