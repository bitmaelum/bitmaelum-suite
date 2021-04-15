// Copyright (c) 2021 BitMaelum Authors
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of
// this software and associated documentation files (the "Software"), to deal in
// the Software without restriction, including without limitation the rights to
// use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
// the Software, and to permit persons to whom the Software is furnished to do so,
// subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
// FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
// COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
// IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
// CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package handlers

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/bitmaelum/bitmaelum-suite/cmd/bm-client/internal"
	"github.com/bitmaelum/bitmaelum-suite/internal/api"
	"github.com/bitmaelum/bitmaelum-suite/internal/message"
	"github.com/bitmaelum/bitmaelum-suite/internal/resolver"
	"github.com/bitmaelum/bitmaelum-suite/internal/vault"
	"github.com/c2h5oh/datasize"
	"github.com/sirupsen/logrus"
)

// ReadMessages will read a specific message blocks
func ReadMessages(info *vault.AccountInfo, routingInfo *resolver.RoutingInfo, box, messageID string, since time.Time) {
	client, err := api.NewAuthenticated(*info.Address, info.GetActiveKey().PrivKey, routingInfo.Routing, internal.JwtErrorFunc)
	if err != nil {
		logrus.Fatal(err)
	}

	// generate a message list of all messages we want do display
	fmt.Print("* Fetching messages from remote server(s)...")
	messageList := queryMessages(client, info, box, messageID, since)
	fmt.Println("")

	if len(messageList) == 0 {
		fmt.Println("*  No messages found to read.")
		return
	}

	readDone := false
	idx := 0

	// iterate list until we quit
	for !readDone {
		// Find key, or iterate all keys when key is not found
		keys := info.Keys
		key, err := info.FindKey(messageList[idx].Header.To.Fingerprint)
		if err == nil {
			keys = []vault.KeyPair{*key}
		}

		// Iterate all keys (or the single key), and see if we can decrypt the message
		var decryptedMsg *message.DecryptedMessage
		for _, key := range keys {
			var err error
			decryptedMsg, err = messageList[idx].Decrypt(key.PrivKey)
			if err == nil {
				break
			}
			decryptedMsg = nil
		}

		// Display message or fail when not decrypted
		if decryptedMsg != nil {
			displayMessage(*decryptedMsg)
		} else {
			fmt.Println("Cannot decrypt message: ", err)
		}

		parseDone := false
		for !parseDone {
			readDone, parseDone = parseCommands(&idx, messageList, decryptedMsg)
		}
	}
}

func parseCommands(idx *int, messageList []message.EncryptedMessage, decryptedMsg *message.DecryptedMessage) (bool, bool) {
	// Build command string
	var cmds []string

	if *idx < len(messageList)-1 {
		cmds = append(cmds, "View (N)ext")
	}
	if *idx > 0 {
		cmds = append(cmds, "View (P)revious")
	}
	if decryptedMsg != nil && len(decryptedMsg.Catalog.Attachments) > 0 {
		cmds = append(cmds, "(S)ave attachments")
	}
	cmds = append(cmds, "(Q)uit")

	if len(messageList) > 1 || len(decryptedMsg.Catalog.Attachments) > 0 {
		fmt.Printf("(%d/%d): %s > ", *idx+1, len(messageList), strings.Join(cmds, ", "))
	} else {
		return true, true
	}

	// Read and parse entry
	reader := bufio.NewReader(os.Stdin)
	ch, _ := reader.ReadByte()
	ch = strings.ToUpper(string(ch))[0]

	// Process commands
	if ch == 'P' && *idx > 0 {
		*idx--
		return false, true
	}
	if ch == 'N' && *idx < len(messageList)-1 {
		*idx++
		return false, true
	}
	if ch == 'S' && decryptedMsg != nil {
		saveAttachments(*decryptedMsg)
	}
	if ch == 'Q' {
		return true, true
	}

	return false, false
}

func queryMessages(client *api.API, info *vault.AccountInfo, boxID, msgID string, since time.Time) []message.EncryptedMessage {
	var ret []message.EncryptedMessage

	// All 4 modes (box, msgid, since and all) are squished inside one single iteration loop. A bit more complex code, but
	// we don't need to duplicate the code 4 times over with just minor tweaks
	mode := getQueryMode(boxID, msgID, since)

	// Make sure we reset since if we don't use that mode
	if mode != "since" {
		since = time.Time{}
	}

	// Get all mailboxes
	mbl, err := client.GetMailboxList(info.Address.Hash())
	if err != nil {
		return nil
	}

	for _, box := range mbl.Boxes {
		// Skip if box mode and not our box id
		if mode == "box" && strconv.Itoa(box.ID) != boxID {
			continue
		}

		mb, err := client.GetMailboxMessages(info.Address.Hash(), strconv.Itoa(box.ID), since)
		if err != nil {
			continue
		}

		for idx := range mb.Messages {
			// skip if we only want specific msg ID's
			if mode == "msgid" && !strings.HasPrefix(mb.Messages[idx].ID, msgID) {
				continue
			}

			em := message.EncryptedMessage{
				ID:      mb.Messages[idx].ID,
				Header:  &mb.Messages[idx].Header,
				Catalog: mb.Messages[idx].Catalog,

				GenerateBlockReader:      client.GenerateAPIBlockReader(info.Address.Hash()),
				GenerateAttachmentReader: client.GenerateAPIAttachmentReader(info.Address.Hash()),
			}

			ret = append(ret, em)
		}
	}

	return ret
}

func getQueryMode(boxID string, msgID string, since time.Time) string {
	if !since.IsZero() {
		return "since"
	}

	if msgID != "" {
		return "msgid"
	}

	if boxID != "" {
		return "box"
	}

	return "all"
}

func displayMessage(msg message.DecryptedMessage) {
	fmt.Printf("--------------------------------------------------------\n")
	fmt.Printf("From       : %s <%s>\n", msg.Catalog.From.Name, msg.Catalog.From.Address)
	fmt.Printf("To         : %s\n", msg.Catalog.To.Address)
	fmt.Printf("Subject    : %s\n", msg.Catalog.Subject)
	fmt.Printf("\n")
	fmt.Printf("Msg ID     : %s\n", msg.ID)
	fmt.Printf("Created at : %s\n", msg.Catalog.CreatedAt)
	fmt.Printf("ThreadID   : %s\n", msg.Catalog.ThreadID)
	fmt.Printf("Flags      : %s\n", msg.Catalog.Flags)
	fmt.Printf("Labels     : %s\n", msg.Catalog.Labels)
	fmt.Println("--------------------------------------------------------")
	for idx, b := range msg.Catalog.Blocks {
		fmt.Printf("Block %02d: %-20s %8s\n", idx, b.Type, datasize.ByteSize(b.Size))
		fmt.Printf("\n")

		content, err := ioutil.ReadAll(msg.Catalog.Blocks[idx].Reader)
		if err != nil {
			logrus.Fatal(err)
		}

		fmt.Println(string(content))
	}

	if len(msg.Catalog.Attachments) > 0 {
		fmt.Println("--------------------------------------------------------")
		for idx, a := range msg.Catalog.Attachments {
			fmt.Printf("Attachment %02d: %30s %8d %s\n", idx, a.FileName, datasize.ByteSize(a.Size), a.MimeType)
		}
	}
	fmt.Println("--------------------------------------------------------")
}

// save attachments
func saveAttachments(msg message.DecryptedMessage) {
	for _, att := range msg.Catalog.Attachments {
		err := saveAttachment(att)
		if err != nil {
			fmt.Println(err)
		}
	}
}

// saveAttachment will save the given attachment to disk
func saveAttachment(att message.AttachmentType) error {
	defer func() {
		// Close stream if it's closeable
		_, ok := att.Reader.(io.Closer)
		if ok {
			_ = att.Reader.(io.Closer).Close()
		}
	}()

	_, ok := os.Stat(att.FileName)
	if ok == nil {
		fmt.Printf("cannot write to %s: file exists\n", att.FileName)
		return errors.New("cannot write to file")
	}

	f, err := os.Create(att.FileName)
	if err != nil {
		fmt.Printf("cannot open file %s: %s\n", att.FileName, err)
		return err
	}

	n, err := io.Copy(f, att.Reader)
	if err != nil || n != int64(att.Size) {
		fmt.Printf("error while writing file %s: %s (%d/%d bytes)\n", att.FileName, err, n, att.Size)
		return err
	}

	_ = f.Close()
	fmt.Printf("saved file: %s\n", att.FileName)

	return nil
}
