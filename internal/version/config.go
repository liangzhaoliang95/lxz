package version

import (
	"os"
	"strings"
)

// Config 版本配置
type Config struct {
	Repository string `json:"repository" yaml:"repository"`
	AutoCheck  bool   `json:"auto_check" yaml:"auto_check"`
	CheckURL   string `json:"check_url"  yaml:"check_url"`
}

// DefaultConfig 默认配置
var DefaultConfig = &Config{
	Repository: "liangzhaoliang95/lxz",
	AutoCheck:  true,
	CheckURL:   "https://api.github.com",
}

// LoadConfig 加载配置
func LoadConfig() *Config {
	config := &Config{}

	// 从环境变量读取
	if repo := os.Getenv("LXZ_REPOSITORY"); repo != "" {
		config.Repository = repo
	}

	if autoCheck := os.Getenv("LXZ_AUTO_CHECK"); autoCheck != "" {
		config.AutoCheck = strings.ToLower(autoCheck) == "true"
	}

	if checkURL := os.Getenv("LXZ_CHECK_URL"); checkURL != "" {
		config.CheckURL = checkURL
	}

	// 如果没有设置，使用默认值
	if config.Repository == "" {
		config.Repository = DefaultConfig.Repository
	}

	if config.CheckURL == "" {
		config.CheckURL = DefaultConfig.CheckURL
	}

	return config
}

// GetRepositoryInfo 获取仓库信息
func GetRepositoryInfo() string {
	config := LoadConfig()
	return config.Repository
}

// GetCheckURL 获取检查URL
func GetCheckURL() string {
	config := LoadConfig()
	return config.CheckURL
}

// IsAutoCheckEnabled 是否启用自动检查
func IsAutoCheckEnabled() bool {
	config := LoadConfig()
	return config.AutoCheck
}
