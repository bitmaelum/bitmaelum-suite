package api

import (
	"fmt"
	"github.com/bitmaelum/bitmaelum-suite/internal/message"
	"github.com/bitmaelum/bitmaelum-suite/pkg/address"
)

type MailboxListBoxIn struct {
	ID    int `json:"id"`
	Total int `json:"total"`
}
type MailboxListIn struct {
	Meta struct {
		Total    int `json:"total"`
		Returned int `json:"returned"`
	} `json:"meta"`
	Boxes []MailboxListBoxIn `json:"boxes"`
}

type MailboxMessagesMessageIn struct {
	Header  message.Header `json:"h"`
	Catalog []byte `json:"c"`
}
type MailboxMessagesIn struct {
	Meta struct {
		Total    int `json:"total"`
		Returned int `json:"returned"`
		Offset   int `json:"offset"`
		Limit    int `json:"limit"`
	} `json:"meta"`
	Messages []MailboxMessagesMessageIn `json:"messages"`
}

func (api *API) GetMailboxList(addr address.HashAddress) (*MailboxListIn, error) {
	in := &MailboxListIn{}

	statusCode, err := api.GetJSON(fmt.Sprintf("/account/%s/boxes", addr.String()), in)
	if err != nil {
		return nil, err
	}

	if statusCode < 200 || statusCode > 299 {
		return nil, errNoSuccess
	}

	return in, nil
}

func (api *API) GetMailboxMessages(addr address.HashAddress, box string) (*MailboxMessagesIn, error) {
	in := &MailboxMessagesIn{}

	statusCode, err := api.GetJSON(fmt.Sprintf("/account/%s/box/%s", addr.String(), box), in)
	if err != nil {
		return nil, err
	}

	if statusCode < 200 || statusCode > 299 {
		return nil, errNoSuccess
	}

	return in, nil
}
