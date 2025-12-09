package version

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// VersionInfo 表示版本信息
type VersionInfo struct {
	Version      string    `json:"version"`
	Commit       string    `json:"commit"`
	Date         string    `json:"date"`
	BuildTime    time.Time `json:"build_time"`
	GoVersion    string    `json:"go_version"`
	Platform     string    `json:"platform"`
	Architecture string    `json:"architecture"`
}

// GitHubRelease 表示GitHub release信息
type GitHubRelease struct {
	TagName     string    `json:"tag_name"`
	Name        string    `json:"name"`
	Body        string    `json:"body"`
	PublishedAt time.Time `json:"published_at"`
	Assets      []Asset   `json:"assets"`
}

// Asset 表示release中的资源文件
type Asset struct {
	Name               string `json:"name"`
	BrowserDownloadURL string `json:"browser_download_url"`
	Size               int64  `json:"size"`
	DownloadCount      int64  `json:"download_count"`
}

// 版本信息变量
var (
	// 这些变量将在编译时通过ldflags注入
	Version   = "dev"
	Commit    = "dev"
	Date      = "unknown"
	BuildTime = "unknown"
	GoVersion = "unknown"
	Platform  = "unknown"
	Arch      = "unknown"
)

// GetVersion 返回当前版本信息
func GetVersion() *VersionInfo {
	buildTime, _ := time.Parse(time.RFC3339, BuildTime)
	return &VersionInfo{
		Version:      Version,
		Commit:       Commit,
		Date:         Date,
		BuildTime:    buildTime,
		GoVersion:    GoVersion,
		Platform:     Platform,
		Architecture: Arch,
	}
}

// IsDevVersion 检查是否为开发版本
func IsDevVersion() bool {
	return Version == "dev" || Version == "unknown"
}

// IsLatestVersion 检查是否为最新版本
func IsLatestVersion() bool {
	if IsDevVersion() {
		return false
	}

	latest, err := GetLatestGitHubRelease()
	if err != nil {
		return false
	}

	return Version == latest.TagName
}

// GetLatestGitHubRelease 从GitHub获取最新的release信息
func GetLatestGitHubRelease() (*GitHubRelease, error) {
	// 从环境变量或配置中获取仓库信息
	repo := getRepositoryInfo()
	if repo == "" {
		// 如果无法获取仓库信息，使用默认仓库
		repo = "liangzhaoliang95/lxz"
	}

	url := fmt.Sprintf("https://api.github.com/repos/%s/releases/latest", repo)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("请求GitHub API失败: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API返回错误状态码: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	var release GitHubRelease
	if err := json.Unmarshal(body, &release); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	slog.Info("GitHub Release", "release", release)

	return &release, nil
}

// CheckForUpdates 检查是否有新版本可用
func CheckForUpdates() (*UpdateInfo, error) {
	if IsDevVersion() {
		return nil, fmt.Errorf("开发版本无法检查更新")
	}

	latest, err := GetLatestGitHubRelease()
	if err != nil {
		return nil, err
	}

	current := parseVersion(Version)
	latestVersion := parseVersion(latest.TagName)

	if latestVersion == nil || current == nil {
		return nil, fmt.Errorf("版本号格式错误")
	}

	slog.Info("CheckForUpdates", "current", current, "latest", latestVersion)

	if latestVersion.GreaterThan(current) {
		return &UpdateInfo{
			CurrentVersion: Version,
			LatestVersion:  latest.TagName,
			ReleaseNotes:   latest.Body,
			DownloadURL:    getDownloadURL(latest, Platform, Arch),
			PublishedAt:    latest.PublishedAt,
		}, nil
	}

	return nil, nil
}

// UpdateInfo 表示更新信息
type UpdateInfo struct {
	CurrentVersion string    `json:"current_version"`
	LatestVersion  string    `json:"latest_version"`
	ReleaseNotes   string    `json:"release_notes"`
	DownloadURL    string    `json:"download_url"`
	PublishedAt    time.Time `json:"published_at"`
}

// SemVer 表示语义化版本
type SemVer struct {
	Major int
	Minor int
	Patch int
}

// parseVersion 解析版本号
func parseVersion(version string) *SemVer {
	// 移除v前缀
	version = strings.TrimPrefix(version, "v")

	parts := strings.Split(version, ".")
	if len(parts) != 3 {
		return nil
	}

	major, err := strconv.Atoi(parts[0])
	if err != nil {
		return nil
	}

	minor, err := strconv.Atoi(parts[1])
	if err != nil {
		return nil
	}

	patch, err := strconv.Atoi(parts[2])
	if err != nil {
		return nil
	}

	return &SemVer{Major: major, Minor: minor, Patch: patch}
}

// GreaterThan 比较版本号
func (v *SemVer) GreaterThan(other *SemVer) bool {
	if v.Major > other.Major {
		return true
	}
	if v.Major < other.Major {
		return false
	}

	if v.Minor > other.Minor {
		return true
	}
	if v.Minor < other.Minor {
		return false
	}

	return v.Patch > other.Patch
}

// getRepositoryInfo 获取仓库信息
func getRepositoryInfo() string {
	return GetRepositoryInfo()
}

// getDownloadURL 根据平台和架构获取下载URL
func getDownloadURL(release *GitHubRelease, platform, arch string) string {
	// 构建期望的文件名
	var expectedName string
	switch platform {
	case "linux":
		switch arch {
		case "amd64":
			expectedName = "lxz-linux-amd64"
		case "arm64":
			expectedName = "lxz-linux-arm64"
		}
	case "darwin":
		switch arch {
		case "amd64":
			expectedName = "lxz-darwin-amd64"
		case "arm64":
			expectedName = "lxz-darwin-arm64"
		}
	case "windows":
		switch arch {
		case "amd64":
			expectedName = "lxz-windows-amd64.exe"
		case "arm64":
			expectedName = "lxz-windows-arm64.exe"
		}
	}

	// 查找匹配的资源文件
	for _, asset := range release.Assets {
		if asset.Name == expectedName {
			return asset.BrowserDownloadURL
		}
	}

	return ""
}

// FormatVersion 格式化版本信息显示
func FormatVersion(short bool) string {
	v := GetVersion()

	if short {
		return v.Version
	}

	var result strings.Builder
	result.WriteString(fmt.Sprintf("Version: %s\n", v.Version))
	result.WriteString(fmt.Sprintf("Commit: %s\n", v.Commit))
	result.WriteString(fmt.Sprintf("Date: %s\n", v.Date))
	result.WriteString(fmt.Sprintf("Go Version: %s\n", v.GoVersion))
	result.WriteString(fmt.Sprintf("Platform: %s/%s", v.Platform, v.Architecture))

	return result.String()
}
