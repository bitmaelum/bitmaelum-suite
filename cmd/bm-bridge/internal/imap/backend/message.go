// Copyright (c) 2022 BitMaelum Authors
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

package imapgw

import (
	"bufio"
	"bytes"
	"io"
	"io/ioutil"
	netmail "net/mail"
	"time"

	common "github.com/bitmaelum/bitmaelum-suite/cmd/bm-bridge/internal"
	"github.com/bitmaelum/bitmaelum-suite/internal"
	bmmessage "github.com/bitmaelum/bitmaelum-suite/internal/message"
	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/backend/backendutil"
	"github.com/emersion/go-message"
	"github.com/emersion/go-message/textproto"
	"github.com/sirupsen/logrus"
)

// Message holds all the information needed to retrieve a full message
type Message struct {
	UID     uint32
	ID      string
	Date    time.Time
	Size    uint32
	Flags   []string
	DMsg    *bmmessage.DecryptedMessage
	User    *User
	FullMsg []byte
}

// Build will create the message to be read
func (m *Message) Build() error {
	mimeMsg := &common.MimeMessage{
		From: &netmail.Address{
			Name:    m.DMsg.Catalog.From.Name,
			Address: common.AddrToEmail(m.DMsg.Catalog.From.Address),
		},

		To: []*netmail.Address{{
			Address: common.AddrToEmail(m.DMsg.Catalog.To.Address),
		}},

		ID:      "<" + m.ID + "@bitmaelum.network>",
		Subject: m.DMsg.Catalog.Subject,
		Date:    m.DMsg.Catalog.CreatedAt,
	}

	for _, block := range m.DMsg.Catalog.Blocks {
		blockContent := make([]byte, block.Size)

		if block.Reader != nil {
			blockContent, _ = ioutil.ReadAll(block.Reader)
		}

		mimeMsg.Blocks = append(mimeMsg.Blocks, block.Type+","+string(blockContent))
	}

	mimeMsg.Attachments = make(map[string][]byte)
	for _, attachment := range m.DMsg.Catalog.Attachments {
		attachContent := make([]byte, attachment.Size)

		r := attachment.Reader
		if attachment.Reader != nil {
			attachContent, _ = ioutil.ReadAll(r)
		}

		// Create an attachment
		mimeMsg.Attachments[attachment.FileName] = internal.Encode(attachContent)
	}

	var err error
	m.FullMsg, err = mimeMsg.EncodeToMime()
	if err != nil {
		return err
	}

	m.Size = uint32(9999 + len(m.FullMsg))

	return nil
}

func (m *Message) entity() (*message.Entity, error) {
	if m.DMsg == nil {
		err := m.Decrypt()
		if err != nil {
			return nil, err
		}
	}

	return message.Read(bytes.NewReader(m.FullMsg))
}

func (m *Message) headerAndBody() (textproto.Header, io.Reader, error) {
	body := bufio.NewReader(bytes.NewReader(m.FullMsg))
	hdr, err := textproto.ReadHeader(body)
	hdr.Add("Message-ID", "<"+m.DMsg.ID+"@bitmaelum.network/>")
	return hdr, body, err
}

func (m *Message) headerOnly() (textproto.Header, io.Reader, error) {
	body := bufio.NewReader(bytes.NewReader(m.FullMsg))

	hdr, err := textproto.ReadHeader(body)
	hdr.Add("Message-ID", "<"+m.DMsg.ID+"@bitmaelum.network/>")
	return hdr, body, err
}

// Fetch will retrieve a message
func (m *Message) Fetch(seqNum uint32, items []imap.FetchItem, user *User) (*imap.Message, error) {

	if m.DMsg == nil {
		err := m.Decrypt()
		if err != nil {
			return nil, err
		}
	}

	fetched := imap.NewMessage(seqNum, items)
	for _, item := range items {
		switch item {
		case imap.FetchEnvelope:
			hdr, _, _ := m.headerOnly()
			fetched.Envelope, _ = backendutil.FetchEnvelope(hdr)
		case imap.FetchBody, imap.FetchBodyStructure:
			hdr, body, _ := m.headerAndBody()
			fetched.BodyStructure, _ = backendutil.FetchBodyStructure(hdr, body, item == imap.FetchBodyStructure)
			logrus.Infof("IMAP: fetching message %s", m.ID)
		case imap.FetchFlags:
			fetched.Flags = m.Flags
		case imap.FetchInternalDate:
			fetched.InternalDate = m.Date
		case imap.FetchRFC822Size:
			fetched.Size = m.Size
		case imap.FetchUid:
			fetched.Uid = m.UID
		default:
			section, err := imap.ParseBodySectionName(item)
			if err != nil {
				break
			}

			hdr, body, _ := m.headerAndBody()

			l, _ := backendutil.FetchBodySection(hdr, body, section)
			fetched.Body[section] = l
		}
	}

	return fetched, nil
}

// Match checks if the message is the same
func (m *Message) Match(seqNum uint32, c *imap.SearchCriteria) (bool, error) {
	e, err := m.entity()
	if err != nil {
		return false, err
	}
	return backendutil.Match(e, seqNum, m.UID, m.Date, m.Flags, c)
}

// Decrypt a message and store it decrypted
func (m *Message) Decrypt() error {
	msg, err := m.User.Client.GetMessage(m.User.Info.Address.Hash(), m.ID)
	if err != nil {
		return err
	}

	// Decrypt message
	em := bmmessage.EncryptedMessage{
		ID:      msg.ID,
		Header:  &msg.Header,
		Catalog: msg.Catalog,

		GenerateBlockReader:      m.User.Client.GenerateAPIBlockReader(m.User.Info.Address.Hash()),
		GenerateAttachmentReader: m.User.Client.GenerateAPIAttachmentReader(m.User.Info.Address.Hash()),
	}

	m.DMsg, err = em.Decrypt(m.User.Info.GetActiveKey().PrivKey)
	if err != nil {
		return err
	}

	m.Date = m.DMsg.Catalog.CreatedAt

	m.Build()

	return nil
}
