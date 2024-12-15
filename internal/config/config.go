package config

import (
	"log"

	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
)

type Config struct {
	Port string `env:"PORT" default:"8082"`
	DSN  string `env:"DATABASE_URL"`
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
