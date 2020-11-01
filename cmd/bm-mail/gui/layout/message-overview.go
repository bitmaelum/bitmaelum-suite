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
	"math/rand"

	"github.com/bitmaelum/bitmaelum-suite/cmd/bm-mail/gui/app"
	"github.com/bitmaelum/bitmaelum-suite/cmd/bm-mail/gui/components"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

// NewMessageOverviewScreen creates a new message overview screen
func NewMessageOverviewScreen(app app.Type) tview.Primitive {
	boxes := tview.NewTreeView()
	boxes.SetBorder(true)
	boxes.SetTitle("Boxes")
	boxes.SetGraphics(true)

	root := tview.NewTreeNode("Accounts").SetColor(tcell.ColorGreen).SetExpanded(true).SetSelectable(false)
	for _, acc := range app.Vault.Store.Accounts {
		accNode := tview.NewTreeNode(acc.Address).SetColor(tcell.ColorYellow)
		root.AddChild(accNode)
	}

	accounts := tview.NewTreeView()
	accounts.SetBorder(true)
	accounts.SetRoot(root)
	accounts.SetTitle("Accounts")
	accounts.SetGraphics(true)

	accounts.SetSelectedFunc(func(node *tview.TreeNode) {
		box := node.GetText()

		tree := tview.NewTreeNode(box).SetColor(tcell.ColorYellow).SetExpanded(true).SetSelectable(false)
		tree.AddChild(tview.NewTreeNode(fmt.Sprintf("incoming (%d)", rand.Intn(100))))
		tree.AddChild(tview.NewTreeNode("sent"))
		tree.AddChild(tview.NewTreeNode(fmt.Sprintf("trash (%d)", rand.Intn(100))))
		boxes.SetRoot(tree)
	})

	messages := tview.NewTextView()
	messages.SetBorder(true)

	menuBar := components.NewMenubar(app.App)
	menuBar.SetSlot(9, "Back", func() {
		app.Pages.SwitchToPage("main_menu")
	})

	messages.SetDoneFunc(func(key tcell.Key) {
		app.Pages.SwitchToPage("main_menu")
	})

	grid := tview.NewGrid().SetColumns(30, 0).SetRows(20, 0, 1)
	grid.AddItem(accounts, 0, 0, 1, 1, 20, 30, true)
	grid.AddItem(boxes, 1, 0, 1, 1, 0, 0, false)
	grid.AddItem(messages, 0, 1, 2, 1, 0, 0, false)
	grid.AddItem(menuBar, 2, 0, 1, 2, 0, 0, false)

	grid.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == 'a' {
			app.App.SetFocus(accounts)
			return nil
		}
		if event.Rune() == 'b' {
			app.App.SetFocus(boxes)
			return nil
		}
		if event.Rune() == 'm' {
			app.App.SetFocus(messages)
			return nil
		}

		return event
	})

	return grid
}
