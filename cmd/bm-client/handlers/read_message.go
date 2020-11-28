// Copyright (c) 2020 BitMaelum Authors
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
	"bytes"
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
	"github.com/bitmaelum/bitmaelum-suite/pkg/bmcrypto"
	"github.com/c2h5oh/datasize"
	"github.com/sirupsen/logrus"
)

type messageEntryType struct {
	Box     string
	ID      string
	Header  message.Header
	Catalog message.Catalog
}

// ReadMessages will read a specific message blocks
func ReadMessages(info *vault.AccountInfo, routingInfo *resolver.RoutingInfo, box, messageID string, since time.Time) {
	client, err := api.NewAuthenticated(*info.Address, &info.PrivKey, routingInfo.Routing, internal.JwtErrorFunc)
	if err != nil {
		logrus.Fatal(err)
	}

	// generate a message list of all messages we want do display
	fmt.Print("* Fetching remote messages...")
	entryList := queryMessageEntries(client, info, box, messageID, since)
	idx := 0
	fmt.Println("")

	if len(entryList) == 0 {
		fmt.Println("*  No messages found to read.")
		return
	}

	done := false

	// iterate list until we quit
	for {
		if done {
			break
		}

		displayMessage(client, info, entryList[idx])

		for {

			// Build command string
			cmds := []string{}

			// display text based on current item
			if idx < len(entryList)-1 {
				cmds = append(cmds, "View (N)ext")
			}
			if idx > 0 {
				cmds = append(cmds, "View (P)revious")
			}
			if len(entryList[idx].Catalog.Attachments) > 0 {
				cmds = append(cmds, "(S)ave attachments")
			}
			cmds = append(cmds, "(Q)uit")

			fmt.Printf("(%d/%d): %s > ", idx+1, len(entryList), strings.Join(cmds, ", "))

			// Read and parse entry
			reader := bufio.NewReader(os.Stdin)

			ch, _ := reader.ReadByte()
			ch = strings.ToUpper(string(ch))[0]

			if ch == 'P' && idx > 0 {
				idx--
				break
			}
			if ch == 'N' && idx < len(entryList)-1 {
				idx++
				break
			}
			if ch == 'S' {
				saveAttachments(client, info, entryList[idx])
			}
			if ch == 'Q' {
				done = true
				break
			}
		}
	}
}

func queryMessageEntries(client *api.API, info *vault.AccountInfo, boxID, msgID string, since time.Time) []messageEntryType {
	ret := []messageEntryType{}

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

		for _, msg := range mb.Messages {
			// skip if we only want specific msg ID's
			if mode == "msgid" && !strings.HasPrefix(msg.ID, msgID) {
				continue
			}

			catalog, err := decryptCatalog(info, msg)
			if err != nil {
				continue
			}

			ret = append(ret, messageEntryType{
				ID:      msg.ID,
				Box:     strconv.Itoa(box.ID),
				Header:  msg.Header,
				Catalog: *catalog,
			})
		}
	}

	return ret
}

// decrypt
func decryptCatalog(info *vault.AccountInfo, msg api.MailboxMessagesMessage) (*message.Catalog, error) {
	key, err := bmcrypto.Decrypt(info.PrivKey, msg.Header.Catalog.TransactionID, msg.Header.Catalog.EncryptedKey)
	if err != nil {
		return nil, err
	}

	// Verify the clientSignature
	if !message.VerifyClientHeader(msg.Header) {
		return nil, errors.New("invalid client signature")
	}

	// Decrypt the catalog
	catalog := &message.Catalog{}
	err = bmcrypto.CatalogDecrypt(key, msg.Catalog, catalog)
	if err != nil {
		return nil, errors.New("cannot decrypt")
	}

	return catalog, nil
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

func displayMessage(client *api.API, info *vault.AccountInfo, entry messageEntryType) {
	fmt.Printf("--------------------------------------------------------\n")
	fmt.Printf("From       : %s <%s>\n", entry.Catalog.From.Name, entry.Catalog.From.Address)
	fmt.Printf("To         : %s <%s>\n", entry.Catalog.To.Name, entry.Catalog.To.Address)
	fmt.Printf("Subject    : %s\n", entry.Catalog.Subject)
	fmt.Printf("\n")
	fmt.Printf("Msg ID     : %s\n", entry.ID)
	fmt.Printf("Created at : %s\n", entry.Catalog.CreatedAt)
	fmt.Printf("ThreadID   : %s\n", entry.Catalog.ThreadID)
	fmt.Printf("Flags      : %s\n", entry.Catalog.Flags)
	fmt.Printf("Labels     : %s\n", entry.Catalog.Labels)
	fmt.Println("--------------------------------------------------------")
	for idx, b := range entry.Catalog.Blocks {
		fmt.Printf("Block %02d: %-20s %8s\n", idx, b.Type, datasize.ByteSize(b.Size))
		fmt.Printf("\n")

		data, err := client.GetMessageBlock(info.Address.Hash(), entry.Box, entry.ID, b.ID)
		if err != nil {
			continue
		}
		bb := bytes.NewBuffer(data)

		r, err := bmcrypto.GetAesDecryptorReader(b.IV, b.Key, bb)
		if err != nil {
			panic(err)
		}

		content, err := ioutil.ReadAll(r)
		if err != nil {
			logrus.Fatal(err)
		}

		fmt.Println(string(content))
	}
	if len(entry.Catalog.Attachments) > 0 {
		fmt.Println("--------------------------------------------------------")
		for idx, a := range entry.Catalog.Attachments {
			fmt.Printf("Attachment %02d: %30s %8d %s\n", idx, a.FileName, datasize.ByteSize(a.Size), a.MimeType)
		}
	}
	fmt.Println("--------------------------------------------------------")
}

// save attachments
func saveAttachments(client *api.API, info *vault.AccountInfo, entry messageEntryType) {
	for _, att := range entry.Catalog.Attachments {
		attReader, err := client.GetMessageAttachment(info.Address.Hash(), entry.Box, entry.ID, att.ID)
		if err != nil {
			continue
		}

		err = saveAttachment(att, attReader)
		if err != nil {
			fmt.Println(err)
		}
	}
}

// saveAttachment will save the given attachment to disk
func saveAttachment(att message.AttachmentType, ar io.ReadCloser) error {
	defer func() {
		_ = ar.Close()
	}()

	_, ok := os.Stat(att.FileName)
	if ok == nil {
		fmt.Printf("cannot write to %s: file exists\n", att.FileName)
		return errors.New("cannot write to file")
	}

	r, err := bmcrypto.GetAesDecryptorReader(att.IV, att.Key, ar)
	if err != nil {
		fmt.Printf("cannot create decryptor to %s: %s\n", att.FileName, err)
		return err
	}

	f, err := os.Create(att.FileName)
	if err != nil {
		fmt.Printf("cannot open file %s: %s\n", att.FileName, err)
		return err
	}

	if att.Compression == "zlib" {
		r, err = message.ZlibDecompress(r)
		if err != nil {
			fmt.Printf("error while creating zlib reader %s: %s\n", att.FileName, err)
			return err
		}
	}

	n, err := io.Copy(f, r)
	if err != nil || n != int64(att.Size) {
		fmt.Printf("error while writing file %s: %s (%d/%d bytes)\n", att.FileName, err, n, att.Size)
		return err
	}

	_ = f.Close()
	fmt.Printf("saved file: %s\n", att.FileName)

	return nil
}
