package imap

import (
	"strconv"
)

func List(c *Conn, tag, cmd string, args []string) error {
	// if args[0] == "" {
	// 	c.Write("*", "LIST (\\Noselect) \"/\" \"\"")
	// }

	if args[1] == "*" {
		mbl, err := c.Client.GetMailboxList(c.Info.Address.Hash())
		if err != nil {
			panic(err)
		}

		for _, box := range mbl.Boxes {
			boxName := strconv.Itoa(box.ID)
			switch box.ID {
			case 1:
				boxName = "Inbox"
			case 2:
				boxName = "Sent"
			case 3:
				boxName = "Trash"
			}
			c.Write("*", "LIST (\\HasNoChildren \\UnMarked) \".\" \""+boxName+"\"")
		}
	}
	c.Write(tag, "OK LIST completed")

	return nil
}
