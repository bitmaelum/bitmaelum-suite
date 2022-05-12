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
	"os"

	"github.com/bitmaelum/bitmaelum-suite/cmd/bm-mail/app"
	"github.com/bitmaelum/bitmaelum-suite/internal/vault"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// ModalDialogBox is the modal box that displays the password form
var ModalDialogBox *tview.Frame

// ModalDialogInput is the input for the password
var ModalDialogInput *tview.Form

// NewUnlockVault Display a modal dialog that opens the vault
func NewUnlockVault() tview.Primitive {

	createModal := func(p tview.Primitive, width, height int) tview.Primitive {
		return tview.NewFlex().
			AddItem(nil, 0, 1, false).
			AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
				AddItem(nil, 0, 1, false).
				AddItem(p, height, 1, true).
				AddItem(nil, 0, 1, false), width, 1, false).
			AddItem(nil, 0, 1, false)
	}

	errorDisplayed := false

	ModalDialogInput = tview.NewForm().
		AddPasswordField("Password:", "", 40, '*', func(text string) {
			vault.VaultPassword = text
		}).
		AddButton("Open vault", func() {
			var err error
			app.MailApp.Vault, err = vault.Open(vault.VaultPath, vault.VaultPassword)
			if err == nil {
				app.MailApp.StaleVault = true
				Refresh()
				app.MailApp.Pages.SwitchToPage("main_menu")
				app.MailApp.App.SetFocus(MainMenu)
				return
			}

			if !errorDisplayed {
				errorDisplayed = true
				ModalDialogBox.AddText("Incorrect vault password specified", false, tview.AlignCenter, tcell.ColorRed)
			}
			// Clear password
			ModalDialogInput.GetFormItem(0).(*tview.InputField).SetText("")
			ModalDialogInput.SetFocus(0)
			app.MailApp.App.SetFocus(ModalDialogInput)
		}).
		AddButton("Quit", func() {
			app.MailApp.App.Stop()
			os.Exit(0)
		}).
		SetButtonsAlign(tview.AlignCenter)

	ModalDialogBox = tview.NewFrame(ModalDialogInput).SetBorders(0, 0, 0, 0, 0, 0)

	ModalDialogBox.SetBorder(true).
		SetBorderPadding(1, 1, 1, 1).
		SetTitle(" Enter your vault password ")

	return createModal(ModalDialogBox, 60, 10)
}
