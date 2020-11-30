package tui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type QueryView struct {
	*tview.InputField
	Query string
}

func NewQueryView() *QueryView {
	qv := &QueryView{
		InputField: tview.NewInputField(),
	}
	qv.SetTitle("Query").SetTitleAlign(tview.AlignLeft)
	qv.SetBorder(true)
	qv.SetFieldBackgroundColor(tcell.ColorBlack)

	return qv
}
