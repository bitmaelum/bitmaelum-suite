// Copyright (c) 2020 BitMaelum Authors
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of
// this software and associated documentation files (the "Software"), to deal in
// the Software without restriction, including without limitation the rights to
// use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
// the Software, and to permit persons to whom the Software is furnished to do so,
// subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
// FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
// COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
// IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
// CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package components

import (
	"fmt"
	"strings"
	"time"

	"github.com/bitmaelum/bitmaelum-suite/cmd/bm-mail/gui/app"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

const (
	slotSize int = 8
)

type menubarSlot struct {
	Text     string
	Selected func()
}

// Menubar is a structure that holds a menu bar in the bottom
type Menubar struct {
	*tview.Box
	DisplayTime bool
	Slots       [10]*menubarSlot
}

// NewMenubar creates a new menu bar
func NewMenubar(app *tview.Application) *Menubar {
	m := &Menubar{
		Box: tview.NewBox(),
	}

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() >= tcell.KeyF1 && event.Key() <= tcell.KeyF10 {
			idx := event.Key() - tcell.KeyF1

			if m.Slots[idx] != nil && m.Slots[idx].Selected != nil {
				m.Slots[idx].Selected()
			}

			return nil
		}

		return event
	})

	// @TODO: We should not always display timer, only when we set displayTimer to true
	go refreshTimer()

	return m
}

func refreshTimer() {
	t := time.NewTicker(1000 * time.Millisecond)

	for {
		<-t.C
		app.App.App.QueueUpdateDraw(func() {})
	}
}

// SetDisplayTime will display or undisplay the time
func (m *Menubar) SetDisplayTime(b bool) *Menubar {
	m.DisplayTime = b
	return m
}

// Draw will draw the menubar on the screen
func (m *Menubar) Draw(screen tcell.Screen) {
	x, y, width, _ := m.GetInnerRect()
	x++

	if m.GetTitle() != "" {
		tview.Print(screen, m.GetTitle()+" |", x, y, width-2, tview.AlignLeft, tcell.ColorYellow)
		x += len(m.GetTitle()) + 3
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

// SetSlot set the given menubar index with a menu item. It will call the selected() function whenever we select
// the given function from the menubar.
func (m *Menubar) SetSlot(idx int, text string, selected func()) *Menubar {
	m.Slots[idx] = &menubarSlot{
		Text:     text,
		Selected: selected,
	}

	return m
}
