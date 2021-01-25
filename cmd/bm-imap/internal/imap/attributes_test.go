package imap

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseAttributes(t *testing.T) {
	s := "(UID RFC822.SIZE FLAGS BODY.PEEK[HEADER.FIELDS (From To Cc Bcc Resent-Message-ID Subject Date Message-ID Priority X-Priority References Newsgroups In-Reply-To Content-Type Reply-To List-Unsubscribe Received Delivery-Date)])"
	ret := ParseAttributes(s)
	assert.Len(t, ret, 4)
	assert.Len(t, ret[3].Headers, 18)
	assert.True(t, ret[3].Peek)
	assert.Equal(t, ret[3].Name, "BODY")
	assert.False(t, ret[2].Peek)

	assert.Equal(t, ret[0].ToString(), "UID")
	assert.Equal(t, ret[1].ToString(), "RFC822.SIZE")
	assert.Equal(t, ret[2].ToString(), "FLAGS")
	assert.Equal(t, ret[3].ToString(), "BODY[HEADER.FIELDS (\"from\" \"to\" \"cc\" \"bcc\" \"resent-message-id\" \"subject\" \"date\" \"message-id\" \"priority\" \"x-priority\" \"references\" \"newsgroups\" \"in-reply-to\" \"content-type\" \"reply-to\" \"list-unsubscribe\" \"received\" \"delivery-date\" )]")
	assert.Equal(t, 0, ret[2].MinRange)
	assert.Equal(t, 0, ret[2].MaxRange)


	s = "(FULL)"
	ret = ParseAttributes(s)
	assert.Len(t, ret, 5)
	assert.Equal(t, ret[0].ToString(), "FLAGS")
	assert.Equal(t, ret[1].ToString(), "INTERNALDATE")
	assert.Equal(t, ret[2].ToString(), "RFC822.SIZE")
	assert.Equal(t, ret[3].ToString(), "ENVELOPE")
	assert.Equal(t, ret[4].ToString(), "BODY")
	assert.Equal(t, 0, ret[2].MinRange)
	assert.Equal(t, 0, ret[2].MaxRange)


	s = "(BODY[HEADER.FIELDS.NOT (From)])"
	ret = ParseAttributes(s)
	assert.Len(t, ret, 1)
	assert.Len(t, ret[0].Headers, 1)
	assert.Equal(t, ret[0].Name, "BODY")
	assert.False(t, ret[0].Peek)
	assert.True(t, ret[0].Not)
	assert.Equal(t, ret[0].ToString(), "BODY[HEADER.FIELDS.NOT (\"from\" )]")
	assert.Equal(t, 0, ret[0].MinRange)
	assert.Equal(t, 0, ret[0].MaxRange)


	s = "(BODY.PEEK[TEXT]<0.8192>)"
	ret = ParseAttributes(s)
	assert.Len(t, ret, 1)
	assert.Len(t, ret[0].Headers, 0)
	assert.Equal(t, ret[0].Section, "TEXT")
	assert.Equal(t, ret[0].Name, "BODY")
	assert.True(t, ret[0].Peek)
	assert.False(t, ret[0].Not)
	assert.Equal(t, 0, ret[0].MinRange)
	assert.Equal(t, 8192, ret[0].MaxRange)
	assert.Equal(t, ret[0].ToString(), "BODY[TEXT]<0.8192>")
}
