package config

import (
	"os"
)

type Config struct {
	DBDriver   string
	DSN        string
	JWTSecret  string
	AdminEmail string
}

func Load() Config {
	driver := getEnv("DB_DRIVER", "sqlite")
	dsn := getEnv("DB_DSN", "file:epl.db?cache=shared&_journal_mode=WAL")
	secret := getEnv("JWT_SECRET", "dev-secret-change")
	adminEmail := getEnv("ADMIN_EMAIL", "admin@epl.local")
	return Config{
		DBDriver:   driver,
		DSN:        dsn,
		JWTSecret:  secret,
		AdminEmail: adminEmail,
	}
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
