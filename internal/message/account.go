package message

// MailBoxInfo represents the .info.json file
type MailBoxInfo struct {
	Name        string
	Description string `json:"description"`
	Quota       int    `json:"quota"`
	Total       int    `json:"total"`
}

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
