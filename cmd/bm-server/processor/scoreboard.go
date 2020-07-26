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
