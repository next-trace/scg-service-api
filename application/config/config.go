// Package config defines the abstract interface (PORT) for configuration management.
package config

// Loader defines the abstract interface for loading configuration from various sources.
type Loader interface {
	// Load loads configuration from the specified path and file name into the provided struct.
	Load(path, fileName string, configStruct interface{}) error

	// LoadWithOptions loads configuration with additional options.
	LoadWithOptions(path, fileName string, configStruct interface{}, options Options) error

	// Reload reloads configuration from the source.
	Reload(path, fileName string, configStruct interface{}, options Options) error

	// FileExists checks if a configuration file exists.
	FileExists(path, fileName string, fileType string) bool
}

// Options defines configuration loading options.
type Options struct {
	// ConfigType specifies the configuration file type (yaml, json, toml, etc.)
	ConfigType string

	// EnvPrefix is the prefix for environment variables.
	EnvPrefix string

	// AllowEnvOverride allows environment variables to override file settings.
	AllowEnvOverride bool

	// RequireConfigFile requires a config file to exist (error if not found).
	RequireConfigFile bool
}

// DefaultOptions returns the default configuration options.
func DefaultOptions() Options {
	return Options{
		ConfigType:        "yaml",
		EnvPrefix:         "",
		AllowEnvOverride:  true,
		RequireConfigFile: false,
	}
}