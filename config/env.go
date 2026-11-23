package config

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort        string
	DatabaseURL    string
	MongoURI       string
	MongoDB        string
	RedisAddr      string
	JWTSecret      string
	JWTAccessTTL   time.Duration
	JWTRefreshTTL  time.Duration
	SwaggerEnabled bool
}

var Env *Config

func LoadEnv() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found, relying on system env")
	}

	accessTTL, err := time.ParseDuration(getEnv("JWT_ACCESS_TTL", "15m"))
	if err != nil {
		log.Fatalf("invalid JWT_ACCESS_TTL: %v", err)
	}

	refreshTTL, err := time.ParseDuration(getEnv("JWT_REFRESH_TTL", "720h"))
	if err != nil {
		log.Fatalf("invalid JWT_REFRESH_TTL: %v", err)
	}

	Env = &Config{
		AppPort:        getEnv("APP_PORT", "8080"),
		DatabaseURL:    getEnv("DATABASE_URL", ""),
		MongoURI:       getEnv("MONGO_URI", ""),
		MongoDB:        getEnv("MONGO_DB", "sewakos"),
		RedisAddr:      getEnv("REDIS_ADDR", "localhost:6379"),
		JWTSecret:      getEnv("JWT_SECRET", ""),
		JWTAccessTTL:   accessTTL,
		JWTRefreshTTL:  refreshTTL,
		SwaggerEnabled: getEnv("SWAGGER_ENABLED", "false") == "true",
	}

	if Env.DatabaseURL == "" || Env.JWTSecret == "" {
		log.Fatal("DATABASE_URL and JWT_SECRET must be set")
	}

	return Env
}

func getEnv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return fallback
}
