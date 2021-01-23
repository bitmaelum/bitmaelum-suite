package internal

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/bitmaelum/bitmaelum-suite/internal/api"
	"github.com/bitmaelum/bitmaelum-suite/internal/vault"
	"github.com/bitmaelum/bitmaelum-suite/pkg/address"
)

type State int

const (
	StateUnauthenticated State = iota
	StateAuthenticated
)

type Conn struct {
	c       net.Conn       // Actual connection
	scanner *bufio.Scanner // Scanner to read data from connection

	state   State        // Current state of the connection
	vault   *vault.Vault // Opened Vault
	account string       // bitmaelum account connected to the connection
	info    *vault.AccountInfo
	client  *api.API

	box string // current selected box
}

func NewConn(c net.Conn, v *vault.Vault) Conn {
	return Conn{
		c:       c,
		scanner: bufio.NewScanner(c),
		state:   StateUnauthenticated,
		vault:   v,
	}
}

func (c *Conn) Handle() {
	defer c.Close()

	c.Write("*", "OK IMAP4rev1 Service Ready")

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

		fmt.Println("TAG: ", tag)
		fmt.Println("CMD: ", cmd)
		fmt.Println("ARGS: ", args)

		switch cmd {
		case "CAPABILITY":
			// c.Write("*", "CAPABILITY IMAP4rev1 STARTTLS LOGINDISABLED AUTH=PLAIN")
			c.Write("*", "CAPABILITY IMAP4rev1 LOGINDISABLED AUTH=PLAIN")
			c.Write(tag, "OK CAPABILITY completed")
		case "LOGOUT":
			return
		case "AUTHENTICATE":
			c.Write("+", "")

			line, ok := c.Read()
			if !ok {
				c.Write(tag, "BAD Cannot read data")
				continue
			}

			b, err := base64.StdEncoding.DecodeString(line)
			if err != nil {
				c.Write(tag, "BAD Cannot read data")
				continue
			}
			creds := strings.Split(string(b), "\x00")
			c.account = creds[1] + "!"

			addr, err := address.NewAddress(c.account)
			if err != nil {
				c.Write(tag, "NO incorrect address format specified")
				continue
			}
			if !c.vault.HasAccount(*addr) {
				c.Write(tag, "NO account not found in vault")
				continue
			}

			_, c.info, c.client, err = GetClientAndInfo(c.account)
			if err != nil {
				panic(err)
			}
			c.state = StateAuthenticated

			fmt.Println("Authenticated as ", c.account)
			c.Write(tag, "OK Authenticated")
		case "LIST":
			if args[0] == "\"\"" {
				c.Write("*", "LIST (\\Noselect) \"/\" \"\"")
			}
			if args[1] == "\"*\"" {
				mbl, err := c.client.GetMailboxList(c.info.Address.Hash())
				if err != nil {
					panic(err)
				}

				for _, box := range mbl.Boxes {
					boxName := strconv.Itoa(box.ID)
					switch box.ID {
					case 1:
						boxName = "Inbox"
					case 2:
						boxName = "Sent"
					case 3:
						boxName = "Trash"
					}
					c.Write("*", "LIST \"/\" \""+boxName+"\"")
				}
			}
			c.Write(tag, "OK LIST completed")
		case "LSUB":
			mbl, err := c.client.GetMailboxList(c.info.Address.Hash())
			if err != nil {
				panic(err)
			}

			for _, box := range mbl.Boxes {
				boxName := strconv.Itoa(box.ID)
				switch box.ID {
				case 1:
					boxName = "Inbox"
				case 2:
					boxName = "Sent"
				case 3:
					boxName = "Trash"
				}
				c.Write("*", "LSUB () \"/\" \""+boxName+"\"")
			}
			c.Write(tag, "OK LSUB completed")
		case "SELECT":
			c.box = args[0]
			switch (strings.ToUpper(args[0])) {
			case "INBOX":
				c.box = "1"
			case "SENT":
				c.box = "2"
			case "TRASH":
				c.box = "3"
			}

			msg, err := c.client.GetMailboxMessages(c.info.Address.Hash(), c.box, time.Time{})
			if err != nil {
				panic(err)
			}


			c.Write("*", fmt.Sprintf("%d EXISTS", len(msg.Messages)))
			c.Write("*", fmt.Sprintf("%d RECENT", len(msg.Messages)))
			c.Write("*", fmt.Sprintf("OK [UNSEEN %d]", len(msg.Messages)))
			c.Write("*", fmt.Sprintf("OK [UIDNEXT %d]", 123))
			c.Write("*", fmt.Sprintf("OK [UIDVALIDITY %d]", 250))
			c.Write("*", "FLAGS (\\Answered \\Flagged \\Deleted \\Seen \\Draft)\r\n")

			c.Write(tag, "OK [READ-WRITE] SELECT completed")
		case "NOOP":
			c.Write(tag, "OK NOOP completed")
		}
	}
}

func (c *Conn) Close() {
	fmt.Printf("Closing connection\n")
	_ = c.c.Close()
}

func (c *Conn) Read() (string, bool) {
	ok := c.scanner.Scan()
	if !ok {
		return "", false
	}

	msg := c.scanner.Text()
	fmt.Printf("IN > %s\n", msg)
	return msg, true
}

func (c *Conn) Write(seq, msg string) {
	fmt.Printf("OUT> %s %s\n", seq, msg)

	_, _ = fmt.Fprintf(c.c, "%s %s\n", seq, msg)
}
