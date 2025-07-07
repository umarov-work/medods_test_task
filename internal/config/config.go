package config

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"
)

type Config struct {
	DbDsn          string
	AccessTokenTTL time.Duration
	JWTSecret      []byte
	WebHook        string
}

var (
	cfg  *Config
	once sync.Once
)

func Load() *Config {
	once.Do(func() {
		ttlStr := getEnv("ACCESS_TOKEN_TTL")
		ttl, err := time.ParseDuration(ttlStr)
		if err != nil {
			panic(fmt.Sprintf("Failed to load config: ACCESS_TOKEN_TTL is incorrect: %s", ttlStr))
		}

		jwtSecret := getEnv("JWT_SECRET")

		cfg = &Config{
			DbDsn:          getEnv("DB_DSN"),
			AccessTokenTTL: ttl,
			JWTSecret:      []byte(jwtSecret),
			WebHook:        getEnv("WEBHOOK"),
		}
	})
	return cfg
}

func getEnv(key string) string {
	value := os.Getenv(key)
	if strings.TrimSpace(value) == "" {
		panic(fmt.Sprintf("Environment variable %s not set", key))
	}
	return value
}
