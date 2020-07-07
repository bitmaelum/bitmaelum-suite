package main

import (
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
	"github.com/sirupsen/logrus"
)

var app *tview.Application
var pages *tview.Pages

var panels []tview.Primitive
var panelIdx = 0

func createWelcomeScreen() {
	modal := tview.NewModal().
		SetText("Welcome to BitMaelum\n\nThis is the first message client for the BitMaelum system. It's highly experimental but should give you a good idea on the functionality.").
		AddButtons([]string{"Okay"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			pages.SwitchToPage("layout")
		})

	pages.AddPage("welcome", modal, false, false)
}

func navigate(dir int) {
	panelIdx += dir

	panelIdx %= len(panels)
	if panelIdx < 0 {
		panelIdx = len(panels) - 1
	}

	app.SetFocus(panels[panelIdx])
}

func createLayout() {
	newPrimitive := func(text string) tview.Primitive {
		tmp := tview.NewTextView().
			SetTextAlign(tview.AlignCenter).
			SetText(text)

		tmp.SetBorder(true)
		return tmp
	}

	main := newPrimitive("BitMaelum Client")
	accounts := newAccountsPanel()
	boxes := newMessageBoxPanel()

	panels = []tview.Primitive{accounts, boxes, main}

	status := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetWrap(false).
		SetText("This is the text we want to display").
		SetBorder(false).
		SetBackgroundColor(tcell.ColorBlue)

	grid := tview.NewGrid().
		SetRows(0, 0, 1).
		SetColumns(30, 0).
		SetBorders(false).
		AddItem(status, 2, 0, 1, 2, 0, 0, false).
		AddItem(accounts, 0, 0, 1, 1, 0, 0, true).
		AddItem(boxes, 1, 0, 1, 1, 0, 0, false).
		AddItem(main, 0, 1, 2, 1, 0, 0, false)

	pages.AddPage("layout", grid, true, false)
}

func newAccountsPanel() tview.Primitive {
	header := tview.NewTextView().SetText("Accounts").SetTextAlign(tview.AlignCenter)

	list := tview.NewList().
		AddItem("jaytaph! (5)", "", '1', nil).
		AddItem("jthijssen!", "", '2', nil).
		AddItem("noxlogic!", "", '3', nil).
		AddItem("jaytaph@phpnl! (1)", "", '4', nil).
		ShowSecondaryText(false)

	grid := tview.NewGrid().
		SetColumns(0).SetRows(1, 0).SetBorders(true).
		AddItem(header, 0, 0, 1, 1, 0, 0, false).
		AddItem(list, 1, 0, 1, 1, 0, 0, true)

	grid.SetBorderColor(tcell.ColorSteelBlue)

	return grid
}

func newMessageBoxPanel() tview.Primitive {
	header := tview.NewTextView().SetText("Message Boxes").SetTextAlign(tview.AlignCenter)

	root := tview.NewTreeNode("Message boxes")

	root.AddChild(tview.NewTreeNode("Inbox (4)"))
	root.AddChild(tview.NewTreeNode("Send"))
	root.AddChild(tview.NewTreeNode("Trash (32)"))

	personalBox := tview.NewTreeNode("Personal")
	root.AddChild(personalBox)
	personalBox.AddChild(tview.NewTreeNode("Gym"))
	personalBox.AddChild(tview.NewTreeNode("Home"))
	personalBox.AddChild(tview.NewTreeNode("Hobbies (1)"))
	root.AddChild(tview.NewTreeNode("Work related"))
	list := tview.NewTreeView().SetRoot(root).SetCurrentNode(root)

	grid := tview.NewGrid().
		SetColumns(0).SetRows(1, 0).SetBorders(true).
		AddItem(header, 0, 0, 1, 1, 0, 0, false).
		AddItem(list, 1, 0, 1, 1, 0, 0, true)

	return grid
}

func main() {
	app = tview.NewApplication()
	pages = tview.NewPages()

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyTab {
			navigate(1)
			return nil
		}
		if event.Key() == tcell.KeyBacktab {
			navigate(-1)
			return nil
		}

		return event
	})

	createLayout()
	createWelcomeScreen()

	pages.SwitchToPage("welcome")

	if err := app.SetRoot(pages, true).SetFocus(pages).Run(); err != nil {
		logrus.Fatal(err)
	}
}
