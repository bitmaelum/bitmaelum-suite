package layout

import (
	"io/ioutil"
	"net/http"

	"github.com/bitmaelum/bitmaelum-suite/cmd/bm-client/gui/app"
	"github.com/bitmaelum/bitmaelum-suite/cmd/bm-client/gui/components"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

func NewHelpScreen(app app.AppType) tview.Primitive {

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
	help.SetDoneFunc(func(key tcell.Key) {
		app.Pages.SwitchToPage("main_menu")
	})


	// Create menu bar
	mb := components.NewMenubar(app.App)
	mb.SetSlot(0, "fooo", nil)
	mb.SetSlot(9, "close", nil)


	// Create a Flex layout that centers the logo and subtitle.
	grid := tview.NewGrid().SetColumns(10, 0, 10).SetRows(1, 0, 2, 1)
	grid.AddItem(help, 1, 1,  1, 1, 10, 70, true)
	grid.AddItem(mb,   3, 0,  1, 3,  0,  0, false)

	return grid
}
