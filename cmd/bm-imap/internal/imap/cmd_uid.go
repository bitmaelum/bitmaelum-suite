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
		info := c.DB.GetBoxInfo(c.Account, c.Box)
		for _, uid := range info.Uids {
			s += strconv.Itoa(uid) + " "
		}
		c.Write("*", s)

		c.Write(tag, "OK UID SEARCH completed")
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
