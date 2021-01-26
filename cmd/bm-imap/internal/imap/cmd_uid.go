package imap

import (
	"strconv"
	"strings"
)

func Uid(c *Conn, tag, cmd string, args []string) error {
	if strings.ToUpper(args[0]) == "FETCH" {
		return UidFetch(c, tag, cmd, args)
	}
	if strings.ToUpper(args[0]) == "SEARCH" {
		s := ""
		for _, idx := range c.Index {
			s += strconv.Itoa(idx.UID) + " "
		}
		c.Write("* SEARCH", s)

		c.Write(tag, "OK UID SEARCH completed")
	}


	if strings.ToUpper(args[0]) == "COPY" {
		// @TODO: add copy
		c.Write(tag, "OK UID COPY completed")
	}
	if strings.ToUpper(args[0]) == "STORE" {
		// @TODO: add store
		c.Write(tag, "OK UID STORE completed")
	}

	return nil
}
