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
	"github.com/bitmaelum/bitmaelum-suite/internal/config"
	"github.com/bitmaelum/bitmaelum-suite/internal/container"
	"github.com/bitmaelum/bitmaelum-suite/internal/message"
	"github.com/bitmaelum/bitmaelum-suite/internal/messages"
	"github.com/bitmaelum/bitmaelum-suite/internal/vault"
	"github.com/bitmaelum/bitmaelum-suite/pkg/address"
	"github.com/sirupsen/logrus"
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

func (fe *Fetcher) sendMailForMessage(mailboxMessage api.MailboxMessagesMessage) error {
	// get the message
	logrus.Infof("processing incoming message %s", mailboxMessage.ID)
	msg, err := fe.Client.GetMessage(fe.Info.Address.Hash(), mailboxMessage.ID)
	if err != nil {
		return err
	}

	// Decrypt message
	em := message.EncryptedMessage{
		ID:      msg.ID,
		Header:  &msg.Header,
		Catalog: msg.Catalog,

		GenerateBlockReader:      fe.Client.GenerateAPIBlockReader(fe.Info.Address.Hash()),
		GenerateAttachmentReader: fe.Client.GenerateAPIAttachmentReader(fe.Info.Address.Hash()),
	}

	dm, err := em.Decrypt(fe.Info.GetActiveKey().PrivKey)
	if err != nil {
		logrus.Debugf("error while decrypting message - %s", err.Error())
		return err
	}

	// Check if there is a mimeparts and destination block because we need that to reconstruct the mime message
	for _, block := range dm.Catalog.Blocks {
		if block.Type == DestinationBlock {
			// Use this block as destination address
			recipientAddress := make([]byte, block.Size)
			if block.Reader != nil {
				recipientAddress, _ = ioutil.ReadAll(block.Reader)
			}

			if err := processMIMEMessage(string(recipientAddress), dm.Catalog, mailboxMessage.ID); err != nil {
				logrus.Infof("error delivering message %s - %v", mailboxMessage.ID, err)
				err = fe.sendPostmasterResponse(dm.Catalog.From.Address, dm.Catalog, err)
				if err != nil {
					logrus.Debugf("error sending postmaster - %v", err)
				}
			} else {
				logrus.Infof("message %s processed, sent to %s", mailboxMessage.ID, recipientAddress)
			}

			// Delete the message
			fe.Client.RemoveMessage(fe.Info.Address.Hash(), mailboxMessage.ID)

			return err
		}
	}

	logrus.Infof("the message %s does not contain a \"%s\" block. ignoring", mailboxMessage.ID, DestinationBlock)
	// Delete the message
	return fe.Client.RemoveMessage(fe.Info.Address.Hash(), mailboxMessage.ID)
}

func processMIMEMessage(toMail string, catalog *message.Catalog, msgID string) error {
	mimeMsg := &MimeMessage{
		From: &mail.Address{
			Name:    catalog.From.Name,
			Address: AddrToEmail(catalog.From.Address),
		},

		To: []*mail.Address{{
			Address: toMail,
		}},

		ID:      "<" + msgID + "@" + config.Bridge.Server.SMTP.Domain + ">",
		Subject: catalog.Subject,
		Date:    catalog.CreatedAt,
	}

	for _, block := range catalog.Blocks {
		if block.Type == DestinationBlock {
			// ignore the "destination" block
			continue
		}

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

		logrus.Tracef("sending mail to SMTP host %s - from: %s - to: %s", mx[0].Host, mimeMsg.From.Address, toMail)
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

func (fe *Fetcher) sendPostmasterResponse(destination string, catalog *message.Catalog, msgError error) error {
	svc := container.Instance.GetResolveService()

	senderInfo, _ := svc.ResolveAddress(fe.Info.Address.Hash())

	// Check to address
	toAddr, err := address.NewAddress(destination)
	if err != nil {
		return err
	}

	recipientInfo, err := svc.ResolveAddress(toAddr.Hash())
	if err != nil {
		return err
	}

	// Setup addressing
	addressing := message.NewAddressing(message.SignedByTypeOrigin)
	addressing.AddSender(fe.Info.Address, nil, fe.Info.Name, fe.Info.GetActiveKey().PrivKey, senderInfo.RoutingInfo.Routing)
	addressing.AddRecipient(toAddr, nil, &recipientInfo.PublicKey)

	var blocks []string
	blocks = append(blocks, "default,I was unable to send mail with subject \""+catalog.Subject+"\". The error was: "+msgError.Error())

	// Compose mail
	envelope, err := message.Compose(addressing, "Unable to deliver mail", blocks, nil)
	if err != nil {
		return err
	}

	// Send mail
	client, err := api.NewAuthenticated(*fe.Info.Address, fe.Info.GetActiveKey().PrivKey, senderInfo.RoutingInfo.Routing, nil)
	if err != nil {
		return err
	}

	err = messages.Send(*client, envelope)
	if err != nil {
		return err
	}

	return nil
}
