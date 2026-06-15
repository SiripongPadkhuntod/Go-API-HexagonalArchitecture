package config

import "os"

// Config keeps runtime settings for the API.
type Config struct {
	Port        string
	Database    DatabaseConfig
	OutboundAPI OutboundAPIConfig
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

type OutboundAPIConfig struct {
	BaseURL string
}

func Load() Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	return Config{
		Port: port,
		Database: DatabaseConfig{
			Host:     envOrDefault("DATABASE_HOST", "localhost"),
			Port:     envOrDefault("DATABASE_PORT", "5439"),
			User:     envOrDefault("DATABASE_USER", "postgres"),
			Password: envOrDefault("DATABASE_PASSWORD", "postgres"),
			Name:     envOrDefault("DATABASE_NAME", "hexagonal_architecture"),
			SSLMode:  envOrDefault("DATABASE_SSLMODE", "disable"),
		},
		OutboundAPI: OutboundAPIConfig{
			BaseURL: os.Getenv("OUTBOUND_API_BASE_URL"),
		},
	}
}

func (c Config) ServerAddress() string {
	return ":" + c.Port
}

func (d DatabaseConfig) URL() string {
	return "postgres://" + d.User + ":" + d.Password + "@" + d.Host + ":" + d.Port + "/" + d.Name + "?sslmode=" + d.SSLMode
}

func envOrDefault(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}
