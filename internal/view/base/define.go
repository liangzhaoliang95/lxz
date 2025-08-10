package base

import "github.com/gdamore/tcell/v2"

var (
	FlexBorderColor     = tcell.ColorNavy        // 边框颜色
	BoarderDefaultColor = tcell.ColorGreenYellow // 默认边框颜色
	ActiveBorderColor   = tcell.ColorGreenYellow // ✅ 获得焦点
	InactiveBorderColor = tcell.ColorWhite       // ❌ 非焦点
)
