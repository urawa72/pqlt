package tui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Tui struct {
	App			*tview.Application
	QueryView	*QueryView
	ResultView	*ResultView
	Info		*InfoView
	Pages		*tview.Pages
	Panels
}

type Panels struct {
	Current	int
	Panels	[]tview.Primitive
}

func New() *Tui {
	// NewClient()
	queryView := NewQueryView()
	resultView := NewResultView()
	infoView := NewInfoView()

	t := &Tui {
		App:		tview.NewApplication(),
		QueryView: 	queryView,
		ResultView:	resultView,
		Info:		infoView,
	}

	t.Panels = Panels{
		Panels: []tview.Primitive{
			queryView,
			resultView,
		},
	}

	return t
}

func (t *Tui) switchPanel(p tview.Primitive) *tview.Application {
	return t.App.SetFocus(p)
}

func (t *Tui) nextPanel() {
	idx := (t.Panels.Current + 1) % len(t.Panels.Panels)
	t.Panels.Current = idx
	t.switchPanel(t.Panels.Panels[t.Panels.Current])
}

func (t *Tui) panelKeybindings(event *tcell.EventKey) {
	switch event.Key() {
	case tcell.KeyTab:
		t.nextPanel()
	}
}

func (t *Tui) queryKeybindings() {
 	t.QueryView.SetDoneFunc(func(key tcell.Key) {
		switch key {
		case tcell.KeyEnter:
			text := t.QueryView.GetText()
			t.QueryView.Query = text
			t.ResultView.UpdateView(t)
		}
	}).SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		t.panelKeybindings(event)
		return event
	})
}

func (t *Tui) resultKeybindings() {
    t.ResultView.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		t.panelKeybindings(event)
		return event
	})
}

func (t *Tui) Run() error {
	grid := tview.NewGrid().
		SetRows(3, 0, 4).
		AddItem(t.QueryView, 0, 0, 1, 1, 0, 0, true).
		AddItem(t.ResultView, 1, 0, 1, 1, 0, 0, true).
		AddItem(t.Info, 2, 0, 1, 1, 0, 0, false)

	t.Pages = tview.NewPages().AddAndSwitchToPage("main", grid, true)

	t.queryKeybindings()
	t.resultKeybindings()

	if err := t.App.SetRoot(t.Pages, true).SetFocus(t.QueryView).Run(); err != nil {
		t.App.Stop()
		return err
	}

	return nil
}
