/**
 * @author  zhaoliang.liang
 * @date  2025/8/5 21:21
 */

package ui

import "github.com/liangzhaoliang95/tview"

func IsInputPrimitive(p tview.Primitive) bool {
	switch p.(type) {
	case *tview.InputField, *tview.TextArea:
		return true
	default:
		return false
	}
}
