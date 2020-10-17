package components

import (
	"fmt"
	"strings"
	"time"

	"github.com/bitmaelum/bitmaelum-suite/cmd/bm-client/gui/app"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

const (
	slotSize int = 8
)

type menubarSlot struct {
	Text     string
	Shortcut tcell.Key
	Selected func()
}

type Menubar struct {
	*tview.Box
	Title       string
	DisplayTime bool
	Slots       [10]*menubarSlot
}

func NewMenubar(app *tview.Application) *Menubar {
	m := &Menubar{
		Box:   tview.NewBox(),
	}

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		for _, i := range m.Slots {
			if i == nil {
				continue
			}
			if event.Key() == i.Shortcut {
				if i.Selected != nil {
					i.Selected()
				}
				break
			}
		}

		return event
	})

	// @TODO: We should not always display timer, only when we set displayTimer to true
	go refreshTimer(m)

	return m
}

func refreshTimer(m *Menubar) {
	t := time.NewTicker(500 * time.Millisecond)

	for {
		select {
		case <-t.C:
			app.App.App.QueueUpdateDraw(func() {})
		}
	}
}

func (m *Menubar) SetTitle(title string) *Menubar {
	m.Title = title
	return m
}

func (m *Menubar) SetDisplayTime(b bool) *Menubar {
	m.DisplayTime = b
	return m
}


func (m *Menubar) Draw(screen tcell.Screen) {
	x, y, width, _ := m.GetInnerRect()
	x++

	if m.Title != "" {
		tview.Print(screen, m.Title+" |", x, y, width-2, tview.AlignLeft, tcell.ColorYellow)
		x += len(m.Title) + 3
	}

	for i, slot := range m.Slots {
		kn := tcell.Key(int(tcell.KeyF1) + i)
		tview.Print(screen, tcell.KeyNames[kn], x, y, width-2, tview.AlignLeft, tcell.ColorBlue)
		x += len(tcell.KeyNames[kn]) + 1

		t := "[yellow:blue:b]       [-:-:-]"
		if slot != nil {
			slot.Text += strings.Repeat(" ", slotSize)
			t = fmt.Sprintf("[yellow:blue:b]%s[-:-:-]", slot.Text[:slotSize])
		}
		tview.Print(screen, t, x, y, width, tview.AlignLeft, tcell.ColorYellow)
		x += slotSize + 2
	}

	if m.DisplayTime {
		t := time.Now().Format("15:03:05")
		tview.Print(screen, t, width-8, y, width-2, tview.AlignLeft, tcell.ColorYellow)
	}
}

func (m *Menubar) SetSlot(idx int, text string, selected func()) *Menubar {
	m.Slots[idx] = &menubarSlot{
		Text:     text,
		Shortcut: tcell.Key(int(tcell.KeyF1) + idx),
		Selected: selected,
	}

	return m
}
