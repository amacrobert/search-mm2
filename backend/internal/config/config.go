package config

import (
	"os"
	"time"
)

type Config struct {
	DatabaseURL    string
	JWTSecret      string
	Port           string
	ScrapeInterval time.Duration
	AdminUsername   string
	AdminPassword  string
}

func Load() *Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	intervalStr := os.Getenv("SCRAPE_INTERVAL")
	if intervalStr == "" {
		intervalStr = "30m"
	}
	interval, err := time.ParseDuration(intervalStr)
	if err != nil {
		interval = 30 * time.Minute
	}

	adminUser := os.Getenv("ADMIN_USERNAME")
	if adminUser == "" {
		adminUser = "admin"
	}

	adminPass := os.Getenv("ADMIN_PASSWORD")
	if adminPass == "" {
		adminPass = "admin"
	}

	return &Config{
		DatabaseURL:    os.Getenv("DATABASE_URL"),
		JWTSecret:      os.Getenv("JWT_SECRET"),
		Port:           port,
		ScrapeInterval: interval,
		AdminUsername:   adminUser,
		AdminPassword:  adminPass,
	}
}
