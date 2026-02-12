package config

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

func loadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}
}

func NewRedisClient() (*redis.Client, error) {
	loadEnv()

	cfg := RedisConfig{
		Host:     os.Getenv("REDIS_HOST"),
		Port:     os.Getenv("REDIS_PORT"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	}

	client := redis.NewClient(&redis.Options{
		Addr:      fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Password:  cfg.Password,
		DB:        cfg.DB,
		TLSConfig: &tls.Config{},
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return client, nil
}
