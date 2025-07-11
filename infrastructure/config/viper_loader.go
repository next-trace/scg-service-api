package config

import (
	"os"
	"path/filepath"
	"strings"

	appconfig "github.com/hbttundar/scg-service-base/application/config"
	"github.com/spf13/viper"
)

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
	vp.AddConfigPath(path)

	// Remove file extension if present
	fileName = removeFileExtension(fileName)

	vp.SetConfigName(fileName)
	vp.SetConfigType(options.ConfigType)

	// Configure environment variables
	if options.AllowEnvOverride {
		vp.AutomaticEnv()
		if options.EnvPrefix != "" {
			vp.SetEnvPrefix(options.EnvPrefix)
		}
		vp.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	}

	// Read configuration file
	if err := vp.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			if options.RequireConfigFile {
				return err
			}
		} else {
			return err
		}
	}

	return vp.Unmarshal(configStruct)
}

// Reload reloads configuration from the source.
func (v *ViperLoader) Reload(path, fileName string, configStruct interface{}, options appconfig.Options) error {
	return v.LoadWithOptions(path, fileName, configStruct, options)
}

// FileExists checks if a configuration file exists.
func (v *ViperLoader) FileExists(path, fileName string, fileType string) bool {
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
