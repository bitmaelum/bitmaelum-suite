package imap

import (
	"strings"
)

func Uid(c *Conn, tag, cmd string, args []string) error {
	if strings.ToUpper(args[0]) == "FETCH" {
		return UidFetch(c, tag, cmd, args)
	}
	if strings.ToUpper(args[0]) == "SEARCH" {
		c.Write(tag, "OK UID STORE completed")
	}
	if strings.ToUpper(args[0]) == "COPY" {
		// @TODO: add copy
		c.Write(tag, "OK UID COPY completed")
	}
	if strings.ToUpper(args[0]) == "STORE" {
		c.Write(tag, "OK UID STORE completed")
	}

	return nil
}
