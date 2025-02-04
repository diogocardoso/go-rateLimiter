package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	RedisAddr      string
	RedisPassword  string
	RedisDB        int
	RateLimitIP    int
	RateLimitToken int
	BlockTime      int
	Port           string
}

func LoadConfig() Config {
	godotenv.Load()

	return Config{
		RedisAddr:      getEnv("REDIS_ADDR", "localhost:6379"),
		RedisPassword:  getEnv("REDIS_PASSWORD", ""),
		RedisDB:        getEnvAsInt("REDIS_DB", 0),
		RateLimitIP:    getEnvAsInt("RATE_LIMIT_IP", 100),
		RateLimitToken: getEnvAsInt("RATE_LIMIT_TOKEN", 200),
		BlockTime:      getEnvAsInt("BLOCK_TIME", 60),
		Port:           getEnv("PORT", "8080"),
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func getEnvAsInt(key string, defaultValue int) int {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	intValue, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return intValue
}
