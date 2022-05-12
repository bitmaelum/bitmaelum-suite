// Copyright (c) 2022 BitMaelum Authors
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

package processor

import "sync"

type scoreboardType struct {
	mux sync.Mutex
	v   map[string]int
}

var scoreboard scoreboardType

// IsInScoreboard will check if the given message ID in the section is present
func IsInScoreboard(section int, msgID string) bool {
	s, found := scoreboard.v[msgID]

	return found && s == section
}

// AddToScoreboard will set the message ID in the section in the scoreboard
func AddToScoreboard(section int, msgID string) {
	scoreboard.mux.Lock()
	scoreboard.v[msgID] = section
	scoreboard.mux.Unlock()
}

// RemoveFromScoreboard will remove the message ID from the section in the scoreboard
func RemoveFromScoreboard(section int, msgID string) {
	scoreboard.mux.Lock()
	delete(scoreboard.v, msgID)
	scoreboard.mux.Unlock()
}

func init() {
	scoreboard.v = make(map[string]int)
}
