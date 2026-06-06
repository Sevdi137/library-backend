package config

import (
    "os"
    "github.com/joho/godotenv"
)

type Config struct {
    Port      string
    DBPath    string
    JWTSecret string
}

func Load() (*Config, error) {
    godotenv.Load()
    return &Config{
        Port:      getEnv("PORT", "8080"),
        DBPath:    getEnv("DB_PATH", "./library.db"),
        JWTSecret: getEnv("JWT_SECRET", "defaultsecret"),
    }, nil
}

func getEnv(key, fallback string) string {
    if val := os.Getenv(key); val != "" {
        return val
    }
    return fallback
}
