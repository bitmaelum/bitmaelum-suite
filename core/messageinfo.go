package core

// MessageList is a message list
type MessageList struct {
	ID    string   `json:"id"`
	Dt    string   `json:"datetime"`
	Flags []string `json:"flags"`
}
