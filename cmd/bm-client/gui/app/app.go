package app

import (
	"github.com/bitmaelum/bitmaelum-suite/internal/vault"
	"github.com/rivo/tview"
)

type AppType struct {
	App    *tview.Application
	Pages  *tview.Pages
	Vault  *vault.Vault
}

var App AppType
