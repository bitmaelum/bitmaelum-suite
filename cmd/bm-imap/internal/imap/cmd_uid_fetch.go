package imap

import (
	"fmt"
	"strings"

	"github.com/bitmaelum/bitmaelum-suite/cmd/bm-imap/internal"
	"github.com/bitmaelum/bitmaelum-suite/internal/message"
)

const TimeFormat = "Sun, 02-Jan-2006 15:04:05 -0700"

func UidFetch(c *Conn, tag, cmd string, args []string) error {
	set := NewSequenceSet(args[1])
	attrs := ParseAttributes(strings.Join(args[2:], " "))

	var i = 1
	info := c.DB.GetBoxInfo(c.Account, c.Box)
	for _, uid := range info.Uids {
		if !set.InSet(uid) {
			continue
		}

		// Get message info based on box, UIDvalidity and UID
		msgInfo, err := c.DB.Fetch(c.Account, c.Box, info.UIDValidity, uid)
		if err != nil {
			continue
		}

		// Fetch actual message from bitmaelum
		msg, err := c.Client.GetMessage(c.Info.Address.Hash(), msgInfo.MessageID)
		if err != nil {
			continue
		}

		// Decrypt message
		em := message.EncryptedMessage{
			ID:      msg.ID,
			Header:  &msg.Header,
			Catalog: msg.Catalog,

			GenerateBlockReader:      c.Client.GenerateAPIBlockReader(c.Info.Address.Hash()),
			GenerateAttachmentReader: c.Client.GenerateAPIAttachmentReader(c.Info.Address.Hash()),
		}

		dMsg, err := em.Decrypt(c.Info.GetActiveKey().PrivKey)
		if err != nil {
			continue
		}

		// fill up found attributes
		ret := executeAttributes(attrs, *msgInfo, *dMsg)

		c.Write("*", fmt.Sprintf("%d FETCH (%s)", i, strings.Join(ret, " ")))
		i++
	}
	c.Write(tag, "OK FETCH completed")

	return nil
}

func executeAttributes(attrs []Attribute, msgInfo internal.MessageInfo, dMsg message.DecryptedMessage) []string {
	var ret []string

	from := addrToEmail(dMsg.Catalog.From.Address)
	to := addrToEmail(dMsg.Catalog.To.Address)

	fromParts := strings.Split(from, "@")
	toParts := strings.Split(to, "@")



	for _, attr := range attrs {
		switch attr.Name {
		case "RFC822":
			switch attr.Section {
			case "SIZE":
				ret = append(ret, fmt.Sprintf("RFC822.SIZE %d", 12345))
			}

		case "ENVELOPE":

			ret = append(ret, fmt.Sprintf("ENVELOPE (\"%s\" \"%s\" %s %s %s %s %s %s %s \"<%s@bitmaelum.network>\")",
				dMsg.Catalog.CreatedAt.Format(TimeFormat),
				dMsg.Catalog.Subject,
				fmt.Sprintf("((\"%s\" NIL \"%s\" \"%s\"))", from, fromParts[0], fromParts[1]),
				fmt.Sprintf("((\"%s\" NIL \"%s\" \"%s\"))", from, fromParts[0], fromParts[1]),
				fmt.Sprintf("((\"%s\" NIL \"%s\" \"%s\"))", from, fromParts[0], fromParts[1]),
				fmt.Sprintf("((\"%s\" NIL \"%s\" \"%s\"))", to, toParts[0], toParts[1]),
				"NIL",
				"NIL",
				"NIL",
				msgInfo.MessageID,
			))

		case "BODYSTRUCTURE":
			ret = append(ret, "BODYSTRUCTURE (\"TEXT\" \"PLAIN\" (\"CHARSET\" \"UTF-8\") NIL NIL \"7bit\" 12345 64 NIL NIL NIL NIL)")

		case "UID":
			ret = append(ret, fmt.Sprintf("UID %d", msgInfo.UID))

		case "FLAGS":
			ret = append(ret, fmt.Sprintf("FLAGS (%s)", strings.Join(msgInfo.Flags, " ")))

		case "INTERNALDATE":
			ret = append(ret, fmt.Sprintf("INTERNALDATE \"%s\"", dMsg.Catalog.CreatedAt.Format(TimeFormat)))

		case "BODY":
			switch attr.Section {
			case "HEADER.FIELDS":
				hdrs := ""
				for _, h := range attr.Headers {
					s := ""
					switch h {
					case "x-priority":
						s = fmt.Sprintf("X-Priority: 3")
					case "content-type":
						s = fmt.Sprintf("Content-Type: text/plain")
					case "from":
						s = fmt.Sprintf("From: %s <%s>", dMsg.Catalog.From.Name, from)
					case "to":
						s = fmt.Sprintf("To: <%s>", to)
					case "reply-to":
						s = fmt.Sprintf("Reply-To: <%s>", to)
					case "subject":
						s = fmt.Sprintf("Subject: %s", dMsg.Catalog.Subject)
					case "date":
						s = fmt.Sprintf("Date: %s", dMsg.Catalog.CreatedAt.Format(TimeFormat))
					case "message-id":
						s = "Message-ID: <" + msgInfo.MessageID + "@bitmaelum.network>\n"
					case "received":
						s = fmt.Sprintf("Received: from imap.bitmaelum.network\n"+
							"        by imap.bitmaelum.network with LMTP\n"+
							"        id %s\n"+
							"        (envelope-from <%s>)\n"+
							"        for <%s>; %s",
							from, msgInfo.MessageID, from, dMsg.Catalog.CreatedAt.Format(TimeFormat))
					case "delivery-date":
						s = "Delivery-Date: " + dMsg.Catalog.CreatedAt.Format(TimeFormat)
					}
					if s != "" {
						hdrs += s + "\n"
					}
				}
				ret = append(ret, fmt.Sprintf("%s {%d}\n%s\n", attr.ToString(), len(hdrs), hdrs))

			case "TEXT":
			case "MIME":
			case "HEADER":
			}


		}
	}

	return ret
}

func addrToEmail(address string) string {
	address = strings.Replace(address, "@", "#", -1)
	address = strings.Replace(address, "!", "", -1)

	return address + "@bitmaelum.network"

}

/*
* 478 FETCH (
    UID 318459
    RFC822.SIZE 7057
	FLAGS (\Seen NonJunk)
	INTERNALDATE "24-Jan-2021 14:10:54 +0100"
	BODYSTRUCTURE
		(
			("text" "plain" ("charset" "UTF-8") NIL NIL "7bit" 194 5 NIL NIL NIL NIL)
			("text" "html" ("charset" "UTF-8") NIL NIL "7bit" 1485 21 NIL NIL NIL NIL)
			"alternative"
			("boundary" "--==_mimepart_600d71d99b2f3_631a04201887" "charset" "UTF-8") NIL NIL NIL)

	ENVELOPE (
		"Sun, 24 Jan 2021 05:10:49 -0800" {70}
Re: [91divoc-ln/harrie-0] Add harrie 6 (work by Joshua Thijssen) (#31)
	(("Renze Nicolai" NIL "notifications" "github.com"))
	(("Renze Nicolai" NIL "notifications" "github.com"))
	(("91divoc-ln/harrie-0" NIL "reply+AAB26MXGCD2AGRWJ3KTRWWN6DFJNTEVBNHHC54BV3M" "reply.github.com"))
	(("91divoc-ln/harrie-0" NIL "harrie-0" "noreply.github.com"))
	(("Joshua Thijssen" NIL "jthijssen" "noxlogic.nl")("Review requested" NIL "review_requested" "noreply.github.com")) NIL "<91divoc-ln/harrie-0/pull/31@github.com>" "<91divoc-ln/harrie-0/pull/31/issue_event/4242487782@github.com>") BODY[HEADER.FIELDS (IMPORTANCE X-PRIORITY REFERENCES CONTENT-TYPE)] {165}
References: <91divoc-ln/harrie-0/pull/31@github.com>
Content-Type: multipart/alternative;
 boundary="--==_mimepart_600d71d99b2f3_631a04201887";
 charset=UTF-8

)



* 1 FETCH (
	UID 19
	RFC822.SIZE 12345
	FLAGS (\Unseen)
	INTERNALDATE "10-Jan-2021 17:24:55 +0000"
	BODYSTRUCTURE
		(
			("TEXT" "PLAIN" ("CHARSET" "UTF-8") NIL NIL "QUOTED-PRINTABLE" 12345 64 NIL NIL NIL NIL)
	ENVELOPE (
		"10-Jan-2021 17:24:55 +0000" "test"
	(("hello@bitmaelum.network" NIL "hello" "bitmaelum.network"))
	(("hello@bitmaelum.network" NIL "hello" "bitmaelum.network"))
	(("hello@bitmaelum.network" NIL "hello" "bitmaelum.network"))
	(("jaytaph@bitmaelum.network" NIL "jaytaph" "bitmaelum.network")) NIL NIL NIL "<e749580d-e059-43f0-b1ec-0d13006773c0@bitmaelum.network>") BODY[HEADER.FIELDS ("importance" "x-priority" "references" "content-type" )] {39}
X-Priority: 3
Content-Type: text/plain

)
 */
