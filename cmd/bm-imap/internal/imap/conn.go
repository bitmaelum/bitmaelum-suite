package imap

import (
	"bufio"
	"fmt"
	"net"
	"strings"

	"github.com/bitmaelum/bitmaelum-suite/cmd/bm-imap/internal"
	"github.com/bitmaelum/bitmaelum-suite/internal/api"
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

	Box string // current selected box
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

	c.Write("*", "OK [CAPABILITY IMAP4rev1 LOGINDISABLED IDLE AUTH=PLAIN] BitMaelum IMAP Service Ready")

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
	fmt.Printf("IN > %s\n", msg)
	return msg, true
}

func (c *Conn) Write(seq, msg string) {
	fmt.Printf("OUT> %s %s\n", seq, msg)
	_, _ = fmt.Fprintf(c.C, "%s %s\r\n", seq, msg)
}

func (c *Conn) UpdateImapDB(list *api.MailboxMessages) (internal.BoxInfo, int) {
	boxInfo := c.DB.GetBoxInfo(c.Account, c.Box)

	boxInfo.Uids = make([]int, len(list.Messages))

	unseen := 0
	for i, msg := range list.Messages {
		info, err := c.DB.FetchByMessageID(c.Account, msg.ID)
		if err != nil {
			boxInfo.HighestUID++
			info, err = c.DB.Store(c.Account, c.Box, 11, boxInfo.HighestUID, msg.ID)
			if err != nil {
				return boxInfo, 0
			}
		}

		// Store this message in the box
		boxInfo.Uids[i] = info.UID

		// Check and count the unseen flags
		for _, f := range info.Flags {
			if f == "\\Unseen" {
				unseen++
			}
		}
	}

	_ = c.DB.StoreBoxInfo(c.Account, boxInfo)

	return boxInfo, unseen
}
