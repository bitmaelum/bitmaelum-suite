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

package layout

import (
	"github.com/bitmaelum/bitmaelum-suite/cmd/bm-mail/gui/app"
	"github.com/bitmaelum/bitmaelum-suite/cmd/bm-mail/gui/components"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

// NewAccountScreen creats a new account screen
func NewAccountScreen(app app.Type) tview.Primitive {

	list := tview.NewList().ShowSecondaryText(false)
	list.SetBorder(true).SetTitle("Accounts")
	for _, acc := range app.Vault.Store.Accounts {
		list.AddItem(acc.Name+" <"+acc.Address.String()+">", "", rune(0), nil)
	}

	list2 := tview.NewList().ShowSecondaryText(false)
	list2.SetBorder(true).SetTitle("Organisations")
	for _, org := range app.Vault.Store.Organisations {
		list2.AddItem(org.FullName+" <...@"+org.Addr+">", "", rune(0), nil)
	}

	menuBar := components.NewMenubar(app.App)
	menuBar.SetSlot(0, "New Acc", func() {})
	menuBar.SetSlot(1, "New Org", func() {})
	menuBar.SetSlot(2, "Bar", func() {})
	menuBar.SetSlot(9, "Back", func() {
		app.Pages.SwitchToPage("main_menu")
	})

	grid := tview.NewGrid().SetColumns(0, 0).SetRows(0, 1)
	grid.AddItem(list, 0, 0, 1, 1, 0, 0, true)
	grid.AddItem(list2, 0, 1, 1, 1, 0, 0, false)
	grid.AddItem(menuBar, 1, 0, 1, 2, 0, 0, true)

	curActiveElement := 0
	elements := []tview.Primitive{list, list2}

	grid.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyTab {
			curActiveElement++
			if curActiveElement >= len(elements) {
				curActiveElement = 0
			}
			p := elements[curActiveElement]

			app.App.SetFocus(p)
			return nil
		}

		return event
	})

	return grid
}
