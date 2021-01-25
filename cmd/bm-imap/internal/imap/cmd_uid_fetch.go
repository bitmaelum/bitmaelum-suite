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
	attrs := ParseAttributes(strings.Join(args[2:], " "), true)

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

		c.Write("*", fmt.Sprintf("%d FETCH (%s)", uid, strings.Join(ret, " ")))
		i++
	}
	c.Write(tag, "OK FETCH completed")

	return nil
}

func generateHeaderMap(msgInfo internal.MessageInfo, dMsg message.DecryptedMessage) map[string]string {
	ret := make(map[string]string)

	ret["x-priority"] = "3"
	ret["date"] = dMsg.Catalog.CreatedAt.Format(TimeFormat)
	ret["subject"] = dMsg.Catalog.Subject
	ret["from"] = addrToEmail(dMsg.Catalog.From.Address)
	ret["content-type"] = "text/plain"
	ret["to"] = addrToEmail(dMsg.Catalog.To.Address)
	ret["message-id"] = "<" + msgInfo.MessageID + "@bitmaelum.network>"

	return ret
}

func executeAttributes(attrs []Attribute, msgInfo internal.MessageInfo, dMsg message.DecryptedMessage) []string {
	var ret []string

	from := addrToEmail(dMsg.Catalog.From.Address)
	to := addrToEmail(dMsg.Catalog.To.Address)
	header := generateHeaderMap(msgInfo, dMsg)

	fromParts := strings.Split(from, "@")
	toParts := strings.Split(to, "@")

	for _, attr := range attrs {
		switch attr.Name {
		case "RFC822":
			switch attr.SubName {
			case "SIZE":
				ret = append(ret, fmt.Sprintf("RFC822.SIZE %d", 200))
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
			ret = append(ret, "BODYSTRUCTURE (\"TEXT\" \"PLAIN\" (\"CHARSET\" \"UTF-8\") NIL NIL \"7bit\" 200 4 NIL NIL NIL NIL)")

		case "UID":
			ret = append(ret, fmt.Sprintf("UID %d", msgInfo.UID))

		case "FLAGS":
			//ret = append(ret, fmt.Sprintf("FLAGS (%s)", strings.Join(msgInfo.Flags, " ")))
			ret = append(ret, "FLAGS (\\Recent)")

		case "INTERNALDATE":
			ret = append(ret, fmt.Sprintf("INTERNALDATE \"%s\"", dMsg.Catalog.CreatedAt.Format(TimeFormat)))

		case "BODY":
			switch attr.Section {
			case "HEADER.FIELDS":
				hdrs := ""
				for _, h := range attr.Headers {
					v, ok := header[h]
					if ok {
						hdrs += strings.ToUpper(h[0:0]) + h[1:] + ": " + v + "\n"
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

				s := string(body) + "\n\n"
				ret = append(ret, fmt.Sprintf("%s {%d}\n%s", attr.ToString(), len(s), s))
			case "MIME":
			case "HEADER":
				hdrs := ""
				for i := range header {
					hdrs += i + ": " + header[i] + "\n"
				}
				hdrs += "\n"

				ret = append(ret, fmt.Sprintf("%s {%d}\n%s", attr.ToString(), len(hdrs), hdrs))
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
