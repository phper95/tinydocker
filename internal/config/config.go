package config

import "os"

type Config struct {
	Server ServerConfig
}
type ServerConfig struct {
	Host string
	Port string
	Mode string // debug, release, test
}

func Load() *Config {
	return &Config{
		Server: ServerConfig{
			Host: getEnv("SERVER_HOST", "0.0.0.0"),
			Port: getEnv("SERVER_PORT", "80"),
			Mode: getEnv("GIN_MODE", "debug"),
		},
	}
}
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
