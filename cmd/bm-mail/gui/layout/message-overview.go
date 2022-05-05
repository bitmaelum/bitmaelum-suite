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

package layout

import (
	"io/ioutil"
	"strings"

	"github.com/bitmaelum/bitmaelum-suite/cmd/bm-mail/app"
	"github.com/bitmaelum/bitmaelum-suite/cmd/bm-mail/gui/components"
	"github.com/bitmaelum/bitmaelum-suite/internal/message"
	"github.com/bitmaelum/bitmaelum-suite/pkg/address"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// NewMessageOverviewScreen creates a new message overview screen
func NewMessageOverviewScreen() tview.Primitive {

	app.MailApp.MessageAccountTree = tview.NewTreeView()
	app.MailApp.MessageAccountTree.SetBorder(true)
	app.MailApp.MessageAccountTree.SetTitle("Accounts")
	app.MailApp.MessageAccountTree.SetGraphics(true)
	app.MailApp.MessageAccountTree.SetSelectedFunc(func(node *tview.TreeNode) {
		addr, err := address.NewAddress(node.GetText())
		if err != nil {
			return
		}
		app.MailApp.CurrentAddr = addr
		app.MailApp.StaleAddr = true
		Refresh()
	})

	app.MailApp.MessageBoxTree = tview.NewTreeView()
	app.MailApp.MessageBoxTree.SetBorder(true)
	app.MailApp.MessageBoxTree.SetTitle("Mailboxes")
	app.MailApp.MessageBoxTree.SetGraphics(true)
	app.MailApp.MessageBoxTree.SetSelectedFunc(func(node *tview.TreeNode) {
		box := "1" // node.GetText()
		app.MailApp.CurrentBox = &box
		app.MailApp.StaleBox = true
		Refresh()
	})

	app.MailApp.MessageList = components.NewMessageList()
	app.MailApp.MessageList.SetBorder(true)
	app.MailApp.MessageList.SetSelectFunc(func(ml components.MessageList, idx int) {
		app.MailApp.MessageView.SetText(createMessageText(ml.Items[idx]))
	})

	menuBar := components.NewMenubar(app.MailApp.App)
	menuBar.SetSlot(9, "Back", func() {
		app.MailApp.Pages.SwitchToPage("main_menu")
	})

	app.MailApp.MessageView = tview.NewTextView()
	app.MailApp.MessageView.SetBorder(true)
	app.MailApp.MessageView.SetText("load a message first")

	// MessageList.SetDoneFunc(func(key tcell.Key) {
	// 	app.Pages.SwitchToPage("main_menu")
	// })

	grid := tview.NewGrid().SetColumns(30, 0).SetRows(20, 0, 1)
	grid.AddItem(app.MailApp.MessageAccountTree, 0, 0, 1, 1, 20, 30, true)
	grid.AddItem(app.MailApp.MessageBoxTree, 1, 0, 1, 1, 0, 0, false)
	grid.AddItem(app.MailApp.MessageList, 0, 1, 1, 1, 0, 0, false)
	grid.AddItem(app.MailApp.MessageView, 1, 1, 1, 1, 0, 0, false)
	grid.AddItem(menuBar, 2, 0, 1, 2, 0, 0, false)

	grid.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == 'a' {
			app.MailApp.App.SetFocus(app.MailApp.MessageAccountTree)
			return nil
		}
		if event.Rune() == 'b' {
			app.MailApp.App.SetFocus(app.MailApp.MessageBoxTree)
			return nil
		}
		if event.Rune() == 'l' {
			app.MailApp.App.SetFocus(app.MailApp.MessageList)
			return nil
		}
		if event.Rune() == 'm' {
			app.MailApp.App.SetFocus(app.MailApp.MessageView)
			return nil
		}

		return event
	})

	return grid
}

func createMessageText(msg message.DecryptedMessage) string {
	ret := ""
	ret += "Subject : " + msg.Catalog.Subject + "\n"
	ret += "From    : " + msg.Catalog.From.Name + "<" + msg.Catalog.From.Address + ">\n"
	ret += "To      : <" + msg.Catalog.To.Address + ">\n"
	ret += "Date    : " + msg.Catalog.CreatedAt.String() + "\n"
	ret += "Flags   : [" + strings.Join(msg.Catalog.Flags, ",") + "]\n"
	ret += "Labels  : [" + strings.Join(msg.Catalog.Labels, ",") + "]\n"

	ret += "\n"

	if msg.Catalog.Blocks[0].Reader != nil {
		b, err := ioutil.ReadAll(msg.Catalog.Blocks[0].Reader)
		if err != nil {
			ret += "cannot decode block"
		} else {
			ret += string(b) + "\n"
		}
	}
	return ret
}
