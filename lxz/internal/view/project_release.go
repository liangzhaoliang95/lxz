/**
 * @author  zhaoliang.liang
 * @date  2025/7/30 17:35
 */

package view

import (
	"github.com/rivo/tview"
	"lxz/internal/ui"
)

// ProjectRelease 项目Release视图
type ProjectRelease struct {
	*tview.Flex
	actions    *ui.KeyActions
	name       string
	fullScreen bool
}
