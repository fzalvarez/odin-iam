package config

import "os"

type Config struct {
	Port        string
	DatabaseURL string
	JwtSecret   string
}

func Load() (*Config, error) {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	cfg := &Config{
		Port:        port,
		DatabaseURL: os.Getenv("DATABASE_URL"),
		JwtSecret:   os.Getenv("JWT_SECRET"),
	}

	return cfg, nil
}
