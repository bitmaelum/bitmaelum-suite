package imap

func Noop(c *Conn, tag, _ string, _ []string) error {
	c.Write(tag, "OK NOOP completed")

	return nil
}
