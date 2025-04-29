package config

import (
	"log"

	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
)

type Config struct {
	Port string `env:"PORT" default:"8082"`
	DSN  string `env:"DATABASE_URL"`

	// Redis configuration
	RedisAddr     string `env:"REDIS_ADDR" default:"redis:6379"`
	RedisPassword string `env:"REDIS_PASSWORD" default:""`
	RedisDB       int    `env:"REDIS_DB" default:"0"`
	TokenTTL      int    `env:"TOKEN_TTL" default:"600"`   // 10 minutes in seconds
	SessionTTL    int    `env:"SESSION_TTL" default:"300"` // 5 minutes in seconds
}

func ParseConfig(cfg *Config) {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file")
	}
	err := env.Parse(cfg)
	if err != nil {
		log.Fatalf("unable to parse ennvironment variables: %e", err)
	}
}
