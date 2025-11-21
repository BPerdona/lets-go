package config

import (
	"os"

	"github.com/joho/godotenv"
)

const defaultPort = "8080"

// ServerConfig represents the configuration for the server
type ServerConfig struct {
	ListenAddr string
}

// NewServerConfig creates a new ServerConfig
func NewServerConfig() *ServerConfig {
	_ = godotenv.Load()

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = defaultPort
	}

	return &ServerConfig{
		ListenAddr: ":" + port,
	}
}
