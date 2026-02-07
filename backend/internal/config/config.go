package config

import "os"

type PostgresConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
	SSLMode  string
}

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

type Config struct {
	Port       string
	Env        string
	Kubeconfig string
	Postgres   PostgresConfig
	Redis      RedisConfig
}

func Load() Config {
	return Config{
		Port:       GetEnv("PORT", "8080"),
		Env:        GetEnv("ENV", ""),
		Kubeconfig: GetEnv("KUBECONFIG", ""),
		Postgres: PostgresConfig{
			Host:     GetEnv("POSTGRES_HOST", ""),
			Port:     GetEnv("POSTGRES_PORT", ""),
			User:     GetEnv("POSTGRES_USER", ""),
			Password: GetEnv("POSTGRES_PASSWORD", ""),
			Database: GetEnv("POSTGRES_DB", ""),
			SSLMode:  GetEnv("POSTGRES_SSLMODE", ""),
		},
		Redis: RedisConfig{
			Addr:     GetEnv("REDIS_ADDR", ""),
			Password: GetEnv("REDIS_PASSWORD", ""),
			DB:       GetEnvInt("REDIS_DB", 0),
		},
	}
}

func GetEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func GetEnvInt(key string, fallback int) int {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	n := 0
	for _, c := range value {
		if c < '0' || c > '9' {
			return fallback
		}
		n = n*10 + int(c-'0')
	}
	return n
}
