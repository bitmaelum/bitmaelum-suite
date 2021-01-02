// Copyright (c) 2021 BitMaelum Authors
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

package gui

import (
	"github.com/bitmaelum/bitmaelum-suite/cmd/bm-mail/gui/app"
	"github.com/bitmaelum/bitmaelum-suite/cmd/bm-mail/gui/layout"
	"github.com/bitmaelum/bitmaelum-suite/internal/vault"
	"github.com/rivo/tview"
)

// Run is the main entrypoint of the gui
func Run(v *vault.Vault) {
	app.App = app.Type{
		App:   tview.NewApplication(),
		Pages: tview.NewPages(),
		Vault: v,
	}

	app.App.Pages.AddPage("main_menu", layout.NewMainMenuScreen(app.App), true, false)
	app.App.Pages.AddPage("accounts", layout.NewAccountScreen(app.App), true, false)
	app.App.Pages.AddPage("message-overview", layout.NewMessageOverviewScreen(app.App), true, false)
	app.App.Pages.AddPage("help", layout.NewHelpScreen(app.App), true, false)
	app.App.Pages.SwitchToPage("main_menu")

	// Setup AFTER we set app.App, otherwise other systems cannot find App.
	if err := app.App.App.SetRoot(app.App.Pages, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}
