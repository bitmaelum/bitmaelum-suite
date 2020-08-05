package message

// Flags represents the .flags.json file which holds all current flags for the given mailbox/message
type Flags struct {
	Flags []string `json:"flags"`
}

// List is a message list
type List struct {
	ID    string   `json:"id"`
	Dt    string   `json:"datetime"`
	Flags []string `json:"flags"`
}
