package layout

 import (
	"fmt"
	"math/rand"

	"github.com/bitmaelum/bitmaelum-suite/cmd/bm-client/gui/app"
	"github.com/bitmaelum/bitmaelum-suite/cmd/bm-client/gui/components"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

func NewMessageOverviewScreen(app app.AppType) tview.Primitive {
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
	grid.AddItem(menuBar,  2, 0,  1, 2,  0,  0, false)

	grid.SetInputCapture(func (event *tcell.EventKey) *tcell.EventKey {
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
