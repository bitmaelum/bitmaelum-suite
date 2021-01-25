package imap

import (
	"fmt"
	"strings"
	"time"
)

func Select(c *Conn, tag, cmd string, args []string) error {
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
		fmt.Println("ERR: ", err)
		return err
	}

	boxInfo, _ := c.UpdateImapDB(msgList)


	c.SeqList = make([]string, len(msgList.Messages))
	for i, msg := range msgList.Messages {
		c.SeqList[i] = msg.ID
	}
	fmt.Println(c.SeqList)

	c.Write("*", fmt.Sprintf("%d EXISTS", len(msgList.Messages)))
	c.Write("*", fmt.Sprintf("%d RECENT", 0))
	c.Write("*", fmt.Sprintf("OK [UNSEEN %d]", 1))
	c.Write("*", fmt.Sprintf("OK [UIDNEXT %d]", boxInfo.HighestUID))
	c.Write("*", fmt.Sprintf("OK [UIDVALIDITY %d]", boxInfo.UIDValidity))
	c.Write("*", "OK [HIGHESTMODSEQ 1234567]")

	c.Write(tag, "OK [READ-WRITE] SELECT completed")

	return nil
}
