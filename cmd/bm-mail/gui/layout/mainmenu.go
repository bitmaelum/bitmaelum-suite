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
	"fmt"

	"github.com/bitmaelum/bitmaelum-suite/cmd/bm-mail/gui/app"
	"github.com/bitmaelum/bitmaelum-suite/cmd/bm-mail/gui/components"
	"github.com/bitmaelum/bitmaelum-suite/internal"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

const (
	subtitle = `BitMaelum - Read your mail privately from the terminal`
)

var shortcuts = []string{
	"Read (new) BitMaelum mail",
	"Compose a new BitMaelum message",
	"List accounts & Organisations",
	"Configure BitMaelum",
	"Get some help",
	"Exit BitMaelum client",
}

// NewMainMenuScreen creates a new main menu screen
func NewMainMenuScreen(app app.Type) tview.Primitive {

	// Convert our ANSI logo into textview
	logo := tview.NewTextView().
		SetDynamicColors(true)
	w := tview.ANSIWriter(logo)
	_, _ = fmt.Fprint(w, internal.GetASCIILogo())

	// in order to center the logo, we need to set it inside a grid
	logoBox := tview.NewGrid()
	logoBox.SetColumns(0, 51, 0)
	logoBox.AddItem(logo, 0, 1, 1, 1, 10, 10, false)

	// Create a frame for the subtitle
	frame := tview.NewFrame(tview.NewBox()).
		AddText(subtitle, true, tview.AlignCenter, tcell.ColorWhite)

	// List with options
	menu := components.NewMainMenu().
		SetMainTextColor(tcell.ColorTeal).
		SetSelectedTextColor(tcell.ColorYellow).
		SetSelectedBackgroundColor(tcell.ColorTeal)

	for _, s := range shortcuts {
		menu.AddItem(s, rune(0))
	}

	menu.SetSelectedFunc(func(idx int, main string, r rune) {
		switch idx {
		case 0:
			// Display account configuration
			app.Pages.SwitchToPage("message-overview")

		case 1:
			// Display account conssfiguration
			app.Pages.SwitchToPage("message-overview")

		case 2:
			// Display account configuration
			app.Pages.SwitchToPage("accounts")
		case 4:
			// Display help
			app.Pages.SwitchToPage("help")
		case 5:
			// Exit application
			app.App.Stop()
		default:
			modal := tview.NewModal().
				SetText("This is a modal").
				AddButtons([]string{"Ok"}).
				SetDoneFunc(func(buttonIndex int, buttonLabel string) {
					app.Pages.HidePage("modal")
				})

			app.Pages.AddPage("modal", modal, false, true)
		}
	})

	// Create a Flex layout that centers the logo and subtitle.
	grid := tview.NewGrid().SetColumns(0, 70, 0).SetRows(1, 10, 4, 0, 1)

	grid.AddItem(logoBox, 1, 1, 1, 1, 10, 70, false)
	grid.AddItem(frame, 2, 1, 1, 1, 0, 0, false)
	grid.AddItem(menu, 3, 1, 1, 1, 0, 0, true)

	return grid
}
