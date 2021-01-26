package imap

import (
	"bufio"
	"errors"
	"fmt"
	"hash/crc32"
	"net"
	"strings"
	"time"

	"github.com/bitmaelum/bitmaelum-suite/cmd/bm-imap/internal"
	"github.com/bitmaelum/bitmaelum-suite/internal/api"
	"github.com/bitmaelum/bitmaelum-suite/internal/message"
	"github.com/bitmaelum/bitmaelum-suite/internal/vault"
	"github.com/mitchellh/go-homedir"
)

type State int

const (
	StateUnauthenticated State = iota
	StateAuthenticated
)

type Conn struct {
	C       net.Conn       // Actual connection
	Scanner *bufio.Scanner // Scanner to read data from connection

	State   State        // Current state of the connection
	Vault   *vault.Vault // Opened Vault
	Account string       // bitmaelum account connected to the connection
	Info    *vault.AccountInfo
	Client  *api.API

	DB *internal.BoltRepo

	Box     string           // current selected box
	Index   []MessageIndex   // message index list for current selected box
	BoxInfo internal.BoxInfo // BoxInfo
}

type MessageIndex struct {
	MessageID string
	UID       int
}

func NewConn(c net.Conn, v *vault.Vault) Conn {
	p, _ := homedir.Dir()

	return Conn{
		C:       c,
		Scanner: bufio.NewScanner(c),
		State:   StateUnauthenticated,
		Vault:   v,
		DB:      internal.NewBolt(p),
	}
}

func (c *Conn) Handle() {
	defer c.Close()

	c.Write("*", "OK [CAPABILITY IMAP4rev1 AUTH=PLAIN] BitMaelum IMAP Service Ready")

	for {
		line, ok := c.Read()
		if !ok {
			fmt.Printf("OK is not ok: '%s'\n", line)
			return
		}

		parts := strings.Split(line, " ")
		tag := parts[0]
		cmd := strings.ToUpper(parts[1])
		var args []string
		if len(parts) > 2 {
			args = parts[2:]
			for i := range args {
				args[i] = strings.Trim(args[i], "\"")
			}
		}

		var err error

		switch cmd {
		case "CAPABILITY":
			err = Capability(c, tag, cmd, args)
		case "AUTHENTICATE":
			err = Authenticate(c, tag, cmd, args)
		case "LOGIN":
			err = Login(c, tag, cmd, args)
		case "LIST":
			err = List(c, tag, cmd, args)
		case "LSUB":
			err = Lsub(c, tag, cmd, args)
		case "SELECT":
			err = Select(c, tag, cmd, args)
		case "UID":
			err = Uid(c, tag, cmd, args)
		case "NOOP":
			err = Noop(c, tag, cmd, args)
		case "EXPUNGE":
			err = Expunge(c, tag, cmd, args)
		case "STATUS":
			err = Status(c, tag, cmd, args)
		case "LOGOUT":
			return
		}

		if err != nil {
			return
		}
	}
}

func (c *Conn) Close() {
	fmt.Printf("Closing connection\n")
	_ = c.C.Close()
}

func (c *Conn) Read() (string, bool) {
	ok := c.Scanner.Scan()
	if !ok {
		return "", false
	}

	msg := c.Scanner.Text()
	fmt.Printf("IN [%s]> %s\n", c.C.RemoteAddr().String(), msg)
	return msg, true
}

func (c *Conn) Write(seq, msg string) {
	fmt.Printf("OUT[%s]> %s %s\n", c.C.RemoteAddr().String(), seq, msg)

	_, _ = fmt.Fprintf(c.C, "%s %s\r\n", seq, msg)
}

func (c *Conn) GetUIDForMessage(msgID string) int {
	boxInfo := c.DB.GetBoxInfo(c.Account, c.Box)

	info, err := c.DB.FetchByMessageID(c.Account, msgID)
	if err == nil {
		return info.UID
	}

	boxInfo.HighestUID++
	info, err = c.DB.Store(c.Account, c.Box, int(crc32.ChecksumIEEE([]byte(c.Box))), boxInfo.HighestUID, msgID)
	if err != nil {
		return boxInfo.HighestUID
	}

	_ = c.DB.StoreBoxInfo(c.Account, boxInfo)
	return info.UID
}

func (c *Conn) updateMessageIndex() {
	var index []MessageIndex

	// fetch a list of all messages from the given box on the server
	msgList, err := c.Client.GetMailboxMessages(c.Info.Address.Hash(), c.Box, time.Time{})
	if err != nil {
		return
	}

	for _, msg := range msgList.Messages {
		// Decrypt message if possible
		em := message.EncryptedMessage{
			ID:      msg.ID,
			Header:  &msg.Header,
			Catalog: msg.Catalog,

			GenerateBlockReader:      c.Client.GenerateAPIBlockReader(c.Info.Address.Hash()),
			GenerateAttachmentReader: c.Client.GenerateAPIAttachmentReader(c.Info.Address.Hash()),
		}

		_, err := em.Decrypt(c.Info.GetActiveKey().PrivKey)
		if err != nil {
			continue
		}

		index = append(index, MessageIndex{
			MessageID: msg.ID,
			UID:       c.GetUIDForMessage(msg.ID),
		})
	}

	c.Index = index
}

func (c *Conn) FindByUID(uid int) (*MessageIndex, error) {
	for _, idx := range c.Index {
		if idx.UID == uid {
			return &idx, nil
		}
	}

	return nil, errors.New("not found")
}

func (c *Conn) FindByMsgID(msgID string) (*MessageIndex, error) {
	for _, idx := range c.Index {
		if idx.MessageID == msgID {
			return &idx, nil
		}
	}

	return nil, errors.New("not found")
}

func (c *Conn) FindBySeq(seq int) (*MessageIndex, error) {
	return &c.Index[seq], nil
}


func (c *Conn) ChangeBox(box string) {
	c.Box = box
	c.BoxInfo = c.DB.GetBoxInfo(c.Account, c.Box)

	c.updateMessageIndex()
}
