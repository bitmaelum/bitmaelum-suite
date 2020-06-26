package messagebox

type MessageList struct {
	Id    string   `json:"id"`
	Dt    string   `json:"datetime"`
	Flags []string `json:"flags"`
}
