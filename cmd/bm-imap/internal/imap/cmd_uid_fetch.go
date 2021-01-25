package imap

import (
	"fmt"
	"io"
	"io/ioutil"
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
				// Check if we have a maxrange, if so, only read until that number
				r := dMsg.Catalog.Blocks[0].Reader
				if attr.MaxRange > 0 {
					r = io.LimitReader(r, int64(attr.MaxRange))
				}

				body, err := ioutil.ReadAll(r)
				if err != nil {
					continue
				}

				// Remove everything below minrange
				if attr.MinRange > 0 {
					body = body[attr.MinRange:]
				}

				ret = append(ret, fmt.Sprintf("%s {%d}\n%s", attr.ToString(), len(body), body))
			case "MIME":
			case "HEADER":
				body := "From: joshua@bitmaelum.network\n"

				ret = append(ret, fmt.Sprintf("%s {%d}\n%s", attr.ToString(), len(body), body))
			}


		}
	}

	return ret
}

func addrToEmail(address string) string {
	address = strings.Replace(address, "@", "_", -1)
	address = strings.Replace(address, "!", "", -1)

	return address + "@bitmaelum.network"

}
