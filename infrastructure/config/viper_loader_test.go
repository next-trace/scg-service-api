package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/hbttundar/scg-service-base/infrastructure/config"
	"github.com/stretchr/testify/assert"
)

type TestConfig struct {
	AppName string `yaml:"app_name"`
	Server  struct {
		Port int `yaml:"port"`
	} `yaml:"server"`
	Database struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
	} `yaml:"database"`
}

func TestLoad(t *testing.T) {
	// Create a temporary directory for test config files
	tempDir, err := os.MkdirTemp("", "config-test")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create a test config file
	configContent := `
app_name: test-app
server:
  port: 8080
database:
  host: localhost
  port: 5432
  username: testuser
  password: testpass
`
	configPath := filepath.Join(tempDir, "config.yaml")
	err = os.WriteFile(configPath, []byte(configContent), 0644)
	assert.NoError(t, err)

	// Verify the file exists
	_, err = os.Stat(configPath)
	assert.NoError(t, err)

	// Test loading from file
	t.Run("Load from file", func(t *testing.T) {
		var cfg TestConfig
		err := config.Load(tempDir, "config", &cfg)
		assert.NoError(t, err)
		assert.Equal(t, "test-app", cfg.AppName)
		assert.Equal(t, 8080, cfg.Server.Port)
		assert.Equal(t, "localhost", cfg.Database.Host)
		assert.Equal(t, 5432, cfg.Database.Port)
		assert.Equal(t, "testuser", cfg.Database.Username)
		assert.Equal(t, "testpass", cfg.Database.Password)
	})

	// Test environment variable override
	t.Run("Environment variable override", func(t *testing.T) {
		os.Setenv("DATABASE_HOST", "db.example.com")
		os.Setenv("SERVER_PORT", "9090")
		defer func() {
			os.Unsetenv("DATABASE_HOST")
			os.Unsetenv("SERVER_PORT")
		}()

		var cfg TestConfig
		err := config.Load(tempDir, "config", &cfg)
		assert.NoError(t, err)
		assert.Equal(t, "test-app", cfg.AppName)
		assert.Equal(t, 9090, cfg.Server.Port) // Should be overridden by env var
		assert.Equal(t, "db.example.com", cfg.Database.Host) // Should be overridden by env var
	})

	// Test non-existent file
	t.Run("Non-existent file", func(t *testing.T) {
		var cfg TestConfig
		err := config.Load(tempDir, "nonexistent", &cfg)
		assert.NoError(t, err) // Should not error, just use defaults/env vars
		assert.Empty(t, cfg.AppName)
	})

	// Test invalid config path
	t.Run("Invalid config path", func(t *testing.T) {
		var cfg TestConfig
		err := config.Load("/nonexistent/path", "config", &cfg)
		assert.NoError(t, err) // Should not error, just use defaults/env vars
		assert.Empty(t, cfg.AppName)
	})
}
