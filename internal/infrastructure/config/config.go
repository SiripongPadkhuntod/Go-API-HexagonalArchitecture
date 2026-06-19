package config

import "os"

// Config keeps runtime settings for the API.
type Config struct {
	Port        string
	Database    DatabaseConfig
	OutboundAPI OutboundAPIConfig
	Storage     StorageConfig
}

type DatabaseConfig struct {
	Driver   string
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string // คือการตั้งค่าการเชื่อมต่อกับฐานข้อมูล เช่น disable, require, verify-ca, verify-full 
}

type OutboundAPIConfig struct {
	BaseURL string // คือการตั้งค่าการเชื่อมต่อกับ outbound API เช่น http://localhost:8081
}

type StorageConfig struct {
	Endpoint   string
	AccessKey  string
	SecretKey  string
	UseSSL     bool
	BucketName string
}

func Load() Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	return Config{
		Port: port,
		Database: DatabaseConfig{
			Driver:   envOrDefault("DATABASE_DRIVER", "postgres"),
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
		Storage: StorageConfig{
			Endpoint:   envOrDefault("MINIO_ENDPOINT", "localhost:9000"),
			AccessKey:  envOrDefault("MINIO_ACCESS_KEY", "root"),
			SecretKey:  envOrDefault("MINIO_SECRET_KEY", "password123"),
			UseSSL:     os.Getenv("MINIO_USE_SSL") == "true",
			BucketName: envOrDefault("MINIO_BUCKET_NAME", "uploads"),
		},
	}
}

func (c Config) ServerAddress() string {
	return ":" + c.Port
}

func (d DatabaseConfig) PostgresURL() string {
	return "postgres://" + d.User + ":" + d.Password + "@" + d.Host + ":" + d.Port + "/" + d.Name + "?sslmode=" + d.SSLMode
}

func (d DatabaseConfig) MySQLURL() string {
	return d.User + ":" + d.Password + "@tcp(" + d.Host + ":" + d.Port + ")/" + d.Name + "?parseTime=true"
}

func envOrDefault(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}
