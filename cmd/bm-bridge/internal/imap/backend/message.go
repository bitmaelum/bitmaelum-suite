package imapgw

import (
	"bufio"
	"bytes"
	"io"
	"io/ioutil"
	"time"

	common "github.com/bitmaelum/bitmaelum-suite/cmd/bm-bridge/internal"
	bmmessage "github.com/bitmaelum/bitmaelum-suite/internal/message"
	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/backend/backendutil"
	"github.com/emersion/go-message"
	"github.com/emersion/go-message/textproto"
	"github.com/jhillyerd/enmime"
)

type Message struct {
	Uid     uint32
	ID      string
	Date    time.Time
	Size    uint32
	Flags   []string
	Builder enmime.MailBuilder
	DMsg    *bmmessage.DecryptedMessage
	User    *User
	FullMsg []byte
}

func (m *Message) Build(realContent bool) []byte {

	mimeMsg := m.Builder

	if len(m.FullMsg) > 0 {
		return m.FullMsg
	}

	for _, block := range m.DMsg.Catalog.Blocks {
		blockContent := make([]byte, block.Size)

		if realContent {
			if block.Reader != nil {
				blockContent, _ = ioutil.ReadAll(block.Reader)
			}
		} else {
			// filling the content with spaces so it will create a dummy message to generate a correct mime message header
			for i := range blockContent {
				blockContent[i] = ' '
			}

		}

		switch block.Type {
		case "html":
			mimeMsg = mimeMsg.HTML(blockContent)
		case "default":
		case "text/plain":
			mimeMsg = mimeMsg.Text(blockContent)
		default:
			mimeMsg = mimeMsg.AddInline(blockContent, block.Type, "inline.dat", block.ID)
		}
	}

	for _, attachment := range m.DMsg.Catalog.Attachments {
		attachContent := make([]byte, attachment.Size)

		if realContent {
			r := attachment.Reader
			attachContent, _ = ioutil.ReadAll(r)
		}

		mimeMsg = mimeMsg.AddAttachment(attachContent, attachment.MimeType, attachment.FileName)
	}

	finalMsg, _ := mimeMsg.Build()

	b := &bytes.Buffer{}
	finalMsg.Encode(b)
	m.FullMsg = b.Bytes()

	m.Size = uint32(9999 + len(m.FullMsg))
	return m.FullMsg
}

func (m *Message) entity() (*message.Entity, error) {
	if m.DMsg == nil {
		err := m.Decrypt()
		if err != nil {
			return nil, err
		}
	}

	msgBuffer := m.Build(true)

	return message.Read(bytes.NewReader(msgBuffer))
}

func (m *Message) headerAndBody() (textproto.Header, io.Reader, error) {
	msgBuffer := m.Build(true)

	body := bufio.NewReader(bytes.NewReader(msgBuffer))
	hdr, err := textproto.ReadHeader(body)
	hdr.Add("Message-ID", "<"+m.DMsg.ID+"@bitmaelum.network/>")
	return hdr, body, err
}

func (m *Message) headerOnly() (textproto.Header, io.Reader, error) {
	msgBuffer := m.Build(true)

	body := bufio.NewReader(bytes.NewReader(msgBuffer))

	hdr, err := textproto.ReadHeader(body)
	hdr.Add("Message-ID", "<"+m.DMsg.ID+"@bitmaelum.network/>")
	return hdr, body, err
}

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
		case imap.FetchFlags:
			fetched.Flags = m.Flags
		case imap.FetchInternalDate:
			fetched.InternalDate = m.Date
		case imap.FetchRFC822Size:
			fetched.Size = m.Size
		case imap.FetchUid:
			fetched.Uid = m.Uid
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

func (m *Message) Match(seqNum uint32, c *imap.SearchCriteria) (bool, error) {
	e, _ := m.entity()
	return backendutil.Match(e, seqNum, m.Uid, m.Date, m.Flags, c)
}

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

	mimeMsg := enmime.Builder().
		From(m.DMsg.Catalog.From.Name, common.AddrToEmail(m.DMsg.Catalog.From.Address)).
		To("", common.AddrToEmail(m.DMsg.Catalog.To.Address)).
		Subject(m.DMsg.Catalog.Subject).
		Date(m.DMsg.Catalog.CreatedAt)

	m.Builder = mimeMsg

	m.Build(true)

	return nil
}
