package config

import (
	"errors"
	"fmt"
	"github.com/liangzhaoliang95/lxz/internal/helper"

	"github.com/liangzhaoliang95/lxz/internal/config/data"
	"github.com/liangzhaoliang95/lxz/internal/slogs"
	"gopkg.in/yaml.v3"
	"io/fs"
	"log/slog"
	"os"
)

type RedisConnConfig struct {
	Name     string `yaml:"name"     json:"name"`
	UserName string `yaml:"username" json:"username"`
	Password string `yaml:"password" json:"password"`
	Host     string `yaml:"host"     json:"host"`
	Port     int64  `yaml:"port"     json:"port"`
}

type RedisConfig struct {
	RedisConnConfig []*RedisConnConfig `yaml:"redisConnConfig" json:"redisConnConfig"`
}

// String()
func (c *RedisConfig) String() string {
	return helper.Prettify(c)
}

// Save database configuration to disk.
func (c *RedisConfig) Save(force bool) error {
	if _, err := os.Stat(AppRedisConfigFile); errors.Is(err, fs.ErrNotExist) {
		return c.SaveFile(AppRedisConfigFile)
	}
	if force {
		return c.SaveFile(AppRedisConfigFile)
	}

	return nil
}

// SaveFile lxz database configuration to disk.
func (c *RedisConfig) SaveFile(path string) error {
	if err := data.EnsureDirPath(path, data.DefaultDirMod); err != nil {
		return err
	}

	if err := data.SaveYAML(path, c); err != nil {
		slog.Error("Unable to save LXZ redis config file", slogs.Error, err)
		return err
	}

	slog.Info("[CONFIG] Saving LXZ redis config to disk", slogs.Path, path)
	return nil
}

func (c *RedisConfig) Merge(fileRead *RedisConfig) {

	if len(fileRead.RedisConnConfig) == 0 {
		slog.Info("[CONFIG] No redis connections found in config, using default connection")
		//c.DBConnections = append(c.DBConnections, &DBConnection{})
	} else {
		c.RedisConnConfig = fileRead.RedisConnConfig
	}
}

// Load loads LXZ redis configuration from file.
func (c *RedisConfig) Load(path string, force bool) error {
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

	var cfg RedisConfig
	if err = yaml.Unmarshal(bb, &cfg); err != nil {
		errs = errors.Join(errs, fmt.Errorf("redis config load failed: %w", err))
	}
	c.Merge(&cfg)

	return errs
}

func NewRedisConfig() *RedisConfig {
	return &RedisConfig{}
}
