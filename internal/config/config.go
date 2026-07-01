package config

import (
	"fmt"
	"os"
)

type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string
	Env        string
	Port       string
}

func Load() Config {
	return Config{
		DBHost:     GetEnv("DB_HOST", "localhost"),
		DBPort:     GetEnv("DB_PORT", "5432"),
		DBUser:     GetEnv("DB_USER", "postgres"),
		DBPassword: GetEnv("DB_PASSWORD", "password"),
		DBName:     GetEnv("DB_NAME", "mydb"),
		DBSSLMode:  GetEnv("DB_SSLMODE", "disable"),
		Env:        GetEnv("ENV", "development"),
		Port:       GetEnv("PORT", "8080"),
	}
}

func (c Config) DSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBName, c.DBSSLMode,
	)
}

func (c Config) ServerAddr() string {
	return ":" + c.Port
}

func GetEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
