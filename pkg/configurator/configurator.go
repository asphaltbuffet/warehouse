package configurator

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/lmittmann/tint"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
)

const (
	WherehouseEnvPrefix   string     = "WHEREHOUSE"
	DefaultConfigFileBase string     = "wherehouse"
	DefaultConfigExt      string     = "toml"
	DefaultLogLevel       slog.Level = slog.LevelInfo
)

type Config struct {
	viper       *viper.Viper
	dir         string
	cfgFile     string
	cfgFileType string
	logger      *slog.Logger
	fs          afero.Fs
}

func New(options ...func(*Config)) (Config, error) {
	w := os.Stderr
	cfg := Config{
		viper: viper.New(),
		logger: slog.New(
			tint.NewHandler(w, &tint.Options{
				Level:      slog.LevelInfo,
				TimeFormat: time.StampMilli,
			}),
		),
	}
	slog.SetDefault(cfg.logger)

	for _, option := range options {
		option(&cfg)
	}

	// set up fs
	if cfg.fs == nil {
		cfg.fs = afero.NewOsFs()
	}
	cfg.viper.SetFs(cfg.fs)

	// set up viper
	cfg.setViperConfigFile()

	cfg.viper.SetEnvPrefix(WherehouseEnvPrefix)

	for k, v := range defaults {
		cfg.viper.SetDefault(string(k), v)
	}

	configDir, err := os.UserConfigDir()
	if err != nil {
		cfg.logger.Error("get default user config dir", "error", tint.Err(err))
		return Config{}, err
	}
	cfg.viper.SetDefault(string(ConfigDirKey), filepath.Join(configDir, "wherehouse"))

	// set config sources
	if cfg.dir != "" {
		cfg.viper.AddConfigPath(cfg.dir)
	}
	cfg.viper.AddConfigPath(".")

	userCfg, err := os.UserConfigDir()
	if err == nil {
		cfg.viper.AddConfigPath(filepath.Join(userCfg, "wherehouse"))
	}

	err = cfg.viper.ReadInConfig()
	if err != nil {
		if !errors.As(err, &viper.ConfigFileNotFoundError{}) {
			// only return error if it's not a missing config file
			cfg.logger.Error("failed to read config file", "error", err, "config", cfg.cfgFile)
			return Config{}, err
		}

		cfg.logger.Warn("no config file found", slog.String("file", cfg.cfgFile), tint.Err(err))
	} else {
		cfg.logger.Debug("starting with config file", "config", cfg.viper.ConfigFileUsed())
	}

	return cfg, nil
}

// WithFile sets the configuration file and type.
//
// If the file is empty, the default file name and type are used.
func WithFile(f string) func(*Config) {
	return func(c *Config) {
		file := filepath.Base(f)
		ext := filepath.Ext(f)
		c.dir = filepath.Dir(f)

		// handle dotfiles
		// foo 		-> "foo" + "" 			(false)
		// foo.bar 	-> "foo.bar" + ".bar" 	(false)
		// .foo.bar -> ".foo.bar" + ".bar" 	(false)
		// .foo.foo -> ".foo.foo" + ".foo" 	(false)
		// .foo 	-> ".foo" + ".foo" 		(true)
		// "" 		-> "." + "" 			(false)
		if file == ext {
			ext = ""
		}

		// remove leading dot from extension
		ext = strings.TrimPrefix(ext, ".")

		switch {
		// filepath.Base returns "." for empty path
		case file == ".":
			// no filename; use defaults
			c.cfgFile = fmt.Sprintf("%s.%s", DefaultConfigFileBase, DefaultConfigExt)
			c.cfgFileType = DefaultConfigExt
		case file != "." && ext == "":
			// filename without extension; use default extension
			c.cfgFile = file
			c.cfgFileType = DefaultConfigExt

		case file != "." && ext != "":
			// filename with extension; set type as well
			c.cfgFile = file
			c.cfgFileType = ext

		default:
		}

		c.logger.Debug("config file set", "path", c.dir, "file", c.cfgFile, "type", c.cfgFileType)
	}
}

func (c *Config) setViperConfigFile() {
	if c.viper == nil {
		panic("viper not initialized")
	}

	if c.cfgFile == "" {
		c.cfgFile = DefaultConfigFileBase + "." + DefaultConfigExt
	}

	c.viper.SetConfigName(c.cfgFile)

	if c.cfgFileType != "" {
		c.viper.SetConfigType(c.cfgFileType)
	}

	c.viper.AddConfigPath(c.dir)
}

// WithFs sets the file system.
//
// If the file system is nil, a new OS file system is used.
func WithFs(fs afero.Fs) func(*Config) {
	return func(c *Config) {
		c.fs = fs
	}
}

// GetConfigFileUsed returns the configuration file used.
//
// If no configuration file is loaded, an empty string is returned. Failure to read a
// configuration file does not cause an error and will still result in an empty string.
func (c Config) GetConfigFileUsed() string {
	return c.viper.ConfigFileUsed()
}

func (c Config) GetConfigDir() string {
	return c.viper.GetString(string(ConfigDirKey))
}

func (c Config) GetLogger() *slog.Logger {
	return c.logger
}

func (c Config) GetFs() afero.Fs {
	return c.fs
}
