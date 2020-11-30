package tui

import "github.com/rivo/tview"

type InfoView struct {
	*tview.TextView
}

func NewInfoView() *InfoView {
	info := &InfoView{
		TextView: tview.NewTextView().SetTextAlign(tview.AlignLeft).SetDynamicColors(true),
	}
	info.SetTitleAlign(tview.AlignLeft)
	navi := "[yellow::b]Enter[white]: Execute query\n[yellow::b]Tab[white]: Switch panel\n[yellow::b]Ctrl-c[white]: Exit"
	info.SetText(navi)
	return info
}
