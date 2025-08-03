package config

import (
	"errors"
	"fmt"
	"lxz/internal/helper"

	"gopkg.in/yaml.v3"
	"io/fs"
	"log/slog"
	"lxz/internal/config/data"
	"lxz/internal/slogs"
	"os"
)

const (
	DatabaseProviderMySQL = "MySQL"
)

var DatabaseProviderList = []string{
	DatabaseProviderMySQL,
}

type DBConnection struct {
	Name      string   `yaml:"name" json:"name"`
	URL       string   `yaml:"url" json:"url"`
	Provider  string   `yaml:"provider" json:"provider"`
	UserName  string   `yaml:"username" json:"username"`
	Password  string   `yaml:"password" json:"password"`
	Host      string   `yaml:"host" json:"host"`
	Port      int64    `yaml:"port" json:"port"`
	DBName    string   `yaml:"dbname" json:"dbname"`
	URLParams string   `yaml:"urlParams" json:"urlParams"`
	Commands  []string `yaml:"commands" json:"commands"`
}

func (d *DBConnection) GetUniqKey() string {
	return fmt.Sprintf("%s@%s", d.Provider, d.Name)
}

type DatabaseConfig struct {
	DefaultPageSize int             `yaml:"defaultPageSize" json:"defaultPageSize"`
	DBConnections   []*DBConnection `yaml:"dbConnections" json:"dbConnections"`
}

// String()
func (c *DatabaseConfig) String() string {
	return helper.Prettify(c)
}

// Save database configuration to disk.
func (c *DatabaseConfig) Save(force bool) error {
	if _, err := os.Stat(AppDatabaseConfigFile); errors.Is(err, fs.ErrNotExist) {
		return c.SaveFile(AppDatabaseConfigFile)
	}
	if force {
		return c.SaveFile(AppDatabaseConfigFile)
	}

	return nil
}

// SaveFile lxz database configuration to disk.
func (c *DatabaseConfig) SaveFile(path string) error {
	if err := data.EnsureDirPath(path, data.DefaultDirMod); err != nil {
		return err
	}

	if err := data.SaveYAML(path, c); err != nil {
		slog.Error("Unable to save LXZ database_config file", slogs.Error, err)
		return err
	}

	slog.Info("[CONFIG] Saving LXZ database_config to disk", slogs.Path, path)
	return nil
}

func (c *DatabaseConfig) Merge(c1 *DatabaseConfig) {
	if c1.DefaultPageSize != 0 {
		c.DefaultPageSize = c1.DefaultPageSize
	}
	if len(c1.DBConnections) == 0 {
		slog.Info("[CONFIG] No database connections found in config, using default connection")
		//c.DBConnections = append(c.DBConnections, &DBConnection{})
	} else {
		c.DBConnections = c1.DBConnections
	}
}

// Load loads LXZ database configuration from file.
func (c *DatabaseConfig) Load(path string, force bool) error {
	if _, err := os.Stat(path); errors.Is(err, fs.ErrNotExist) {
		if err = c.Save(force); err != nil {
			return err
		}
	}
	bb, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	var errs error

	var cfg DatabaseConfig
	if err = yaml.Unmarshal(bb, &cfg); err != nil {
		errs = errors.Join(errs, fmt.Errorf("database_config load failed: %w", err))
	}
	c.Merge(&cfg)

	return errs
}

func NewDatabaseConfig() *DatabaseConfig {
	return &DatabaseConfig{
		DefaultPageSize: 300,
	}
}
