// Copyright (c) 2022 BitMaelum Authors
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
	"github.com/bitmaelum/bitmaelum-suite/cmd/bm-mail/app"
	"github.com/bitmaelum/bitmaelum-suite/cmd/bm-mail/gui/components"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// AccountList is the list of all accounts
var AccountList *tview.List

// OrganisationList is the list of all organisations
var OrganisationList *tview.List

// NewAccountScreen creats a new account screen
func NewAccountScreen() tview.Primitive {
	AccountList = tview.NewList().ShowSecondaryText(false)
	AccountList.SetBorder(true).SetTitle("Accounts")

	OrganisationList = tview.NewList().ShowSecondaryText(false)
	OrganisationList.SetBorder(true).SetTitle("Organisations")

	menuBar := components.NewMenubar(app.MailApp.App)
	menuBar.SetSlot(0, "New Acc", func() {})
	menuBar.SetSlot(1, "New Org", func() {})
	menuBar.SetSlot(2, "Bar", func() {})
	menuBar.SetSlot(9, "Back", func() {
		app.MailApp.Pages.SwitchToPage("main_menu")
	})

	grid := tview.NewGrid().SetColumns(0, 0).SetRows(0, 1)
	grid.AddItem(AccountList, 0, 0, 1, 1, 0, 0, true)
	grid.AddItem(OrganisationList, 0, 1, 1, 1, 0, 0, false)
	grid.AddItem(menuBar, 1, 0, 1, 2, 0, 0, true)

	curActiveElement := 0
	elements := []tview.Primitive{AccountList, OrganisationList}

	grid.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyTab {
			curActiveElement++
			if curActiveElement >= len(elements) {
				curActiveElement = 0
			}
			p := elements[curActiveElement]

			app.MailApp.App.SetFocus(p)
			return nil
		}

		return event
	})

	return grid
}
