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

package components

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type mainMenuItem struct {
	Text     string
	Shortcut rune
}

// MainMenu is our main menu
type MainMenu struct {
	*tview.Box

	items       []mainMenuItem // Items in the main menu
	currentItem int            // Current selected item

	mainTextColor           tcell.Color // main color
	selectedTextColor       tcell.Color // selected item color
	selectedBackgroundColor tcell.Color // selected item background color

	selected func(index int, mainText string, shortcut rune) // function when an item is selected
}

// NewMainMenu creates an empty main menu
func NewMainMenu() *MainMenu {
	return &MainMenu{
		Box:                     tview.NewBox(),
		mainTextColor:           tview.Styles.PrimaryTextColor,
		selectedTextColor:       tview.Styles.PrimitiveBackgroundColor,
		selectedBackgroundColor: tview.Styles.PrimaryTextColor,
	}
}

// Draw will draw the main menu on screen
func (m *MainMenu) Draw(screen tcell.Screen) {
	m.Box.Draw(screen)

	// Determine the dimensions.
	x, y, width, height := m.GetInnerRect()
	bottomLimit := y + height
	_, totalHeight := screen.Size()
	if bottomLimit > totalHeight {
		bottomLimit = totalHeight
	}

	// Draw the list items.
	for index, item := range m.items {
		if y >= bottomLimit {
			break
		}

		// Main text.
		tview.Print(screen, item.Text, x, y, width, tview.AlignCenter, m.mainTextColor)

		// Background color of selected text.
		if index == m.currentItem && m.HasFocus() {
			textWidth := width
			for bx := 0; bx < textWidth; bx++ {
				mc, c, style, _ := screen.GetContent(x+bx, y)
				fg, _, _ := style.Decompose()
				if fg == m.mainTextColor {
					fg = m.selectedTextColor
				}
				style = style.Background(m.selectedBackgroundColor).Foreground(fg)
				screen.SetContent(x+bx, y, mc, c, style)
			}
		}

		y++

		if y >= bottomLimit {
			break
		}
	}
}

// SetMainTextColor sets the main text color
func (m *MainMenu) SetMainTextColor(c tcell.Color) *MainMenu {
	m.mainTextColor = c
	return m
}

// SetSelectedTextColor is the selected text color
func (m *MainMenu) SetSelectedTextColor(c tcell.Color) *MainMenu {
	m.selectedTextColor = c
	return m
}

// SetSelectedBackgroundColor is the selected background color
func (m *MainMenu) SetSelectedBackgroundColor(c tcell.Color) *MainMenu {
	m.selectedBackgroundColor = c
	return m
}

// SetSelectedFunc sets the function that is called when an item is selected
func (m *MainMenu) SetSelectedFunc(handler func(int, string, rune)) *MainMenu {
	m.selected = handler
	return m
}

// AddItem adds a new item to the main menu
func (m *MainMenu) AddItem(text string, shortcut rune) *MainMenu {
	m.items = append(m.items, mainMenuItem{
		Text:     text,
		Shortcut: shortcut,
	})

	return m
}

// Clear will remove all menu items
func (m *MainMenu) Clear() *MainMenu {
	m.items = nil

	return m
}

// InputHandler will be called when the main menu is active and we have an input. It will deal with navigating through
// the main menu
func (m *MainMenu) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return m.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		if event.Key() == tcell.KeyEscape {
			return
		} else if len(m.items) == 0 {
			return
		}

		// previousItem := m.currentItem

		switch key := event.Key(); key {
		case tcell.KeyTab, tcell.KeyDown, tcell.KeyRight:
			m.currentItem++
		case tcell.KeyBacktab, tcell.KeyUp, tcell.KeyLeft:
			m.currentItem--
		case tcell.KeyHome:
			m.currentItem = 0
		case tcell.KeyEnd:
			m.currentItem = len(m.items) - 1
		case tcell.KeyEnter:
			if m.currentItem >= 0 && m.currentItem < len(m.items) {
				item := m.items[m.currentItem]
				if m.selected != nil {
					m.selected(m.currentItem, item.Text, item.Shortcut)
				}
			}
		case tcell.KeyRune:
			doRune(m, event.Rune())
		}

		if m.currentItem < 0 {
			m.currentItem = 0
		} else if m.currentItem >= len(m.items) {
			m.currentItem = len(m.items) - 1
		}
	})
}

func doRune(m *MainMenu, ch rune) {
	if ch != ' ' {
		// It's not a space bar. Is it a shortcut?
		var found bool
		for index, item := range m.items {
			if item.Shortcut == ch {
				// We have a shortcut.
				found = true
				m.currentItem = index
				break
			}
		}
		if !found {
			return
		}
	}

	// call the selected function if any
	item := m.items[m.currentItem]
	if m.selected != nil {
		m.selected(m.currentItem, item.Text, item.Shortcut)
	}
}
