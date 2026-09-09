package config

import (
	"os"
)

type Config struct {
	Port        string
	JWTSecret   string
	NodeAPIURL  string
}

func LoadConfig() *Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "interseguro_super_secret_key_2024"
	}

	nodeAPIURL := os.Getenv("NODE_API_URL")
	if nodeAPIURL == "" {
		nodeAPIURL = "http://api-node:4000"
	}

	return &Config{
		Port:       port,
		JWTSecret:  jwtSecret,
		NodeAPIURL: nodeAPIURL,
	}
}
