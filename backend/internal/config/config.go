package config

import "os"

type Config struct {
	AppEnv           string
	HTTPAddr         string
	DatabaseURL      string
	RedisURL         string
	JWTIssuer        string
	JWTAccessSecret  string
	JWTRefreshSecret string
	S3Endpoint       string
	S3Bucket         string
	S3AccessKey      string
	S3SecretKey      string
	CDNBaseURL       string
	GeocodingAPIKey  string
	EmailAPIKey      string
}

func Load() Config {
	return Config{
		AppEnv:           env("APP_ENV", "local"),
		HTTPAddr:         env("HTTP_ADDR", ":8080"),
		DatabaseURL:      env("DATABASE_URL", "postgres://adotapet:adotapet@localhost:5432/adotapet?sslmode=disable"),
		RedisURL:         env("REDIS_URL", "redis://localhost:6379/0"),
		JWTIssuer:        env("JWT_ISSUER", "adotapet"),
		JWTAccessSecret:  env("JWT_ACCESS_SECRET", "change-me-access"),
		JWTRefreshSecret: env("JWT_REFRESH_SECRET", ""),
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
