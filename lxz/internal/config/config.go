// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of K9s

package config

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v3"
	"io/fs"
	"log/slog"
	"lxz/internal/config/data"
	"lxz/internal/config/json"
	"lxz/internal/slogs"
	"os"
)

// Config tracks K9s configuration options.
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
		slog.Error("Unable to save K9s config file", slogs.Error, err)
		return err
	}

	slog.Info("[CONFIG] Saving K9s config to disk", slogs.Path, path)
	return nil
}

func (c *Config) Merge(c1 *Config) {

}

// Load loads K9s configuration from file.
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
