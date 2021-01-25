package imap

import (
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"github.com/bitmaelum/bitmaelum-suite/cmd/bm-imap/internal"
	"github.com/bitmaelum/bitmaelum-suite/pkg/address"
)

func Login(c *Conn, tag, _ string, args []string) error {
	err := login(c, args[0], args[1])
	if err != nil {
		c.Write(tag, err.Error())

		return err
	}

	fmt.Println("Logged in as ", c.Account)
	c.Write(tag, "OK Login")

	return nil
}

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

	err = login(c, creds[0], creds[1])
	if err != nil {
		c.Write(tag, err.Error())
	}

	fmt.Println("Authenticated as ", c.Account)
	c.Write(tag, "OK Authenticated")

	return nil
}

func login(c *Conn, user, _ string) error {
	c.Account = user + "!"

	addr, err := address.NewAddress(c.Account)
	if err != nil {
		return errors.New("NO incorrect address format specified")
	}
	if !c.Vault.HasAccount(*addr) {
		return errors.New("NO account not found in vault")
	}

	_, c.Info, c.Client, err = internal.GetClientAndInfo(c.Account)
	if err != nil {
		return errors.New("NO error")
	}
	c.State = StateAuthenticated

	return nil
}
