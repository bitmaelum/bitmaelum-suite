package layout

import (
	"github.com/bitmaelum/bitmaelum-suite/cmd/bm-client/gui/app"
	"github.com/bitmaelum/bitmaelum-suite/cmd/bm-client/gui/components"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

func NewMessageOverviewScreen(app app.AppType) tview.Primitive {
	root := tview.NewTreeNode("Accounts").SetColor(tcell.ColorGreen).SetExpanded(true)
	for _, acc := range app.Vault.Store.Accounts {
		accNode := tview.NewTreeNode(acc.Address).SetColor(tcell.ColorYellow)
		root.AddChild(accNode)
	}

	accounts := tview.NewTreeView()
	accounts.SetBorder(true)
	accounts.SetRoot(root)
	accounts.SetTitle("Accounts")
	accounts.SetGraphics(true)

	boxroot := tview.NewTreeNode("joshua!").SetColor(tcell.ColorYellow).SetExpanded(true)
	boxroot.AddChild(tview.NewTreeNode("incoming (52)"))
	boxroot.AddChild(tview.NewTreeNode("sent"))
	boxroot.AddChild(tview.NewTreeNode("trashcan (4)"))

	boxes := tview.NewTreeView()
	boxes.SetBorder(true)
	boxes.SetRoot(boxroot)
	boxes.SetTitle("Accounts")
	boxes.SetGraphics(true)



	messages := tview.NewTextView()
	messages.SetBorder(true)


	menuPages := tview.NewPages()
	searchInput := tview.NewInputField()
	searchInput.
		SetLabel("Search for: ").
		SetFieldWidth(40).
		SetAcceptanceFunc(nil).
		SetDoneFunc(func(key tcell.Key) {
			// search = searchInput.GetText()
			// menuPages.HidePage("search")
			// pipe = selectService(pipe, list.Model[list.GetCurrentItem()], filter, logView, serviceView, infoTable, confirmDialog)
			// app.SetFocus(list)
		})

	menuBar := components.NewMenubar(app.App)
	menuPages.AddPage("menu", menuBar, true, true)
	menuPages.AddPage("search", searchInput, true, false)
	//menuPages.SwitchToPage("search")
	menuBar.SetSlot(9, "Back", func() {})

	messages.SetDoneFunc(func(key tcell.Key) {
		app.Pages.SwitchToPage("main_menu")
	})


	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(tview.NewFlex().
			SetDirection(tview.FlexColumn).

			AddItem(tview.NewFlex().
				SetDirection(tview.FlexRow).
				AddItem(accounts, 15, 1, true).
				AddItem(boxes, 0, 1, false),
				30, 1, true,
			).
			AddItem(messages, 0, 1, false),
			0, 1, false,
		).
		AddItem(menuPages, 3, 1, false)

	return flex
}
