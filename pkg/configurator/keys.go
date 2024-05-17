package configurator

// ConfigKey is a type for configuration keys.
type ConfigKey string

const (
	// General configuration keys.
	ConfigDirKey        ConfigKey = "config-dir"    // Configuration key for application configuration files.
	DatabaseFilenameKey ConfigKey = "database.file" // Configuration key for database filename
)

func (k ConfigKey) String() string {
	return string(k)
}
