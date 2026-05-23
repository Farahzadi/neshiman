package config

import "os"

type Config struct {
	DBURL       string
	Port        string
	JWTSecret   string
	CORSOrigins string
}

func Load() *Config {
	return &Config{
		DBURL:       getEnv("DB_URL", "postgres://postgres:postgres@localhost:5432/neshiman?sslmode=disable"),
		Port:        getEnv("PORT", "8080"),
		JWTSecret:   getEnv("JWT_SECRET", "dev-secret-change-in-production"),
		CORSOrigins: getEnv("CORS_ORIGINS", "*"),
	}
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
