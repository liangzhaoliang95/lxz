/**
 * @author  zhaoliang.liang
 * @date  2025/8/5 21:21
 */

package ui

import "github.com/rivo/tview"

func IsInputPrimitive(p tview.Primitive) bool {
	switch p.(type) {
	case *tview.InputField, *tview.TextArea, *tview.TextView:
		return true
	default:
		return false
	}
}
