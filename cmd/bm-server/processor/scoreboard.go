package processor

var scoreboard = map[string]int{}

// IsInScoreboard will check if the given message ID in the section is present
func IsInScoreboard(section int, msgID string) bool {
	s, found := scoreboard[msgID]

	return found && s == section
}

// AddToScoreboard will set the message ID in the section in the scoreboard
func AddToScoreboard(section int, msgID string) {
	scoreboard[msgID] = section
}

// RemoveFromScoreboard will remove the message ID from the section in the scoreboard
func RemoveFromScoreboard(section int, msgID string) {
	delete(scoreboard, msgID)
}
