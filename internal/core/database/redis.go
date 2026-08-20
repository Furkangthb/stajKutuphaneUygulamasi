package database

import (
	"context"
	"log"
	"os"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

func RedisConnect() {
	addr := os.Getenv("REDIS_ADDR")
	if addr == "" {
		addr = "stajKutuphaneRedis:6379"
	}
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})
	if err := RedisClient.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("Redis bağlantı hatası: %v", err)
	}
	log.Println("Redis bağlantı başarılı.")
}
