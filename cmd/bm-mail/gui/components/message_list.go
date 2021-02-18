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

package components

import (
	"strings"
	"time"

	"github.com/bitmaelum/bitmaelum-suite/internal/message"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type MessageList struct {
	*tview.Box

	Items []message.DecryptedMessage

	selectFunc func(ml MessageList, idx int) // Called when we have selected an entry

	scrollOffset  int // the amount of messages in the beginning we DONT show
	selectedIndex int // current selected index in Items[]
}

// NewMessageList creates a new message list
func NewMessageList() *MessageList {
	m := &MessageList{
		Box: tview.NewBox(),
	}

	return m
}

func (m *MessageList) SetSelectFunc(f func(ml MessageList, idx int)) {
	m.selectFunc = f
}

var cols = []int{8, 40, 20, 20, 20, 10, 0}

// Draw will draw the message list on the screen
func (m *MessageList) Draw(screen tcell.Screen) {
	m.Box.DrawForSubclass(screen, m)

	// Determine the dimensions.
	x, y, width, height := m.GetInnerRect()

	bottomLimit := y + height
	_, totalHeight := screen.Size()
	if bottomLimit > totalHeight {
		bottomLimit = totalHeight
	}

	tcell.StyleDefault = tcell.StyleDefault.Background(tcell.ColorWhite)
	tview.Print(screen, strings.Repeat(" ", width), x, y, width, tview.AlignLeft, tcell.ColorBlack)
	m.printColumns(screen, x, y, tcell.ColorBlack, "Flags", "Subject", "Name", "Address", "Received", "ID")
	tcell.StyleDefault = tcell.StyleDefault.Background(tcell.ColorDefault)
	y++

	for index, msg := range m.Items {
		// Don't display any messages before our scroll offset
		if index < m.scrollOffset {
			continue
		}

		// WHen we hit the end of our window, stop displaying messages
		if y >= bottomLimit {
			break
		}

		columns := []string{
			"[A-!-]",
			msg.Catalog.Subject,
			msg.Catalog.From.Name,
			msg.Catalog.From.Address,
			msg.Catalog.CreatedAt.Format(time.RFC822),
			msg.ID[:8],
		}
		m.printColumns(screen, x, y, tcell.ColorWhite, columns...)

		// Inverse colors on selected entry
		if index == m.selectedIndex {
			for bx := 0; bx < width; bx++ {
				m, c, style, _ := screen.GetContent(x+bx, y)
				style = style.Background(tcell.ColorYellow).Foreground(tcell.ColorBlack)
				screen.SetContent(x+bx, y, m, c, style)
			}
		}

		y++
	}
}

func (m *MessageList) printColumns(screen tcell.Screen, x, y int, col tcell.Color, items ...string) {
	_, _, width, _ := m.GetInnerRect()

	for idx, item := range items {
		size := cols[idx]
		if size == 0 {
			size = width - x
		}
		tview.Print(screen, item, x, y, size, tview.AlignLeft, col)
		x += cols[idx] + 1
	}
}

func (m *MessageList) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return m.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		_, _, _, height := m.GetInnerRect()

		switch key := event.Key(); key {
		case tcell.KeyTab, tcell.KeyDown, tcell.KeyRight:
			m.selectedIndex++
		case tcell.KeyBacktab, tcell.KeyUp, tcell.KeyLeft:
			m.selectedIndex--
		case tcell.KeyHome:
			m.selectedIndex = 0
		case tcell.KeyEnd:
			m.selectedIndex = len(m.Items) - 1
		case tcell.KeyPgDn:
			m.selectedIndex += height
		case tcell.KeyPgUp:
			m.selectedIndex -= height
		case tcell.KeyEnter:
			if m.selectFunc != nil {
				m.selectFunc(*m, m.selectedIndex)
			}
			return
		}

		// Limit our selected entry within our item range
		if m.selectedIndex < 0 {
			m.selectedIndex = 0
		}
		if m.selectedIndex > len(m.Items)-1 {
			m.selectedIndex = len(m.Items) - 1
		}

		// Adjust offset to keep the current selection in view.
		if m.selectedIndex < m.scrollOffset {
			m.scrollOffset = m.selectedIndex
		} else {
			if m.selectedIndex-m.scrollOffset >= height-1 {
				m.scrollOffset = m.selectedIndex + 2 - height
			}
		}
	})
}
