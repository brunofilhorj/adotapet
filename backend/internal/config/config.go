package config

import (
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv           string
	HTTPAddr         string
	DatabaseURL      string
	RedisURL         string
	JWTIssuer        string
	JWTAccessSecret  string
	JWTRefreshSecret string
	JWTAccessTTL     time.Duration
	JWTRefreshTTL    time.Duration
	S3Endpoint       string
	S3Bucket         string
	S3AccessKey      string
	S3SecretKey      string
	CDNBaseURL       string
	GeocodingAPIKey  string
	EmailAPIKey      string
}

func Load() Config {
	_ = godotenv.Load()

	return Config{
		AppEnv:           env("APP_ENV", "local"),
		HTTPAddr:         env("HTTP_ADDR", ":8080"),
		DatabaseURL:      env("DATABASE_URL", ""),
		RedisURL:         env("REDIS_URL", "redis://localhost:6379/0"),
		JWTIssuer:        env("JWT_ISSUER", "adotapet"),
		JWTAccessSecret:  env("JWT_ACCESS_SECRET", ""),
		JWTRefreshSecret: env("JWT_REFRESH_SECRET", ""),
		JWTAccessTTL:     durationEnv("JWT_ACCESS_TTL", 15*time.Minute),
		JWTRefreshTTL:    durationEnv("JWT_REFRESH_TTL", 30*24*time.Hour),
		S3Endpoint:       env("S3_ENDPOINT", ""),
		S3Bucket:         env("S3_BUCKET", ""),
		S3AccessKey:      env("S3_ACCESS_KEY", ""),
		S3SecretKey:      env("S3_SECRET_KEY", ""),
		CDNBaseURL:       env("CDN_BASE_URL", ""),
		GeocodingAPIKey:  env("GEOCODING_API_KEY", ""),
		EmailAPIKey:      env("EMAIL_PROVIDER_API_KEY", ""),
	}
}

func env(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

func durationEnv(key string, fallback time.Duration) time.Duration {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	duration, err := time.ParseDuration(value)
	if err != nil {
		return fallback
	}
	return duration
}
