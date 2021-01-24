package imap

func Expunge(c *Conn, tag, _ string, _ []string) error {
	c.Write(tag, "OK EXPUNGE completed")

	return nil
}
