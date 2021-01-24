package imap

import (
	"fmt"
	"strings"
	"time"
)

func Status(c *Conn, tag, cmd string, args []string) error {
	c.Box = args[0]
	switch strings.ToUpper(args[0]) {
	case "INBOX":
		c.Box = "1"
	case "SENT":
		c.Box = "2"
	case "TRASH":
		c.Box = "3"
	}

	msgList, err := c.Client.GetMailboxMessages(c.Info.Address.Hash(), c.Box, time.Time{})
	if err != nil {
		return err
	}

	_, unseen := c.UpdateImapDB(msgList)

	c.Write("*", fmt.Sprintf("STATUS %s (UNSEEN %d)", args[0], unseen))

	c.Write(tag, "OK STATUS completed")

	return nil
}
