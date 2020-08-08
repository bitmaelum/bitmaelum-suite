package api

import (
	"fmt"
	"github.com/bitmaelum/bitmaelum-suite/internal/message"
	"github.com/bitmaelum/bitmaelum-suite/pkg/address"
)

// Message is a standard structure that returns a message header + catalog
type Message struct {
	ID      string         `json:"id"`
	Header  message.Header `json:"h"`
	Catalog []byte         `json:"c"`
}

// GetMessage retrieves a message header + catalog from a message box
func (api *API) GetMessage(addr address.HashAddress, box, messageID string) (*Message, error) {
	in := &Message{}

	statusCode, err := api.GetJSON(fmt.Sprintf("/account/%s/box/%s/%s", addr.String(), box, messageID), in)
	if err != nil {
		return nil, err
	}

	if statusCode < 200 || statusCode > 299 {
		return nil, errNoSuccess
	}

	return in, nil
}

// GetMessageBlock retrieves a message block
func (api *API) GetMessageBlock(addr address.HashAddress, box, messageID, blockID string) ([]byte, error) {
	body, statusCode, err := api.Get(fmt.Sprintf("/account/%s/box/%s/%s/%s", addr.String(), box, messageID, blockID))
	if err != nil {
		return nil, err
	}

	if statusCode < 200 || statusCode > 299 {
		return nil, errNoSuccess
	}

	return body, nil
}

// GetMessageAttachment retrieves a message attachment
func (api *API) GetMessageAttachment(addr address.HashAddress, box, messageID, attachmentID string) ([]byte, error) {
	body, statusCode, err := api.Get(fmt.Sprintf("/account/%s/box/%s/%s/%s", addr.String(), box, messageID, attachmentID))
	if err != nil {
		return nil, err
	}

	if statusCode < 200 || statusCode > 299 {
		return nil, errNoSuccess
	}

	return body, nil
}
