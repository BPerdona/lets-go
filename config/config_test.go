package config

import (
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	// Load .env if it exists
	_ = godotenv.Load("../.env")

	// Define default values for tests if they don't exist
	if os.Getenv("SERVER_PORT") == "" {
		os.Setenv("SERVER_PORT", "8081")
	}

	// Run tests
	code := m.Run()

	// Cleanup
	os.Exit(code)
}

func TestGetEnv_ServerPort(t *testing.T) {
	port := os.Getenv("SERVER_PORT")
	assert.NotEmpty(t, port)
}

func TestNewServerConfig_ListenAddr(t *testing.T) {
	port := os.Getenv("SERVER_PORT")
	config := NewServerConfig()
	assert.Equal(t, config.ListenAddr, ":"+port)
}
