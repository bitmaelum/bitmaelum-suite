package imap

func Capability(c *Conn, tag, _ string, _ []string) error {
	// c.Write("*", "CAPABILITY IMAP4rev1 STARTTLS LOGINDISABLED AUTH=PLAIN")
	c.Write("*", "CAPABILITY IMAP4rev1 LOGINDISABLED AUTH=PLAIN")
	c.Write(tag, "OK CAPABILITY completed")

	return nil
}
