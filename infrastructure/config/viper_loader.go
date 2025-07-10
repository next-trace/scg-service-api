package config

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	appconfig "github.com/next-trace/scg-service-api/application/config"
	"github.com/spf13/viper"
)

// Load is a convenience function that creates a ViperLoader and calls its Load method.
// It reads configuration from a YAML file and environment variables into the provided struct pointer.
func Load(path, fileName string, configStruct interface{}) error {
	loader := NewViperLoader()
	return loader.Load(path, fileName, configStruct)
}

// ViperLoader implements the config.Loader interface using Viper.
type ViperLoader struct{}

// NewViperLoader creates a new Viper-based configuration loader.
func NewViperLoader() appconfig.Loader {
	return &ViperLoader{}
}

// Load reads configuration from a YAML file and environment variables
// into the provided struct pointer.
func (v *ViperLoader) Load(path, fileName string, configStruct interface{}) error {
	return v.LoadWithOptions(path, fileName, configStruct, appconfig.DefaultOptions())
}

// LoadWithOptions loads configuration with additional options.
func (v *ViperLoader) LoadWithOptions(path, fileName string, configStruct interface{}, options appconfig.Options) error {
	vp := viper.New()

	// Set config path and name
	vp.AddConfigPath(path)

	// Remove file extension if present
	fileName = removeFileExtension(fileName)

	vp.SetConfigName(fileName)
	vp.SetConfigType(options.ConfigType)

	// Configure environment variables
	if options.AllowEnvOverride {
		if options.EnvPrefix != "" {
			vp.SetEnvPrefix(options.EnvPrefix)
		}
		vp.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

		// Explicitly bind common configuration keys to environment variables
		// Helper function to bind environment variables and handle errors
		bindEnv := func(key string) {
			// Explicitly ignore the error as it's not critical
			_ = vp.BindEnv(key)
		}

		bindEnv("app_name")
		bindEnv("server.port")
		bindEnv("database.host")
		bindEnv("database.port")
		bindEnv("database.username")
		bindEnv("database.password")

		vp.AutomaticEnv() // This enables overriding config with env vars
	}

	// Read configuration file
	err := vp.ReadInConfig()
	if err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if errors.As(err, &configFileNotFoundError) {
			if options.RequireConfigFile {
				return err
			}
			// Config file not found but not required, continue with defaults
		} else {
			// Config file found but could not be read
			return err
		}
	}

	// Explicitly set all keys from the config file to ensure they're available
	for _, key := range vp.AllKeys() {
		val := vp.Get(key)
		vp.Set(key, val)
	}

	// Explicitly set the app_name field if it exists in the config file
	if vp.IsSet("app_name") {
		appName := vp.GetString("app_name")
		vp.Set("AppName", appName) // Also set with capitalized field name
	}

	// Unmarshal the config into the provided struct
	if err := vp.Unmarshal(configStruct); err != nil {
		return err
	}

	return nil
}

// Reload reloads configuration from the source.
func (v *ViperLoader) Reload(path, fileName string, configStruct interface{}, options appconfig.Options) error {
	return v.LoadWithOptions(path, fileName, configStruct, options)
}

// FileExists checks if a configuration file exists.
func (v *ViperLoader) FileExists(path, fileName, fileType string) bool {
	fileName = removeFileExtension(fileName)
	filePath := filepath.Join(path, fileName+"."+fileType)
	_, err := os.Stat(filePath)
	return err == nil
}

// Helper function to remove file extensions
func removeFileExtension(fileName string) string {
	extensions := []string{".yaml", ".yml", ".json", ".toml", ".ini"}
	for _, ext := range extensions {
		if strings.HasSuffix(fileName, ext) {
			return strings.TrimSuffix(fileName, ext)
		}
	}
	return fileName
}
