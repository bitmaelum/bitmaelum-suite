package imap

import (
	"fmt"
	"strings"
)

func Select(c *Conn, tag, cmd string, args []string) error {
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

	c.Write("*", fmt.Sprintf("%d EXISTS", len(c.Index)))
	c.Write("*", fmt.Sprintf("%d RECENT", 0))
	c.Write("*", fmt.Sprintf("OK [UNSEEN %d]", 1))
	c.Write("*", fmt.Sprintf("OK [UIDNEXT %d]", c.BoxInfo.HighestUID))
	c.Write("*", fmt.Sprintf("OK [UIDVALIDITY %d]", c.BoxInfo.UIDValidity))
	c.Write("*", "OK [HIGHESTMODSEQ 1234567]")

	c.Write(tag, "OK [READ-WRITE] SELECT completed")

	return nil
}

// func createSequenceList(c *Conn, msgList *api.MailboxMessages) []string {
// 	var seqList []string
//
// 	for _, msg := range msgList.Messages {
// 		// Fetch actual message from bitmaelum
// 		msg, err := c.Client.GetMessage(c.Info.Address.Hash(), msg.ID)
// 		if err != nil {
// 			continue
// 		}
//
// 		// Decrypt message
// 		em := message.EncryptedMessage{
// 			ID:      msg.ID,
// 			Header:  &msg.Header,
// 			Catalog: msg.Catalog,
//
// 			GenerateBlockReader:      c.Client.GenerateAPIBlockReader(c.Info.Address.Hash()),
// 			GenerateAttachmentReader: c.Client.GenerateAPIAttachmentReader(c.Info.Address.Hash()),
// 		}
//
// 		_, err = em.Decrypt(c.Info.GetActiveKey().PrivKey)
// 		if err != nil {
// 			continue
// 		}
//
// 		seqList = append(seqList, msg.ID)
// 	}
//
// 	return seqList
// }
