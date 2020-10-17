package components

import (
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

type mainMenuItem struct {
	Text     string
	Shortcut rune
}

type MainMenu struct {
	*tview.Box

	items []mainMenuItem
	currentItem int

	mainTextColor tcell.Color
	selectedTextColor tcell.Color
	selectedBackgroundColor tcell.Color

	selected func(index int, mainText string, shortcut rune)
}

func NewMainMenu() *MainMenu {
	return &MainMenu{
		Box:                     tview.NewBox(),
		mainTextColor:           tview.Styles.PrimaryTextColor,
		selectedTextColor:       tview.Styles.PrimitiveBackgroundColor,
		selectedBackgroundColor: tview.Styles.PrimaryTextColor,
	}
}

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

func (m *MainMenu) SetMainTextColor(c tcell.Color) *MainMenu {
	m.mainTextColor = c
	return m
}

func (m *MainMenu) SetSelectedTextColor(c tcell.Color) *MainMenu {
	m.selectedTextColor = c
	return m
}

func (m *MainMenu) SetSelectedBackgroundColor(c tcell.Color) *MainMenu {
	m.selectedBackgroundColor = c
	return m
}

func (m *MainMenu) SetSelectedFunc(handler func(int, string, rune)) *MainMenu {
	m.selected = handler
	return m
}

func (m *MainMenu) AddItem(text string, shortcut rune) *MainMenu {
	m.items = append(m.items, mainMenuItem{
		Text:     text,
		Shortcut: shortcut,
	})

	return m
}

func (m *MainMenu) Clear() *MainMenu {
	m.items = nil

	return m
}

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
			ch := event.Rune()
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
					break
				}
			}
			item := m.items[m.currentItem]
			if m.selected != nil {
				m.selected(m.currentItem, item.Text, item.Shortcut)
			}
		}

		if m.currentItem < 0 {
			m.currentItem = 0
		} else if m.currentItem >= len(m.items) {
			m.currentItem = len(m.items) - 1
		}
	})
}
