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

package bitmaelumClient

import (
	"io"
	"io/ioutil"
	"os"
	"path"
	"sort"
	"strconv"
	"time"

	"github.com/bitmaelum/bitmaelum-suite/internal/api"
	"github.com/bitmaelum/bitmaelum-suite/internal/mailbox"
	"github.com/bitmaelum/bitmaelum-suite/internal/message"
	"github.com/bitmaelum/bitmaelum-suite/internal/messages"
	"github.com/bitmaelum/bitmaelum-suite/pkg/address"
	"github.com/bitmaelum/bitmaelum-suite/pkg/bmcrypto"
	"github.com/pkg/errors"
)

type Message struct {
	ID          string                   `json:"id"`
	Date        time.Time                `json:"date"`
	Subject     string                   `json:"subject"`
	FromAddress string                   `json:"from_address"`
	FromName    string                   `json:"from_name"`
	Attachments []message.AttachmentType `json:"attachments"`
	Blocks      []message.BlockType      `json:"blocks"`
	SignedBy    message.SignedByType     `json:"signed_by"`
}

func (b *BitMaelumClient) SendSimpleMessage(to, subject, body string) error {
	// Setup blocks
	blocks := map[string]string{"default": body}

	return b.SendMessage(to, subject, blocks, nil)
}

func (b *BitMaelumClient) SendMessage(to, subject string, blocks map[string]string, attachments []string) error {
	senderInfo, _ := b.resolverService.ResolveAddress(b.user.Address.Hash())

	// Check recipient
	toAddr, err := address.NewAddress(to)
	if err != nil {
		return errors.Wrap(err, "check recipient address")
	}

	recipientInfo, err := b.resolverService.ResolveAddress(toAddr.Hash())
	if err != nil {
		return errors.Wrap(err, "resolve recipient address")
	}

	// Setup addressing
	addressing := message.NewAddressing(message.SignedByTypeOrigin)
	addressing.AddSender(b.user.Address, nil, b.user.Name, *b.user.PrivateKey, senderInfo.RoutingInfo.Routing)
	addressing.AddRecipient(toAddr, nil, &recipientInfo.PublicKey)

	/*
		// Setup attachments
		var mAttachments []string
		for fileName, fileData := range attachments {
			// We write the attachments temporary to disk so we can use it later on message.Compose,
			// however this needs to be improved so we don't need to write them to disk
			fName := filepath.Join(os.TempDir(), fileName)
			err = ioutil.WriteFile(fName, fileData, 0644)
			if err != nil {
				return err
			}

			defer os.Remove(fName)

			mAttachments = append(mAttachments, fName)
		}
	*/

	// Setup blocks
	var mBlocks []string
	for t, d := range blocks {
		mBlocks = append(mBlocks, t+","+d)
	}

	// Compose mail
	envelope, err := message.Compose(addressing, subject, mBlocks, attachments)
	if err != nil {
		return errors.Wrap(err, "composing mail")
	}

	// Send mail
	client, err := api.NewAuthenticated(*b.user.Address, *b.user.PrivateKey, senderInfo.RoutingInfo.Routing, nil)
	if err != nil {
		return errors.Wrap(err, "setting api")
	}

	err = messages.Send(*client, envelope)
	if err != nil {
		return errors.Wrap(err, "sending message")
	}

	return nil
}

func (b *BitMaelumClient) ListMessages(since time.Time, boxID int) ([]Message, error) {

	client, err := api.NewAuthenticated(*b.user.Address, *b.user.PrivateKey, b.user.RoutingInfo.Routing, nil)
	if err != nil {
		return nil, err
	}

	msgList, err := client.GetMailboxMessages(b.user.Address.Hash(), strconv.Itoa(boxID), since)
	if err != nil {
		return nil, err
	}

	// Sort messages first
	msort := mailbox.NewMessageSort(*b.user.PrivateKey, msgList.Messages, mailbox.SortDate, true)
	sort.Sort(&msort)

	var messages []Message
	for _, msg := range msgList.Messages {
		var newMessage Message

		key, _ := bmcrypto.Decrypt(*b.user.PrivateKey, msg.Header.Catalog.TransactionID, msg.Header.Catalog.EncryptedKey)
		cat := &message.Catalog{}
		err = bmcrypto.CatalogDecrypt(key, msg.Catalog, cat)
		if err != nil {
			// The message could not be decrypted
			newMessage = Message{
				ID:          msg.ID,
				Subject:     "[UNABLE TO DECRYPT]",
				FromAddress: msg.Header.From.Addr.String(),
				SignedBy:    msg.Header.From.SignedBy,
			}
		} else {
			newMessage = Message{
				ID:          msg.ID,
				Subject:     cat.Subject,
				FromAddress: cat.From.Address,
				FromName:    cat.From.Name,
				Date:        cat.CreatedAt,
				Attachments: cat.Attachments,
				Blocks:      cat.Blocks,
				SignedBy:    msg.Header.From.SignedBy,
			}
		}

		messages = append(messages, newMessage)
	}

	return messages, nil
}

func (b *BitMaelumClient) ReadBlock(msgID, blockID string) ([]byte, error) {
	client, err := api.NewAuthenticated(*b.user.Address, *b.user.PrivateKey, b.user.RoutingInfo.Routing, nil)
	if err != nil {
		return nil, err
	}

	msg, err := client.GetMessage(b.user.Address.Hash(), msgID)
	if err != nil {
		return nil, err
	}

	// Decrypt message
	em := message.EncryptedMessage{
		ID:      msg.ID,
		Header:  &msg.Header,
		Catalog: msg.Catalog,

		GenerateBlockReader: client.GenerateAPIBlockReader(b.user.Address.Hash()),
	}

	decMsg, err := em.Decrypt(*b.user.PrivateKey)
	if err != nil {
		return nil, err
	}

	for _, block := range decMsg.Catalog.Blocks {
		if block.ID != blockID {
			continue
		}

		if block.Reader != nil {
			return ioutil.ReadAll(block.Reader)
		}

	}

	return nil, errors.New("block not found")
}

func (b *BitMaelumClient) SaveAttachment(msgID, attachmentID, path string, overwrite bool) (interface{}, error) {
	client, err := api.NewAuthenticated(*b.user.Address, *b.user.PrivateKey, b.user.RoutingInfo.Routing, nil)
	if err != nil {
		return nil, err
	}

	msg, err := client.GetMessage(b.user.Address.Hash(), msgID)
	if err != nil {
		return nil, err
	}

	// Decrypt message
	em := message.EncryptedMessage{
		ID:      msg.ID,
		Header:  &msg.Header,
		Catalog: msg.Catalog,

		GenerateAttachmentReader: client.GenerateAPIAttachmentReader(b.user.Address.Hash()),
	}

	decMsg, err := em.Decrypt(*b.user.PrivateKey)
	if err != nil {
		return nil, err
	}

	for _, attachment := range decMsg.Catalog.Attachments {
		if attachment.ID != attachmentID {
			continue
		}

		return saveAttachment(attachment, path, overwrite)
	}

	return nil, errors.New("attachment not found")
}

func saveAttachment(att message.AttachmentType, savePath string, overwrite bool) (interface{}, error) {
	defer func() {
		// Close stream if it's closeable
		_, ok := att.Reader.(io.Closer)
		if ok {
			_ = att.Reader.(io.Closer).Close()
		}
	}()

	destFile := path.Join(savePath, att.FileName)

	_, ok := os.Stat(destFile)
	if ok == nil {
		if !overwrite {
			return nil, errors.New("cannot write, file exists")
		}
	}

	f, err := os.Create(destFile)
	if err != nil {
		return nil, err
	}

	n, err := io.Copy(f, att.Reader)
	if err != nil || n != int64(att.Size) {
		return nil, err
	}

	_ = f.Close()

	return map[string]interface{}{
		"path": destFile,
	}, nil
}
