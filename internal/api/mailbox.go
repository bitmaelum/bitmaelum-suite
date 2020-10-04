package api

import (
	"fmt"

	"github.com/bitmaelum/bitmaelum-suite/internal/message"
	"github.com/bitmaelum/bitmaelum-suite/pkg/address"
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
	Header  message.Header `json:"h"`
	Catalog []byte         `json:"c"`
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
func (api *API) GetMailboxList(addr address.Hash) (*MailboxList, error) {
	in := &MailboxList{}

	resp, statusCode, err := api.GetJSON(fmt.Sprintf("/account/%s/boxes", addr.String()), in)
	if err != nil {
		return nil, err
	}

	if statusCode < 200 || statusCode > 299 {
		return nil, getErrorFromResponse(resp)
	}

	return in, nil
}

// GetMailboxMessages returns a list of message within a specific mailbox
func (api *API) GetMailboxMessages(addr address.Hash, box string) (*MailboxMessages, error) {
	in := &MailboxMessages{}

	body, statusCode, err := api.GetJSON(fmt.Sprintf("/account/%s/box/%s", addr.String(), box), in)
	if err != nil {
		return nil, err
	}

	if statusCode < 200 || statusCode > 299 {
		return nil, getErrorFromResponse(body)
	}

	return in, nil
}
