package imap

import (
	"strconv"
)

func Lsub(c *Conn, tag, cmd string, args []string) error {
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
		c.Write("*", "LSUB () \"/\" \""+boxName+"\"")
	}
	c.Write(tag, "OK LSUB completed")

	return nil
}
