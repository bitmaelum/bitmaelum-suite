package imap

import (
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"github.com/bitmaelum/bitmaelum-suite/cmd/bm-imap/internal"
	"github.com/bitmaelum/bitmaelum-suite/pkg/address"
)

func Authenticate(c *Conn, tag, _ string, _ []string) error {
	c.Write("+", "")

	line, ok := c.Read()
	if !ok {
		c.Write(tag, "BAD Cannot read data")
		return errors.New("error")
	}

	b, err := base64.StdEncoding.DecodeString(line)
	if err != nil {
		c.Write(tag, "BAD Cannot read data")
		return errors.New("error")
	}
	creds := strings.Split(string(b), "\x00")
	c.Account = creds[1] + "!"

	addr, err := address.NewAddress(c.Account)
	if err != nil {
		c.Write(tag, "NO incorrect address format specified")
		return errors.New("error")
	}
	if !c.Vault.HasAccount(*addr) {
		c.Write(tag, "NO account not found in vault")
		return errors.New("error")
	}

	_, c.Info, c.Client, err = internal.GetClientAndInfo(c.Account)
	if err != nil {
		return errors.New("error")
	}
	c.State = StateAuthenticated

	fmt.Println("Authenticated as ", c.Account)
	c.Write(tag, "OK Authenticated")

	return nil
}
