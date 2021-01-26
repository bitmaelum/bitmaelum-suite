package imap

import (
	"fmt"
	"strings"
)

func Status(c *Conn, tag, cmd string, args []string) error {
	box := args[0]
	switch strings.ToUpper(args[0]) {
	case "INBOX":
		box = "1"
	case "SENT":
		box = "2"
	case "TRASH":
		box = "3"
	}
	c.ChangeBox(box)

	c.Write("*", fmt.Sprintf("STATUS %s (UNSEEN %d)", args[0], len(c.Index)))

	c.Write(tag, "OK STATUS completed")

	return nil
}
