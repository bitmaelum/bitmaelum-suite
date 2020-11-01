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
	"io/ioutil"
	"net/http"

	"github.com/bitmaelum/bitmaelum-suite/cmd/bm-mail/gui/app"
	"github.com/bitmaelum/bitmaelum-suite/cmd/bm-mail/gui/components"
	"github.com/rivo/tview"
)

// NewHelpScreen creates a new help screen
func NewHelpScreen(app app.Type) tview.Primitive {

	// Load text into help window
	res, _ := http.Get("http://www.gutenberg.org/cache/epub/17192/pg17192.txt")
	b, _ := ioutil.ReadAll(res.Body)
	defer func() { _ = res.Body.Close() }()

	help := tview.NewTextView()
	_, _ = help.Write(b)

	help.SetBorder(true).
		SetBackgroundColor(tview.Styles.ContrastBackgroundColor).
		SetBorderPadding(1, 1, 1, 1)
	help.SetTitle(" BitMaelum Help ")

	// Create menu bar
	menuBar := components.NewMenubar(app.App)
	menuBar.SetSlot(0, "fooo", func() {
		app.Pages.SwitchToPage("accounts")
	})
	menuBar.SetSlot(9, "Back", func() {
		app.Pages.SwitchToPage("main_menu")
	})

	// Create a Flex layout that centers the logo and subtitle.
	grid := tview.NewGrid().SetColumns(10, 0, 10).SetRows(1, 0, 2, 1)
	grid.AddItem(help, 1, 1, 1, 1, 10, 70, true)
	grid.AddItem(menuBar, 3, 0, 1, 3, 0, 0, true)

	return grid
}
