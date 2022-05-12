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

package gui

import (
	"fmt"

	"github.com/bitmaelum/bitmaelum-suite/cmd/bm-mail/app"
	"github.com/bitmaelum/bitmaelum-suite/cmd/bm-mail/gui/layout"
	"github.com/bitmaelum/bitmaelum-suite/internal/api"
	"github.com/bitmaelum/bitmaelum-suite/internal/vault"
	"github.com/rivo/tview"
)

// Run is the main entrypoint of the gui
func Run() {
	app.MailApp = &app.BmMailAppType{
		App:          tview.NewApplication(),
		Pages:        tview.NewPages(),
		MailBoxLists: make(map[string]*api.MailboxList),
	}

	app.MailApp.Pages.AddPage("unlock_vault", layout.NewUnlockVault(), true, false)
	app.MailApp.Pages.AddPage("main_menu", layout.NewMainMenuScreen(), true, false)
	app.MailApp.Pages.AddPage("accounts", layout.NewAccountScreen(), true, false)
	app.MailApp.Pages.AddPage("message-overview", layout.NewMessageOverviewScreen(), true, false)
	app.MailApp.Pages.AddPage("help", layout.NewHelpScreen(), true, false)
	app.MailApp.App.SetRoot(app.MailApp.Pages, true)

	// Check if the vault can be opened. If not, display our unlock modal
	v, err := vault.Open(vault.VaultPath, vault.VaultPassword)
	if err != nil {
		// Switch unlock page
		app.MailApp.Pages.SwitchToPage("unlock_vault")
		app.MailApp.App.SetFocus(layout.ModalDialogInput)
	} else {
		app.MailApp.Vault = v
		app.MailApp.StaleVault = true
		layout.Refresh()
		app.MailApp.Pages.SwitchToPage("main_menu")
		app.MailApp.App.SetFocus(layout.MainMenuGrid)
	}

	if err := app.MailApp.App.EnableMouse(true).Run(); err != nil {
		panic(err)
	}

	fmt.Print("\nThank you for using BitMaelum, the privacy-first email alternative network.\n\n")
}
