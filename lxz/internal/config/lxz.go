/**
 * @author  zhaoliang.liang
 * @date  2025/7/23 13:31
 */

package config

type LXZ struct {
	RefreshRate   int    `json:"refreshRate"   yaml:"refreshRate"`
	ScreenDumpDir string `json:"screenDumpDir" yaml:"screenDumpDir,omitempty"`
	Logger        Logger `json:"logger"        yaml:"logger"`
	UI            UI     `json:"ui"            yaml:"ui"`
}

// NewLXZ create a new LXZ configuration.
func NewLXZ() *LXZ {
	return &LXZ{
		RefreshRate:   defaultRefreshRate, // 刷新率
		ScreenDumpDir: AppDumpsDir,        // 截图目录
		Logger:        NewLogger(),        // 日志配置
		UI:            UI{},               // UI配置
	}
}

// IsHeadless returns headless setting.
func (k *LXZ) IsHeadless() bool {
	if IsBoolSet(k.UI.manualHeadless) {
		return true
	}

	return k.UI.Headless
}

// IsSplashless returns splashless setting.
func (k *LXZ) IsSplashless() bool {
	if IsBoolSet(k.UI.manualSplashless) {
		return true
	}

	return k.UI.Splashless
}

// Override overrides lxz config from cli args.
func (k *LXZ) Override(lxzFlags *Flags) {
	// 可以使用将命令行配置覆盖到k上
}
