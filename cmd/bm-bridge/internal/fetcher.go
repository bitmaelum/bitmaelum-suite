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

package common

import (
	"crypto/tls"
	"io/ioutil"
	"net"
	"net/mail"
	"net/smtp"
	"os"
	"strings"
	"time"

	"github.com/bitmaelum/bitmaelum-suite/internal"
	"github.com/bitmaelum/bitmaelum-suite/internal/api"
	bmmessage "github.com/bitmaelum/bitmaelum-suite/internal/message"
	"github.com/bitmaelum/bitmaelum-suite/internal/vault"
)

// Fetcher struct
type Fetcher struct {
	Account string
	Info    *vault.AccountInfo
	Client  *api.API
	Vault   *vault.Vault
}

// CheckMail will check if there is new mail waiting to be sent to the outside
func (fe *Fetcher) CheckMail() {
	msgList, err := fe.Client.GetMailboxMessages(fe.Info.Address.Hash(), "1", time.Time{})
	if err != nil {
		return
	}

	for _, message := range msgList.Messages {
		fe.sendMailForMessage(message)
	}
}

func (fe *Fetcher) sendMailForMessage(message api.MailboxMessagesMessage) error {
	// get the message
	msg, err := fe.Client.GetMessage(fe.Info.Address.Hash(), message.ID)
	if err != nil {
		return err
	}

	// Decrypt message
	em := bmmessage.EncryptedMessage{
		ID:      msg.ID,
		Header:  &msg.Header,
		Catalog: msg.Catalog,

		GenerateBlockReader:      fe.Client.GenerateAPIBlockReader(fe.Info.Address.Hash()),
		GenerateAttachmentReader: fe.Client.GenerateAPIAttachmentReader(fe.Info.Address.Hash()),
	}

	dm, err := em.Decrypt(fe.Info.GetActiveKey().PrivKey)
	if err != nil {
		return err
	}

	// Check if there is a mimeparts and destination block because we need that to reconstruct the mime message
	for _, block := range dm.Catalog.Blocks {
		if block.Type == "destination" {
			// Use this block as destination address
			recipientAddress := make([]byte, block.Size)
			if block.Reader != nil {
				recipientAddress, _ = ioutil.ReadAll(block.Reader)
			}

			if err := processMIMEMessage(string(recipientAddress), dm.Catalog, message.ID); err == nil {
				// Delete the message
				fe.Client.RemoveMessage(fe.Info.Address.Hash(), message.ID)
			}

			return err
		}
	}

	// Delete the message
	return fe.Client.RemoveMessage(fe.Info.Address.Hash(), message.ID)
}

func processMIMEMessage(toMail string, catalog *bmmessage.Catalog, msgID string) error {
	mimeMsg := &MimeMessage{
		From: &mail.Address{
			Name:    catalog.From.Name,
			Address: AddrToEmail(catalog.From.Address),
		},

		To: []*mail.Address{{
			Address: toMail,
		}},

		ID:      "<" + msgID + "@bitmaelum.network>",
		Subject: catalog.Subject,
		Date:    catalog.CreatedAt,
	}

	for _, block := range catalog.Blocks {
		blockContent := make([]byte, block.Size)

		if block.Reader != nil {
			blockContent, _ = ioutil.ReadAll(block.Reader)
		}

		mimeMsg.Blocks = append(mimeMsg.Blocks, block.Type+","+string(blockContent))
	}

	mimeMsg.Attachments = make(map[string][]byte)
	for _, attachment := range catalog.Attachments {
		attachContent := make([]byte, attachment.Size)

		r := attachment.Reader
		if attachment.Reader != nil {
			attachContent, _ = ioutil.ReadAll(r)
		}

		// Create an attachment
		mimeMsg.Attachments[attachment.FileName] = internal.Encode(attachContent)
	}

	var err error
	if mime, err := mimeMsg.EncodeToMime(); err == nil {
		// Send the mail to the main MX
		parts := strings.Split(toMail, "@")
		mx, err := net.LookupMX(parts[1])
		if err != nil || len(mx) == 0 {
			return err
		}

		return sendSMTP(mx[0].Host, mimeMsg.From.Address, []string{toMail}, mime)
	}

	return err
}

func sendSMTP(host, from string, to []string, msg []byte) error {

	c, err := smtp.Dial(host + ":25")
	if err != nil {
		return err
	}
	defer c.Close()
	hostname, _ := os.Hostname()
	c.Hello(hostname)

	if ok, _ := c.Extension("STARTTLS"); ok {
		config := &tls.Config{ServerName: host}
		if err = c.StartTLS(config); err != nil {
			return err
		}
	}

	if err = c.Mail(from); err != nil {
		return err
	}
	for _, addr := range to {
		if err = c.Rcpt(addr); err != nil {
			return err
		}
	}
	w, err := c.Data()
	if err != nil {
		return err
	}
	_, err = w.Write(msg)
	if err != nil {
		return err
	}
	err = w.Close()
	if err != nil {
		return err
	}
	return c.Quit()
}
