package gui

import (
	"github.com/bitmaelum/bitmaelum-suite/cmd/bm-mail/gui/app"
	"github.com/bitmaelum/bitmaelum-suite/cmd/bm-mail/gui/layout"
	"github.com/bitmaelum/bitmaelum-suite/internal/vault"
	"github.com/rivo/tview"
)

func Run(v *vault.Vault) {
	app.App = app.AppType{
		App:    tview.NewApplication(),
		Pages:  tview.NewPages(),
		Vault:  v,
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
