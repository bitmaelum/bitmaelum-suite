// Copyright (c) 2020 BitMaelum Authors
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of
// this software and associated documentation files (the "Software"), to deal in
// the Software without restriction, including without limitation the rights to
// use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
// the Software, and to permit persons to whom the Software is furnished to do so,
// subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
// FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
// COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
// IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
// CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package api

import (
	"fmt"
	"time"

	"github.com/bitmaelum/bitmaelum-suite/internal/message"
	"github.com/bitmaelum/bitmaelum-suite/pkg/hash"
)

// MailboxListBox is a structure that holds a given mailbox and the total messages inside
type MailboxListBox struct {
	ID    int `json:"id"`
	Total int `json:"total"`
}

// MailboxList is a list of mailboxes
type MailboxList struct {
	Meta struct {
		Total    int `json:"total"`
		Returned int `json:"returned"`
	} `json:"meta"`
	Boxes []MailboxListBox `json:"boxes"`
}

// MailboxMessagesMessage is a message (header + catalog) within a mailbox
type MailboxMessagesMessage struct {
	ID      string         `json:"id"`
	Header  message.Header `json:"header"`
	Catalog []byte         `json:"catalog"`
}

// MailboxMessages returns a list of mailbox messages
type MailboxMessages struct {
	Meta struct {
		Total    int `json:"total"`
		Returned int `json:"returned"`
		Offset   int `json:"offset"`
		Limit    int `json:"limit"`
	} `json:"meta"`
	Messages []MailboxMessagesMessage `json:"messages"`
}

// GetMailboxList returns a list of mailboxes
func (api *API) GetMailboxList(addr hash.Hash) (*MailboxList, error) {
	in := &MailboxList{}

	resp, statusCode, err := api.GetJSON(fmt.Sprintf("/account/%s/boxes", addr.String()), in)
	if err != nil {
		return nil, err
	}

	if statusCode < 200 || statusCode > 299 {
		return nil, GetErrorFromResponse(resp)
	}

	return in, nil
}

// GetMailboxMessages returns a list of message within a specific mailbox
func (api *API) GetMailboxMessages(addr hash.Hash, box string, since time.Time) (*MailboxMessages, error) {
	in := &MailboxMessages{}

	// Add since query string if needed
	qs := ""
	if !since.IsZero() {
		qs = fmt.Sprintf("since=%d", since.Unix())
	}

	body, statusCode, err := api.GetJSON(fmt.Sprintf("/account/%s/box/%s?%s", addr.String(), box, qs), in)
	if err != nil {
		return nil, err
	}

	if statusCode < 200 || statusCode > 299 {
		return nil, GetErrorFromResponse(body)
	}

	return in, nil
}
