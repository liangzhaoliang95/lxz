/**
 * @author  zhaoliang.liang
 * @date  2025/8/4 16:02
 */

package dialog

import (
	"github.com/gdamore/tcell/v2"
	"github.com/liangzhaoliang95/lxz/internal/ui"
	"github.com/liangzhaoliang95/tview"
)

type ForceDrawFn func() *tview.Application

type LoadingDialog struct {
	Pages        *ui.Pages
	Message      string
	modalLoading *tview.ModalLoading
	forceDrawFn  ForceDrawFn
}

func (_this *LoadingDialog) Hide() {
	dismissConfirm(_this.Pages)
	_this.forceDrawFn()
}

func ShowLoadingDialog(pages *ui.Pages, Message string, forceDrawFn ForceDrawFn) *LoadingDialog {

	l := &LoadingDialog{
		Pages:        pages,
		modalLoading: tview.NewModalLoading(""),
		Message:      Message,
		forceDrawFn:  forceDrawFn,
	}

	if l.Message == "" {
		l.Message = "⏳ Loading..."
	}
	l.modalLoading.SetText(l.Message)
	l.modalLoading.SetTextColor(tcell.ColorRed)
	l.modalLoading.SetBorder(true)
	l.modalLoading.SetBackgroundColor(tcell.ColorGrey)

	l.Pages.AddPage(confirmKey, l.modalLoading, false, false)
	l.Pages.ShowPage(confirmKey)
	l.forceDrawFn()
	return l
}
