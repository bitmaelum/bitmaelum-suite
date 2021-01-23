package imap

import (
	"fmt"
	"strings"
	"time"
)

var (
	FlagsAll  = []string{"FLAGS", "INTERNALDATE", "RFC822.SIZE", "ENVELOPE"}
	FlagsFast = []string{"FLAGS", "INTERNALDATE", "RFC822.SIZE"}
	FlagsFull  = []string{"FLAGS", "INTERNALDATE", "RFC822.SIZE", "ENVELOPE", "BODY"}
)

func UidFetch(c *Conn, tag, cmd string, args []string) error {
	set := NewSequenceSet(args[1])

	attrs := fetchAttrs(args[2])

	var i = 1
	info := c.DB.GetBoxInfo(c.Account, c.Box)
	for _, uid := range info.Uids {
		if set.InSet(uid) {

			msgInfo, err := c.DB.Fetch(c.Account, c.Box, info.UIDValidity, uid)
			if err != nil {
				continue
			}

			ret := []string{}

			for _, attr := range attrs {
				switch attr {
				case "RFC822.SIZE":
					ret = append(ret, fmt.Sprintf("RFC822.SIZE %d", 3700))
				case "UID":
					ret = append(ret, fmt.Sprintf("UID %d", uid))
				case "FLAGS":
					ret = append(ret, "FLAGS (\\Unseen)")
					// ret = append(ret, fmt.Sprintf("FLAGS (%s)", strings.Join(msgInfo.Flags, " ")))
				case "BODY.PEEK":
					s := "From: <joshua@bitmaelum.network>\n"
					s += "Reply-To: <joshua@bitmaelum.network>\n"
					s += "To: <joshua@bitmaelum.network>\n"
					// s += "Cc: \n"
					// s += "Bcc: \n"
					s += "Subject: second message on IMAP (" + msgInfo.MessageID + ")\n"
					s += "Date: " + time.Now().Format(time.RFC822) + "\n"
					s += "Message-ID: <" + msgInfo.MessageID + "@bitmaelum.network>\n"
					// s += "Priority: \n"
					// s += "X-Priority: \n"
					// s += "References: \n"
					// s += "Newsgroups: \n"
					// s += "In-Reply-To: \n"
					s += "Content-Type: plain/text\n"
					// s += "Reply-To: \n"
					// s += "List-Unsubscribe: \n"
					s += "Received: from imap.bitmaelum.network\n"
					s += "        by imap.bitmaelum.network with LMTP\n"
					s += "        id GL77MveLDGCvXwAA2ul2EA\n"
					s += "        (envelope-from <joshua@bitmaelum.network>)\n"
					s += "        for <joshua@bitmaelum.network>; " + time.Now().Format(time.RFC822) + "\n"
					s += "Delivery-Date: "+ time.Now().Format(time.RFC822) + ""
					ret = append(ret, fmt.Sprintf("BODY[HEADER.FIELDS (\"From\" \"To\" \"Cc\" \"Bcc\" \"Resent-Message-ID\" \"Subject\" \"Date\" \"Message-ID\" \"Priority\" \"X-Priority\" \"References\" \"Newsgroups\" \"In-Reply-To\" \"Content-Type\" \"Reply-To\" \"List-Unsubscribe\" \"Received\" \"Delivery-Date\")] {%d}\n%s\n", len(s), s))
				}
			}

			c.Write("*", fmt.Sprintf("%d FETCH (%s)", i, strings.Join(ret, " ")))
			i++
		}
	}
	c.Write(tag, "OK FETCH completed")

	return nil
}

func fetchAttrs(s string) []string {
	s = strings.Trim(s, "()")

	if s == "FLAGS" {
		return []string{"FLAGS", "UID"}
	}

	if s == "ALL" {
		return FlagsAll
	}
	if s == "FULL" {
		return FlagsFull
	}
	if s == "FAST" {
		return FlagsFast
	}

	return []string{
		"UID",
		"RFC822.SIZE",
		"FLAGS",
		"BODY.PEEK",
		//"BODY.PEEK[HEADERS.FIELDS(From To Cc Bcc Resent-Message-ID Subject Date Message-ID Priority X-Priority References Newsgroups In-Reply-To Content-Type Reply-To List-Unsubscribe Received Delivery-Date)]",
	}
}

/*
(
	UID
	RFC822.SIZE
	FLAGS
	BODY.PEEK[HEADER.FIELDS (From To Cc Bcc Resent-Message-ID Subject Date Message-ID Priority X-Priority References Newsgroups In-Reply-To Content-Type Reply-To List-Unsubscribe Received Delivery-Date)]
)
*/
