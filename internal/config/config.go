package config

import (
	"errors"
	"fmt"
	"github.com/liangzhaoliang95/lxz/internal/config/data"
	"github.com/liangzhaoliang95/lxz/internal/config/json"
	"github.com/liangzhaoliang95/lxz/internal/helper"
	"github.com/liangzhaoliang95/lxz/internal/slogs"
	"gopkg.in/yaml.v3"
	"io/fs"
	"log/slog"
	"os"
)

// Config tracks LXZ configuration options.
type Config struct {
	LXZ *LXZ `yaml:"lxz" json:"lxz"`
}

// NewConfig creates a new default config.
func NewConfig() *Config {
	return &Config{
		LXZ: NewLXZ(),
	}
}

// Save configuration to disk.
func (c *Config) Save(force bool) error {
	if _, err := os.Stat(AppConfigFile); errors.Is(err, fs.ErrNotExist) {
		return c.SaveFile(AppConfigFile)
	}

	return nil
}

// SaveFile lxz configuration to disk.
func (c *Config) SaveFile(path string) error {
	if err := data.EnsureDirPath(path, data.DefaultDirMod); err != nil {
		return err
	}

	if err := data.SaveYAML(path, c); err != nil {
		slog.Error("Unable to save LXZ config file", slogs.Error, err)
		return err
	}

	slog.Info("[CONFIG] Saving LXZ config to disk", slogs.Path, path)
	return nil
}

func (c *Config) Merge(c1 *Config) {

}

// Load loads LXZ configuration from file.
func (c *Config) Load(path string, force bool) error {
	if _, err := os.Stat(path); errors.Is(err, fs.ErrNotExist) {
		if err := c.Save(force); err != nil {
			return err
		}
	}
	bb, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	var errs error
	if err := data.JSONValidator.Validate(json.LXZSchema, bb); err != nil {
		errs = errors.Join(errs, fmt.Errorf("lxz config file %q load failed:\n%w", path, err))
	}

	var cfg Config
	if err := yaml.Unmarshal(bb, &cfg); err != nil {
		errs = errors.Join(errs, fmt.Errorf("main config.yaml load failed: %w", err))
	}
	c.Merge(&cfg)

	return errs
}

// Validate the configuration.
func (c *Config) Validate(contextName, clusterName string) {

}

// string
func (c *Config) String() string {
	return helper.Prettify(c)
}
