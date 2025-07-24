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
		RefreshRate: defaultRefreshRate,

		ScreenDumpDir: AppDumpsDir,
		Logger:        NewLogger(),
		UI:            UI{},
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
